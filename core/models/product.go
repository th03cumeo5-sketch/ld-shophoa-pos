package models

import "time"

// Product ánh xạ bảng DM_HH (Danh mục hàng hóa / hoa / phụ liệu)
type Product struct {
	MaVT      string    `json:"mavt"`       // Mã vật tư / hàng hóa (PK)
	TenVT     string    `json:"tenvt"`      // Tên hàng hóa
	DVT       string    `json:"dvt"`        // Đơn vị tính (Cành, Bó, Giỏ, Chậu...)
	GiaBan    float64   `json:"giaban"`     // Đơn giá bán lẻ tiêu chuẩn
	GiaBan1   float64   `json:"giaban1"`    // Đơn giá bán phụ
	GiaSi     float64   `json:"giasi"`      // Đơn giá bán sỉ
	GiaVon    float64   `json:"giavon"`     // Giá vốn
	LoaiThue  string    `json:"loaithue"`   // Loại thuế (T00, T05, T08, T10, KCT)
	VatRa     float64   `json:"vat_ra"`     // Thuế suất VAT bán (%)
	MaNhom    string    `json:"manhom"`     // Nhóm hàng
	MaNganh   string    `json:"manganh"`    // Ngành hàng
	TinhTrang string    `json:"tinhtrang"`  // Tình trạng kinh doanh
	PLU       string    `json:"plu"`        // Mã barcode / cân điện tử
	Checked   bool      `json:"checked"`    // Đang kinh doanh
	UpdateDate time.Time `json:"update_date"`
}

// ProductSearchResult kết quả tìm kiếm nhanh sản phẩm trên POS
type ProductSearchResult struct {
	MaVT     string  `json:"mavt"`
	TenVT    string  `json:"tenvt"`
	DVT      string  `json:"dvt"`
	GiaBan   float64 `json:"giaban"`
	LoaiThue string  `json:"loaithue"`
	VatRa    float64 `json:"vat_ra"`
	TonKho   float64 `json:"tonkho"` // Tồn kho tức thời (nếu có)
}

// Category đại diện cho nhóm hàng trong DM_NHOM
type Category struct {
	MaNhom  string `json:"manhom"`
	TenNhom string `json:"tennhom"`
	Count   int    `json:"count"`
}
