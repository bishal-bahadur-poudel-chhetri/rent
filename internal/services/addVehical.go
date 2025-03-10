package services

import (
	"renting/internal/models"
	"renting/internal/repositories"
)

type VehicleService struct {
	vehicleRepo *repositories.VehicleRepository
}

func NewVehicleService(vehicleRepo *repositories.VehicleRepository) *VehicleService {
	return &VehicleService{
		vehicleRepo: vehicleRepo,
	}
}

// RegisterVehicle registers a new vehicle using the repository
func (s *VehicleService) RegisterVehicle(vehicle models.VehicleRequest) (int, error) {
	return s.vehicleRepo.RegisterVehicle(vehicle)
}

// ListVehicles lists vehicles with filters and pagination
func (s *VehicleService) ListVehicles(filter models.VehicleFilter) ([]models.VehicleResponse, error) {
	return s.vehicleRepo.ListVehicles(filter)
}
