# TECHNICAL SPECIFICATION DOCUMENT (TSD) - GOLANG
## POSz - Point of Sale System Migration to Go

---

## 1. SYSTEM ARCHITECTURE

### 1.1 Architecture Overview
```
┌─────────────────────────────────────────────────────────────┐
│                     POS Client Layer                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   POS UI     │  │   Receipt    │  │   Barcode    │      │
│  │   (Desktop)  │  │   Printer    │  │   Scanner    │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
│         │                  │                  │              │
│         └──────────────────┴──────────────────┘              │
│                            │                                 │
└────────────────────────────┼─────────────────────────────────┘
                             │
┌────────────────────────────┼─────────────────────────────────┐
│                     API Gateway Layer                        │
│                    (Go + Gin/Echo/Fiber)                     │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  REST API Endpoints                                    │  │
│  │  - /api/pos/sale         - /api/pos/return            │  │
│  │  - /api/pos/payment      - /api/pos/promotion         │  │
│  │  - /api/products/*       - /api/customers/*           │  │
│  └───────────────────────────────────────────────────────┘  │
└────────────────────────────┼─────────────────────────────────┘
                             │
┌────────────────────────────┼─────────────────────────────────┐
│                  Business Logic Layer (Go)                   │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  Sale Engine │  │   Discount   │  │  Promotion   │      │
│  │              │  │   Calculator │  │   Engine     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │    Return    │  │   Payment    │  │     Tax      │      │
│  │   Handler    │  │   Processor  │  │  Calculator  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└────────────────────────────┼─────────────────────────────────┘
                             │
┌────────────────────────────┼─────────────────────────────────┐
│                    Data Access Layer (Go)                    │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  GORM / sqlx / database/sql                            │  │
│  │  - Connection Pool                                     │  │
│  │  - Transaction Management                              │  │
│  │  - Query Builder                                       │  │
│  └───────────────────────────────────────────────────────┘  │
└────────────────────────────┼─────────────────────────────────┘
                             │
┌────────────────────────────┼─────────────────────────────────┐
│                    Database Layer                            │
│                    (SQL Server)                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │  BANLE   │  │ BANLE_CT │  │ TRAHANG  │  │  DMVT    │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 Technology Stack

**Backend:**
- **Language:** Go 1.21+
- **Web Framework:** Gin (high performance) OR Echo OR Fiber
- **ORM/Database:** 
  - GORM (recommended for productivity)
  - OR sqlx (for performance)
  - Driver: github.com/denisenkom/go-mssqldb
- **Authentication:** JWT (golang-jwt/jwt)
- **Decimal Math:** shopspring/decimal (CRITICAL for financial calculations)
- **Validation:** go-playground/validator
- **Testing:** testify + httptest

**Database:**
- **RDBMS:** Microsoft SQL Server (existing schema)
- **Driver:** go-mssqldb

**Build & Tools:**
- **Build:** Go modules
- **CI/CD:** GitHub Actions / GitLab CI
- **Container:** Docker (optional)
- **Linting:** golangci-lint

---

## 2. PROJECT STRUCTURE

```
posz/
├── go.mod
├── go.sum
├── .env
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── models/
│   │   ├── sale.go              # BANLE, BANLE_CT
│   │   ├── return.go            # TRAHANG, TRAHANG_CT
│   │   ├── promotion.go         # KHUYENMAI, KHUYENMAI_CT
│   │   ├── product.go           # DM_VT, DM_VT_BO
│   │   ├── customer.go          # DM_KH_NCC
│   │   └── common.go            # Common types
│   ├── repository/
│   │   ├── sale_repo.go         # Sale data access
│   │   ├── product_repo.go      # Product data access
│   │   ├── promotion_repo.go    # Promotion data access
│   │   └── customer_repo.go     # Customer data access
│   ├── service/
│   │   ├── sale_service.go      # Sale business logic
│   │   ├── discount_service.go  # Discount calculation
│   │   ├── tax_service.go       # Tax calculation
│   │   ├── promotion_service.go # Promotion logic
│   │   ├── return_service.go    # Return processing
│   │   └── payment_service.go   # Payment processing
│   ├── handler/
│   │   ├── sale_handler.go      # HTTP handlers for sale
│   │   ├── return_handler.go    # HTTP handlers for return
│   │   ├── product_handler.go   # HTTP handlers for products
│   │   └── auth_handler.go      # Authentication handlers
│   ├── middleware/
│   │   ├── auth.go              # JWT authentication
│   │   ├── logger.go            # Request logging
│   │   └── recovery.go          # Panic recovery
│   ├── database/
│   │   └── db.go                # Database connection
│   └── utils/
│       ├── guid.go              # GUID generation
│       ├── decimal.go           # Decimal helpers
│       └── response.go          # HTTP response helpers
├── pkg/                         # Public packages (if any)
└── test/
    └── integration/
        └── sale_test.go
```

---

## 3. DEPENDENCIES (go.mod)

```go
module github.com/yourorg/posz

go 1.21

require (
    // Web framework
    github.com/gin-gonic/gin v1.9.1
    
    // Database
    github.com/denisenkom/go-mssqldb v0.12.3
    gorm.io/gorm v1.25.5
    gorm.io/driver/sqlserver v1.5.2
    // OR for raw SQL
    github.com/jmoiron/sqlx v1.3.5
    
    // Decimal math (CRITICAL)
    github.com/shopspring/decimal v1.3.1
    
    // UUID
    github.com/google/uuid v1.4.0
    
    // JWT
    github.com/golang-jwt/jwt/v5 v5.1.0
    
    // Password hashing
    golang.org/x/crypto v0.15.0
    
    // Configuration
    github.com/spf13/viper v1.17.0
    github.com/joho/godotenv v1.5.1
    
    // Validation
    github.com/go-playground/validator/v10 v10.16.0
    
    // Logging
    github.com/sirupsen/logrus v1.9.3
    
    // Testing
    github.com/stretchr/testify v1.8.4
)
```

---

## 4. DATABASE CONNECTION

### 4.1 Connection Setup (GORM)

**File: `internal/database/db.go`**

```go
package database

import (
    "fmt"
    "log"
    "gorm.io/driver/sqlserver"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(dsn string) error {
    var err error
    
    DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
    })
    
    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }
    
    sqlDB, err := DB.DB()
    if err != nil {
        return err
    }
    
    // Connection pool settings
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    log.Println("Database connection established")
    return nil
}

func Close() error {
    sqlDB, err := DB.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}
```

### 4.2 Connection String Format

```
sqlserver://username:password@hostname:1433?database=SOFTZ_SHOPHOA&encrypt=disable
```

Or in `.env`:
```
DB_DSN=sqlserver://sa:YourPassword@localhost:1433?database=SOFTZ_SHOPHOA&encrypt=disable
JWT_SECRET=your-secret-key-here
SERVER_PORT=8080
```

---

## 5. DATA MODELS

### 5.1 Sale Models

**File: `internal/models/sale.go`**

```go
package models

import (
    "time"
    "github.com/google/uuid"
    "github.com/shopspring/decimal"
)

// Banle represents BANLE table (Retail Sales Header)
type Banle struct {
    KHOA               uuid.UUID        `gorm:"column:KHOA;type:uniqueidentifier;primaryKey"`
    LCT                *string          `gorm:"column:LCT;type:nvarchar(5)"`
    NGAY               *time.Time       `gorm:"column:NGAY;type:datetime"`
    SCT                *string          `gorm:"column:SCT;type:nvarchar(20)"`
    LOC                *string          `gorm:"column:LOC;type:nvarchar(2)"`
    MAKHO              *string          `gorm:"column:MAKHO;type:nvarchar(2)"`
    QUAY               *string          `gorm:"column:QUAY;type:nvarchar(15)"`
    CA                 *string          `gorm:"column:CA;type:nvarchar(5)"`
    THUNGAN            *string          `gorm:"column:THUNGAN;type:nvarchar(15)"`
    MADT               *string          `gorm:"column:MADT;type:nvarchar(15)"`
    SOLUONG            decimal.Decimal  `gorm:"column:SOLUONG;type:float"`
    THANHTIEN_CHUA_VAT decimal.Decimal  `gorm:"column:THANHTIEN_CHUA_VAT;type:float"`
    TIEN_THUE          decimal.Decimal  `gorm:"column:TIEN_THUE;type:float"`
    THANHTOAN          decimal.Decimal  `gorm:"column:THANHTOAN;type:float"`
    
    // Tax breakdown
    TIEN_THUE_0        decimal.Decimal  `gorm:"column:TIEN_THUE_0;type:float"`
    TIEN_THUE_5        decimal.Decimal  `gorm:"column:TIEN_THUE_5;type:float"`
    TIEN_THUE_8        decimal.Decimal  `gorm:"column:TIEN_THUE_8;type:float"`
    TIEN_THUE_10       decimal.Decimal  `gorm:"column:TIEN_THUE_10;type:float"`
    TIEN_THUE_KHAC     decimal.Decimal  `gorm:"column:TIEN_THUE_KHAC;type:float"`
    
    THANHTIEN_CHUA_VAT_0    decimal.Decimal `gorm:"column:THANHTIEN_CHUA_VAT_0;type:float"`
    THANHTIEN_CHUA_VAT_5    decimal.Decimal `gorm:"column:THANHTIEN_CHUA_VAT_5;type:float"`
    THANHTIEN_CHUA_VAT_8    decimal.Decimal `gorm:"column:THANHTIEN_CHUA_VAT_8;type:float"`
    THANHTIEN_CHUA_VAT_10   decimal.Decimal `gorm:"column:THANHTIEN_CHUA_VAT_10;type:float"`
    THANHTIEN_CHUA_VAT_KHAC decimal.Decimal `gorm:"column:THANHTIEN_CHUA_VAT_KHAC;type:float"`
    
    THANHTOAN_0   decimal.Decimal `gorm:"column:THANHTOAN_0;type:float"`
    THANHTOAN_5   decimal.Decimal `gorm:"column:THANHTOAN_5;type:float"`
    THANHTOAN_8   decimal.Decimal `gorm:"column:THANHTOAN_8;type:float"`
    THANHTOAN_10  decimal.Decimal `gorm:"column:THANHTOAN_10;type:float"`
    THANHTOAN_KHAC decimal.Decimal `gorm:"column:THANHTOAN_KHAC;type:float"`
    THANHTOAN_KCT decimal.Decimal `gorm:"column:THANHTOAN_KCT;type:float"`
    
    // Payment
    PTTT       *string         `gorm:"column:PTTT;type:nvarchar(2)"`
    TTOAN1     decimal.Decimal `gorm:"column:TTOAN1;type:float"`
    PTTT2      *string         `gorm:"column:PTTT2;type:nvarchar(2)"`
    TTOAN2     decimal.Decimal `gorm:"column:TTOAN2;type:float"`
    CHUATHOI   decimal.Decimal `gorm:"column:CHUATHOI;type:float"`
    SO_CHUANCHI *string        `gorm:"column:SO_CHUANCHI;type:nvarchar(20)"`
    
    // Status
    PRINTED    *bool           `gorm:"column:PRINTED;type:bit"`
    PRINT_NO   *int            `gorm:"column:PRINT_NO;type:int"`
    
    // Audit trail
    CREATE_BY   *int           `gorm:"column:CREATE_BY;type:int"`
    CREATE_DATE *time.Time     `gorm:"column:CREATE_DATE;type:datetime"`
    UPDATE_BY   *int           `gorm:"column:UPDATE_BY;type:int"`
    UPDATE_DATE *time.Time     `gorm:"column:UPDATE_DATE;type:datetime"`
    DELETE_BY   *int           `gorm:"column:DELETE_BY;type:int"`
    DELETE_DATE *time.Time     `gorm:"column:DELETE_DATE;type:datetime"`
    
    // Relationships
    Items []BanleCT `gorm:"foreignKey:KHOA;references:KHOA"`
}

func (Banle) TableName() string {
    return "BANLE"
}

// BanleCT represents BANLE_CT table (Retail Sales Detail)
type BanleCT struct {
    KHOACT             uuid.UUID       `gorm:"column:KHOACT;type:uniqueidentifier;primaryKey"`
    KHOA               uuid.UUID       `gorm:"column:KHOA;type:uniqueidentifier"`
    STT                *int            `gorm:"column:STT;type:int"`
    
    // Product info
    MAVT               *string         `gorm:"column:MAVT;type:nvarchar(15)"`
    TENVT              *string         `gorm:"column:TENVT;type:nvarchar(200)"`
    DVT                *string         `gorm:"column:DVT;type:nvarchar(20)"`
    MABO               *string         `gorm:"column:MABO;type:nvarchar(15)"`
    
    // Quantity & Price
    SOLUONG            decimal.Decimal `gorm:"column:SOLUONG;type:float"`
    DONGIA             decimal.Decimal `gorm:"column:DONGIA;type:float"`
    DONGIA_CHUA_VAT    decimal.Decimal `gorm:"column:DONGIA_CHUA_VAT;type:float"`
    
    // Multi-level discounts
    TL_CK              decimal.Decimal `gorm:"column:TL_CK;type:float"`
    TL_CK2             decimal.Decimal `gorm:"column:TL_CK2;type:float"`
    TL_CK3             decimal.Decimal `gorm:"column:TL_CK3;type:float"`
    TL_CK4             decimal.Decimal `gorm:"column:TL_CK4;type:float"`
    TL_CK5             decimal.Decimal `gorm:"column:TL_CK5;type:float"`
    TL_CK6             decimal.Decimal `gorm:"column:TL_CK6;type:float"`
    TL_CK7             decimal.Decimal `gorm:"column:TL_CK7;type:float"`
    TL_CK_MAX          decimal.Decimal `gorm:"column:TL_CK_MAX;type:float"`
    
    CK_BY              *int            `gorm:"column:CK_BY;type:int"`
    
    // Tax
    LOAITHUE           *string         `gorm:"column:LOAITHUE;type:nvarchar(15)"`
    THUE_SUAT          decimal.Decimal `gorm:"column:THUE_SUAT;type:float"`
    
    // Amounts
    THANHTIEN_CHUA_VAT decimal.Decimal `gorm:"column:THANHTIEN_CHUA_VAT;type:float"`
    TIEN_THUE          decimal.Decimal `gorm:"column:TIEN_THUE;type:float"`
    THANHTOAN          decimal.Decimal `gorm:"column:THANHTOAN;type:float"`
    
    // Audit
    CREATE_BY          *int            `gorm:"column:CREATE_BY;type:int"`
    CREATE_DATE        *time.Time      `gorm:"column:CREATE_DATE;type:datetime"`
}

func (BanleCT) TableName() string {
    return "BANLE_CT"
}
```

### 5.2 Product Model

**File: `internal/models/product.go`**

```go
package models

import (
    "github.com/shopspring/decimal"
    "github.com/google/uuid"
    "time"
)

type DmVT struct {
    MAVT        string          `gorm:"column:MAVT;type:nvarchar(15);primaryKey"`
    TENVT       *string         `gorm:"column:TENVT;type:nvarchar(200)"`
    TENVT2      *string         `gorm:"column:TENVT2;type:nvarchar(200)"`
    TENTAT      *string         `gorm:"column:TENTAT;type:nvarchar(50)"`
    DVT         *string         `gorm:"column:DVT;type:nvarchar(20)"`
    DVT1        *string         `gorm:"column:DVT1;type:nvarchar(20)"`
    QD1         *int            `gorm:"column:QD1;type:int"`
    MANHOM      *string         `gorm:"column:MANHOM;type:nvarchar(4)"`
    GIABAN      decimal.Decimal `gorm:"column:GIABAN;type:float"`
    LOAITHUE    *string         `gorm:"column:LOAITHUE;type:nvarchar(15)"`
    THUE_SUAT   decimal.Decimal `gorm:"column:THUE_SUAT;type:float"`
    TL_CK_MAX   decimal.Decimal `gorm:"column:TL_CK_MAX;type:float"`
    BARCODE     *string         `gorm:"column:BARCODE;type:nvarchar(30)"`
    PLU         *string         `gorm:"column:PLU;type:nvarchar(20)"`
    BO          *bool           `gorm:"column:BO;type:bit"`
    TINHTRANG   *string         `gorm:"column:TINHTRANG;type:nvarchar(20)"`
    CREATE_BY   *int            `gorm:"column:CREATE_BY;type:int"`
    CREATE_DATE *time.Time      `gorm:"column:CREATE_DATE;type:datetime"`
}

func (DmVT) TableName() string {
    return "DM_VT"
}

type DmVTBO struct {
    KHOACT      uuid.UUID       `gorm:"column:KHOACT;type:uniqueidentifier;primaryKey"`
    MABO        string          `gorm:"column:MABO;type:nvarchar(15)"`
    MAVT        string          `gorm:"column:MAVT;type:nvarchar(15)"`
    SOLUONG     decimal.Decimal `gorm:"column:SOLUONG;type:float"`
    STT         *int            `gorm:"column:STT;type:int"`
    CREATE_BY   *int            `gorm:"column:CREATE_BY;type:int"`
    CREATE_DATE *time.Time      `gorm:"column:CREATE_DATE;type:datetime"`
}

func (DmVTBO) TableName() string {
    return "DM_VT_BO"
}
```

---

## 6. REPOSITORY LAYER

### 6.1 Sale Repository

**File: `internal/repository/sale_repo.go`**

```go
package repository

import (
    "context"
    "fmt"
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "yourorg/posz/internal/models"
)

type SaleRepository struct {
    db *gorm.DB
}

func NewSaleRepository(db *gorm.DB) *SaleRepository {
    return &SaleRepository{db: db}
}

// SaveDraftSale lưu thông tin đơn hàng đang soạn hoặc treo đơn (F10) vào cụm bảng tạm BANLE_TAM & BANLE_TAM_CT
func (r *SaleRepository) SaveDraftSale(ctx context.Context, banle *models.BanleTam, items []models.BanleTamCT) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // Upsert header vào BANLE_TAM
        if err := tx.Save(banle).Error; err != nil {
            return fmt.Errorf("failed to save banle_tam: %w", err)
        }
        
        // Xóa chi tiết cũ nếu có và thêm lại chi tiết mới vào BANLE_TAM_CT
        if err := tx.Where("KHOA = ?", banle.Khoa).Delete(&models.BanleTamCT{}).Error; err != nil {
            return fmt.Errorf("failed to clean banle_tam_ct: %w", err)
        }
        if len(items) > 0 {
            if err := tx.Create(&items).Error; err != nil {
                return fmt.Errorf("failed to insert banle_tam_ct items: %w", err)
            }
        }
        return nil
    })
}

// CompleteSalePayment thực hiện thanh toán chính thức bằng cách gọi Stored Procedure POS_PAID(@KHOA)
// Thủ tục này sẽ tự động chuyển dữ liệu từ BANLE_TAM -> BANLE, BANLE_TAM_CT -> BANLE_CT, xóa bảng tạm và merge vào ehoadon
func (r *SaleRepository) CompleteSalePayment(ctx context.Context, khoa uuid.UUID) error {
    err := r.db.WithContext(ctx).Exec("EXEC POS_PAID @KHOA = ?", khoa).Error
    if err != nil {
        return fmt.Errorf("POS_PAID failed for khoa %s: %w", khoa.String(), err)
    }
    return nil
}

// AllocateRetailBillNumber cấp số hóa đơn bán lẻ POS thông qua Stored Procedure POS_ALLOC_BILL_NO
func (r *SaleRepository) AllocateRetailBillNumber(ctx context.Context, makho, quay string, ngay time.Time) (string, error) {
    var sct string
    row := r.db.WithContext(ctx).Raw("DECLARE @sct nvarchar(20); EXEC POS_ALLOC_BILL_NO @MAKHO = ?, @QUAY = ?, @NGAY = ?, @SCT = @sct OUTPUT; SELECT @sct;", 
        makho, quay, ngay).Row()
    if err := row.Scan(&sct); err != nil {
        return "", fmt.Errorf("POS_ALLOC_BILL_NO error: %w", err)
    }
    return sct, nil
}

// AllocateSCT cấp số chứng từ hệ thống khác thông qua Stored Procedure ALLOC_SCT (dùng atomic SYS_SEQ_NEXTVAL)
func (r *SaleRepository) AllocateSCT(ctx context.Context, loai, makho string, ngay time.Time) (string, error) {
    var sct string
    row := r.db.WithContext(ctx).Raw("DECLARE @sct nvarchar(50); EXEC ALLOC_SCT @SCT = @sct OUTPUT, @LOAI = ?, @NGAY = ?, @MAKHO = ?; SELECT @sct;", 
        loai, ngay, makho).Row()
    if err := row.Scan(&sct); err != nil {
        return "", fmt.Errorf("ALLOC_SCT error: %w", err)
    }
    return sct, nil
}

// GetByKHOA retrieves a sale by its KHOA
func (r *SaleRepository) GetByKHOA(ctx context.Context, khoa uuid.UUID) (*models.Banle, error) {
    var banle models.Banle
    err := r.db.WithContext(ctx).
        Where("KHOA = ? AND ISNULL(DELETE_BY, 0) = 0", khoa).
        First(&banle).Error
    
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, err
    }
    
    return &banle, nil
}

// GetItemsByKHOA retrieves all items for a sale
func (r *SaleRepository) GetItemsByKHOA(ctx context.Context, khoa uuid.UUID) ([]models.BanleCT, error) {
    var items []models.BanleCT
    err := r.db.WithContext(ctx).
        Where("KHOA = ?", khoa).
        Order("STT").
        Find(&items).Error
    
    if err != nil {
        return nil, err
    }
    
    return items, nil
}

// CreateItem adds a new item to a sale
func (r *SaleRepository) CreateItem(ctx context.Context, item *models.BanleCT) error {
    return r.db.WithContext(ctx).Create(item).Error
}

// UpdateItem updates an existing item
func (r *SaleRepository) UpdateItem(ctx context.Context, item *models.BanleCT) error {
    return r.db.WithContext(ctx).Save(item).Error
}

// UpdateHeaderTotals updates the header totals
func (r *SaleRepository) UpdateHeaderTotals(ctx context.Context, khoa uuid.UUID, totals *HeaderTotals) error {
    return r.db.WithContext(ctx).
        Model(&models.Banle{}).
        Where("KHOA = ?", khoa).
        Updates(map[string]interface{}{
            "SOLUONG":                totals.SOLUONG,
            "THANHTIEN_CHUA_VAT":     totals.THANHTIEN_CHUA_VAT,
            "TIEN_THUE":              totals.TIEN_THUE,
            "THANHTOAN":              totals.THANHTOAN,
            "TIEN_THUE_0":            totals.TIEN_THUE_0,
            "TIEN_THUE_5":            totals.TIEN_THUE_5,
            "TIEN_THUE_8":            totals.TIEN_THUE_8,
            "TIEN_THUE_10":           totals.TIEN_THUE_10,
            "TIEN_THUE_KHAC":         totals.TIEN_THUE_KHAC,
            "THANHTIEN_CHUA_VAT_0":   totals.THANHTIEN_CHUA_VAT_0,
            "THANHTIEN_CHUA_VAT_5":   totals.THANHTIEN_CHUA_VAT_5,
            "THANHTIEN_CHUA_VAT_8":   totals.THANHTIEN_CHUA_VAT_8,
            "THANHTIEN_CHUA_VAT_10":  totals.THANHTIEN_CHUA_VAT_10,
            "THANHTIEN_CHUA_VAT_KHAC": totals.THANHTIEN_CHUA_VAT_KHAC,
            "THANHTOAN_0":            totals.THANHTOAN_0,
            "THANHTOAN_5":            totals.THANHTOAN_5,
            "THANHTOAN_8":            totals.THANHTOAN_8,
            "THANHTOAN_10":           totals.THANHTOAN_10,
            "THANHTOAN_KHAC":         totals.THANHTOAN_KHAC,
            "THANHTOAN_KCT":          totals.THANHTOAN_KCT,
        }).Error
}

type HeaderTotals struct {
    SOLUONG                decimal.Decimal
    THANHTIEN_CHUA_VAT     decimal.Decimal
    TIEN_THUE              decimal.Decimal
    THANHTOAN              decimal.Decimal
    TIEN_THUE_0            decimal.Decimal
    TIEN_THUE_5            decimal.Decimal
    TIEN_THUE_8            decimal.Decimal
    TIEN_THUE_10           decimal.Decimal
    TIEN_THUE_KHAC         decimal.Decimal
    THANHTIEN_CHUA_VAT_0   decimal.Decimal
    THANHTIEN_CHUA_VAT_5   decimal.Decimal
    THANHTIEN_CHUA_VAT_8   decimal.Decimal
    THANHTIEN_CHUA_VAT_10  decimal.Decimal
    THANHTIEN_CHUA_VAT_KHAC decimal.Decimal
    THANHTOAN_0            decimal.Decimal
    THANHTOAN_5            decimal.Decimal
    THANHTOAN_8            decimal.Decimal
    THANHTOAN_10           decimal.Decimal
    THANHTOAN_KHAC         decimal.Decimal
    THANHTOAN_KCT          decimal.Decimal
}
```

---

## 7. SERVICE LAYER (BUSINESS LOGIC)

### 7.1 Discount Service

**File: `internal/service/discount_service.go`**

```go
package service

import (
    "github.com/shopspring/decimal"
    "yourorg/posz/internal/models"
)

type DiscountService struct{}

func NewDiscountService() *DiscountService {
    return &DiscountService{}
}

// CalculateFinalDiscount returns the MAX of all discount sources
// This is CORE business logic for multi-level discounts
func (s *DiscountService) CalculateFinalDiscount(item *models.BanleCT) decimal.Decimal {
    discounts := []decimal.Decimal{
        item.TL_CK2,  // VIP
        item.TL_CK3,  // Group
        item.TL_CK4,  // Promotion
        item.TL_CK5,  // Manual
        item.TL_CK6,  // Additional
        item.TL_CK7,  // Special
    }
    
    max := decimal.Zero
    for _, d := range discounts {
        if d.GreaterThan(max) {
            max = d
        }
    }
    
    return max
}

// CalculatePriceAfterDiscount calculates discounted price
func (s *DiscountService) CalculatePriceAfterDiscount(pricePreTax, discountPercent decimal.Decimal) decimal.Decimal {
    hundred := decimal.NewFromInt(100)
    return pricePreTax.Mul(decimal.NewFromInt(1).Sub(discountPercent.Div(hundred)))
}

// RequiresAuthorization checks if discount needs approval
func (s *DiscountService) RequiresAuthorization(discount, maxAllowed, userLimit decimal.Decimal) bool {
    return discount.GreaterThan(maxAllowed) || discount.GreaterThan(userLimit)
}
```

### 7.2 Tax Service

**File: `internal/service/tax_service.go`**

```go
package service

import (
    "github.com/shopspring/decimal"
    "yourorg/posz/internal/models"
)

type TaxService struct{}

func NewTaxService() *TaxService {
    return &TaxService{}
}

// CalculateLineTax calculates tax for a line item
func (s *TaxService) CalculateLineTax(thanhtienChuaVat, thueSuat decimal.Decimal) decimal.Decimal {
    hundred := decimal.NewFromInt(100)
    tax := thanhtienChuaVat.Mul(thueSuat).Div(hundred)
    return tax.Round(2) // Vietnamese currency - 2 decimal places
}

// CalculatePreTaxAmount calculates pre-tax from price including tax
func (s *TaxService) CalculatePreTaxAmount(priceWithTax, thueSuat decimal.Decimal) decimal.Decimal {
    hundred := decimal.NewFromInt(100)
    divisor := decimal.NewFromInt(1).Add(thueSuat.Div(hundred))
    return priceWithTax.Div(divisor).Round(2)
}

// CalculateTaxBreakdown calculates tax breakdown by rate
func (s *TaxService) CalculateTaxBreakdown(items []models.BanleCT) *TaxBreakdown {
    breakdown := &TaxBreakdown{}
    
    for _, item := range items {
        loaithue := ""
        if item.LOAITHUE != nil {
            loaithue = *item.LOAITHUE
        }
        
        thanhtien := item.THANHTIEN_CHUA_VAT
        tienThue := item.TIEN_THUE
        thanhtoan := item.THANHTOAN
        
        switch loaithue {
        case "T00":
            breakdown.THANHTIEN_CHUA_VAT_0 = breakdown.THANHTIEN_CHUA_VAT_0.Add(thanhtien)
            breakdown.TIEN_THUE_0 = breakdown.TIEN_THUE_0.Add(tienThue)
            breakdown.THANHTOAN_0 = breakdown.THANHTOAN_0.Add(thanhtoan)
        case "T05":
            breakdown.THANHTIEN_CHUA_VAT_5 = breakdown.THANHTIEN_CHUA_VAT_5.Add(thanhtien)
            breakdown.TIEN_THUE_5 = breakdown.TIEN_THUE_5.Add(tienThue)
            breakdown.THANHTOAN_5 = breakdown.THANHTOAN_5.Add(thanhtoan)
        case "T08":
            breakdown.THANHTIEN_CHUA_VAT_8 = breakdown.THANHTIEN_CHUA_VAT_8.Add(thanhtien)
            breakdown.TIEN_THUE_8 = breakdown.TIEN_THUE_8.Add(tienThue)
            breakdown.THANHTOAN_8 = breakdown.THANHTOAN_8.Add(thanhtoan)
        case "T10":
            breakdown.THANHTIEN_CHUA_VAT_10 = breakdown.THANHTIEN_CHUA_VAT_10.Add(thanhtien)
            breakdown.TIEN_THUE_10 = breakdown.TIEN_THUE_10.Add(tienThue)
            breakdown.THANHTOAN_10 = breakdown.THANHTOAN_10.Add(thanhtoan)
        case "KCT":
            breakdown.THANHTOAN_KCT = breakdown.THANHTOAN_KCT.Add(thanhtoan)
        default:
            breakdown.THANHTIEN_CHUA_VAT_KHAC = breakdown.THANHTIEN_CHUA_VAT_KHAC.Add(thanhtien)
            breakdown.TIEN_THUE_KHAC = breakdown.TIEN_THUE_KHAC.Add(tienThue)
            breakdown.THANHTOAN_KHAC = breakdown.THANHTOAN_KHAC.Add(thanhtoan)
        }
    }
    
    return breakdown
}

type TaxBreakdown struct {
    THANHTIEN_CHUA_VAT_0    decimal.Decimal
    THANHTIEN_CHUA_VAT_5    decimal.Decimal
    THANHTIEN_CHUA_VAT_8    decimal.Decimal
    THANHTIEN_CHUA_VAT_10   decimal.Decimal
    THANHTIEN_CHUA_VAT_KHAC decimal.Decimal
    TIEN_THUE_0             decimal.Decimal
    TIEN_THUE_5             decimal.Decimal
    TIEN_THUE_8             decimal.Decimal
    TIEN_THUE_10            decimal.Decimal
    TIEN_THUE_KHAC          decimal.Decimal
    THANHTOAN_0             decimal.Decimal
    THANHTOAN_5             decimal.Decimal
    THANHTOAN_8             decimal.Decimal
    THANHTOAN_10            decimal.Decimal
    THANHTOAN_KHAC          decimal.Decimal
    THANHTOAN_KCT           decimal.Decimal
}
```

### 7.3 Sale Service

**File: `internal/service/sale_service.go`**

```go
package service

import (
    "context"
    "fmt"
    "time"
    "github.com/google/uuid"
    "github.com/shopspring/decimal"
    "yourorg/posz/internal/models"
    "yourorg/posz/internal/repository"
)

type SaleService struct {
    saleRepo        *repository.SaleRepository
    discountService *DiscountService
    taxService      *TaxService
}

func NewSaleService(
    saleRepo *repository.SaleRepository,
    discountService *DiscountService,
    taxService *TaxService,
) *SaleService {
    return &SaleService{
        saleRepo:        saleRepo,
        discountService: discountService,
        taxService:      taxService,
    }
}

// StartSale creates a new sale transaction
func (s *SaleService) StartSale(ctx context.Context, req *StartSaleRequest) (*StartSaleResponse, error) {
    khoa := uuid.New()
    ngay := time.Now().UTC()
    
    sct, err := s.saleRepo.AllocateSCT(ctx, req.Location, ngay)
    if err != nil {
        return nil, fmt.Errorf("failed to allocate SCT: %w", err)
    }
    
    lct := "BANLE"
    banle := &models.Banle{
        KHOA:        khoa,
        LCT:         &lct,
        NGAY:        &ngay,
        SCT:         &sct,
        LOC:         &req.Location,
        QUAY:        &req.Counter,
        CA:          &req.Shift,
        THUNGAN:     stringPtr(fmt.Sprintf("%d", req.CashierUID)),
        CREATE_BY:   &req.CashierUID,
        CREATE_DATE: &ngay,
    }
    
    if err := s.saleRepo.CreateSale(ctx, banle, nil); err != nil {
        return nil, err
    }
    
    return &StartSaleResponse{
        TransactionID: khoa,
        SCT:           sct,
        NGAY:          ngay,
    }, nil
}

// AddItem adds an item to the sale
func (s *SaleService) AddItem(ctx context.Context, req *AddItemRequest) (*AddItemResponse, error) {
    // Get existing items to determine STT
    items, err := s.saleRepo.GetItemsByKHOA(ctx, req.TransactionID)
    if err != nil {
        return nil, err
    }
    
    stt := len(items) + 1
    
    // Calculate pre-tax price
    dongiaC huaVat := s.taxService.CalculatePreTaxAmount(req.DonGia, req.ThueSuat)
    
    // Create item
    item := &models.BanleCT{
        KHOACT:          uuid.New(),
        KHOA:            req.TransactionID,
        STT:             &stt,
        MAVT:            &req.MaVT,
        TENVT:           &req.TenVT,
        DVT:             &req.DVT,
        SOLUONG:         req.SoLuong,
        DONGIA:          req.DonGia,
        DONGIA_CHUA_VAT: dongiaC huaVat,
        LOAITHUE:        &req.LoaiThue,
        THUE_SUAT:       req.ThueSuat,
        CREATE_BY:       &req.UserID,
        CREATE_DATE:     timePtr(time.Now().UTC()),
    }
    
    // Calculate amounts (no discount initially)
    s.calculateItemAmounts(item)
    
    // Save item
    if err := s.saleRepo.CreateItem(ctx, item); err != nil {
        return nil, err
    }
    
    // Recalculate header totals
    if err := s.recalculateHeaderTotals(ctx, req.TransactionID); err != nil {
        return nil, err
    }
    
    return &AddItemResponse{
        Item: item,
    }, nil
}

// calculateItemAmounts calculates all amounts for an item
func (s *SaleService) calculateItemAmounts(item *models.BanleCT) {
    // Price after discount
    dongiaSauCK := s.discountService.CalculatePriceAfterDiscount(
        item.DONGIA_CHUA_VAT,
        item.TL_CK,
    )
    
    // Line total pre-tax
    thanhtienChuaVat := item.SOLUONG.Mul(dongiaSauCK).Round(2)
    
    // Tax
    tienThue := s.taxService.CalculateLineTax(thanhtienChuaVat, item.THUE_SUAT)
    
    // Total
    thanhtoan := thanhtienChuaVat.Add(tienThue)
    
    item.THANHTIEN_CHUA_VAT = thanhtienChuaVat
    item.TIEN_THUE = tienThue
    item.THANHTOAN = thanhtoan
}

// recalculateHeaderTotals recalculates all totals in header from items
func (s *SaleService) recalculateHeaderTotals(ctx context.Context, khoa uuid.UUID) error {
    items, err := s.saleRepo.GetItemsByKHOA(ctx, khoa)
    if err != nil {
        return err
    }
    
    // Sum quantities and amounts
    totals := &repository.HeaderTotals{
        SOLUONG:            decimal.Zero,
        THANHTIEN_CHUA_VAT: decimal.Zero,
        TIEN_THUE:          decimal.Zero,
        THANHTOAN:          decimal.Zero,
    }
    
    for _, item := range items {
        totals.SOLUONG = totals.SOLUONG.Add(item.SOLUONG)
        totals.THANHTIEN_CHUA_VAT = totals.THANHTIEN_CHUA_VAT.Add(item.THANHTIEN_CHUA_VAT)
        totals.TIEN_THUE = totals.TIEN_THUE.Add(item.TIEN_THUE)
        totals.THANHTOAN = totals.THANHTOAN.Add(item.THANHTOAN)
    }
    
    // Tax breakdown
    breakdown := s.taxService.CalculateTaxBreakdown(items)
    totals.TIEN_THUE_0 = breakdown.TIEN_THUE_0
    totals.TIEN_THUE_5 = breakdown.TIEN_THUE_5
    totals.TIEN_THUE_8 = breakdown.TIEN_THUE_8
    totals.TIEN_THUE_10 = breakdown.TIEN_THUE_10
    totals.TIEN_THUE_KHAC = breakdown.TIEN_THUE_KHAC
    totals.THANHTIEN_CHUA_VAT_0 = breakdown.THANHTIEN_CHUA_VAT_0
    totals.THANHTIEN_CHUA_VAT_5 = breakdown.THANHTIEN_CHUA_VAT_5
    totals.THANHTIEN_CHUA_VAT_8 = breakdown.THANHTIEN_CHUA_VAT_8
    totals.THANHTIEN_CHUA_VAT_10 = breakdown.THANHTIEN_CHUA_VAT_10
    totals.THANHTIEN_CHUA_VAT_KHAC = breakdown.THANHTIEN_CHUA_VAT_KHAC
    totals.THANHTOAN_0 = breakdown.THANHTOAN_0
    totals.THANHTOAN_5 = breakdown.THANHTOAN_5
    totals.THANHTOAN_8 = breakdown.THANHTOAN_8
    totals.THANHTOAN_10 = breakdown.THANHTOAN_10
    totals.THANHTOAN_KHAC = breakdown.THANHTOAN_KHAC
    totals.THANHTOAN_KCT = breakdown.THANHTOAN_KCT
    
    return s.saleRepo.UpdateHeaderTotals(ctx, khoa, totals)
}

// Helper functions
func stringPtr(s string) *string { return &s }
func timePtr(t time.Time) *time.Time { return &t }

// Request/Response DTOs
type StartSaleRequest struct {
    Location   string `json:"location" binding:"required"`
    Counter    string `json:"counter" binding:"required"`
    Shift      string `json:"shift" binding:"required"`
    CashierUID int    `json:"cashier_uid" binding:"required"`
}

type StartSaleResponse struct {
    TransactionID uuid.UUID `json:"transaction_id"`
    SCT           string    `json:"sct"`
    NGAY          time.Time `json:"ngay"`
}

type AddItemRequest struct {
    TransactionID uuid.UUID       `json:"transaction_id" binding:"required"`
    MaVT          string          `json:"mavt" binding:"required"`
    TenVT         string          `json:"tenvt" binding:"required"`
    DVT           string          `json:"dvt" binding:"required"`
    SoLuong       decimal.Decimal `json:"soluong" binding:"required"`
    DonGia        decimal.Decimal `json:"dongia" binding:"required"`
    ThueSuat      decimal.Decimal `json:"thue_suat"`
    LoaiThue      string          `json:"loai_thue"`
    UserID        int             `json:"user_id"`
}

type AddItemResponse struct {
    Item *models.BanleCT `json:"item"`
}
```

---

## 8. API HANDLERS (GIN FRAMEWORK)

**File: `internal/handler/sale_handler.go`**

```go
package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "yourorg/posz/internal/service"
)

type SaleHandler struct {
    saleService *service.SaleService
}

func NewSaleHandler(saleService *service.SaleService) *SaleHandler {
    return &SaleHandler{saleService: saleService}
}

// StartSale handles POST /api/pos/sale/start
func (h *SaleHandler) StartSale(c *gin.Context) {
    var req service.StartSaleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }
    
    resp, err := h.saleService.StartSale(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    resp,
    })
}

// AddItem handles POST /api/pos/sale/add-item
func (h *SaleHandler) AddItem(c *gin.Context) {
    var req service.AddItemRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }
    
    resp, err := h.saleService.AddItem(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    resp,
    })
}

// RegisterRoutes registers all sale routes
func (h *SaleHandler) RegisterRoutes(r *gin.RouterGroup) {
    sale := r.Group("/sale")
    {
        sale.POST("/start", h.StartSale)
        sale.POST("/add-item", h.AddItem)
        // ... other routes
    }
}
```

---

## 9. MAIN APPLICATION

**File: `cmd/api/main.go`**

```go
package main

import (
    "log"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "yourorg/posz/internal/database"
    "yourorg/posz/internal/repository"
    "yourorg/posz/internal/service"
    "yourorg/posz/internal/handler"
)

func main() {
    // Load .env
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found")
    }
    
    // Connect to database
    dsn := os.Getenv("DB_DSN")
    if err := database.Connect(dsn); err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer database.Close()
    
    // Initialize repositories
    saleRepo := repository.NewSaleRepository(database.DB)
    
    // Initialize services
    discountService := service.NewDiscountService()
    taxService := service.NewTaxService()
    saleService := service.NewSaleService(saleRepo, discountService, taxService)
    
    // Initialize handlers
    saleHandler := handler.NewSaleHandler(saleService)
    
    // Setup Gin
    r := gin.Default()
    
    // API routes
    api := r.Group("/api/pos")
    {
        saleHandler.RegisterRoutes(api)
        // ... other handlers
    }
    
    // Start server
    port := os.Getenv("SERVER_PORT")
    if port == "" {
        port = "8080"
    }
    
    log.Printf("Server starting on port %s", port)
    if err := r.Run(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
```

---

## 10. TESTING

**File: `internal/service/discount_service_test.go`**

```go
package service_test

import (
    "testing"
    "github.com/shopspring/decimal"
    "github.com/stretchr/testify/assert"
    "yourorg/posz/internal/models"
    "yourorg/posz/internal/service"
)

func TestCalculateFinalDiscount(t *testing.T) {
    svc := service.NewDiscountService()
    
    item := &models.BanleCT{
        TL_CK2: decimal.NewFromInt(5),   // VIP
        TL_CK3: decimal.NewFromInt(10),  // Group
        TL_CK4: decimal.NewFromInt(15),  // Promo (MAX)
        TL_CK5: decimal.NewFromInt(7),   // Manual
    }
    
    result := svc.CalculateFinalDiscount(item)
    assert.Equal(t, decimal.NewFromInt(15), result)
}

func TestCalculatePriceAfterDiscount(t *testing.T) {
    svc := service.NewDiscountService()
    
    price := decimal.NewFromInt(100000)
    discount := decimal.NewFromInt(10)
    
    result := svc.CalculatePriceAfterDiscount(price, discount)
    expected := decimal.NewFromInt(90000)
    
    assert.True(t, result.Equal(expected))
}
```

---

## 11. DEPLOYMENT

### Build
```bash
go build -o posz ./cmd/api
```

### Run
```bash
export DB_DSN="sqlserver://sa:password@localhost:1433?database=SOFTZ_SHOPHOA"
export JWT_SECRET="your-secret"
./posz
```

### Docker
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o posz ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/posz .
EXPOSE 8080
CMD ["./posz"]
```

---

**Document Version:** 1.0 (Golang)  
**Last Updated:** 2026-10-14  
**Author:** Technical Architect (AI-assisted)  
**Status:** Ready for Implementation
