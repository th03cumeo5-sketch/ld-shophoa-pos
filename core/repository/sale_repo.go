package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"shophoa/core/models"
)

type SaleRepository struct {
	db *sql.DB
}

func NewSaleRepository(db *sql.DB) *SaleRepository {
	return &SaleRepository{db: db}
}

// formatSQLDateTime format thời gian chuẩn tương thích SQL Server datetime (tránh lỗi 22008 fractional second overflow)
func formatSQLDateTime(t time.Time) string {
	if t.IsZero() {
		t = time.Now()
	}
	return t.Format("2006-01-02 15:04:05")
}

// AllocateRetailBillNumber gọi Stored Procedure POS_ALLOC_BILL_NO để cấp số hóa đơn bán lẻ
func (r *SaleRepository) AllocateRetailBillNumber(ctx context.Context, makho, quay string, ngay time.Time) (string, error) {
	query := `
		SET NOCOUNT ON;
		DECLARE @sct nvarchar(20);
		EXEC [dbo].[POS_ALLOC_BILL_NO] ?, ?, ?, @sct OUTPUT;
		SELECT @sct;
	`
	var sct string
	err := r.db.QueryRowContext(ctx, query, makho, quay, formatSQLDateTime(ngay)).Scan(&sct)
	if err != nil {
		return "", fmt.Errorf("POS_ALLOC_BILL_NO failed: %w", err)
	}
	return sct, nil
}

// AllocateSCT gọi Stored Procedure ALLOC_SCT để cấp số chứng từ hệ thống (NMUA, XBAN, THU, CHI, NTBL...)
func (r *SaleRepository) AllocateSCT(ctx context.Context, loai, makho string, ngay time.Time) (string, error) {
	query := `
		SET NOCOUNT ON;
		DECLARE @sct nvarchar(50);
		EXEC [dbo].[ALLOC_SCT] @sct OUTPUT, ?, ?, ?;
		SELECT @sct;
	`
	var sct string
	err := r.db.QueryRowContext(ctx, query, loai, formatSQLDateTime(ngay), makho).Scan(&sct)
	if err != nil {
		return "", fmt.Errorf("ALLOC_SCT failed: %w", err)
	}
	return sct, nil
}

// SaveDraftSale lưu hóa đơn bán hàng đang soạn thảo hoặc treo đơn (F10) vào cụm bảng BANLE_TAM & BANLE_TAM_CT
func (r *SaleRepository) SaveDraftSale(ctx context.Context, header *models.BanleTam, items []models.BanleTamCT) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	if header.Khoa == uuid.Nil {
		header.Khoa = uuid.New()
	}
	if header.LCT == "" {
		header.LCT = "BLE"
	}

	khoaStr := header.Khoa.String()
	ngayStr := formatSQLDateTime(header.Ngay)
	nowStr := formatSQLDateTime(time.Now())

	// 1. Dọn dẹp dữ liệu cũ của phiếu này trong bảng tạm (nếu là cập nhật)
	delCT := "DELETE FROM BANLE_TAM_CT WHERE KHOA = ?"
	if _, err := tx.ExecContext(ctx, delCT, khoaStr); err != nil {
		return fmt.Errorf("failed to clean BANLE_TAM_CT: %w", err)
	}
	delHeader := "DELETE FROM BANLE_TAM WHERE KHOA = ?"
	if _, err := tx.ExecContext(ctx, delHeader, khoaStr); err != nil {
		return fmt.Errorf("failed to clean BANLE_TAM: %w", err)
	}

	// 2. Chèn Header vào BANLE_TAM
	insertHeader := `
		INSERT INTO BANLE_TAM (
			KHOA, LCT, NGAY, SCT, LOC, MAKHO, QUAY, CA, THUNGAN, MADT, MAVIP,
			SOLUONG, SOTIEN, CHIETKHAU_MH, TL_CK, CHIETKHAU, THANHTIEN_CHUA_VAT,
			TIEN_THUE, TIEN_THUE_CL, THANHTOAN, PTTT, TTOAN1, TTOAN2, CHUATHOI,
			PRINTED, PRINT_NO, DGIAI, CREATE_BY, CREATE_DATE, UPDATE_BY, UPDATE_DATE
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?
		)
	`
	// Xử lý các trường có Foreign Key: nếu rỗng thì truyền nil (NULL trong SQL Server)
	var madtVal interface{} = nil
	if header.MaDT != "" {
		madtVal = header.MaDT
	}
	var mavipVal interface{} = nil
	if header.MaVIP != "" {
		mavipVal = header.MaVIP
	}

	_, err = tx.ExecContext(ctx, insertHeader,
		khoaStr,
		header.LCT,
		ngayStr,
		header.SCT,
		header.LOC,
		header.MaKho,
		header.Quay,
		header.Ca,
		header.ThuNgan,
		madtVal,
		mavipVal,
		header.SoLuong,
		header.SoTien,
		header.ChietKhauMH,
		header.TL_CK,
		header.ChietKhau,
		header.ThanhTienChuaVAT,
		header.TienThue,
		header.TienThueCL,
		header.ThanhToan,
		header.PTTT,
		header.TToan1,
		header.TToan2,
		header.ChuaThoi,
		header.Printed,
		header.PrintNo,
		header.DGIAI,
		header.CreateBy,
		nowStr,
		header.UpdateBy,
		nowStr,
	)
	if err != nil {
		return fmt.Errorf("failed to insert BANLE_TAM: %w", err)
	}

	// 3. Chèn các dòng chi tiết vào BANLE_TAM_CT
	insertItem := `
		INSERT INTO BANLE_TAM_CT (
			KHOACT, KHOA, STT, MAVT, TENVT, DVT, SOLUONG, DONGIA, DONGIA_SAUCK,
			SOTIEN, TL_CK, TL_CK2, TL_CK3, TL_CK4, TL_CK5, CHIETKHAU,
			LOAITHUE, THUE_SUAT, TIEN_THUE, TIEN_THUE_CL, THANHTIEN_CHUA_VAT,
			THANHTOAN, GHICHU, LOC
		) VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?,
			?, ?, ?
		)
	`
	for i := range items {
		item := &items[i]
		if item.KhoaCT == uuid.Nil {
			item.KhoaCT = uuid.New()
		}
		item.Khoa = header.Khoa
		item.STT = i + 1

		_, err := tx.ExecContext(ctx, insertItem,
			item.KhoaCT.String(),
			khoaStr,
			item.STT,
			item.MaVT,
			item.TenVT,
			item.DVT,
			item.SoLuong,
			item.DonGia,
			item.DonGiaSauCK,
			item.SoTien,
			item.TL_CK,
			item.TL_CK2,
			item.TL_CK3,
			item.TL_CK4,
			item.TL_CK5,
			item.ChietKhau,
			item.LoaiThue,
			item.ThueSuat,
			item.TienThue,
			item.TienThueCL,
			item.ThanhTienChuaVAT,
			item.ThanhToan,
			item.GhiChu,
			header.LOC,
		)
		if err != nil {
			return fmt.Errorf("failed to insert BANLE_TAM_CT row %d: %w", i+1, err)
		}
	}

	return tx.Commit()
}

// CompleteSalePayment gọi Stored Procedure POS_PAID(@KHOA)
func (r *SaleRepository) CompleteSalePayment(ctx context.Context, khoa uuid.UUID) error {
	query := "EXEC [dbo].[POS_PAID] ?"
	_, err := r.db.ExecContext(ctx, query, khoa.String())
	if err != nil {
		return fmt.Errorf("executing POS_PAID failed for khoa %s: %w", khoa.String(), err)
	}
	return nil
}

// ListDraftSales lấy danh sách các đơn hàng lưu tạm (F10) tại một quầy
func (r *SaleRepository) ListDraftSales(ctx context.Context, makho, quay string) ([]models.BanleTam, error) {
	query := `
		SELECT KHOA, SCT, NGAY, ISNULL(MADT, ''), ISNULL(SOLUONG, 0), ISNULL(THANHTOAN, 0), ISNULL(DGIAI, '')
		FROM BANLE_TAM
		WHERE MAKHO = ? AND QUAY = ?
		ORDER BY NGAY DESC
	`
	rows, err := r.db.QueryContext(ctx, query, makho, quay)
	if err != nil {
		return nil, fmt.Errorf("failed to list draft sales: %w", err)
	}
	defer rows.Close()

	var list []models.BanleTam
	for rows.Next() {
		var b models.BanleTam
		var khoaStr string
		var ngayStr string
		if err := rows.Scan(&khoaStr, &b.SCT, &ngayStr, &b.MaDT, &b.SoLuong, &b.ThanhToan, &b.DGIAI); err != nil {
			return nil, fmt.Errorf("failed to scan draft sale: %w", err)
		}
		b.Khoa, _ = uuid.Parse(khoaStr)
		b.Ngay, _ = time.Parse("2006-01-02 15:04:05", ngayStr)
		list = append(list, b)
	}
	return list, nil
}
