# DATABASE SCHEMA REFERENCE
## POSz - Complete Database Structure

---

## TABLE OF CONTENTS
1. [Core Transaction Tables](#1-core-transaction-tables)
2. [Master Data Tables](#2-master-data-tables)
3. [System Tables](#3-system-tables)
4. [Views](#4-views)
5. [Stored Procedures](#5-stored-procedures)
6. [Functions](#6-functions)
7. [Triggers](#7-triggers)
8. [Field Naming Conventions](#8-field-naming-conventions)

---

## 1. CORE TRANSACTION TABLES

### 1.1 BANLE & BANLE_TAM (Retail Sales Header)

**Purpose:** 
- `BANLE_TAM`: Lưu trữ toàn bộ hóa đơn đang bán dở, giỏ hàng đang quét, và các hóa đơn tạm (treo đơn `F10`) tại các quầy thu ngân.
- `BANLE`: Bảng chính thức lưu trữ các hóa đơn đã thanh toán thành công (được move qua từ `BANLE_TAM` thông qua Stored Procedure `POS_PAID`).
- **Schema:** `BANLE_TAM` và `BANLE` có cấu trúc cột hoàn toàn giống nhau 100%.

**Primary Key:** KHOA (uniqueidentifier)

| Column | Type | Nullable | Description |
|--------|------|----------|-------------|
| KHOA | uniqueidentifier | NO | Transaction ID (GUID), Primary Key |
| LCT | nvarchar(5) | YES | Transaction type code ('BANLE') |
| NGAY | datetime | YES | Transaction date & time |
| SCT | nvarchar(20) | YES | Receipt number (format: LOC-YYYYMMDD-NNNN) |
| LOC | nvarchar(2) | YES | Location/Store code, FK to DM_LOCATION |
| MAKHO | nvarchar(2) | YES | Warehouse code |
| QUAY | nvarchar(15) | YES | Counter/POS terminal code |
| CA | nvarchar(5) | YES | Shift code (CA1, CA2, CA3) |
| THUNGAN | nvarchar(15) | YES | Cashier user ID |
| MAVIP | nvarchar(20) | YES | VIP customer code (reference only) |
| MADT | nvarchar(15) | YES | Customer code, FK to DM_KH |
| SOLUONG | float | YES | Total quantity of all items |
| SOTIEN | float | YES | Subtotal (sum of line amounts pre-tax) |
| SOTIEN1 | float | YES | Alternate subtotal calculation |
| CHIETKHAU_MH | float | YES | Total item-level discount amount |
| CK_BY | int | YES | User ID who authorized discounts |
| TL_CK | float | YES | Header-level discount % |
| TL_CK2 | float | YES | Secondary discount % (split payment calc) |
| TL_CK_TAM | float | YES | Temporary discount holder |
| CHIETKHAU | float | YES | Header-level discount amount |
| THUE | float | YES | **DEPRECATED** - Use TIEN_THUE instead |
| THANHTIEN_CHUA_VAT | float | YES | Total pre-tax amount (after all discounts) |
| THANHTIEN_CHUA_VAT_CL | float | YES | Adjusted pre-tax amount (rounding) |
| TIEN_THUE | float | YES | Total tax amount (sum of all tax rates) |
| TIEN_THUE_0 | float | YES | Tax amount for 0% rate |
| TIEN_THUE_5 | float | YES | Tax amount for 5% rate |
| TIEN_THUE_8 | float | YES | Tax amount for 8% rate |
| TIEN_THUE_10 | float | YES | Tax amount for 10% rate |
| TIEN_THUE_KHAC | float | YES | Tax amount for other rates |
| TIEN_THUE_CHUA_CL | float | YES | Tax before rounding adjustment |
| TIEN_THUE_CL | float | YES | Tax rounding adjustment amount |
| THANHTIEN_CHUA_VAT_0 | float | YES | Pre-tax amount for 0% tax items |
| THANHTIEN_CHUA_VAT_5 | float | YES | Pre-tax amount for 5% tax items |
| THANHTIEN_CHUA_VAT_8 | float | YES | Pre-tax amount for 8% tax items |
| THANHTIEN_CHUA_VAT_10 | float | YES | Pre-tax amount for 10% tax items |
| THANHTIEN_CHUA_VAT_KHAC | float | YES | Pre-tax amount for other tax rates |
| THANHTOAN | float | YES | **Grand total to pay** (with tax, after discounts) |
| THANHTOAN_0 | float | YES | Total for 0% tax items |
| THANHTOAN_5 | float | YES | Total for 5% tax items |
| THANHTOAN_8 | float | YES | Total for 8% tax items |
| THANHTOAN_10 | float | YES | Total for 10% tax items |
| THANHTOAN_KHAC | float | YES | Total for other tax rates |
| THANHTOAN_KCT | float | YES | Total for tax-exempt items (KCT) |
| PTTT | nvarchar(2) | YES | Payment method 1 (TM, CK, CN, TANG, BIEU) |
| TTOAN1 | float | YES | Payment amount 1 |
| PTTT2 | nvarchar(2) | YES | Payment method 2 (split payment) |
| TTOAN2 | float | YES | Payment amount 2 |
| CHUATHOI | float | YES | Change to return = (TTOAN1 + TTOAN2) - THANHTOAN |
| SO_CHUANCHI | nvarchar(20) | YES | Check/Transfer reference number (for CK payment) |
| TTTT | nvarchar(10) | YES | Payment status |
| SPTT | nvarchar(10) | YES | **DEPRECATED** - Payment sub-type (legacy) |
| PRINTED | bit | YES | Has receipt been printed? (0=No, 1=Yes) |
| PRINT_NO | int | YES | Number of times receipt has been printed |
| CN_TENDV | nvarchar(200) | YES | Customer company name (for CN - credit sales) |
| CN_MST | nvarchar(30) | YES | Customer tax ID |
| CN_DIACHI | nvarchar(200) | YES | Customer address |
| CN_DIACHI_HD | nvarchar(200) | YES | Customer invoice address |
| CN_LIENHE | nvarchar(100) | YES | Contact person |
| CN_DTHOAI | nvarchar(100) | YES | Contact phone |
| CN_EMAIL | nvarchar(100) | YES | Contact email |
| CN_MATK | nvarchar(50) | YES | Customer bank account number |
| CREATE_BY | int | YES | User ID who created this transaction |
| CREATE_DATE | datetime | YES | Creation timestamp |
| UPDATE_BY | int | YES | User ID who last updated |
| UPDATE_DATE | datetime | YES | Last update timestamp |
| DELETE_BY | int | YES | User ID who deleted (soft delete) |
| DELETE_DATE | datetime | YES | Deletion timestamp |
| DRC_STATUS | nvarchar(1) | YES | Data replication/sync status |
| QUAYKE | nvarchar(100) | YES | Shelf/Display location |
| DGIAI | ntext | YES | Notes/Description |
| ROWGUID | uniqueidentifier | NO | Row GUID (for replication) |

**Calculated Fields (not in DB, computed in views/queries):**
- LK_TENDT: Customer name (from DM_KH)
- LK_TENKHO: Warehouse name (from DM_LOCATION)
- LK_TENTK: Account holder name
- LK_DAIDIEN: Representative name
- LK_NGANHANG: Bank name
- LK_PTTT: Payment method description
- LK_PTTT2: Secondary payment method description
- LK_TTTT: Payment status description
- CALC_TONGCK: Calculated total discount
- CALC_TIENTHOI: Calculated change amount

**Business Rules:**
1. KHOA must be unique GUID (use NEWID() or uuid v4)
2. SCT format: `{LOC}-{YYYYMMDD}-{NNNN}` (e.g., "01-20241014-0001")
3. THANHTOAN must equal sum of THANHTOAN_0 + THANHTOAN_5 + THANHTOAN_8 + THANHTOAN_10 + THANHTOAN_KHAC + THANHTOAN_KCT
4. TTOAN1 + TTOAN2 >= THANHTOAN (payment must cover total)
5. CHUATHOI = (TTOAN1 + TTOAN2) - THANHTOAN
6. SOLUONG = sum of BANLE_CT.SOLUONG
7. If DELETE_BY IS NOT NULL, transaction is considered deleted (soft delete)
8. Once PRINTED = 1, transaction cannot be deleted (business rule, not DB constraint)

**Indexes:**
```sql
CREATE INDEX IX_BANLE_NGAY ON BANLE(NGAY);
CREATE INDEX IX_BANLE_LOC_NGAY ON BANLE(LOC, NGAY);
CREATE INDEX IX_BANLE_QUAY_CA ON BANLE(QUAY, CA, NGAY);
CREATE INDEX IX_BANLE_MADT ON BANLE(MADT);
CREATE INDEX IX_BANLE_SCT ON BANLE(SCT);
CREATE INDEX IX_BANLE_CREATE_DATE ON BANLE(CREATE_DATE);
CREATE INDEX IX_BANLE_DRC_STATUS ON BANLE(DRC_STATUS);
```

---

### 1.2 BANLE_CT (Retail Sales Detail)

**Purpose:** Stores line items for each POS sale

**Primary Key:** KHOACT (uniqueidentifier)
**Foreign Key:** KHOA → BANLE(KHOA)

| Column | Type | Nullable | Description |
|--------|------|----------|-------------|
| KHOACT | uniqueidentifier | NO | Detail line ID (GUID), Primary Key |
| KHOA | uniqueidentifier | YES | FK to BANLE.KHOA |
| STT | int | YES | Line number (display order, 1, 2, 3...) |
| RSTT | int | YES | Reverse line number (for UI display) |
| MAVT | nvarchar(15) | YES | Product code, FK to DM_VT |
| TENVT | nvarchar(200) | YES | Product name (snapshot from DM_VT) |
| TENVT2 | nvarchar(200) | YES | Short product name |
| TENTAT | nvarchar(50) | YES | Product abbreviation |
| DVT | nvarchar(20) | YES | Unit of measure |
| MABO | nvarchar(15) | YES | Bundle product code (if this item is part of a bundle) |
| SOLUONG | float | YES | Quantity sold |
| DONGIA | float | YES | Unit price (including tax) |
| DONGIA_CHUA_VAT | float | YES | Unit price before tax = DONGIA / (1 + THUE_SUAT/100) |
| DONGIA_SAUCK | float | YES | Unit price after discount = DONGIA_CHUA_VAT * (1 - TL_CK/100) |
| TL_CK | float | YES | **Applied discount %** = MAX(TL_CK2...TL_CK7) |
| TL_CK2 | float | YES | VIP customer discount % |
| TL_CK3 | float | YES | Product group discount % |
| TL_CK4 | float | YES | Promotion discount % |
| TL_CK5 | float | YES | Manual discount % (entered by cashier) |
| TL_CK6 | float | YES | Additional discount % |
| TL_CK7 | float | YES | Special discount % |
| TL_CK_MAX | float | YES | Maximum allowed discount % (from DM_VT) |
| TL_CK_THEM | float | YES | Excess discount beyond max (requires approval) |
| CHIETKHAU | float | YES | Discount amount = SOLUONG * DONGIA_CHUA_VAT * TL_CK/100 |
| CK_BY | int | YES | User ID who authorized the discount |
| LOAITHUE | nvarchar(15) | YES | Tax type code (T00, T05, T08, T10, KCT, TTT) |
| THUE_SUAT | float | YES | Tax rate % |
| TIEN_THUE | float | YES | Tax amount = THANHTIEN_CHUA_VAT * THUE_SUAT/100 |
| TIEN_THUE_CHUA_CL | float | YES | Tax before rounding adjustment |
| TIEN_THUE_CL | float | YES | Tax rounding adjustment |
| THANHTIEN_CHUA_VAT | float | YES | Line total pre-tax = SOLUONG * DONGIA_SAUCK |
| THANHTIEN_CHUA_VAT_CL | float | YES | Adjusted pre-tax total (rounding) |
| SOTIEN | float | YES | **DEPRECATED** - Use THANHTOAN instead |
| SOTIEN1 | float | YES | Alternate total calculation |
| THANHTOAN | float | YES | **Line amount to pay** = THANHTIEN_CHUA_VAT + TIEN_THUE |
| TRA_DATE | datetime | YES | Return deadline date (e.g., 30 days from BANLE.NGAY) |
| TRA_BY | int | YES | User ID who processed return (if returned) |
| GHICHU | ntext | YES | Line item notes |
| CREATE_BY | int | YES | Created by user ID |
| CREATE_DATE | datetime | YES | Creation timestamp |
| UPDATE_BY | int | YES | Last updated by user ID |
| UPDATE_DATE | datetime | YES | Last update timestamp |
| ROWGUID | uniqueidentifier | NO | Row GUID |

**Calculated Fields:**
- LK_TENVT: Product name (from DM_VT, if different from snapshot)
- LK_DVT: Unit name
- LK_TENVT2: Short name
- CALC_SOTIEN_SAUCK: Amount after discount

**Business Rules:**
1. KHOACT must be unique GUID
2. KHOA must reference existing BANLE.KHOA
3. STT should be sequential within same KHOA (1, 2, 3...)
4. TL_CK = MAX(TL_CK2, TL_CK3, TL_CK4, TL_CK5, TL_CK6, TL_CK7)
5. If TL_CK > TL_CK_MAX, then TL_CK_THEM = TL_CK - TL_CK_MAX, requires approval
6. THANHTOAN can be 0 for free items (promotional gifts) - mark as ItemTypeID=15 for eInvoice
7. DONGIA_CHUA_VAT formula: `DONGIA / (1 + THUE_SUAT/100)`
8. THANHTIEN_CHUA_VAT formula: `SOLUONG * DONGIA_SAUCK`
9. TIEN_THUE formula: `THANHTIEN_CHUA_VAT * THUE_SUAT/100`, rounded to 2 decimals
10. THANHTOAN formula: `THANHTIEN_CHUA_VAT + TIEN_THUE`

**Indexes:**
```sql
CREATE INDEX IX_BANLECT_KHOA ON BANLE_CT(KHOA);
CREATE INDEX IX_BANLECT_MAVT ON BANLE_CT(MAVT);
CREATE INDEX IX_BANLECT_STT ON BANLE_CT(KHOA, STT);
```

---

### 1.3 TRAHANG (Returns Header)

**Purpose:** Stores return/refund transaction headers

**Primary Key:** KHOA (uniqueidentifier)
**Foreign Key:** KHOA2 → BANLE(KHOA)

| Column | Type | Nullable | Description |
|--------|------|----------|-------------|
| KHOA | uniqueidentifier | NO | Return transaction ID, Primary Key |
| KHOA2 | uniqueidentifier | YES | **FK to original sale** BANLE.KHOA |
| LCT | nvarchar(5) | YES | Transaction type ('TRAH' for return) |
| NGAY | datetime | YES | Return date & time |
| SCT | nvarchar(20) | YES | Return receipt number (format: LOC-YYYYMMDD-TNNNN) |
| LOC | nvarchar(2) | YES | Location code |
| MAKHO | nvarchar(2) | YES | Warehouse code |
| QUAY | nvarchar(15) | YES | Counter code |
| CA | nvarchar(5) | YES | Shift code |
| THUNGAN | nvarchar(15) | YES | Cashier user ID |
| MADT | nvarchar(15) | YES | Customer code |
| SOLUONG | float | YES | Total quantity returned (**negative value**) |
| THANHTIEN_CHUA_VAT | float | YES | Total pre-tax amount (**negative**) |
| TIEN_THUE | float | YES | Total tax amount (**negative**) |
| THANHTOAN | float | YES | **Refund amount (negative)** |
| TIEN_THUE_0 | float | YES | Tax at 0% (**negative**) |
| TIEN_THUE_5 | float | YES | Tax at 5% (**negative**) |
| TIEN_THUE_8 | float | YES | Tax at 8% (**negative**) |
| TIEN_THUE_10 | float | YES | Tax at 10% (**negative**) |
| TIEN_THUE_KHAC | float | YES | Tax at other rates (**negative**) |
| THANHTIEN_CHUA_VAT_0 | float | YES | Pre-tax for 0% (**negative**) |
| THANHTIEN_CHUA_VAT_5 | float | YES | Pre-tax for 5% (**negative**) |
| THANHTIEN_CHUA_VAT_8 | float | YES | Pre-tax for 8% (**negative**) |
| THANHTIEN_CHUA_VAT_10 | float | YES | Pre-tax for 10% (**negative**) |
| THANHTIEN_CHUA_VAT_KHAC | float | YES | Pre-tax for other rates (**negative**) |
| THANHTOAN_0 | float | YES | Total for 0% (**negative**) |
| THANHTOAN_5 | float | YES | Total for 5% (**negative**) |
| THANHTOAN_8 | float | YES | Total for 8% (**negative**) |
| THANHTOAN_10 | float | YES | Total for 10% (**negative**) |
| THANHTOAN_KHAC | float | YES | Total for other rates (**negative**) |
| THANHTOAN_KCT | float | YES | Tax-exempt total (**negative**) |
| PTTT | nvarchar(2) | YES | Refund method (should match original PTTT) |
| DGIAI | ntext | YES | Return reason/notes |
| CREATE_BY | int | YES | User ID who created |
| CREATE_DATE | datetime | YES | Creation timestamp |
| UPDATE_BY | int | YES | Last updated by |
| UPDATE_DATE | datetime | YES | Last update timestamp |
| DELETE_BY | int | YES | Deleted by (soft delete) |
| DELETE_DATE | datetime | YES | Deletion timestamp |
| DRC_STATUS | nvarchar(1) | YES | Sync status |
| PRINTED | bit | YES | Receipt printed? |
| ROWGUID | uniqueidentifier | NO | Row GUID |

**Business Rules:**
1. KHOA2 must reference existing BANLE.KHOA (the original sale)
2. All amount fields are **negative** (return is reverse of sale)
3. SOLUONG cannot exceed original BANLE_CT.SOLUONG for each item
4. SCT format: `{LOC}-{YYYYMMDD}-T{NNNN}` where 'T' indicates return
5. Return must be within allowed return period (check TRA_DATE in BANLE_CT)
6. PTTT (refund method) should generally match original payment method

**Indexes:**
```sql
CREATE INDEX IX_TRAHANG_KHOA2 ON TRAHANG(KHOA2);
CREATE INDEX IX_TRAHANG_NGAY ON TRAHANG(NGAY);
CREATE INDEX IX_TRAHANG_SCT ON TRAHANG(SCT);
```

---

### 1.4 TRAHANG_CT (Returns Detail)

**Purpose:** Stores line items for return transactions

**Primary Key:** KHOACT (uniqueidentifier)
**Foreign Key:** KHOA → TRAHANG(KHOA)

| Column | Type | Nullable | Description |
|--------|------|----------|-------------|
| KHOACT | uniqueidentifier | NO | Detail line ID, Primary Key |
| KHOA | uniqueidentifier | YES | FK to TRAHANG.KHOA |
| STT | int | YES | Line number |
| MAVT | nvarchar(15) | YES | Product code |
| TENVT | nvarchar(200) | YES | Product name (snapshot) |
| DVT | nvarchar(20) | YES | Unit |
| SOLUONG | float | YES | Quantity returned (**negative**) |
| DONGIA_CHUA_VAT | float | YES | Original unit price pre-tax |
| TL_CK | float | YES | Original discount % (preserved from sale) |
| LOAITHUE | nvarchar(15) | YES | Tax type |
| THUE_SUAT | float | YES | Tax rate % |
| THANHTIEN_CHUA_VAT | float | YES | Line pre-tax total (**negative**) |
| TIEN_THUE | float | YES | Tax amount (**negative**) |
| THANHTOAN | float | YES | Refund amount (**negative**) |
| GHICHU | ntext | YES | Line notes |
| CREATE_BY | int | YES | Created by |
| CREATE_DATE | datetime | YES | Creation timestamp |
| ROWGUID | uniqueidentifier | NO | Row GUID |

**Business Rules:**
1. All quantity and amount fields are **negative**
2. Pricing and discount match original sale (from BANLE_CT)
3. SOLUONG (absolute value) <= original BANLE_CT.SOLUONG

**Indexes:**
```sql
CREATE INDEX IX_TRAHANGCT_KHOA ON TRAHANG_CT(KHOA);
CREATE INDEX IX_TRAHANGCT_MAVT ON TRAHANG_CT(MAVT);
```

---

### 1.5 KHUYENMAI (Promotions Header)

**Purpose:** Stores promotion campaigns

**Primary Key:** KHOA (uniqueidentifier)

| Column | Type | Nullable | Description |
|--------|------|----------|-------------|
| KHOA | uniqueidentifier | NO | Promotion ID, Primary Key |
| LCT | nvarchar(5) | YES | Transaction type ('KM' for promotion) |
| NGAY | datetime | YES | Promotion created date |
| SCT | nvarchar(20) | YES | Promotion code (unique identifier) |
| LYDO | nvarchar(200) | YES | Promotion name/title |
| TUNGAY | datetime | YES | Start date |
| DENNGAY | datetime | YES | End date |
| DENNGAY2 | datetime | YES | Extended end date (if applicable) |
| TUGIO | datetime | YES | Start time of day (time part only) |
| DENGIO | datetime | YES | End time of day (time part only) |
| HTKM | nvarchar(1) | YES | Promotion type: '%'=Discount, 'G'=Fixed price, 'T'=Free |
| TL_CK | float | YES | Discount % (if HTKM = '%') |
| GIABAN | float | YES | Fixed sale price (if HTKM = 'G') |
| NOIDUNG | ntext | YES | Promotion description/details |
| DGIAI | ntext | YES | Additional notes |
| CHECKED | bit | YES | Approved? (0=Pending, 1=Approved) |
| CHECKED_BY | int | YES | User ID who approved |
| CHECKED_DATE | datetime | YES | Approval date |
| CHECKED_DESC | ntext | YES | Approval notes |
| CAPDO | int | YES | Priority level (higher = higher priority) |
| CREATE_BY | int | YES | Created by |
| CREATE_DATE | datetime | YES | Creation timestamp |
| UPDATE_BY | int | YES | Updated by |
| UPDATE_DATE | datetime | YES | Last update timestamp |
| DELETE_BY | int | YES | Deleted by |
| DELETE_DATE | datetime | YES | Deletion timestamp |
| DRC_STATUS | nvarchar(1) | YES | Sync status |
| LOC | nvarchar(2) | YES | Applicable location (NULL = all locations) |
| ROWGUID | uniqueidentifier | NO | Row GUID |

**Business Rules:**
1. TUNGAY <= DENNGAY
2. TUGIO <= DENGIO (time of day validation)
3. HTKM must be '%', 'G', or 'T'
4. If HTKM = '%', TL_CK must be provided
5. If HTKM = 'G', GIABAN must be provided
6. If CHECKED = 0, promotion is not active
7. Promotion applies if: CHECKED = 1 AND current_date BETWEEN TUNGAY AND DENNGAY AND current_time BETWEEN TUGIO AND DENGIO

**Indexes:**
```sql
CREATE INDEX IX_KHUYENMAI_TUNGAY_DENNGAY ON KHUYENMAI(TUNGAY, DENNGAY);
CREATE INDEX IX_KHUYENMAI_SCT ON KHUYENMAI(SCT);
CREATE INDEX IX_KHUYENMAI_CHECKED ON KHUYENMAI(CHECKED);
```

---

### 1.6 KHUYENMAI_CT (Promotions Detail)

**Purpose:** Specifies which products are included in each promotion

**Primary Key:** KHOACT (uniqueidentifier)
**Foreign Key:** KHOA → KHUYENMAI(KHOA)

| Column | Type | Nullable | Description |
|--------|------|----------|-------------|
| KHOACT | uniqueidentifier | NO | Detail line ID, Primary Key |
| KHOA | uniqueidentifier | YES | FK to KHUYENMAI.KHOA |
| STT | int | YES | Line number |
| MAVT | nvarchar(15) | YES | Product code, FK to DM_VT |
| TENVT | nvarchar(200) | YES | Product name (snapshot) |
| GIABAN | float | YES | Promotion price (if different from header) |
| GIABAN1 | float | YES | Alternate unit price |
| GIANHAPVAT | float | YES | Purchase price with VAT (reference) |
| TL_LAI | float | YES | Profit margin % at promotion price |
| TL_LAI1 | float | YES | Alternate unit profit margin |
| TL_CK | float | YES | Product-specific discount % (overrides header) |
| B1 | bit | YES | Flag for alternate unit |
| TINHTRANG | nvarchar(20) | YES | Status (ACTIVE, INACTIVE) |
| CREATE_BY | int | YES | Created by |
| CREATE_DATE | datetime | YES | Creation timestamp |
| UPDATE_BY | int | YES | Updated by |
| UPDATE_DATE | datetime | YES | Last update timestamp |
| ROWGUID | uniqueidentifier | NO | Row GUID |

**Business Rules:**
1. If KHUYENMAI_CT is empty for a KHUYENMAI, promotion applies to **all products**
2. If KHUYENMAI_CT has records, promotion applies only to listed MAVT
3. Product-specific TL_CK or GIABAN overrides header values
4. TINHTRANG = 'INACTIVE' means this product is excluded from promotion

**Indexes:**
```sql
CREATE INDEX IX_KHUYENMCT_KHOA ON KHUYENMAI_CT(KHOA);
CREATE INDEX IX_KHUYENMCT_MAVT ON KHUYENMAI_CT(MAVT);
```

---

## 2. MASTER DATA TABLES

### 2.1 DM_VT (Product Master)

**Purpose:** Master list of all products/items

**Primary Key:** MAVT (nvarchar(15))

| Column | Type | Nullable | Description |
|--------|------|----------|-------------|
| MAVT | nvarchar(15) | NO | Product code, Primary Key |
| TENVT | nvarchar(200) | YES | Product name |
| TENVT2 | nvarchar(200) | YES | Short name |
| TENTAT | nvarchar(50) | YES | Abbreviation |
| DVT | nvarchar(20) | YES | Primary unit of measure |
| DVT1 | nvarchar(20) | YES | Alternate unit |
| QD1 | int | YES | Conversion factor (DVT1 to DVT) |
| MANHOM | nvarchar(4) | YES | Product group code, FK to DM_NHOM |
| MANCC | nvarchar(15) | YES | Default supplier code, FK to DM_NCC |
| GIABAN | float | YES | Standard retail price (including tax) |
| GIABAN1 | float | YES | Alternate unit price |
| GIANHAP | float | YES | Latest purchase price (excluding VAT) |
| GIANHAPVAT | float | YES | Latest purchase price (including VAT) |
| GIAVON | float | YES | Average cost price |
| LOAITHUE | nvarchar(15) | YES | Default tax type, FK to SYS_LOAITHUE |
| THUE_SUAT | float | YES | Default tax rate % |
| TL_CK_MAX | float | YES | Maximum discount % allowed |
| TL_LAI | float | YES | Target profit margin % |
| BARCODE | nvarchar(30) | YES | Primary barcode (EAN-13, EAN-8, etc.) |
| PLU | nvarchar(20) | YES | PLU code (for fresh products with weight) |
| TINHTRANG | nvarchar(20) | YES | Status (ACTIVE, DISCONTINUED, etc.) |
| BO | bit | YES | Is bundle product? (1=Yes, 0=No) |
| GHICHU | ntext | YES | Notes/Description |
| XUATXU | nvarchar(100) | YES | Origin/Manufacturer |
| MALOAI | nvarchar(5) | YES | Category code |
| CREATE_BY | int | YES | Created by |
| CREATE_DATE | datetime | YES | Creation timestamp |
| UPDATE_BY | int | YES | Updated by |
| UPDATE_DATE | datetime | YES | Last update timestamp |
| ROWGUID | uniqueidentifier | NO | Row GUID |

**Business Rules:**
1. MAVT must be unique
2. GIABAN = GIAVON * (1 + TL_LAI/100) / (1 - TL_CK_MAX/100) (approximate)
3. If BO = 1, product is a bundle; components defined in DM_VT_BO
4. BARCODE can be non-unique (multiple products can share barcode in edge cases)
5. QD1 defines conversion: 1 DVT1 = QD1 * DVT
6. TINHTRANG <> 'ACTIVE' means product should not appear in POS by default

**Indexes:**
```sql
CREATE INDEX IX_DMVT_BARCODE ON DM_VT(BARCODE);
CREATE INDEX IX_DMVT_MANHOM ON DM_VT(MANHOM);
CREATE INDEX IX_DMVT_MANCC ON DM_VT(MANCC);
CREATE INDEX IX_DMVT_TINHTRANG ON DM_VT(TINHTRANG);
CREATE INDEX IX_DMVT_TENVT ON DM_VT(TENVT);
```

---

### 2.2 DM_VT_BO (Product Bundle Components)

**Purpose:** Defines components of bundle products

**Primary Key:** KHOACT (uniqueidentifier)
**Foreign Keys:** 
- MABO → DM_VT(MAVT)
- MAVT → DM_VT(MAVT)

| Column | Type | Nullable | Description |
|--------|------|----------|-------------|
| KHOACT | uniqueidentifier | NO | Component ID, Primary Key |
| MABO | nvarchar(15) | YES | Bundle product code (DM_VT.MAVT where BO=1) |
| MAVT | nvarchar(15) | YES | Component product code |
| SOLUONG | float | YES | Component quantity per bundle |
| STT | int | YES | Display order |
| CREATE_BY | int | YES | Created by |
| CREATE_DATE | datetime | YES | Creation timestamp |
| ROWGUID | uniqueidentifier | NO | Row GUID |

**Business Rules:**
1. MABO must reference DM_VT where BO = 1
2. MAVT must reference DM_VT where BO = 0 (components cannot be bundles themselves)
3. When MABO is scanned at POS, auto-expand to add all MAVT components with their SOLUONG

**Indexes:**
```sql
CREATE INDEX IX_DMVTBO_MABO ON DM_VT_BO(MABO);
CREATE INDEX IX_DMVTBO_MAVT ON DM_VT_BO(MAVT);
```

---

### 2.3 DM_KH_NCC (Customer/Supplier Master)

**Purpose:** Master list of customers and suppliers

**Primary Key:** KHOA (nvarchar(100))
**Unique Key:** MADT

| Column | Type | Nullable | Description |
|--------|------|----------|-------------|
| KHOA | nvarchar(100) | NO | Unique ID (GUID string), Primary Key |
| LOAI | int | YES | Type: 0=Customer, 1=Supplier, 3=Both |
| MADT | nvarchar(15) | NO | Code (unique) |
| TENDT | nvarchar(200) | YES | Name |
| TENTAT | nvarchar(50) | YES | Short name |
| PLOAI | nvarchar(2) | YES | Sub-type code |
| NGUNG_LIENLAC | bit | YES | Inactive/Do not contact (1=Yes) |
| DCHI | nvarchar(200) | YES | Address |
| MAQG | nvarchar(5) | YES | Country code |
| MATINH | nvarchar(11) | YES | Province code |
| MAHUYEN | nvarchar(17) | YES | District code |
| DTHOAI | nvarchar(100) | YES | Phone |
| FAX | nvarchar(100) | YES | Fax |
| EMAIL | nvarchar(100) | YES | Email |
| MST | nvarchar(30) | YES | Tax ID |
| TENTK | nvarchar(200) | YES | Bank account holder name |
| MATK | nvarchar(50) | YES | Bank account number |
| TEN_NGANHANG | nvarchar(200) | YES | Bank name |
| DC_NGANHANG | nvarchar(200) | YES | Bank address |
| HAN_TTOAN | int | YES | Payment term (days) |
| SODU | float | YES | Account balance (negative = customer owes) |
| SODU_NGAY | datetime | YES | Balance as of date |
| SODU_LIMIT | float | YES | Credit limit (negative value) |
| MADT2 | nvarchar(15) | YES | Secondary code |
| MANV | nvarchar(15) | YES | Assigned salesperson |
| LOAI_VIP | nvarchar(15) | YES | VIP level code (reference only) |
| MAKV | nvarchar(5) | YES | Region code |
| GHICHU | ntext | YES | Notes |
| CREATE_BY | int | YES | Created by |
| CREATE_DATE | datetime | YES | Creation timestamp |
| UPDATE_BY | int | YES | Updated by |
| UPDATE_DATE | datetime | YES | Last update timestamp |
| ROWGUID | uniqueidentifier | NO | Row GUID |

**Business Rules:**
1. MADT must be unique
2. LOAI determines if record appears in DM_KH view (0, 3) or DM_NCC view (1, 3)
3. NGUNG_LIENLAC = 1 means record is inactive (excluded from views)
4. SODU < 0 means customer owes money; SODU > 0 means customer has credit
5. SODU_LIMIT is negative (e.g., -10000000 means customer can owe up to 10M)
6. For POS, primarily use customer records (LOAI IN (0, 3))

**Views:**
```sql
CREATE VIEW DM_KH AS
SELECT * FROM DM_KH_NCC 
WHERE LOAI IN (0, 3) AND ISNULL(NGUNG_LIENLAC, 0) = 0;

CREATE VIEW DM_NCC AS
SELECT * FROM DM_KH_NCC 
WHERE LOAI IN (1, 3) AND ISNULL(NGUNG_LIENLAC, 0) = 0;
```

**Indexes:**
```sql
CREATE INDEX IX_DMKH_MADT ON DM_KH_NCC(MADT);
CREATE INDEX IX_DMKH_LOAI ON DM_KH_NCC(LOAI);
CREATE INDEX IX_DMKH_NGUNG_LIENLAC ON DM_KH_NCC(NGUNG_LIENLAC);
CREATE INDEX IX_DMKH_TENDT ON DM_KH_NCC(TENDT);
```

---

### 2.4 DM_LOCATION (Store/Location Master)

**Purpose:** Defines stores/locations

**Primary Key:** LOC (nvarchar(2))

| Column | Type | Nullable | Description |
|--------|------|----------|-------------|
| LOC | nvarchar(2) | NO | Location code, Primary Key |
| TENKHO | nvarchar(100) | YES | Store/Location name |
| MAKHO | nvarchar(2) | YES | Default warehouse code |
| DIACHI | nvarchar(200) | YES | Address |
| DTHOAI | nvarchar(100) | YES | Phone |
| HoaDon_Mau | nvarchar(50) | YES | eInvoice template code (for BKAV integration) |
| HoaDon_KyHieu | nvarchar(50) | YES | eInvoice serial (for BKAV integration) |
| GHICHU | ntext | YES | Notes |
| CREATE_BY | int | YES | Created by |
| CREATE_DATE | datetime | YES | Creation timestamp |
| UPDATE_BY | int | YES | Updated by |
| UPDATE_DATE | datetime | YES | Last update timestamp |
| ROWGUID | uniqueidentifier | NO | Row GUID |

**Business Rules:**
1. LOC is typically 2-digit code (01, 02, 03...)
2. Used in BANLE.LOC to identify which store made the sale
3. HoaDon_Mau and HoaDon_KyHieu are for eInvoice module (out of scope for core POS, but preserved)

**Indexes:**
```sql
CREATE INDEX IX_LOCATION_TENKHO ON DM_LOCATION(TENKHO);
```

---

### 2.5 DM_QUAYTN (POS Counter/Terminal Master)

**Purpose:** Defines POS counters/terminals

**Primary Key:** MAQUAY (nvarchar(15))
**Foreign Key:** LOC → DM_LOCATION(LOC)

| Column | Type | Nullable | Description |
|--------|------|----------|-------------|
| MAQUAY | nvarchar(15) | NO | Counter code, Primary Key |
| TENQUAY | nvarchar(100) | YES | Counter name |
| QUAYKE | nvarchar(100) | YES | Shelf/Display location description |
| LOC | nvarchar(2) | YES | Location code, FK to DM_LOCATION |
| PRINTER_NAME | nvarchar(100) | YES | Receipt printer name (Windows printer name) |
| PRINTER_PORT | nvarchar(20) | YES | Printer COM port or network address |
| ACTIVE | bit | YES | Is counter active? (1=Yes, 0=No) |
| CREATE_BY | int | YES | Created by |
| CREATE_DATE | datetime | YES | Creation timestamp |
| UPDATE_BY | int | YES | Updated by |
| UPDATE_DATE | datetime | YES | Last update timestamp |
| ROWGUID | uniqueidentifier | NO | Row GUID |

**Business Rules:**
1. MAQUAY must be unique
2. Each counter belongs to one location (LOC)
3. ACTIVE = 0 means counter is disabled (cannot be selected in POS)

**Indexes:**
```sql
CREATE INDEX IX_QUAYTN_LOC ON DM_QUAYTN(LOC);
CREATE INDEX IX_QUAYTN_ACTIVE ON DM_QUAYTN(ACTIVE);
```

---

## 3. SYSTEM TABLES

### 3.1 SYS_USER (User Master)

**Purpose:** System users (cashiers, supervisors, managers)

**Primary Key:** UID (int IDENTITY)

| Column | Type | Nullable | Description |
|--------|------|----------|-------------|
| UID | int IDENTITY(1,1) | NO | User ID, Primary Key |
| USERNAME | nvarchar(50) | NO | Login username (unique) |
| PASSWORD | nvarchar(100) | YES | Hashed password |
| FULLNAME | nvarchar(100) | YES | Full name |
| ROLE | nvarchar(20) | YES | Role: CASHIER, SUPERVISOR, MANAGER, ADMIN |
| CAN_DISCOUNT_MAX | float | YES | Maximum discount % this user can authorize |
| ACTIVE | bit | YES | Is active? (1=Yes, 0=No) |
| LOCKED | bit | YES | Account locked? (1=Locked, 0=Active) |
| CREATE_DATE | datetime | YES | Account creation date |
| UPDATE_DATE | datetime | YES | Last update timestamp |
| LAST_LOGIN | datetime | YES | Last login timestamp |
| ROWGUID | uniqueidentifier | NO | Row GUID |

**Business Rules:**
1. USERNAME must be unique
2. PASSWORD should be hashed (bcrypt, argon2, etc.)
3. ACTIVE = 0 or LOCKED = 1 means user cannot log in
4. ROLE determines access permissions:
   - CASHIER: Basic POS operations, discount up to configured limit
   - SUPERVISOR: Extended POS ops, higher discount authorization
   - MANAGER: Full POS + reporting, any discount
   - ADMIN: System configuration

**Indexes:**
```sql
CREATE UNIQUE INDEX IX_SYSUSER_USERNAME ON SYS_USER(USERNAME);
CREATE INDEX IX_SYSUSER_ACTIVE ON SYS_USER(ACTIVE);
```

---

### 3.2 SYS_CONFIG (System Configuration)

**Purpose:** System-wide configuration settings

**Primary Key:** CONFIG_KEY (nvarchar(50))

| Column | Type | Nullable | Description |
|--------|------|----------|-------------|
| CONFIG_KEY | nvarchar(50) | NO | Configuration key, Primary Key |
| CONFIG_VALUE | ntext | YES | Configuration value (can be JSON) |
| DESCRIPTION | nvarchar(200) | YES | Description |
| UPDATE_BY | int | YES | Last updated by user ID |
| UPDATE_DATE | datetime | YES | Last update timestamp |

**Key Configuration Keys:**
- `RECEIPT_FORMAT`: Receipt template (HTML/text format)
- `STAMP_PREFIX`: Barcode prefix for weight-embedded barcodes (e.g., '22,23,24,25')
- `COMPANY_INFO`: Business info for receipt header (JSON: name, address, tax ID, etc.)
- `TAX_CONFIG`: Tax calculation rules (JSON)
- `RETURN_POLICY`: Return policy settings (JSON: max days, restrictions)
- `DISCOUNT_APPROVAL_LEVELS`: Discount authorization matrix (JSON)
- `FLEXIBLE`: Legacy flexible config (INI-style format) - see fnGetFlexible function

**Indexes:**
```sql
-- Primary key index on CONFIG_KEY is sufficient
```

---

### 3.3 SYS_LOAITHUE (Tax Types)

**Purpose:** Defines tax types and rates

**Primary Key:** MALT (nvarchar(15))

| Column | Type | Nullable | Description |
|--------|------|----------|-------------|
| MALT | nvarchar(15) | NO | Tax code (T00, T05, T08, T10, KCT, TTT), Primary Key |
| TENLT | nvarchar(200) | YES | Tax name/description |
| VAT_VAO | float | YES | Input VAT % (for purchases) |
| VAT_RA | float | YES | Output VAT % (for sales) |
| CREATE_BY | int | YES | Created by |
| CREATE_DATE | datetime | YES | Creation timestamp |
| UPDATE_BY | int | YES | Updated by |
| UPDATE_DATE | datetime | YES | Last update timestamp |
| ROWGUID | uniqueidentifier | NO | Row GUID |

**Standard Tax Codes:**
- `T00`: 0% VAT
- `T05`: 5% VAT
- `T08`: 8% VAT
- `T10`: 10% VAT
- `KCT`: Không chịu thuế (Tax-exempt)
- `TTT`: Thuế thu trước (Pre-paid tax)

**Business Rules:**
1. MALT must match values in BANLE_CT.LOAITHUE, DM_VT.LOAITHUE
2. VAT_RA is used for retail sales calculations
3. KCT means tax-exempt (VAT_RA = 0)

---

## 4. VIEWS

### 4.1 DM_VT_FULL (Complete Product View)
```sql
CREATE VIEW DM_VT_FULL AS
SELECT 
    v.*,
    n.TENNHOM,
    ng.TENNGANH,
    ncc.TENDT AS TEN_NCC
FROM DM_VT v
LEFT JOIN DM_NHOM n ON v.MANHOM = n.MANHOM
LEFT JOIN DM_NGANH ng ON n.MANGANH = ng.MANGANH
LEFT JOIN DM_KH_NCC ncc ON v.MANCC = ncc.MADT
WHERE ISNULL(v.UPDATE_BY, 0) <> -1;  -- Exclude system-deleted records
```

### 4.2 T_BANLE (Active Sales Transactions)
```sql
CREATE VIEW T_BANLE AS
SELECT * FROM BANLE
WHERE ISNULL(DELETE_BY, 0) = 0;
```

### 4.3 V_BANLE_SUMMARY (Sales Summary View)
```sql
CREATE VIEW V_BANLE_SUMMARY AS
SELECT
    b.KHOA,
    b.NGAY,
    b.SCT,
    b.LOC,
    b.QUAY,
    b.CA,
    b.THUNGAN,
    b.SOLUONG,
    b.THANHTOAN,
    b.PTTT,
    u.FULLNAME AS THUNGAN_NAME,
    l.TENKHO AS LOC_NAME,
    q.TENQUAY AS QUAY_NAME
FROM BANLE b
LEFT JOIN SYS_USER u ON b.THUNGAN = u.USERNAME
LEFT JOIN DM_LOCATION l ON b.LOC = l.LOC
LEFT JOIN DM_QUAYTN q ON b.QUAY = q.MAQUAY
WHERE ISNULL(b.DELETE_BY, 0) = 0;
```

---

## 5. STORED PROCEDURES

### 5.1 ALLOC_SCT (Allocate General Transaction Number)
Cấp số chứng từ hệ thống (NMUA, XBAN, THU, CHI, NTBL, CKHO, GTGT...) dựa trên bảng `SYS_SEQUENCE`:
```sql
CREATE PROCEDURE [dbo].[ALLOC_SCT]
	@SCT	nvarchar(50) out,
	@LOAI	nvarchar(50),	-- Loai chung tu ('NMUA','XBAN','THU','CHI','NTBL','CKHO'...)
	@NGAY	datetime,	
	@MAKHO	nvarchar(5) = ''
as
declare	@S	nvarchar(20),
	@THANG	nvarchar(2),
	@NAM	nvarchar(4),
	@STT	int

if @MAKHO <> ''
	set @MAKHO = @MAKHO + '/'
set @NAM   = right('00' + cast(year (@NGAY) as nvarchar), 2)
set @THANG = right('00'   + cast(month(@NGAY) as nvarchar), 2)

-- Mapping tiền tố theo loại chứng từ
if @LOAI = 'NMUA'	set @SCT = 'PN/'	else
if @LOAI = 'NTRA'	set @SCT = 'NT/'	else
if @LOAI = 'NKHAC'	set @SCT = 'NK/'	else
if @LOAI = 'NDC'	set @SCT = 'NDC/'	else
if @LOAI = 'XBAN'	set @SCT = 'PX/'	else
if @LOAI = 'XTRA'	set @SCT = 'XT/'	else
if @LOAI = 'XKHAC'	set @SCT = 'XK/'	else
if @LOAI = 'XDC'	set @SCT = 'XDC/'	else
if @LOAI = 'CKHO'	set @SCT = 'DKH/'	else
if @LOAI = 'KKE'	set @SCT = 'KK/'	else
if @LOAI = 'NTBL'	set @SCT = 'NTL/'	else
if @LOAI = 'NTBL2'	set @SCT = 'NTL/'	else
if @LOAI = 'CHI'	set @SCT = 'PC/'	else
if @LOAI = 'THU'	set @SCT = 'PT/'	else
if @LOAI = 'GTGT'	set @SCT = 'HD/'	else
			set @SCT = '??/'

-- Sequence alloc atomic counter
set @S = @MAKHO + @LOAI + '/' + @NAM + @THANG
exec SYS_SEQ_NEXTVAL @S, @STT out	-- e.g. MAKHO/LOAI/YYMM

-- Build SCT format: [TIENTO]/[MAKHO/][YY][MM]/[000000]
set @SCT = @SCT + @MAKHO + @NAM + @THANG + '/' +
	right ('000000'	+ cast(@STT as nvarchar), 6)

return(0)
```

### 5.2 POS_ALLOC_BILL_NO (Allocate POS Retail Bill Number)
Cấp số hóa đơn bán lẻ POS tại quầy, lấy số lớn nhất từ cả bảng `BANLE` và `BANLE_TAM`:
```sql
CREATE PROCEDURE [dbo].[POS_ALLOC_BILL_NO]
	@MAKHO	nvarchar(2),
	@QUAY	nvarchar(2),
	@NGAY	datetime,
	@SCT	nvarchar(20) out
AS
declare	@THANG	int,
	@NAM	int,
	@N	int,
	@S	int

set @THANG = month(@NGAY)
set @NAM = year(@NGAY)

select	@S = isnull(max(SCT), '')
  from	(
	select	SCT
	  from	BANLE
	 where	LCT = 'BLE'
	   and	MAKHO = @MAKHO
	   and	QUAY = @QUAY
	   and	month(NGAY) = @THANG
	   and	year(NGAY) = @NAM
  	union 
	select	SCT
	  from	BANLE_TAM
	 where	LCT = 'BLE'
	   and	MAKHO = @MAKHO
	   and	QUAY = @QUAY
	   and	month(NGAY) = @THANG
	   and	year(NGAY) = @NAM
	) a

if isnull(@S, '')=''
	set @S = '0'

set @N = cast(@S as int) + 1
set @SCT = right('00000' + cast(@N as nvarchar), 5)

return (0)
```

### 5.3 POS_PAID (Move Paid Receipt from BANLE_TAM to BANLE)
Chuyển hóa đơn đã hoàn tất thanh toán từ bảng tạm sang bảng chính thức:
```sql
CREATE PROCEDURE [dbo].[POS_PAID]
	@KHOA	uniqueidentifier
AS
begin tran

-- 1. Chuyển Header sang BANLE
insert into BANLE
select	*
  from	BANLE_TAM
 where	KHOA = @KHOA
if @@error <> 0 goto FAILURE

-- 2. Chuyển Details sang BANLE_CT
insert into BANLE_CT (KHOA, STT, MAVT, TENVT, DVT, MABO, SOLUONG, SOLUONG2,
    DONGIA_REF, DONGIA, DONGIA1, SOTIEN, SOTIEN1, SOTIEN_DC, TL_CK, TL_CK2..7,
    CHIETKHAU, LOAITHUE, THUE_SUAT, TIEN_THUE, THANHTOAN...)
select	KHOA, STT, MAVT, TENVT, DVT, MABO, SOLUONG, SOLUONG2,
    DONGIA_REF, DONGIA, DONGIA1, SOTIEN, SOTIEN1, SOTIEN_DC, TL_CK, TL_CK2..7,
    CHIETKHAU, LOAITHUE, THUE_SUAT, TIEN_THUE, THANHTOAN...
  from	BANLE_TAM_CT
 where	KHOA = @KHOA
order by STT
if @@error <> 0 goto FAILURE

-- 3. Xóa dữ liệu tạm
delete BANLE_TAM_CT where KHOA = @KHOA
delete BANLE_TAM    where KHOA = @KHOA

-- 4. Đồng bộ vào bảng hóa đơn điện tử ehoadon
merge into ehoadon as hd
using (
	select KHOA, SCT, LCT, MAKHO, LOC, CREATE_BY
	  from BANLE where KHOA = @KHOA
) as a ON (hd.khoaphieugoc = a.khoa)
when matched then
	update set KhoaPhieuGoc = a.khoa, SctPhieuGoc = a.SCT, LctPhieuGoc = a.LCT, MaKhoPhieuGoc = a.MAKHO, LOC = a.LOC
when not matched then
	insert (KhoaPhieuGoc, SctPhieuGoc, LctPhieuGoc, MaKhoPhieuGoc, LOC, CREATE_BY, UPDATE_BY, CREATE_DATE, UPDATE_DATE)
	values (a.khoa, a.SCT, a.LCT, a.MAKHO, a.LOC, a.CREATE_BY, a.CREATE_BY, getdate(), getdate());

if @@trancount > 0 commit tran
return (0)

FAILURE:
	if @@trancount > 0 rollback tran
	return 1
```

### 5.4 SYS_SEQ_NEXTVAL (Atomic Sequence Generator)
Cơ chế sinh sequence tuần tự an toàn tuyệt đối chống race-condition:
```sql
CREATE PROCEDURE dbo.SYS_SEQ_NEXTVAL 
	@SEQ_NAME	nvarchar(50),
	@SEQ_VAL	numeric out,
	@UPDATE		bit = 1
AS
declare	@SEQ_CUR_VAL numeric
	
select	@SEQ_CUR_VAL = SEQ_CURVAL
  from	SYS_SEQUENCE
 where	SEQ_NAME = @SEQ_NAME

if @@rowcount = 0
begin
	set @SEQ_VAL = 1
	if @UPDATE = 1
	begin
		begin transaction	
		insert into SYS_SEQUENCE values (@SEQ_NAME, 1)
		if @@error <> 0
			exec SYS_SEQ_NEXTVAL @SEQ_NAME, @SEQ_VAL out, @UPDATE
		commit	
	end
end
else
begin
	set @SEQ_VAL = @SEQ_CUR_VAL + 1
	if @UPDATE = 1
	begin
		update	SYS_SEQUENCE
		   set	SEQ_CURVAL = @SEQ_VAL
		 where	SEQ_NAME = @SEQ_NAME and SEQ_CURVAL = @SEQ_CUR_VAL
		if @@rowcount = 0
			exec SYS_SEQ_NEXTVAL @SEQ_NAME, @SEQ_VAL out, @UPDATE
	end
end
return (0)
```

### 5.5 POS_GIABAN (Get Sale Price)
```sql
CREATE PROCEDURE POS_GIABAN
    @MAVT nvarchar(15),
    @MAVIP nvarchar(20) = NULL,
    @NGAY datetime = NULL,
    @GIABAN float OUTPUT,
    @TL_CK float OUTPUT
AS
BEGIN
    IF @NGAY IS NULL SET @NGAY = GETDATE();
    
    -- Default price from product master
    SELECT @GIABAN = GIABAN, @TL_CK = 0
    FROM DM_VT
    WHERE MAVT = @MAVT;
    
    -- Check for active promotion
    DECLARE @PROMO_PRICE float;
    DECLARE @PROMO_DISCOUNT float;
    
    SELECT TOP 1
        @PROMO_PRICE = CASE WHEN km.HTKM = 'G' THEN km.GIABAN ELSE NULL END,
        @PROMO_DISCOUNT = CASE WHEN km.HTKM = '%' THEN km.TL_CK ELSE NULL END
    FROM KHUYENMAI km
    LEFT JOIN KHUYENMAI_CT kmc ON km.KHOA = kmc.KHOA AND kmc.MAVT = @MAVT
    WHERE km.CHECKED = 1
      AND @NGAY BETWEEN km.TUNGAY AND km.DENNGAY
      AND CAST(@NGAY AS time) BETWEEN ISNULL(km.TUGIO, '00:00') AND ISNULL(km.DENGIO, '23:59')
      AND (kmc.MAVT IS NOT NULL OR NOT EXISTS (SELECT 1 FROM KHUYENMAI_CT WHERE KHOA = km.KHOA))
      AND ISNULL(km.DELETE_BY, 0) = 0
    ORDER BY km.CAPDO DESC, km.CREATE_DATE DESC;
    
    IF @PROMO_PRICE IS NOT NULL SET @GIABAN = @PROMO_PRICE;
    IF @PROMO_DISCOUNT IS NOT NULL SET @TL_CK = @PROMO_DISCOUNT;
    
    -- VIP discount logic would go here (out of scope)
END;
```

---

## 6. FUNCTIONS

### 6.1 fnGetFlexible (Get Config Value)
```sql
CREATE FUNCTION fnGetFlexible (
    @SECTION nvarchar(50),
    @ITEM nvarchar(50)
)  
RETURNS nvarchar(500) AS  
BEGIN 
    DECLARE @STR nvarchar(3000), @N int;
    
    SELECT @STR = CONFIG_VALUE 
    FROM SYS_CONFIG 
    WHERE CONFIG_KEY = 'FLEXIBLE';
    
    SET @SECTION = N'[' + @SECTION + N']';
    SET @ITEM = @ITEM + N'=';
    
    SET @N = CHARINDEX(@SECTION, @STR);
    IF @N > 0
    BEGIN
        SET @STR = SUBSTRING(@STR, @N + LEN(@SECTION) + 2, 3000);
        SET @N = CHARINDEX('[', @STR);
        IF @N = 0 SET @N = LEN(@STR);
        ELSE SET @N = @N - 5;
        SET @STR = SUBSTRING(@STR, 1, @N);
        
        SET @N = CHARINDEX(@ITEM, @STR);
        IF @N > 0
        BEGIN
            SET @STR = SUBSTRING(@STR, @N + LEN(@ITEM), 3000);
            SET @N = CHARINDEX(CHAR(13), @STR);
            IF @N = 0 SET @N = LEN(@STR);
            SET @STR = SUBSTRING(@STR, 1, @N);
        END
        ELSE SET @STR = '';
    END
    ELSE SET @STR = '';
    
    SET @STR = REPLACE(@STR, CHAR(13), '');
    
    RETURN RTRIM(LTRIM(@STR));
END;
```

### 6.2 fnStripToneMark (Remove Vietnamese Tone Marks)
Used for search/filtering - converts Vietnamese text to non-accented for search purposes.

### 6.3 fnNum_Text (Convert Number to Vietnamese Words)
Used for receipt printing - converts monetary amount to Vietnamese words.

---

## 7. TRIGGERS

**Note:** Current schema does not rely heavily on triggers. Consider adding these for data integrity:

### 7.1 TR_BANLE_AfterInsert (Audit Trail)
```sql
CREATE TRIGGER TR_BANLE_AfterInsert
ON BANLE
AFTER INSERT
AS
BEGIN
    -- Log to audit table if needed
    -- Update related tables (inventory, customer balance, etc.)
END;
```

### 7.2 TR_BANLE_PreventDeleteIfPrinted
```sql
CREATE TRIGGER TR_BANLE_PreventDeleteIfPrinted
ON BANLE
INSTEAD OF DELETE
AS
BEGIN
    IF EXISTS (SELECT 1 FROM deleted WHERE PRINTED = 1)
    BEGIN
        RAISERROR('Cannot delete printed transactions', 16, 1);
        ROLLBACK TRANSACTION;
        RETURN;
    END
    
    -- Perform soft delete
    UPDATE BANLE
    SET DELETE_BY = USER_ID(),  -- Replace with actual user context
        DELETE_DATE = GETDATE()
    WHERE KHOA IN (SELECT KHOA FROM deleted);
END;
```

---

## 8. FIELD NAMING CONVENTIONS

### 8.1 Common Prefixes
- `MA*`: Code fields (MAVT, MADT, MAKHO, etc.)
- `TEN*`: Name fields (TENVT, TENDT, TENKHO, etc.)
- `TL_*`: Rate/Percentage fields (TL_CK, TL_LAI, etc.)
- `LK_*`: Lookup/Linked fields (calculated, not in DB)
- `CALC_*`: Calculated fields (runtime computation)

### 8.2 Common Suffixes
- `*_BY`: User ID who performed action
- `*_DATE`: Datetime field
- `*_VAT`: Including VAT amount
- `*_CHUA_VAT`: Excluding VAT amount (pre-tax)
- `*_SAUCK`: After discount (sau chiết khấu)

### 8.3 Common Fields
- `KHOA`: Primary key (GUID) for transaction tables
- `KHOACT`: Primary key (GUID) for detail/child tables
- `KHOA2`: Foreign key reference to parent transaction
- `LCT`: Transaction type code (loại chứng từ)
- `SCT`: Receipt/Document number (số chứng từ)
- `NGAY`: Transaction date/time
- `LOC`: Location code
- `STT`: Sequence/Line number
- `CREATE_BY, CREATE_DATE`: Audit trail - creation
- `UPDATE_BY, UPDATE_DATE`: Audit trail - modification
- `DELETE_BY, DELETE_DATE`: Audit trail - soft deletion
- `DRC_STATUS`: Data replication/sync status
- `ROWGUID`: Row GUID (for SQL Server replication)

---

**Document Version:** 1.0  
**Last Updated:** 2026-10-14  
**Author:** Database Architect (AI-assisted)  
**Status:** Reference Document
