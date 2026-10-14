package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"shophoa/core/models"
	"shophoa/core/repository"
)

type SaleService struct {
	saleRepo     *repository.SaleRepository
	productRepo  *repository.ProductRepository
	discountServ *DiscountService
	taxServ      *TaxService
}

func NewSaleService(
	saleRepo *repository.SaleRepository,
	productRepo *repository.ProductRepository,
	discountServ *DiscountService,
	taxServ *TaxService,
) *SaleService {
	return &SaleService{
		saleRepo:     saleRepo,
		productRepo:  productRepo,
		discountServ: discountServ,
		taxServ:      taxServ,
	}
}

// PrepareNewSale khởi tạo một hóa đơn mới, tự động cấp số hóa đơn SCT từ Stored Procedure POS_ALLOC_BILL_NO
func (s *SaleService) PrepareNewSale(ctx context.Context, makho, quay, ca, thungan, loc string) (*models.BanleTam, error) {
	now := time.Now()
	sct, err := s.saleRepo.AllocateRetailBillNumber(ctx, makho, quay, now)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate bill number: %w", err)
	}

	return &models.BanleTam{
		Khoa:       uuid.New(),
		LCT:        "BLE",
		Ngay:       now,
		SCT:        sct,
		LOC:        loc,
		MaKho:      makho,
		Quay:       quay,
		Ca:         ca,
		ThuNgan:    thungan,
		Printed:    false,
		CreateDate: now,
		UpdateDate: now,
	}, nil
}

// RecalculateAndSaveDraft tính toán chiết khấu, thuế, tổng tiền và lưu vào cụm bảng BANLE_TAM / BANLE_TAM_CT
func (s *SaleService) RecalculateAndSaveDraft(ctx context.Context, header *models.BanleTam, items []models.BanleTamCT) error {
	// 1. Tính toán chiết khấu và thuế cho từng dòng
	for i := range items {
		item := &items[i]
		s.discountServ.CalculateLineItemDiscount(item, 50.0) // Trần chiết khấu 50%
		s.taxServ.CalculateTaxForLine(item)
	}

	// 2. Tính chiết khấu hóa đơn
	s.discountServ.CalculateHeaderDiscount(header)

	// 3. Tổng hợp toàn bộ hóa đơn
	s.taxServ.CalculateTotals(header, items)

	// 4. Lưu vào bảng tạm
	return s.saleRepo.SaveDraftSale(ctx, header, items)
}

// CompletePayment thực hiện quy trình thanh toán chính thức
// 1. Cập nhật số tiền khách đưa và hình thức thanh toán vào BANLE_TAM
// 2. Đánh dấu PRINTED = 1
// 3. Gọi Stored Procedure POS_PAID(@KHOA) chuyển dữ liệu sang BANLE/BANLE_CT và merge sang ehoadon
func (s *SaleService) CompletePayment(ctx context.Context, header *models.BanleTam, items []models.BanleTamCT, pttt string, ttoan1, ttoan2 float64) error {
	header.PTTT = pttt
	header.TToan1 = ttoan1
	header.TToan2 = ttoan2
	header.Printed = true
	header.PrintNo++

	// Cập nhật lại toàn bộ số liệu lần cuối vào bảng tạm
	if err := s.RecalculateAndSaveDraft(ctx, header, items); err != nil {
		return fmt.Errorf("failed to update draft before payment: %w", err)
	}

	// Kiểm tra số tiền thanh toán
	if (ttoan1 + ttoan2) < header.ThanhToan {
		return fmt.Errorf("insufficient payment: paid %.0f < required %.0f", ttoan1+ttoan2, header.ThanhToan)
	}

	// Chuyển sang bảng chính thức qua Stored Procedure POS_PAID
	if err := s.saleRepo.CompleteSalePayment(ctx, header.Khoa); err != nil {
		return fmt.Errorf("POS_PAID execution failed: %w", err)
	}

	return nil
}
