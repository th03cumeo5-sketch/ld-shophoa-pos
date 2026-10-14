package repository

import (
	"context"
	"database/sql"
	"fmt"

	"shophoa/core/models"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// SearchCustomers tìm kiếm khách hàng theo SĐT, Tên hoặc Mã khách hàng
func (r *CustomerRepository) SearchCustomers(ctx context.Context, keyword string, limit int) ([]models.Customer, error) {
	if limit <= 0 {
		limit = 10
	}

	query := `
		SELECT TOP (?)
			MADT, ISNULL(TENDT, ''), ISNULL(DTHOAI, ''), 
			ISNULL(DCHI, ''), ISNULL(MST, ''), ISNULL(LOAI_VIP, ''), ISNULL(SODU, 0)
		FROM DM_KH_NCC
		WHERE (NGUNG_LIENLAC = 0 OR NGUNG_LIENLAC IS NULL)
		  AND (MADT LIKE ? OR TENDT LIKE ? OR DTHOAI LIKE ?)
		ORDER BY TENDT
	`

	kw := "%" + keyword + "%"
	rows, err := r.db.QueryContext(ctx, query, limit, kw, kw, kw)
	if err != nil {
		return nil, fmt.Errorf("failed to search customers: %w", err)
	}
	defer rows.Close()

	var list []models.Customer
	for rows.Next() {
		var c models.Customer
		if err := rows.Scan(&c.MaDT, &c.TenDT, &c.DienThoai, &c.DiaChi, &c.MST, &c.LoaiVIP, &c.SoDu); err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}
		list = append(list, c)
	}

	return list, nil
}

// GetByMaDT lấy thông tin một khách hàng theo mã đối tượng
func (r *CustomerRepository) GetByMaDT(ctx context.Context, madt string) (*models.Customer, error) {
	query := `
		SELECT MADT, ISNULL(TENDT, ''), ISNULL(DTHOAI, ''), 
		       ISNULL(DCHI, ''), ISNULL(MST, ''), ISNULL(LOAI_VIP, ''), ISNULL(SODU, 0)
		FROM DM_KH_NCC
		WHERE MADT = ?
	`

	var c models.Customer
	err := r.db.QueryRowContext(ctx, query, madt).Scan(
		&c.MaDT, &c.TenDT, &c.DienThoai, &c.DiaChi, &c.MST, &c.LoaiVIP, &c.SoDu,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get customer %s: %w", madt, err)
	}

	return &c, nil
}
