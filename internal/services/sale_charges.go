package services

import (
	"fmt"
	"renting/internal/models"
	"renting/internal/repositories"
)

type SaleChargeService struct {
	saleRepo *repositories.SaleChargeRepository
}

func NewSaleChargeService(saleRepo *repositories.SaleChargeRepository) *SaleChargeService {
	return &SaleChargeService{
		saleRepo: saleRepo,
	}
}

func (s *SaleChargeService) AddSalesCharge(saleID int, charge models.SalesCharge) error {
	// Wrap the single charge in a slice to match the repository's expected input
	err := s.saleRepo.AddSalesCharges(saleID, []models.SalesCharge{charge})
	if err != nil {
		return fmt.Errorf("failed to add sales charge: %v", err)
	}
	return nil
}
