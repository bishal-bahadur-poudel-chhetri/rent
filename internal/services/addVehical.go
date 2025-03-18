// service/vehicle_service.go
package services

import (
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
