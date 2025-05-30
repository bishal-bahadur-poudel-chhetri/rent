// service/vehicle_service.go
package services

import (
	"fmt"
	"renting/internal/models"
	"renting/internal/repositories"
)

type VehicleService struct {
	vehicleRepo *repositories.VehicleRepository
}

func NewVehicleService(vehicleRepo *repositories.VehicleRepository) *VehicleService {
	return &VehicleService{vehicleRepo: vehicleRepo}
}

func (s *VehicleService) GetVehicles(filters models.VehicleFilter, includeBookingDetails bool, includeSaleid bool) ([]models.VehicleResponse, error) {
	return s.vehicleRepo.GetVehicles(filters, includeBookingDetails, includeSaleid)
}

func (s *SaleService) AddCharge(req models.AddChargeRequest) error {
	// Validate charge type
	switch req.ChargeType {
	case "discount", "wash", "damage":
		// Valid charge types
	default:
		return fmt.Errorf("invalid charge type: %s", req.ChargeType)
	}

	// Validate amount
	if req.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	// Add the charge
	err := s.saleRepo.AddCharge(req.SaleID, req.ChargeType, req.Amount, req.Remark)
	if err != nil {
		return fmt.Errorf("failed to add charge: %v", err)
	}

	return nil
}
