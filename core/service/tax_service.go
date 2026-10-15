package service

import (
	"math"

	"shophoa/core/models"
)

type TaxService struct{}

func NewTaxService() *TaxService {
	return &TaxService{}
}

// CalculateTaxForLine tính tiền thuế cho một dòng mặt hàng (theo nguyên tắc giá bán lẻ đã gồm VAT)
func (s *TaxService) CalculateTaxForLine(item *models.BanleTamCT) {
	switch item.LoaiThue {
	case "T00", "KCT":
		item.ThueSuat = 0
	case "T05":
		item.ThueSuat = 5
	case "T08":
		item.ThueSuat = 8
	case "T10":
		item.ThueSuat = 10
	default:
		if item.ThueSuat <= 0 {
			item.ThueSuat = 10 // Mặc định 10% nếu không chỉ định
		}
	}

	if item.ThueSuat > 0 {
		// Bóc tách VAT từ số tiền thực trả của dòng
		item.ThanhTienChuaVAT = math.Round(item.ThanhToan / (1.0 + item.ThueSuat/100.0))
		item.TienThue = item.ThanhToan - item.ThanhTienChuaVAT
	} else {
		item.ThanhTienChuaVAT = item.ThanhToan
		item.TienThue = 0
	}
}

// CalculateTotals tính tổng cộng toàn bộ hóa đơn từ các dòng chi tiết
func (s *TaxService) CalculateTotals(header *models.BanleTam, items []models.BanleTamCT) {
	var totalQty float64
	var totalPreDiscount float64
	var totalLineDiscount float64
	var totalPreTax float64
	var totalTax float64
	var finalPay float64

	for i := range items {
		item := &items[i]
		totalQty += item.SoLuong
		totalPreDiscount += item.SoTien
		totalLineDiscount += item.ChietKhau
		totalPreTax += item.ThanhTienChuaVAT
		totalTax += item.TienThue
		finalPay += item.ThanhToan
	}

	header.SoLuong = totalQty
	header.SoTien = totalPreDiscount
	header.ChietKhauMH = totalLineDiscount

	header.ThanhToan = math.Max(0, finalPay-header.ChietKhau)
	if finalPay > 0 && header.ChietKhau > 0 {
		ratio := header.ThanhToan / finalPay
		header.ThanhTienChuaVAT = math.Round(totalPreTax * ratio)
		header.TienThue = header.ThanhToan - header.ThanhTienChuaVAT
	} else {
		header.ThanhTienChuaVAT = totalPreTax
		header.TienThue = totalTax
	}

	// Cân đối tiền thừa
	totalPaid := header.TToan1 + header.TToan2
	if totalPaid >= header.ThanhToan {
		header.ChuaThoi = totalPaid - header.ThanhToan
	} else {
		header.ChuaThoi = 0
	}
}

