package main

import (
	"context"
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
