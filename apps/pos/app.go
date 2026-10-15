package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
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
	userRepo     *repository.UserRepository
	discountServ *service.DiscountService
	taxServ      *service.TaxService
	saleServ     *service.SaleService

	// Profiles & Session State
	profilesConfig *ProfilesConfig
	currentProfile DBProfile
	session        *models.SessionInfo
	isLoggedIn     bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	cfg := LoadProfilesConfig()
	var initial DBProfile
	for _, p := range cfg.Profiles {
		if p.ID == cfg.ActiveProfileID {
			initial = p
			break
		}
	}
	if initial.ID == "" && len(cfg.Profiles) > 0 {
		initial = cfg.Profiles[0]
	}

	return &App{
		profilesConfig: cfg,
		currentProfile: initial,
		discountServ:   service.NewDiscountService(),
		taxServ:        service.NewTaxService(),
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Thử khởi tạo kết nối CSDL mặc định gần nhất
	if a.currentProfile.ID != "" {
		_ = a.initDatabaseConnection(a.currentProfile)
	}
}

// shutdown is called when the app terminates
func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

// initDatabaseConnection thiết lập kết nối CSDL và các repositories tương ứng
func (a *App) initDatabaseConnection(p DBProfile) error {
	if a.db != nil {
		a.db.Close()
		a.db = nil
	}

	sqlDB, err := ConnectProfile(p)
	if err != nil {
		return fmt.Errorf("kết nối CSDL thất bại: %w", err)
	}

	a.db = sqlDB
	a.prodRepo = repository.NewProductRepository(sqlDB)
	a.custRepo = repository.NewCustomerRepository(sqlDB)
	a.saleRepo = repository.NewSaleRepository(sqlDB)
	a.userRepo = repository.NewUserRepository(sqlDB)
	a.saleServ = service.NewSaleService(a.saleRepo, a.prodRepo, a.discountServ, a.taxServ)
	a.currentProfile = p

	return nil
}

// =========================================================================
// DATABASE PROFILES & LOGIN / AUTH METHODS
// =========================================================================

// GetDBProfiles trả về danh sách cấu hình kết nối CSDL và cấu hình gần nhất
func (a *App) GetDBProfiles() map[string]interface{} {
	if a.profilesConfig == nil {
		a.profilesConfig = LoadProfilesConfig()
	}
	return map[string]interface{}{
		"active_id": a.profilesConfig.ActiveProfileID,
		"profiles":  a.profilesConfig.Profiles,
	}
}

// SaveDBProfile thêm mới hoặc cập nhật một cấu hình kết nối CSDL
func (a *App) SaveDBProfile(p DBProfile) error {
	if p.Name == "" {
		p.Name = fmt.Sprintf("%s (%s)", p.Server, p.Database)
	}
	if p.ID == "" {
		p.ID = fmt.Sprintf("prof_%d", time.Now().UnixNano())
	}

	found := false
	for i, existing := range a.profilesConfig.Profiles {
		if existing.ID == p.ID {
			a.profilesConfig.Profiles[i] = p
			found = true
			break
		}
	}
	if !found {
		a.profilesConfig.Profiles = append(a.profilesConfig.Profiles, p)
	}

	a.profilesConfig.ActiveProfileID = p.ID
	return SaveProfilesConfig(a.profilesConfig)
}

// DeleteDBProfile xóa bỏ một cấu hình kết nối CSDL
func (a *App) DeleteDBProfile(id string) error {
	if len(a.profilesConfig.Profiles) <= 1 {
		return fmt.Errorf("không thể xóa vì hệ thống cần giữ lại ít nhất 1 cấu hình kết nối")
	}

	var updated []DBProfile
	for _, p := range a.profilesConfig.Profiles {
		if p.ID != id {
			updated = append(updated, p)
		}
	}
	a.profilesConfig.Profiles = updated

	if a.profilesConfig.ActiveProfileID == id {
		a.profilesConfig.ActiveProfileID = updated[0].ID
	}

	return SaveProfilesConfig(a.profilesConfig)
}

// TestDBConnection kiểm tra thử kết nối CSDL
func (a *App) TestDBConnection(p DBProfile) (string, error) {
	testDB, err := ConnectProfile(p)
	if err != nil {
		return "", err
	}
	defer testDB.Close()

	var dbName, dbVersion string
	err = testDB.QueryRow("SELECT DB_NAME(), @@VERSION").Scan(&dbName, &dbVersion)
	if err != nil {
		return "Kết nối thành công nhưng không lấy được tên CSDL", nil
	}

	return fmt.Sprintf("✓ Kết nối thành công tới database [%s] trên máy chủ [%s]!", dbName, p.Server), nil
}

// Login đăng nhập vào hệ thống POS với tài khoản SYS_USER và cấu hình kết nối CSDL được chọn
func (a *App) Login(username, password, ca, quay, profileId string) (*models.SessionInfo, error) {
	// 1. Tìm profile kết nối
	var targetProfile DBProfile
	for _, p := range a.profilesConfig.Profiles {
		if p.ID == profileId {
			targetProfile = p
			break
		}
	}
	if targetProfile.ID == "" {
		if len(a.profilesConfig.Profiles) > 0 {
			targetProfile = a.profilesConfig.Profiles[0]
		} else {
			return nil, fmt.Errorf("chưa có cấu hình kết nối CSDL nào")
		}
	}

	// 2. Chuyển đổi kết nối CSDL nếu chưa kết nối hoặc profile khác
	if a.db == nil || a.currentProfile.ID != targetProfile.ID {
		if err := a.initDatabaseConnection(targetProfile); err != nil {
			return nil, fmt.Errorf("không thể kết nối tới máy chủ CSDL: %w", err)
		}
	}

	// 3. Xác thực tài khoản người dùng
	user, err := a.userRepo.Authenticate(a.ctx, username, password)
	if err != nil {
		return nil, err
	}

	// 4. Lấy thông tin kho
	makho := targetProfile.MaKho
	if makho == "" {
		makho = "01"
	}
	if quay == "" {
		quay = targetProfile.Quay
		if quay == "" {
			quay = "01"
		}
	}
	if ca == "" {
		ca = "1"
	}

	tenKho := a.userRepo.GetWarehouseName(a.ctx, makho)

	// 5. Lưu phiên làm việc
	a.session = &models.SessionInfo{
		Username:  user.Username,
		FullName:  user.FullName,
		MaKho:     makho,
		TenKho:    tenKho,
		Quay:      quay,
		Ca:        ca,
		LOC:       user.LOC,
		ProfileID: targetProfile.ID,
	}
	a.isLoggedIn = true

	// Ghi nhớ profile đã đăng nhập thành công
	a.profilesConfig.ActiveProfileID = targetProfile.ID
	_ = SaveProfilesConfig(a.profilesConfig)

	return a.session, nil
}

// Logout đăng xuất khỏi phiên làm việc hiện tại
func (a *App) Logout() bool {
	a.session = nil
	a.isLoggedIn = false
	return true
}

// GetSessionInfo trả về thông tin phiên làm việc hiện tại
func (a *App) GetSessionInfo() *models.SessionInfo {
	return a.session
}

// GetWarehouses lấy danh sách kho từ CSDL đang kết nối
func (a *App) GetWarehouses() ([]models.Warehouse, error) {
	if a.userRepo == nil {
		return nil, fmt.Errorf("chưa kết nối CSDL")
	}
	return a.userRepo.GetWarehouses(a.ctx)
}

// =========================================================================
// POS BUSINESS METHODS (HOOKED TO BACKEND)
// =========================================================================

// SearchProducts tìm kiếm sản phẩm theo từ khóa / barcode và nhóm hàng
func (a *App) SearchProducts(keyword string, category string) ([]models.ProductSearchResult, error) {
	if a.prodRepo == nil {
		return nil, fmt.Errorf("chưa kết nối CSDL")
	}
	return a.prodRepo.SearchProductsWithCategory(a.ctx, keyword, category, 24)
}

// GetCategories lấy danh sách các nhóm hàng chính
func (a *App) GetCategories() ([]models.Category, error) {
	if a.prodRepo == nil {
		return nil, fmt.Errorf("chưa kết nối CSDL")
	}
	return a.prodRepo.GetTopCategories(a.ctx)
}

// GetCatalogProducts lấy danh sách sản phẩm mẫu cho lưới catalog bên phải
func (a *App) GetCatalogProducts(limit int) ([]models.ProductSearchResult, error) {
	if a.prodRepo == nil {
		return nil, fmt.Errorf("chưa kết nối CSDL")
	}
	return a.prodRepo.SearchProductsWithCategory(a.ctx, "", "", limit)
}

// SearchCustomers tìm kiếm khách hàng
func (a *App) SearchCustomers(keyword string) ([]models.Customer, error) {
	if a.custRepo == nil {
		return nil, fmt.Errorf("chưa kết nối CSDL")
	}
	return a.custRepo.SearchCustomers(a.ctx, keyword, 15)
}

// PrepareNewSale khởi tạo một hóa đơn mới, lấy SCT từ Stored Procedure POS_ALLOC_BILL_NO
func (a *App) PrepareNewSale() (*models.BanleTam, error) {
	if a.saleServ == nil {
		return nil, fmt.Errorf("chưa kết nối CSDL")
	}

	makho := "01"
	quay := "01"
	ca := "1"
	thungan := "ntd"
	loc := "01"

	if a.session != nil {
		makho = a.session.MaKho
		quay = a.session.Quay
		ca = a.session.Ca
		thungan = a.session.Username
		loc = a.session.LOC
	}

	return a.saleServ.PrepareNewSale(a.ctx, makho, quay, ca, thungan, loc)
}

// CalculateCart tính toán lại chiết khấu, thuế VAT và tổng tiền cho giỏ hàng
func (a *App) CalculateCart(sale models.SaleDTO) (*models.SaleDTO, error) {
	if a.saleServ == nil {
		return nil, fmt.Errorf("chưa kết nối CSDL")
	}
	a.saleServ.RecalculateSale(&sale.Header, sale.Items)
	return &sale, nil
}

// SaveDraftSale lưu hóa đơn đang soạn / lưu tạm vào BANLE_TAM & BANLE_TAM_CT
func (a *App) SaveDraftSale(sale models.SaleDTO) error {
	if a.saleServ == nil {
		return fmt.Errorf("chưa kết nối CSDL")
	}
	if a.session != nil {
		sale.Header.MaKho = a.session.MaKho
		sale.Header.Quay = a.session.Quay
		sale.Header.Ca = a.session.Ca
		sale.Header.ThuNgan = a.session.Username
		sale.Header.LOC = a.session.LOC
	}
	return a.saleServ.RecalculateAndSaveDraft(a.ctx, &sale.Header, sale.Items)
}

// GetDraftSale mở lại một đơn hàng lưu tạm
func (a *App) GetDraftSale(khoa string) (*models.SaleDTO, error) {
	if a.saleServ == nil {
		return nil, fmt.Errorf("chưa kết nối CSDL")
	}
	u, err := uuid.Parse(khoa)
	if err != nil {
		return nil, fmt.Errorf("khoa không hợp lệ: %w", err)
	}
	return a.saleServ.GetDraftSale(a.ctx, u)
}

// DeleteDraftSale xóa bỏ một đơn tạm
func (a *App) DeleteDraftSale(khoa string) error {
	if a.saleServ == nil {
		return fmt.Errorf("chưa kết nối CSDL")
	}
	u, err := uuid.Parse(khoa)
	if err != nil {
		return fmt.Errorf("khoa không hợp lệ: %w", err)
	}
	return a.saleServ.DeleteDraftSale(a.ctx, u)
}

// CompletePayment thanh toán và chuyển từ bảng tạm sang bảng chính qua POS_PAID
func (a *App) CompletePayment(sale models.SaleDTO, pttt string, ttoan1, ttoan2 float64) error {
	if a.saleServ == nil {
		return fmt.Errorf("chưa kết nối CSDL")
	}
	if a.session != nil {
		sale.Header.MaKho = a.session.MaKho
		sale.Header.Quay = a.session.Quay
		sale.Header.Ca = a.session.Ca
		sale.Header.ThuNgan = a.session.Username
		sale.Header.LOC = a.session.LOC
	}
	return a.saleServ.CompletePayment(a.ctx, &sale.Header, sale.Items, pttt, ttoan1, ttoan2)
}

// ListDraftSales lấy danh sách các đơn tạm đang treo tại quầy hiện tại
func (a *App) ListDraftSales() ([]models.BanleTam, error) {
	if a.saleRepo == nil {
		return nil, fmt.Errorf("chưa kết nối CSDL")
	}
	makho := "01"
	quay := "01"
	if a.session != nil {
		makho = a.session.MaKho
		quay = a.session.Quay
	}
	return a.saleRepo.ListDraftSales(a.ctx, makho, quay)
}
