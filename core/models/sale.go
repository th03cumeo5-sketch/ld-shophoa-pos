package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// BanleTam Header của hóa đơn đang soạn / lưu tạm (F10) trong bảng BANLE_TAM
type BanleTam struct {
	Khoa              uuid.UUID `json:"khoa"` // GUID PK
	LCT               string    `json:"lct"`  // Loại chứng từ: 'BLE'
	Ngay              time.Time `json:"ngay"` // Ngày giờ tạo
	SCT               string    `json:"sct"`  // Số chứng từ do POS_ALLOC_BILL_NO cấp
	LOC               string    `json:"loc"`  // Chi nhánh (Store)
	MaKho             string    `json:"makho"`
	Quay              string    `json:"quay"`    // Quầy thu ngân
	Ca                string    `json:"ca"`      // Ca làm việc
	ThuNgan           string    `json:"thungan"` // Username thu ngân
	MaDT              string    `json:"madt"`    // Mã khách hàng (DM_KH_NCC)
	MaVIP             string    `json:"mavip"`   // Thẻ VIP
	SoLuong           float64   `json:"soluong"` // Tổng số lượng
	SoTien            float64   `json:"sotien"`  // Tổng tiền hàng trước chiết khấu
	ChietKhauMH       float64   `json:"chietkhau_mh"` // Tổng chiết khấu từng mặt hàng
	TL_CK             float64   `json:"tl_ck"`        // % chiết khấu toàn hóa đơn
	ChietKhau         float64   `json:"chietkhau"`    // Số tiền chiết khấu toàn hóa đơn
	ThanhTienChuaVAT  float64   `json:"thanhtien_chua_vat"`
	TienThue          float64   `json:"tien_thue"`    // Tổng tiền thuế VAT
	TienThueCL        float64   `json:"tien_thue_cl"` // Tiền thuế chênh lệch làm tròn
	ThanhToan         float64   `json:"thanhtoan"`    // Khách phải trả
	PTTT              string    `json:"pttt"`         // TM (Tiền mặt), CK (Chuyển khoản), THE (Thẻ)
	TToan1            float64   `json:"ttoan1"`       // Tiền khách đưa hình thức 1
	TToan2            float64   `json:"ttoan2"`       // Tiền khách đưa hình thức 2
	ChuaThoi          float64   `json:"chuathoi"`     // Tiền thừa trả lại khách
	Printed           bool      `json:"printed"`      // Đã in bill hay chưa
	PrintNo           int       `json:"print_no"`     // Số lần in
	DGIAI             string    `json:"dgiai"`        // Ghi chú đơn tạm / Diễn giải
	CreateBy          int       `json:"create_by"`
	CreateDate        time.Time `json:"create_date"`
	UpdateBy          int       `json:"update_by"`
	UpdateDate        time.Time `json:"update_date"`
}

// BanleTamCT Chi tiết mặt hàng của hóa đơn tạm trong bảng BANLE_TAM_CT
type BanleTamCT struct {
	KhoaCT            uuid.UUID `json:"khoact"` // GUID PK
	Khoa              uuid.UUID `json:"khoa"`   // FK -> BanleTam.Khoa
	STT               int       `json:"stt"`    // Thứ tự dòng (1, 2, 3...)
	MaVT              string    `json:"mavt"`   // Mã hàng hoa
	TenVT             string    `json:"tenvt"`  // Tên hàng hoa
	DVT               string    `json:"dvt"`    // ĐVT
	SoLuong           float64   `json:"soluong"`
	DonGia            float64   `json:"dongia"`       // Đơn giá bán
	DonGiaSauCK       float64   `json:"dongia_sauck"` // Đơn giá sau khi trừ chiết khấu dòng
	SoTien            float64   `json:"sotien"`       // SOLUONG * DONGIA
	TL_CK             float64   `json:"tl_ck"`        // % chiết khấu tổng hợp
	TL_CK2            float64   `json:"tl_ck2"`       // % chiết khấu VIP
	TL_CK3            float64   `json:"tl_ck3"`       // % chiết khấu Nhóm
	TL_CK4            float64   `json:"tl_ck4"`       // % chiết khấu Khuyến mãi
	TL_CK5            float64   `json:"tl_ck5"`       // % chiết khấu Thu ngân
	ChietKhau         float64   `json:"chietkhau"`    // Tiền chiết khấu dòng
	LoaiThue          string    `json:"loaithue"`     // T00, T05, T08, T10, KCT
	ThueSuat          float64   `json:"thue_suat"`    // 0, 5, 8, 10
	TienThue          float64   `json:"tien_thue"`    // Tiền thuế
	TienThueCL        float64   `json:"tien_thue_cl"` // Cân đối làm tròn thuế
	ThanhTienChuaVAT  float64   `json:"thanhtien_chua_vat"`
	ThanhToan         float64   `json:"thanhtoan"`    // Thành tiền cuối cùng của dòng
	GhiChu            string    `json:"ghichu"`       // Ghi chú dòng (đổi trả, yêu cầu đặc biệt)
}

// SaleDTO Đối tượng truyền tải dữ liệu đầy đủ của một đơn hàng (Header + Chi tiết các dòng)
type SaleDTO struct {
	Header BanleTam     `json:"header"`
	Items  []BanleTamCT `json:"items"`
}

// UnmarshalJSON cho BanleTam tự động xử lý chuỗi rỗng / GUID không hợp lệ / định dạng ngày tháng linh hoạt
func (b *BanleTam) UnmarshalJSON(data []byte) error {
	type Alias BanleTam
	aux := &struct {
		Khoa       string `json:"khoa"`
		Ngay       string `json:"ngay"`
		CreateDate string `json:"create_date"`
		UpdateDate string `json:"update_date"`
		*Alias
	}{
		Alias: (*Alias)(b),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Khoa == "" || aux.Khoa == "null" {
		b.Khoa = uuid.New()
	} else {
		u, err := uuid.Parse(aux.Khoa)
		if err != nil {
			b.Khoa = uuid.New()
		} else {
			b.Khoa = u
		}
	}

	if aux.Ngay != "" {
		if t, err := time.Parse(time.RFC3339, aux.Ngay); err == nil {
			b.Ngay = t
		} else if t, err := time.Parse("2006-01-02T15:04:05.000Z", aux.Ngay); err == nil {
			b.Ngay = t
		} else {
			b.Ngay = time.Now()
		}
	} else if b.Ngay.IsZero() {
		b.Ngay = time.Now()
	}

	if aux.CreateDate != "" {
		if t, err := time.Parse(time.RFC3339, aux.CreateDate); err == nil {
			b.CreateDate = t
		}
	}
	if aux.UpdateDate != "" {
		if t, err := time.Parse(time.RFC3339, aux.UpdateDate); err == nil {
			b.UpdateDate = t
		}
	}

	return nil
}

// UnmarshalJSON cho BanleTamCT tự động sinh UUID mới nếu client gửi mã dạng chuỗi tạm thời (ct-...) hoặc rỗng
func (c *BanleTamCT) UnmarshalJSON(data []byte) error {
	type Alias BanleTamCT
	aux := &struct {
		KhoaCT string `json:"khoact"`
		Khoa   string `json:"khoa"`
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.KhoaCT == "" || aux.KhoaCT == "null" {
		c.KhoaCT = uuid.New()
	} else {
		u, err := uuid.Parse(aux.KhoaCT)
		if err != nil {
			c.KhoaCT = uuid.New()
		} else {
			c.KhoaCT = u
		}
	}

	if aux.Khoa == "" || aux.Khoa == "null" {
		c.Khoa = uuid.Nil
	} else {
		u, err := uuid.Parse(aux.Khoa)
		if err != nil {
			c.Khoa = uuid.Nil
		} else {
			c.Khoa = u
		}
	}

	return nil
}
