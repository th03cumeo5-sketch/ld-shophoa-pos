# BUSINESS REQUIREMENTS DOCUMENT (BRD)
## POSz - Point of Sale System Migration to Rust

---

## 1. EXECUTIVE SUMMARY

### 1.1 Project Overview
Migration of existing POSz (Point of Sale) system from Delphi/Pascal to Rust, maintaining all core business logic and functionality for retail operations.

### 1.2 Scope
**IN SCOPE:**
- Point of Sale (POS) operations
- Retail sales (Bán lẻ - BANLE)
- Return/Refund processing (Trả hàng - TRAHANG)
- Promotion management (Khuyến mãi)
- Multi-discount logic
- Payment processing (multiple payment methods)
- Receipt printing
- Customer information (basic)
- Inventory query (read-only from other modules)

**OUT OF SCOPE:**
- VIP/Loyalty card management (handled by NXTz module)
- Inventory management/Stock taking (handled by NXTz module)
- HR/Payroll/Fingerprint attendance (handled by NXTz module)
- Purchase orders (handled by NXTz module)
- eInvoice integration with BKAV (separate workflow)

---

## 2. BUSINESS PROCESSES

### 2.1 Core Business Flow: RETAIL SALES (Bán lẻ)

#### 2.1.1 Process Overview
```
[Start POS Session] → [Scan/Input Items] → [Apply Discounts] → [Apply Promotions] 
→ [Calculate Total] → [Payment] → [Print Receipt] → [Close Transaction]
```

#### 2.1.2 Detailed Steps

**Step 1: Start POS Session**
- Cashier logs in with User ID
- Select Counter/Quay (QUAY)
- Select Shift/Ca làm việc (CA)
- Select Location/Store (LOC)
- Default warehouse (MAKHO) is auto-selected based on location

**Step 2: Scan/Input Items**
- Scan barcode OR input product code (MAVT)
- Support multiple barcode formats:
  - Standard barcode (from DM_VT)
  - EAN-13, EAN-8
  - Stamp prefix barcode (weight-embedded for fresh products)
- System displays:
  - Product name (TENVT)
  - Unit price (DONGIA)
  - Available stock (for reference only)
  - Current quantity in cart
- Quantity can be:
  - Scanned from barcode (weight products)
  - Manually entered
  - Default = 1

**Step 3: Apply Discounts (Multi-level)**
The system supports **7 levels of discount** tracked in fields TL_CK, TL_CK2...TL_CK7:

- **TL_CK**: Final applied discount (max of all discount sources)
- **TL_CK2**: VIP customer discount (from DM_VIP logic)
- **TL_CK3**: Product group discount
- **TL_CK4**: Promotion discount (from KHUYENMAI)
- **TL_CK5**: Manual discount (by cashier, requires permission)
- **TL_CK6**: Additional discount
- **TL_CK7**: Special discount

**Discount Logic:**
```
TL_CK = MAX(TL_CK2, TL_CK3, TL_CK4, TL_CK5, TL_CK6, TL_CK7)
TL_CK_MAX = Maximum allowed discount for the product
TL_CK_THEM = Additional discount beyond max (requires approval)
```

**Discount Calculation:**
```
DONGIA_SAUCK = DONGIA * (1 - TL_CK/100)
THANHTIEN_CHUA_VAT = SOLUONG * DONGIA_SAUCK
TIEN_THUE = THANHTIEN_CHUA_VAT * THUE_SUAT/100
THANHTOAN = THANHTIEN_CHUA_VAT + TIEN_THUE
```

**Step 4: Apply Promotions (Khuyến mãi)**
Types of promotions:
1. **Discount by percentage** (Chiết khấu %)
2. **Fixed price** (Giá bán cố định)
3. **Free items** (Tặng hàng - represented as 100% discount)

Promotion conditions:
- Time-based (TUNGAY, DENNGAY, TUGIO, DENGIO)
- Product-specific (MAVT)
- Automatic application (HTKM field)

**Free items logic:**
- Free items are added with THANHTOAN = 0
- When generating eInvoice, system checks: `if THANHTOAN = 0 then ItemTypeID = 15`
- This marks the item as promotional/free in invoice

**Step 5: Calculate Total**
```
Header totals (BANLE table):
- SOTIEN = Sum of THANHTIEN_CHUA_VAT from all items
- CHIETKHAU_MH = Sum of discount amounts
- THUE = Sum of TIEN_THUE from all items
- THANHTOAN = SOTIEN + THUE - CHIETKHAU_MH (after header-level discount)
- SOLUONG = Total item count

Tax breakdown by rate:
- TIEN_THUE_0, TIEN_THUE_5, TIEN_THUE_8, TIEN_THUE_10, TIEN_THUE_KHAC
- THANHTIEN_CHUA_VAT_0, _5, _8, _10, _KHAC
- THANHTOAN_0, _5, _8, _10, _KHAC, _KCT (không chịu thuế)
```

Tax types (LOAITHUE):
- `T00`: 0% VAT
- `T05`: 5% VAT
- `T08`: 8% VAT
- `T10`: 10% VAT
- `KCT`: Không chịu thuế (tax-exempt)
- `TTT`: Thuế thu trước (pre-paid tax)

**Step 6: Payment Processing**
Multiple payment methods supported (PTTT, PTTT2):

| Code | Description | Notes |
|------|-------------|-------|
| TM | Tiền mặt (Cash) | Default |
| CK | Chuyển khoản (Bank transfer) | Requires account info |
| CN | Công nợ (Credit) | Requires customer info |
| TANG | Tặng (Gift) | No payment |
| BIEU | Biếu (Complimentary) | No payment |

**Split payment:**
- TTOAN1: Primary payment amount (PTTT)
- TTOAN2: Secondary payment amount (PTTT2)
- TTOAN1 + TTOAN2 = THANHTOAN
- CHUATHOI: Change to return = (TTOAN1 + TTOAN2) - THANHTOAN

**Step 7: Temporary Bill & Final Payment (BANLE_TAM → BANLE)**
- **In-Progress & Temporary Bill (Lưu tạm):**
  - Mọi thao tác quét hàng, giỏ hàng đang soạn, và hóa đơn tạm (khi bấm F10/Save2) được lưu vào cụm bảng tạm: **`BANLE_TAM`** (header) và **`BANLE_TAM_CT`** (details).
  - Thu ngân có thể treo đơn (`PRINTED = 0`, ghi chú `DGIAI`), mở đơn mới, và phục hồi lại đơn tạm bất cứ lúc nào.
- **Receipt Number Allocation (Cấp số chứng từ):**
  - **POS Bán lẻ:** Gọi Stored Procedure **`POS_ALLOC_BILL_NO`** (`@MAKHO, @QUAY, @NGAY, @SCT OUT`). Procedure tự động lấy `MAX(SCT)` từ cả `BANLE` và `BANLE_TAM` của quầy và tháng/năm hiện tại rồi sinh số tuần tự chuẩn 5 chữ số (ví dụ: `00001`, `00002`...).
  - **Chứng từ hệ thống khác (Phiếu thu/chi, nhập/xuất, kho, GTGT...):** Gọi Stored Procedure **`ALLOC_SCT`** (`@SCT OUT, @LOAI, @NGAY, @MAKHO`). Procedure dùng cơ chế atomic counter trên bảng `SYS_SEQUENCE` thông qua procedure **`SYS_SEQ_NEXTVAL`**.
- **Final Payment & Move to Main Tables (Thanh toán chính thức):**
  - Khi thu ngân bấm Thanh toán & In (`F4` / `CmdPrint`), hệ thống cập nhật `PRINTED = 1` và gọi Stored Procedure **`POS_PAID(@KHOA)`**.
  - **`POS_PAID` thực hiện trong 1 Transaction ACID duy nhất:**
    1. Chuyển toàn bộ dữ liệu từ `BANLE_TAM` sang bảng chính **`BANLE`**.
    2. Chuyển toàn bộ chi tiết từ `BANLE_TAM_CT` sang bảng chính **`BANLE_CT`**.
    3. Xóa dữ liệu tương ứng khỏi `BANLE_TAM_CT` và `BANLE_TAM`.
    4. Đồng bộ / Merge thông tin vào bảng **`ehoadon`** để phục vụ xuất hóa đơn điện tử.
- **Print & Audit:**
  - Gửi lệnh in bill ra máy in nhiệt (K80/K58).
  - Cập nhật số lần in (`spUPDATE_PRINT_NO`).

---

### 2.2 RETURN/REFUND Process (Trả hàng)

#### 2.2.1 Return Types
1. **Immediate Return** (Same day, same cashier)
2. **Return with Receipt** (Any day, requires original receipt)
3. **Partial Return** (Return some items from original sale)

#### 2.2.2 Return Process Flow
```
[Scan Receipt Barcode] → [Select Items to Return] → [Input Quantity] 
→ [Verify Return Policy] → [Calculate Refund] → [Process Refund] → [Print Return Receipt]
```

#### 2.2.3 Return Business Rules
- **Return within period**: Check TRA_DATE against current date
- **Return quantity validation**: Cannot exceed original quantity
- **Price preservation**: Use original sale price (DONGIA from original BANLE_CT)
- **Discount preservation**: Keep original discount (TL_CK)
- **Link to original**: KHOA2 references original BANLE.KHOA
- **Negative quantities**: Store as negative in TRAHANG_CT

#### 2.2.4 Return Calculation
```
Return amount = Original amount with negative sign
THANHTIEN_CHUA_VAT = -1 * (SOLUONG * DONGIA_CHUA_VAT)
TIEN_THUE = -1 * (THANHTIEN_CHUA_VAT * THUE_SUAT / 100)
THANHTOAN = THANHTIEN_CHUA_VAT + TIEN_THUE
```

#### 2.2.5 Refund Methods
- Cash refund (if original payment was cash)
- Bank transfer (if paid by card/transfer)
- Credit note (for account customers)

---

### 2.3 PROMOTION Management (Khuyến mãi)

#### 2.3.1 Promotion Structure
**Header (KHUYENMAI):**
- Promotion code (SCT)
- Effective period (TUNGAY, DENNGAY)
- Time of day (TUGIO, DENGIO)
- Promotion type (HTKM):
  - Discount percentage (%)
  - Fixed price
  - Free items
- Discount rate (TL_CK) or Fixed price (GIABAN)

**Details (KHUYENMAI_CT):**
- Applicable products (MAVT)
- Product-specific pricing
- Optional: by product group, supplier, etc.

#### 2.3.2 Promotion Application Logic
1. **Time validation**: Current datetime within [TUNGAY TUGIO, DENNGAY DENGIO]
2. **Product matching**: MAVT in KHUYENMAI_CT OR promotion applies to all
3. **Priority**: If multiple promotions apply, use best for customer
4. **Recording**: Store promotion discount in TL_CK4

---

### 2.4 DATA MASTER References

#### 2.4.1 Product Master (DM_VT)
**Key fields:**
- MAVT: Product code (primary key)
- TENVT: Product name
- TENVT2: Short name
- TENTAT: Abbreviation
- DVT: Unit of measure
- LOAITHUE: Tax type
- THUE_SUAT: Default tax rate
- GIABAN: Retail price
- TL_CK_MAX: Maximum discount allowed
- BO: Boolean - is this a product bundle?

**Bundle products (DM_VT_BO):**
- MABO: Bundle code
- MAVT: Component product code
- Automatic expansion at POS when scanned

#### 2.4.2 Customer/Supplier (DM_KH_NCC)
**Types (LOAI):**
- 0: Customer only
- 1: Supplier only
- 3: Both customer and supplier

**For POS, only customer fields used:**
- MADT: Customer code
- TENDT: Customer name
- LOAI_VIP: VIP level (reference only, logic in NXTz)
- MST: Tax code (for invoice)
- DCHI: Address
- DTHOAI: Phone
- EMAIL: Email

#### 2.4.3 Location/Store (DM_LOCATION)
- LOC: Location code
- TENKHO: Store name
- MAKHO: Default warehouse
- HoaDon_Mau, HoaDon_KyHieu: eInvoice template (reference)

#### 2.4.4 Counter/POS Terminal (DM_QUAYTN)
- QUAYKE: Counter code
- MAQUAY: POS terminal ID

#### 2.4.5 System Configuration (SYS_CONFIG)
**Key settings:**
- Receipt format (RECEIPT_FORMAT)
- Barcode stamp prefix (STAMP_PREFIX)
- Tax configuration
- Business info for receipt header

---

## 3. BUSINESS RULES & CONSTRAINTS

### 3.1 Pricing Rules
1. **Price priority:**
   - Promotion price (if active)
   - VIP price (if customer is VIP)
   - Standard retail price (GIABAN)

2. **Negative pricing not allowed**

3. **Zero-price items:**
   - Only allowed for promotional/free items
   - Must be marked as ItemTypeID = 15 for eInvoice

### 3.2 Discount Rules
1. **Maximum discount enforcement:**
   - TL_CK <= TL_CK_MAX (unless authorized)
   - TL_CK_THEM requires manager approval

2. **Discount authorization levels:**
   - Cashier: Up to configured limit
   - Supervisor: Up to TL_CK_MAX
   - Manager: Any discount (TL_CK_MAX + TL_CK_THEM)

3. **Discount tracking:**
   - Record who applied: CK_BY (User ID)
   - Record all discount sources in TL_CK2..TL_CK7

### 3.3 Inventory Rules
1. **Stock check:**
   - POS performs read-only stock query
   - Warning if stock < order quantity (but not blocking)
   - Actual stock update handled by backend/NXTz

2. **No negative stock allowed** (enforced at warehouse level)

### 3.4 Transaction Rules
1. **Unique transaction ID:** KHOA (GUID)

2. **Sequential receipt numbering:**
   - Format: LOC + Date + Sequence
   - E.g., "01-20241014-0001"

3. **Audit trail:**
   - CREATE_BY, CREATE_DATE: Creation
   - UPDATE_BY, UPDATE_DATE: Modifications
   - DELETE_BY, DELETE_DATE: Soft delete
   - DRC_STATUS: Sync status

4. **No deletion after print:**
   - Once PRINTED = 1, cannot delete
   - Can only void/cancel with proper authorization

5. **Return policy:**
   - Maximum return period: configured per store
   - Return quantity <= original quantity
   - Original receipt required (link via KHOA2)

### 3.5 Tax Rules
1. **Tax calculation:**
   - Always calculated from pre-tax amount
   - Round to 2 decimal places
   - Breakdown by tax rate required for eInvoice

2. **Tax-exempt items:**
   - LOAITHUE = 'KCT'
   - No tax amount calculated
   - Separate total for reporting

---

## 4. USER ROLES & PERMISSIONS

### 4.1 Cashier (Thu ngân)
**Permissions:**
- Process sales
- Apply standard discounts (within limit)
- Process returns (same day)
- Reprint receipts
- View customer info

**Restrictions:**
- Cannot exceed discount limit
- Cannot void printed transactions
- Cannot access management reports

### 4.2 Supervisor (Giám sát)
**Permissions:**
- All Cashier permissions
- Apply higher discounts (up to TL_CK_MAX)
- Process returns (any day)
- Void transactions (with reason)
- Access shift reports

### 4.3 Manager (Quản lý)
**Permissions:**
- All Supervisor permissions
- Apply any discount (TL_CK_MAX + TL_CK_THEM)
- Modify pricing
- Manage promotions
- Full report access
- System configuration

---

## 5. REPORTING REQUIREMENTS

### 5.1 Cashier Reports
1. **Shift Summary (Báo cáo ca)**
   - Total sales by payment method
   - Number of transactions
   - Average transaction value
   - Discounts given

2. **Daily Sales Report**
   - Sales by time period
   - Top selling products
   - Returns summary

### 5.2 Management Reports
1. **Sales Analysis**
   - By product, category, supplier
   - By cashier, counter, location
   - Trend analysis

2. **Discount Report**
   - Total discounts by type
   - By cashier/authorization level
   - Unusual discount patterns

3. **Promotion Effectiveness**
   - Sales during promotion period
   - Promotion ROI

---

## 6. NON-FUNCTIONAL REQUIREMENTS

### 6.1 Performance
- **Transaction processing:** < 2 seconds per item scan
- **Payment processing:** < 3 seconds
- **Receipt printing:** < 5 seconds
- **Concurrent users:** Support 50+ POS terminals per server

### 6.2 Reliability
- **Uptime:** 99.5% during business hours
- **Offline capability:** Queue transactions for 2 hours without server
- **Data integrity:** No transaction loss

### 6.3 Usability
- **Keyboard shortcuts:** All functions accessible via keyboard
- **Barcode scanner:** Primary input method
- **Touch screen:** Optional support
- **Response time:** Immediate visual feedback

### 6.4 Security
- **Authentication:** User login required
- **Authorization:** Role-based access control
- **Audit trail:** All transactions logged
- **Data encryption:** Sensitive data encrypted at rest

### 6.5 Compliance
- **Tax compliance:** VAT calculation per Vietnamese regulations
- **eInvoice ready:** Data format compatible with BKAV integration
- **Receipt format:** Comply with Ministry of Finance requirements

---

## 7. DATA RETENTION & ARCHIVAL

### 7.1 Retention Policy
- **Active transactions:** 1 year in operational database
- **Historical data:** 7 years in archive
- **Audit logs:** 3 years

### 7.2 Archival Process
- Monthly archival of transactions older than 1 year
- Maintain summary reports indefinitely
- Purge deleted records after 1 year

---

## 8. INTEGRATION POINTS

### 8.1 Internal Integration
- **NXTz module:** Inventory updates, VIP management, purchase/stock
- **Accounting module:** Journal entries, financial reports
- **eInvoice module:** Invoice data export to BKAV

### 8.2 External Integration
- **Receipt printer:** ESC/POS protocol
- **Barcode scanner:** Standard HID/COM port
- **Cash drawer:** Via printer kick-drawer command
- **Customer display:** Serial/USB pole display

### 8.3 Data Sync
- **Real-time:** Transaction commits
- **Batch:** End-of-day reconciliation
- **Manual:** Return reconciliation

---

## 9. ACCEPTANCE CRITERIA

### 9.1 Functional Criteria
- [ ] Complete a sale transaction from scan to receipt
- [ ] Apply all 7 levels of discount correctly
- [ ] Process split payment (2 payment methods)
- [ ] Process full and partial returns
- [ ] Apply time-based promotions automatically
- [ ] Print receipt with correct tax breakdown
- [ ] Handle bundle products (auto-expansion)
- [ ] Process free items (100% discount)

### 9.2 Data Validation
- [ ] All sales match tax calculation rules
- [ ] Discounts within authorized limits
- [ ] Return quantities do not exceed original
- [ ] Sequential receipt numbering with no gaps
- [ ] All transactions have complete audit trail

### 9.3 Performance Criteria
- [ ] Scan-to-display: < 1 second
- [ ] Complete transaction: < 30 seconds (10 items)
- [ ] Receipt print: < 5 seconds
- [ ] System handles 50 concurrent terminals

### 9.4 Integration Criteria
- [ ] eInvoice export format validated
- [ ] Inventory sync confirmed
- [ ] Receipt printer compatible

---

## 10. GLOSSARY

| Vietnamese | English | Description |
|------------|---------|-------------|
| Bán lẻ | Retail Sale | Point of sale transaction |
| Trả hàng | Return/Refund | Product return |
| Khuyến mãi | Promotion | Sales promotion |
| Chiết khấu | Discount | Price discount |
| Thu ngân | Cashier | POS operator |
| Quầy | Counter | POS terminal |
| Ca | Shift | Work shift |
| Kho | Warehouse | Storage location |
| Vật tư | Product/Item | Merchandise |
| Đơn vị tính | Unit | Unit of measure |
| Thuế VAT | VAT | Value-added tax |
| Hóa đơn | Invoice | Tax invoice |
| Mã vạch | Barcode | Product barcode |

---

**Document Version:** 1.0  
**Last Updated:** 2026-10-14  
**Author:** System Analyst (AI-assisted)  
**Status:** Draft for Review
