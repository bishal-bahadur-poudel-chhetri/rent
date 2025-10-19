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

func (s *SaleChargeService) UpdateSalesCharge(chargeID int, charge models.SalesCharge) error {
	err := s.saleRepo.UpdateSalesCharge(chargeID, charge)
	if err != nil {
		return fmt.Errorf("failed to update sales charge: %v", err)
	}
	return nil
}

func (s *SaleChargeService) DeleteSalesCharge(chargeID int) error {
	err := s.saleRepo.DeleteSalesCharge(chargeID)
	if err != nil {
		return fmt.Errorf("failed to delete sales charge: %v", err)
	}
	return nil
}
