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

func (s *SaleChargeService) GetSalesChargesBySaleID(saleID int) ([]models.SalesCharge, error) {
	charges, err := s.saleRepo.GetSalesChargesBySaleID(saleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales charges: %v", err)
	}
	return charges, nil
}

func (s *SaleChargeService) UpdateSalesCharge(chargeID int, charge models.SalesCharge) error {
	// Validate charge type
	switch charge.ChargeType {
	case "discount", "wash", "damage":
		// Valid charge types
	default:
		return fmt.Errorf("invalid charge type: %s", charge.ChargeType)
	}

	// Validate amount
	if charge.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	err := s.saleRepo.UpdateSalesCharge(chargeID, charge)
	if err != nil {
		return fmt.Errorf("failed to update sales charge: %v", err)
	}
	return nil
}

func (s *SaleChargeService) DeleteSalesCharge(chargeID int, saleID int) error {
	err := s.saleRepo.DeleteSalesCharge(chargeID, saleID)
	if err != nil {
		return fmt.Errorf("failed to delete sales charge: %v", err)
	}
	return nil
}
