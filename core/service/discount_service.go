package service

import (
	"math"

	"shophoa/core/models"
)

type DiscountService struct{}

func NewDiscountService() *DiscountService {
	return &DiscountService{}
}

// CalculateLineItemDiscount tính toán chiết khấu đa cấp cho 1 dòng mặt hàng
func (s *DiscountService) CalculateLineItemDiscount(item *models.BanleTamCT, maxDiscount float64) {
	// 1. TL_CK là giá trị lớn nhất trong các nguồn chiết khấu
	maxCK := math.Max(item.TL_CK2, math.Max(item.TL_CK3, math.Max(item.TL_CK4, item.TL_CK5)))
	if maxCK < 0 {
		maxCK = 0
	}

	// 2. Kiểm tra chặn trần TL_CK_MAX (nếu có quy định và không có cấp quyền vượt trần)
	if maxDiscount > 0 && maxCK > maxDiscount {
		maxCK = maxDiscount
	}

	item.TL_CK = maxCK

	// 3. Tính đơn giá sau chiết khấu và thành tiền
	// Làm tròn số học chuẩn VND
	item.SoTien = math.Round(item.SoLuong * item.DonGia)
	item.ChietKhau = math.Round(item.SoTien * (item.TL_CK / 100.0))
	item.DonGiaSauCK = math.Round(item.DonGia * (1.0 - item.TL_CK/100.0))
	item.ThanhTienChuaVAT = item.SoTien - item.ChietKhau
}

// CalculateHeaderDiscount tính chiết khấu trên tổng hóa đơn
func (s *DiscountService) CalculateHeaderDiscount(header *models.BanleTam) {
	if header.TL_CK > 0 {
		header.ChietKhau = math.Round(header.ThanhTienChuaVAT * (header.TL_CK / 100.0))
	} else {
		header.ChietKhau = 0
	}
}
