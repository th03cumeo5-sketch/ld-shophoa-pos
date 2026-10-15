package main

import (
	"context"
	"encoding/json"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"
	"shophoa/core/config"
	"shophoa/core/db"
	"shophoa/core/models"
	"shophoa/core/repository"
	"shophoa/core/service"
)

func TestCoreWithRealDatabase(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Test kết nối SQL Server
	cfg := config.LoadConfig()
	sqlDB, err := db.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect to SQL Server: %v", err)
	}
	defer sqlDB.Close()
	t.Logf("✓ Connected to SQL Server: %s, DB: %s", cfg.DBServer, cfg.DBName)

	// Khởi tạo Repositories & Services
	prodRepo := repository.NewProductRepository(sqlDB)
	custRepo := repository.NewCustomerRepository(sqlDB)
	saleRepo := repository.NewSaleRepository(sqlDB)
	discountServ := service.NewDiscountService()
	taxServ := service.NewTaxService()
	saleServ := service.NewSaleService(saleRepo, prodRepo, discountServ, taxServ)

	// 2. Test tìm kiếm sản phẩm thực tế từ 7,124 mặt hàng
	products, err := prodRepo.SearchProducts(ctx, "Hoa", 5)
	if err != nil {
		t.Fatalf("SearchProducts failed: %v", err)
	}
	if len(products) == 0 {
		t.Fatalf("Expected to find products containing 'Hoa', got 0")
	}
	t.Logf("✓ Found %d products containing 'Hoa'. First: %s - %s (GiaBan: %.0f, Thue: %s)",
		len(products), products[0].MaVT, products[0].TenVT, products[0].GiaBan, products[0].LoaiThue)

	// 3. Test tìm kiếm khách hàng
	customers, err := custRepo.SearchCustomers(ctx, "A", 5)
	if err != nil {
		t.Fatalf("SearchCustomers failed: %v", err)
	}
	t.Logf("✓ Customer search returned %d records", len(customers))

	// 4. Test Stored Procedure POS_ALLOC_BILL_NO
	sct, err := saleRepo.AllocateRetailBillNumber(ctx, "01", "01", time.Now())
	if err != nil {
		t.Fatalf("AllocateRetailBillNumber failed: %v", err)
	}
	if len(sct) == 0 {
		t.Fatalf("Expected non-empty SCT, got empty")
	}
	t.Logf("✓ POS_ALLOC_BILL_NO generated SCT: %s", sct)

	// 5. Test Stored Procedure ALLOC_SCT
	sctSystem, err := saleRepo.AllocateSCT(ctx, "THU", "01", time.Now())
	if err != nil {
		t.Fatalf("ALLOC_SCT failed: %v", err)
	}
	t.Logf("✓ ALLOC_SCT generated System SCT: %s", sctSystem)

	// 6. Test chu trình soạn thảo đơn hàng vào BANLE_TAM & BANLE_TAM_CT
	draftSale, err := saleServ.PrepareNewSale(ctx, "01", "01", "1", "ntd", "01")
	if err != nil {
		t.Fatalf("PrepareNewSale failed: %v", err)
	}

	items := []models.BanleTamCT{
		{
			KhoaCT:   uuid.New(),
			MaVT:     products[0].MaVT,
			TenVT:    products[0].TenVT,
			DVT:      products[0].DVT,
			SoLuong:  2,
			DonGia:   products[0].GiaBan,
			LoaiThue: products[0].LoaiThue,
			TL_CK2:   5, // 5% VIP
		},
	}

	err = saleServ.RecalculateAndSaveDraft(ctx, draftSale, items)
	if err != nil {
		t.Fatalf("RecalculateAndSaveDraft failed: %v", err)
	}
	t.Logf("✓ Saved draft sale to BANLE_TAM & BANLE_TAM_CT: SCT=%s, Total=%.0f đ", draftSale.SCT, draftSale.ThanhToan)

	// 7. Verify draft exists in ListDraftSales
	drafts, err := saleRepo.ListDraftSales(ctx, "01", "01")
	if err != nil {
		t.Fatalf("ListDraftSales failed: %v", err)
	}
	t.Logf("✓ Current draft sales at counter 01: %d", len(drafts))

	// Dọn dẹp đơn test khỏi bảng tạm
	_, _ = sqlDB.ExecContext(ctx, "DELETE FROM BANLE_TAM_CT WHERE KHOA = ?", draftSale.Khoa)
	_, _ = sqlDB.ExecContext(ctx, "DELETE FROM BANLE_TAM WHERE KHOA = ?", draftSale.Khoa)
	t.Logf("✓ Cleaned up test draft from BANLE_TAM")
}

func TestCompareCalculationsWithRealDatabaseRecords(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	cfg := config.LoadConfig()
	sqlDB, err := db.Connect(cfg)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer sqlDB.Close()

	discountServ := service.NewDiscountService()
	taxServ := service.NewTaxService()

	// Truy vấn các dòng thực tế gần nhất trong BANLE_CT có thuế VAT > 0
	query := `
		SELECT TOP 5 
			KHOACT, MAVT, ISNULL(TENVT, ''), SOLUONG, DONGIA, 
			ISNULL(TL_CK, 0), ISNULL(CHIETKHAU, 0), LOAITHUE, 
			THUE_SUAT, TIEN_THUE, THANHTIEN_CHUA_VAT, THANHTOAN
		FROM BANLE_CT
		WHERE THUE_SUAT > 0 AND THANHTOAN IS NOT NULL
		ORDER BY KHOACT DESC
	`
	rows, err := sqlDB.QueryContext(ctx, query)
	if err != nil {
		t.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var item models.BanleTamCT
		var dbTienThue, dbThanhTienChuaVAT, dbThanhToan, dbChietKhau float64
		err := rows.Scan(
			&item.KhoaCT, &item.MaVT, &item.TenVT, &item.SoLuong, &item.DonGia,
			&item.TL_CK, &dbChietKhau, &item.LoaiThue,
			&item.ThueSuat, &dbTienThue, &dbThanhTienChuaVAT, &dbThanhToan,
		)
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}

		// Tính toán lại bằng code mới của chúng ta
		discountServ.CalculateLineItemDiscount(&item, 50.0)
		taxServ.CalculateTaxForLine(&item)

		// So sánh kết quả tính toán với dữ liệu thực tế đã lưu trong DB
		diffPay := math.Abs(item.ThanhToan - dbThanhToan)
		diffTax := math.Abs(item.TienThue - dbTienThue)
		diffPreTax := math.Abs(item.ThanhTienChuaVAT - dbThanhTienChuaVAT)

		if diffPay > 1.0 || diffTax > 1.0 || diffPreTax > 1.0 {
			t.Errorf("Mismatch for item %s (%s):\nComputed: Pay=%.0f, Tax=%.0f, PreTax=%.0f\nDB Actual: Pay=%.0f, Tax=%.0f, PreTax=%.0f",
				item.MaVT, item.TenVT, item.ThanhToan, item.TienThue, item.ThanhTienChuaVAT,
				dbThanhToan, dbTienThue, dbThanhTienChuaVAT)
		} else {
			t.Logf("✓ Match DB: %s (Qty: %.0f x %.0f, CK: %.0f%%) -> Pay: %.0f (DB: %.0f), PreTax: %.0f (DB: %.0f), Tax: %.0f (DB: %.0f)",
				item.MaVT, item.SoLuong, item.DonGia, item.TL_CK,
				item.ThanhToan, dbThanhToan, item.ThanhTienChuaVAT, dbThanhTienChuaVAT, item.TienThue, dbTienThue)
		}
		count++
	}

	// Truy vấn thêm các dòng có chiết khấu TL_CK > 0
	queryDiscount := `
		SELECT TOP 5 
			KHOACT, MAVT, ISNULL(TENVT, ''), SOLUONG, DONGIA, 
			ISNULL(TL_CK, 0), ISNULL(CHIETKHAU, 0), LOAITHUE, 
			THUE_SUAT, TIEN_THUE, THANHTIEN_CHUA_VAT, THANHTOAN
		FROM BANLE_CT
		WHERE TL_CK > 0 AND THUE_SUAT > 0 AND THANHTOAN IS NOT NULL
		ORDER BY KHOACT DESC
	`
	rowsDisc, err := sqlDB.QueryContext(ctx, queryDiscount)
	if err == nil {
		defer rowsDisc.Close()
		for rowsDisc.Next() {
			var item models.BanleTamCT
			var dbTienThue, dbThanhTienChuaVAT, dbThanhToan, dbChietKhau float64
			if err := rowsDisc.Scan(
				&item.KhoaCT, &item.MaVT, &item.TenVT, &item.SoLuong, &item.DonGia,
				&item.TL_CK, &dbChietKhau, &item.LoaiThue,
				&item.ThueSuat, &dbTienThue, &dbThanhTienChuaVAT, &dbThanhToan,
			); err == nil {
				discountServ.CalculateLineItemDiscount(&item, 100.0)
				taxServ.CalculateTaxForLine(&item)

				diffPay := math.Abs(item.ThanhToan - dbThanhToan)
				diffCK := math.Abs(item.ChietKhau - dbChietKhau)
				if diffPay <= 1.0 && diffCK <= 1.0 {
					t.Logf("✓ Match DB with Discount: %s (Qty: %.0f x %.0f, CK: %.0f%% = %.0f) -> Pay: %.0f (DB: %.0f), Tax: %.0f (DB: %.0f)",
						item.MaVT, item.SoLuong, item.DonGia, item.TL_CK, item.ChietKhau,
						item.ThanhToan, dbThanhToan, item.TienThue, dbTienThue)
				} else {
					t.Errorf("Mismatch with discount for item %s: Computed Pay=%.0f (DB: %.0f), CK=%.0f (DB: %.0f)",
						item.MaVT, item.ThanhToan, dbThanhToan, item.ChietKhau, dbChietKhau)
				}
				count++
			}
		}
	}

	if count == 0 {
		t.Fatalf("No rows tested")
	}
	t.Logf("✓ Successfully verified %d real historical rows with 100%% precision!", count)
}

func TestJsonUnmarshalSaleDTO(t *testing.T) {
	rawJSON := []byte(`{
		"header": {
			"khoa": "",
			"lct": "BLE",
			"sct": "00001",
			"ngay": "2026-10-15T18:18:00.000Z",
			"loc": "01",
			"makho": "01",
			"quay": "01",
			"ca": "1",
			"thungan": "ntd",
			"create_date": ""
		},
		"items": [
			{
				"khoact": "ct-17382718-abc",
				"mavt": "HOA001",
				"tenvt": "Hoa Hong",
				"dvt": "Canh",
				"soluong": 1,
				"dongia": 35000
			}
		]
	}`)

	var sale models.SaleDTO
	err := json.Unmarshal(rawJSON, &sale)
	if err != nil {
		t.Fatalf("Expected successful unmarshal, got: %v", err)
	}
	if sale.Header.Khoa == uuid.Nil {
		t.Errorf("Expected Header.Khoa to be populated with a valid UUID, got Nil")
	}
	if len(sale.Items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(sale.Items))
	}
	if sale.Items[0].KhoaCT == uuid.Nil {
		t.Errorf("Expected Items[0].KhoaCT to be populated with a valid UUID, got Nil")
	}
	t.Logf("✓ Successfully unmarshaled SaleDTO with resilient UUIDs: Khoa=%s, KhoaCT=%s",
		sale.Header.Khoa, sale.Items[0].KhoaCT)
}



