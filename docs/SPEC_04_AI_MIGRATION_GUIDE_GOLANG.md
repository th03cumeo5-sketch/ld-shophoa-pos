# AI AGENTS MIGRATION GUIDE - GOLANG
## POSz System - Delphi/Pascal to Go Migration

---

## DOCUMENT PURPOSE

This document provides AI coding agents with **practical, actionable guidance** for migrating the POSz (Point of Sale) system from Delphi/Pascal to **Go (Golang)**. It includes:

1. **Step-by-step implementation order**
2. **Go-specific code patterns and examples**
3. **Common pitfalls and solutions**
4. **Testing strategies**
5. **Validation checklists**

**Target Audience:** AI Agents (Claude, GPT, Copilot, etc.) performing code generation and migration tasks.

---

## TABLE OF CONTENTS

1. [Migration Overview](#1-migration-overview)
2. [Prerequisites & Setup](#2-prerequisites--setup)
3. [Phase 1: Database Layer](#3-phase-1-database-layer)
4. [Phase 2: Business Logic](#4-phase-2-business-logic)
5. [Phase 3: API Layer](#5-phase-3-api-layer)
6. [Phase 4: Integration](#6-phase-4-integration)
7. [Testing Guidelines](#7-testing-guidelines)
8. [Common Patterns](#8-common-patterns)
9. [Troubleshooting](#9-troubleshooting)

---

## 1. MIGRATION OVERVIEW

### 1.1 What to Migrate

**IN SCOPE:**
- ✅ POS sale processing (BANLE, BANLE_CT)
- ✅ Return/refund handling (TRAHANG, TRAHANG_CT)
- ✅ Discount calculation (multi-level TL_CK logic)
- ✅ Promotion application (KHUYENMAI)
- ✅ Tax calculation (breakdown by rate)
- ✅ Payment processing (split payment)
- ✅ Receipt generation
- ✅ Product lookup (barcode, code, name)
- ✅ Customer lookup

**OUT OF SCOPE:**
- ❌ VIP/Loyalty management (NXTz module)
- ❌ Inventory management (NXTz module)
- ❌ HR/Payroll/Attendance (NXTz module)
- ❌ Purchase orders (NXTz module)
- ❌ eInvoice BKAV integration (separate workflow)

### 1.2 Why Go?

- **Performance:** Native compilation, efficient concurrency
- **Simplicity:** Clear syntax, easy to learn
- **Standard library:** Excellent HTTP and database support
- **Deployment:** Single binary, easy to distribute
- **Tooling:** Built-in formatting, testing, profiling

### 1.3 Success Criteria

- [ ] All core POS functions work identically to Delphi version
- [ ] Performance: Sub-second response for all operations
- [ ] Data integrity: 100% match on financial calculations
- [ ] Backward compatibility: Can read existing database without migration
- [ ] Zero data loss during transition

---

## 2. PREREQUISITES & SETUP

### 2.1 Environment Setup

```bash
# Install Go (version 1.21 or later)
# Download from https://go.dev/dl/

# Verify installation
go version

# Create project
mkdir posz
cd posz
go mod init github.com/yourorg/posz

# Initialize directory structure
mkdir -p cmd/api internal/{config,models,repository,service,handler,middleware,database,utils} test/integration
```

### 2.2 Required Dependencies

```bash
# Add dependencies
go get github.com/gin-gonic/gin@v1.9.1
go get github.com/denisenkom/go-mssqldb@v0.12.3
go get gorm.io/gorm@v1.25.5
go get gorm.io/driver/sqlserver@v1.5.2
go get github.com/shopspring/decimal@v1.3.1
go get github.com/google/uuid@v1.4.0
go get github.com/golang-jwt/jwt/v5@v5.1.0
go get golang.org/x/crypto@v0.15.0
go get github.com/spf13/viper@v1.17.0
go get github.com/joho/godotenv@v1.5.1
go get github.com/go-playground/validator/v10@v10.16.0
go get github.com/sirupsen/logrus@v1.9.3
go get github.com/stretchr/testify@v1.8.4
```

### 2.3 Project Structure

Create the following structure:
```
posz/
├── go.mod
├── go.sum
├── .env
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── models/
│   │   ├── sale.go
│   │   ├── return.go
│   │   ├── promotion.go
│   │   ├── product.go
│   │   └── customer.go
│   ├── repository/
│   │   ├── sale_repo.go
│   │   ├── product_repo.go
│   │   ├── promotion_repo.go
│   │   └── customer_repo.go
│   ├── service/
│   │   ├── sale_service.go
│   │   ├── discount_service.go
│   │   ├── tax_service.go
│   │   └── promotion_service.go
│   ├── handler/
│   │   ├── sale_handler.go
│   │   ├── return_handler.go
│   │   └── product_handler.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── logger.go
│   ├── database/
│   │   └── db.go
│   └── utils/
│       ├── guid.go
│       └── decimal.go
└── test/
    └── integration/
```

### 2.4 Configuration File

**File: `.env`**
```
DB_DSN=sqlserver://sa:YourPassword@localhost:1433?database=SOFTZ_SHOPHOA&encrypt=disable
JWT_SECRET=your-secret-key-here
SERVER_PORT=8080
GIN_MODE=release
LOG_LEVEL=info
```

---

## 3. PHASE 1: DATABASE LAYER

### 3.1 Database Connection

**File: `internal/database/db.go`**

```go
package database

import (
    "fmt"
    "log"
    "time"
    "gorm.io/driver/sqlserver"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect establishes database connection with connection pooling
func Connect(dsn string) error {
    var err error
    
    DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
        // Disable default transaction for better performance
        SkipDefaultTransaction: true,
        // Prepare statements for better performance
        PrepareStmt: true,
    })
    
    if err != nil {
        return fmt.Errorf("failed to connect to database: %w", err)
    }
    
    // Get underlying SQL DB for connection pool settings
    sqlDB, err := DB.DB()
    if err != nil {
        return err
    }
    
    // Connection pool settings
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    sqlDB.SetConnMaxIdleTime(10 * time.Minute)
    
    // Test connection
    if err := sqlDB.Ping(); err != nil {
        return fmt.Errorf("failed to ping database: %w", err)
    }
    
    log.Println("✓ Database connection established")
    return nil
}

// Close closes the database connection
func Close() error {
    sqlDB, err := DB.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}
```

### 3.2 Define Models

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
    KHOA               uuid.UUID       `gorm:"column:KHOA;type:uniqueidentifier;primaryKey"`
    LCT                *string         `gorm:"column:LCT;type:nvarchar(5)"`
    NGAY               *time.Time      `gorm:"column:NGAY;type:datetime"`
    SCT                *string         `gorm:"column:SCT;type:nvarchar(20)"`
    LOC                *string         `gorm:"column:LOC;type:nvarchar(2)"`
    MAKHO              *string         `gorm:"column:MAKHO;type:nvarchar(2)"`
    QUAY               *string         `gorm:"column:QUAY;type:nvarchar(15)"`
    CA                 *string         `gorm:"column:CA;type:nvarchar(5)"`
    THUNGAN            *string         `gorm:"column:THUNGAN;type:nvarchar(15)"`
    MADT               *string         `gorm:"column:MADT;type:nvarchar(15)"`
    
    // Amounts - Use decimal.Decimal for financial calculations
    SOLUONG            decimal.Decimal `gorm:"column:SOLUONG;type:float"`
    THANHTIEN_CHUA_VAT decimal.Decimal `gorm:"column:THANHTIEN_CHUA_VAT;type:float"`
    TIEN_THUE          decimal.Decimal `gorm:"column:TIEN_THUE;type:float"`
    THANHTOAN          decimal.Decimal `gorm:"column:THANHTOAN;type:float"`
    
    // Tax breakdown
    TIEN_THUE_0        decimal.Decimal `gorm:"column:TIEN_THUE_0;type:float"`
    TIEN_THUE_5        decimal.Decimal `gorm:"column:TIEN_THUE_5;type:float"`
    TIEN_THUE_8        decimal.Decimal `gorm:"column:TIEN_THUE_8;type:float"`
    TIEN_THUE_10       decimal.Decimal `gorm:"column:TIEN_THUE_10;type:float"`
    TIEN_THUE_KHAC     decimal.Decimal `gorm:"column:TIEN_THUE_KHAC;type:float"`
    
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
    PRINTED    *bool      `gorm:"column:PRINTED;type:bit"`
    PRINT_NO   *int       `gorm:"column:PRINT_NO;type:int"`
    
    // Audit trail
    CREATE_BY   *int       `gorm:"column:CREATE_BY;type:int"`
    CREATE_DATE *time.Time `gorm:"column:CREATE_DATE;type:datetime"`
    UPDATE_BY   *int       `gorm:"column:UPDATE_BY;type:int"`
    UPDATE_DATE *time.Time `gorm:"column:UPDATE_DATE;type:datetime"`
    DELETE_BY   *int       `gorm:"column:DELETE_BY;type:int"`
    DELETE_DATE *time.Time `gorm:"column:DELETE_DATE;type:datetime"`
    
    // Relationships
    Items []BanleCT `gorm:"foreignKey:KHOA;references:KHOA"`
}

// TableName specifies the table name for GORM
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
    
    // Multi-level discounts (CRITICAL BUSINESS LOGIC)
    TL_CK              decimal.Decimal `gorm:"column:TL_CK;type:float"`        // Final (MAX)
    TL_CK2             decimal.Decimal `gorm:"column:TL_CK2;type:float"`       // VIP
    TL_CK3             decimal.Decimal `gorm:"column:TL_CK3;type:float"`       // Group
    TL_CK4             decimal.Decimal `gorm:"column:TL_CK4;type:float"`       // Promo
    TL_CK5             decimal.Decimal `gorm:"column:TL_CK5;type:float"`       // Manual
    TL_CK6             decimal.Decimal `gorm:"column:TL_CK6;type:float"`       // Additional
    TL_CK7             decimal.Decimal `gorm:"column:TL_CK7;type:float"`       // Special
    TL_CK_MAX          decimal.Decimal `gorm:"column:TL_CK_MAX;type:float"`    // Max allowed
    
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

// Helper functions for pointer types
func StringPtr(s string) *string       { return &s }
func IntPtr(i int) *int                { return &i }
func BoolPtr(b bool) *bool             { return &b }
func TimePtr(t time.Time) *time.Time   { return &t }
```

### 3.3 Repository Implementation

**File: `internal/repository/sale_repo.go`**

```go
package repository

import (
    "context"
    "fmt"
    "time"
    "github.com/google/uuid"
    "github.com/shopspring/decimal"
    "gorm.io/gorm"
    "yourorg/posz/internal/models"
)

type SaleRepository struct {
    db *gorm.DB
}

func NewSaleRepository(db *gorm.DB) *SaleRepository {
    return &SaleRepository{db: db}
}

// CreateSale creates a new sale with items in a transaction
func (r *SaleRepository) CreateSale(ctx context.Context, banle *models.Banle, items []models.BanleCT) error {
    return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // Insert header
        if err := tx.Create(banle).Error; err != nil {
            return fmt.Errorf("failed to create banle: %w", err)
        }
        
        // Insert items if any
        if len(items) > 0 {
            if err := tx.Create(&items).Error; err != nil {
                return fmt.Errorf("failed to create banle items: %w", err)
            }
        }
        
        return nil
    })
}

// GetByKHOA retrieves sale by KHOA
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
    
    return items, err
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

// CompleteSalePayment thực hiện thanh toán chính thức bằng Stored Procedure POS_PAID(@KHOA)
// Thủ tục này tự động chuyển BANLE_TAM -> BANLE, BANLE_TAM_CT -> BANLE_CT, xóa bảng tạm và merge ehoadon
func (r *SaleRepository) CompleteSalePayment(ctx context.Context, khoa uuid.UUID) error {
    err := r.db.WithContext(ctx).Exec("EXEC POS_PAID @KHOA = ?", khoa).Error
    if err != nil {
        return fmt.Errorf("POS_PAID failed for khoa %s: %w", khoa.String(), err)
    }
    return nil
}

// CreateItem adds a new item
func (r *SaleRepository) CreateItem(ctx context.Context, item *models.BanleCT) error {
    return r.db.WithContext(ctx).Create(item).Error
}

// UpdateItem updates an existing item
func (r *SaleRepository) UpdateItem(ctx context.Context, item *models.BanleCT) error {
    return r.db.WithContext(ctx).Save(item).Error
}

// UpdateHeaderTotals updates header totals (all amounts)
func (r *SaleRepository) UpdateHeaderTotals(ctx context.Context, khoa uuid.UUID, totals *HeaderTotals) error {
    updates := map[string]interface{}{
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
    }
    
    return r.db.WithContext(ctx).
        Model(&models.Banle{}).
        Where("KHOA = ?", khoa).
        Updates(updates).Error
}

// HeaderTotals contains all totals for a sale header
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

## 4. PHASE 2: BUSINESS LOGIC

### 4.1 Discount Service

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

// CalculateFinalDiscount returns MAX of all discount sources
// CORE BUSINESS LOGIC for multi-level discounts
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

// CalculatePriceAfterDiscount calculates price after discount percentage
func (s *DiscountService) CalculatePriceAfterDiscount(pricePreTax, discountPercent decimal.Decimal) decimal.Decimal {
    hundred := decimal.NewFromInt(100)
    discountFactor := decimal.NewFromInt(1).Sub(discountPercent.Div(hundred))
    return pricePreTax.Mul(discountFactor)
}

// RequiresAuthorization checks if discount needs approval
func (s *DiscountService) RequiresAuthorization(discount, maxAllowed, userLimit decimal.Decimal) bool {
    return discount.GreaterThan(maxAllowed) || discount.GreaterThan(userLimit)
}
```

### 4.2 Tax Service

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
// CRITICAL: Round to 2 decimal places (Vietnamese currency)
func (s *TaxService) CalculateLineTax(thanhtienChuaVat, thueSuat decimal.Decimal) decimal.Decimal {
    hundred := decimal.NewFromInt(100)
    tax := thanhtienChuaVat.Mul(thueSuat).Div(hundred)
    return tax.Round(2)
}

// CalculatePreTaxAmount calculates pre-tax amount from price including tax
func (s *TaxService) CalculatePreTaxAmount(priceWithTax, thueSuat decimal.Decimal) decimal.Decimal {
    hundred := decimal.NewFromInt(100)
    divisor := decimal.NewFromInt(1).Add(thueSuat.Div(hundred))
    return priceWithTax.Div(divisor).Round(2)
}

// CalculateTaxBreakdown calculates tax breakdown by rate for entire transaction
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

---

## 5. TESTING

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
    
    tests := []struct {
        name     string
        item     *models.BanleCT
        expected decimal.Decimal
    }{
        {
            name: "should return max discount from all sources",
            item: &models.BanleCT{
                TL_CK2: decimal.NewFromInt(5),   // VIP
                TL_CK3: decimal.NewFromInt(10),  // Group
                TL_CK4: decimal.NewFromInt(15),  // Promo (MAX)
                TL_CK5: decimal.NewFromInt(7),   // Manual
            },
            expected: decimal.NewFromInt(15),
        },
        {
            name: "should return zero when all discounts are zero",
            item: &models.BanleCT{
                TL_CK2: decimal.Zero,
                TL_CK3: decimal.Zero,
                TL_CK4: decimal.Zero,
                TL_CK5: decimal.Zero,
            },
            expected: decimal.Zero,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := svc.CalculateFinalDiscount(tt.item)
            assert.True(t, result.Equal(tt.expected), "expected %v, got %v", tt.expected, result)
        })
    }
}

func TestCalculatePriceAfterDiscount(t *testing.T) {
    svc := service.NewDiscountService()
    
    price := decimal.NewFromInt(100000)
    discount := decimal.NewFromInt(10)
    
    result := svc.CalculatePriceAfterDiscount(price, discount)
    expected := decimal.NewFromInt(90000)
    
    assert.True(t, result.Equal(expected), "expected %v, got %v", expected, result)
}
```

---

## 6. CRITICAL GO-SPECIFIC PATTERNS

### 6.1 Decimal Math (CRITICAL!)

**❌ WRONG - Never use float64 for money:**
```go
price := 100000.0
quantity := 2.5
total := price * quantity  // WRONG! Floating point errors
```

**✅ CORRECT - Always use shopspring/decimal:**
```go
price := decimal.NewFromInt(100000)
quantity := decimal.NewFromFloat(2.5)
total := price.Mul(quantity)  // CORRECT!
```

### 6.2 Pointer Handling

Go requires pointers for nullable database fields:

```go
// Creating with pointers
func StringPtr(s string) *string { return &s }

banle := &models.Banle{
    LOC: StringPtr("01"),  // nullable nvarchar in DB
    // ...
}

// Reading from pointers
if banle.LOC != nil {
    location := *banle.LOC
}
```

### 6.3 Error Handling

Go uses explicit error returns:

```go
func (s *SaleService) StartSale(ctx context.Context, req *StartSaleRequest) (*StartSaleResponse, error) {
    sct, err := s.saleRepo.AllocateSCT(ctx, req.Location, time.Now())
    if err != nil {
        return nil, fmt.Errorf("failed to allocate SCT: %w", err)
    }
    
    // ... continue
}
```

### 6.4 Context Management

Always pass `context.Context` for cancellation:

```go
func (r *SaleRepository) GetByKHOA(ctx context.Context, khoa uuid.UUID) (*models.Banle, error) {
    var banle models.Banle
    err := r.db.WithContext(ctx).  // Pass context to DB
        Where("KHOA = ?", khoa).
        First(&banle).Error
    
    return &banle, err
}
```

---

## 7. DEPLOYMENT

### 7.1 Build

```bash
# Development
go run cmd/api/main.go

# Production
go build -o posz ./cmd/api
./posz
```

### 7.2 Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o posz ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/posz .
EXPOSE 8080
CMD ["./posz"]
```

---

## 8. VALIDATION CHECKLIST

- [ ] Use `decimal.Decimal` for ALL financial calculations
- [ ] Handle NULL database fields with pointers
- [ ] Pass `context.Context` to all repository methods
- [ ] Use transactions for multi-table operations
- [ ] Implement proper error wrapping with `fmt.Errorf`
- [ ] Test discount calculation: TL_CK = MAX(TL_CK2...TL_CK7)
- [ ] Test tax calculation matches Delphi to the cent
- [ ] Verify connection pooling is configured
- [ ] Check GORM table names match SQL Server
- [ ] Validate UUID generation and parsing

---

**Document Version:** 1.0 (Golang)  
**Last Updated:** 2026-10-14  
**Author:** AI Migration Architect  
**Status:** Ready for Implementation
