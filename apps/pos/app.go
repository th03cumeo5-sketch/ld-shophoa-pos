package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"shophoa/core/config"
	"shophoa/core/db"
	"shophoa/core/models"
	"shophoa/core/repository"
	"shophoa/core/service"
)

// App struct
type App struct {
	ctx          context.Context
	db           *sql.DB
	prodRepo     *repository.ProductRepository
	custRepo     *repository.CustomerRepository
	saleRepo     *repository.SaleRepository
	discountServ *service.DiscountService
	taxServ      *service.TaxService
	saleServ     *service.SaleService

	// Session State
	MaKho   string
	Quay    string
	Ca      string
	ThuNgan string
	LOC     string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		MaKho:   "01",
		Quay:    "01",
		Ca:      "1",
		ThuNgan: "ntd",
		LOC:     "01",
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Kết nối CSDL SQL Server (Shared Memory / ODBC)
	cfg := config.LoadConfig()
	sqlDB, err := db.Connect(cfg)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}
	a.db = sqlDB

	a.prodRepo = repository.NewProductRepository(sqlDB)
	a.custRepo = repository.NewCustomerRepository(sqlDB)
	a.saleRepo = repository.NewSaleRepository(sqlDB)
	a.discountServ = service.NewDiscountService()
	a.taxServ = service.NewTaxService()
	a.saleServ = service.NewSaleService(a.saleRepo, a.prodRepo, a.discountServ, a.taxServ)

	fmt.Println("✓ POS App Backend initialized with SQL Server connection")
}

// shutdown is called when the app terminates
func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

// SearchProducts tìm kiếm sản phẩm theo từ khóa / barcode và nhóm hàng
func (a *App) SearchProducts(keyword string, category string) ([]models.ProductSearchResult, error) {
	if a.prodRepo == nil {
		return nil, fmt.Errorf("database not connected")
	}
	return a.prodRepo.SearchProductsWithCategory(a.ctx, keyword, category, 24)
}

// GetCategories lấy danh sách các nhóm hàng chính
func (a *App) GetCategories() ([]models.Category, error) {
	if a.prodRepo == nil {
		return nil, fmt.Errorf("database not connected")
	}
	return a.prodRepo.GetTopCategories(a.ctx)
}

// GetCatalogProducts lấy danh sách sản phẩm mẫu cho lưới catalog bên phải
func (a *App) GetCatalogProducts(limit int) ([]models.ProductSearchResult, error) {
	if a.prodRepo == nil {
		return nil, fmt.Errorf("database not connected")
	}
	return a.prodRepo.SearchProductsWithCategory(a.ctx, "", "", limit)
}

// SearchCustomers tìm kiếm khách hàng
func (a *App) SearchCustomers(keyword string) ([]models.Customer, error) {
	if a.custRepo == nil {
		return nil, fmt.Errorf("database not connected")
	}
	return a.custRepo.SearchCustomers(a.ctx, keyword, 15)
}

// PrepareNewSale khởi tạo một hóa đơn mới, lấy SCT từ Stored Procedure POS_ALLOC_BILL_NO
func (a *App) PrepareNewSale() (*models.BanleTam, error) {
	if a.saleServ == nil {
		return nil, fmt.Errorf("database not connected")
	}
	return a.saleServ.PrepareNewSale(a.ctx, a.MaKho, a.Quay, a.Ca, a.ThuNgan, a.LOC)
}

// CalculateCart tính toán lại chiết khấu, thuế VAT và tổng tiền cho giỏ hàng
func (a *App) CalculateCart(sale models.SaleDTO) (*models.SaleDTO, error) {
	if a.saleServ == nil {
		return nil, fmt.Errorf("database not connected")
	}
	a.saleServ.RecalculateSale(&sale.Header, sale.Items)
	return &sale, nil
}

// SaveDraftSale lưu hóa đơn đang soạn / lưu tạm vào BANLE_TAM & BANLE_TAM_CT
func (a *App) SaveDraftSale(sale models.SaleDTO) error {
	if a.saleServ == nil {
		return fmt.Errorf("database not connected")
	}
	return a.saleServ.RecalculateAndSaveDraft(a.ctx, &sale.Header, sale.Items)
}

// GetDraftSale mở lại một đơn hàng lưu tạm
func (a *App) GetDraftSale(khoa string) (*models.SaleDTO, error) {
	if a.saleServ == nil {
		return nil, fmt.Errorf("database not connected")
	}
	u, err := uuid.Parse(khoa)
	if err != nil {
		return nil, fmt.Errorf("invalid khoa guid: %w", err)
	}
	return a.saleServ.GetDraftSale(a.ctx, u)
}

// DeleteDraftSale xóa bỏ một đơn tạm
func (a *App) DeleteDraftSale(khoa string) error {
	if a.saleServ == nil {
		return fmt.Errorf("database not connected")
	}
	u, err := uuid.Parse(khoa)
	if err != nil {
		return fmt.Errorf("invalid khoa guid: %w", err)
	}
	return a.saleServ.DeleteDraftSale(a.ctx, u)
}

// CompletePayment thanh toán và chuyển từ bảng tạm sang bảng chính qua POS_PAID
func (a *App) CompletePayment(sale models.SaleDTO, pttt string, ttoan1, ttoan2 float64) error {
	if a.saleServ == nil {
		return fmt.Errorf("database not connected")
	}
	return a.saleServ.CompletePayment(a.ctx, &sale.Header, sale.Items, pttt, ttoan1, ttoan2)
}

// ListDraftSales lấy danh sách các đơn tạm đang treo
func (a *App) ListDraftSales() ([]models.BanleTam, error) {
	if a.saleRepo == nil {
		return nil, fmt.Errorf("database not connected")
	}
	return a.saleRepo.ListDraftSales(a.ctx, a.MaKho, a.Quay)
}

// GetSessionInfo trả về thông tin ca làm việc hiện tại
func (a *App) GetSessionInfo() map[string]string {
	return map[string]string{
		"makho":   a.MaKho,
		"quay":    a.Quay,
		"ca":      "Ca Sáng",
		"thungan": a.ThuNgan,
		"loc":     "Cửa hàng Q2 (Kho: K01)",
		"time":    time.Now().Format("15:04:05"),
	}
}

