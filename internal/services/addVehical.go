// service/vehicle_service.go
package services

import (
	"log"
	"renting/internal/models"
	"renting/internal/repositories"
)

type VehicleService struct {
	vehicleRepo *repositories.VehicleRepository
}

func NewVehicleService(vehicleRepo *repositories.VehicleRepository) *VehicleService {
	return &VehicleService{vehicleRepo: vehicleRepo}
}

func (s *VehicleService) GetVehicles(filters models.VehicleFilter, includeBookingDetails bool) ([]models.VehicleResponse, error) {
	log.Printf("Fetching vehicles with filters: %+v, includeBookingDetails: %v", filters, includeBookingDetails)
	return s.vehicleRepo.GetVehicles(filters, includeBookingDetails)
}
