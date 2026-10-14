package models

// Customer ánh xạ bảng DM_KH_NCC (Khách hàng & Nhà cung cấp)
type Customer struct {
	MaDT    string  `json:"madt"`     // Mã đối tượng / khách hàng (PK)
	TenDT   string  `json:"tendt"`    // Tên khách hàng
	DienThoai string `json:"dienthoai"` // SĐT liên hệ
	DiaChi  string  `json:"diachi"`   // Địa chỉ
	MST     string  `json:"mst"`      // Mã số thuế
	LoaiVIP string  `json:"loai_vip"` // Hạng VIP (VIP Vàng, Bạc, Kim Cương...)
	SoDu    float64 `json:"sodu"`     // Số dư công nợ
}
