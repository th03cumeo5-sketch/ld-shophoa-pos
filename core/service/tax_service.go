package service

import (
	"math"

	"shophoa/core/models"
)

type TaxService struct{}

func NewTaxService() *TaxService {
	return &TaxService{}
}

// CalculateTaxForLine tính tiền thuế cho một dòng mặt hàng
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
		item.TienThue = math.Round(item.ThanhTienChuaVAT * (item.ThueSuat / 100.0))
	} else {
		item.TienThue = 0
	}

	item.ThanhToan = item.ThanhTienChuaVAT + item.TienThue
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
	header.ThanhTienChuaVAT = totalPreTax - header.ChietKhau
	header.TienThue = totalTax
	header.ThanhToan = finalPay - header.ChietKhau

	// Cân đối tiền thừa
	totalPaid := header.TToan1 + header.TToan2
	if totalPaid >= header.ThanhToan {
		header.ChuaThoi = totalPaid - header.ThanhToan
	} else {
		header.ChuaThoi = 0
	}
}
