package services

import (
	"renting/internal/models"
	"renting/internal/repositories"
	"time"
)

type PaymentService struct {
	paymentRepo *repositories.PaymentRepository
}

func NewPaymentService(paymentRepo *repositories.PaymentRepository) *PaymentService {
	return &PaymentService{paymentRepo: paymentRepo}
}

type SaleFilter struct {
	SaleID        *int
	PaymentStatus *string
	StartDate     *time.Time
	EndDate       *time.Time
	CustomerName  *string
	SaleStatus    *string
	VerifiedBy    *string
}

func (s *PaymentService) GetPaymentsWithSales(filter SaleFilter, limit int, offset int) ([]models.PaymentWithSale, error) {
	return s.paymentRepo.GetPaymentsWithSales(
		repositories.SaleFilter{
			SaleID:        filter.SaleID,
			PaymentStatus: filter.PaymentStatus,
			StartDate:     filter.StartDate,
			EndDate:       filter.EndDate,
			CustomerName:  filter.CustomerName,
			SaleStatus:    filter.SaleStatus,
			VerifiedBy:    filter.VerifiedBy,
		},
		limit,
		offset,
	)
}
