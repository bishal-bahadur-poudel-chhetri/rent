package services

import (
	"fmt"
	"renting/internal/models"
	"renting/internal/repositories"
)

type VehicleServicingService struct {
	servicingRepo *repositories.VehicleServicingRepository
}

func NewVehicleServicingService(servicingRepo *repositories.VehicleServicingRepository) *VehicleServicingService {
	return &VehicleServicingService{servicingRepo: servicingRepo}
}

// InitializeServicingRecord creates a new servicing record for a vehicle
func (s *VehicleServicingService) InitializeServicingRecord(vehicleID int, initialKm float64, servicingInterval float64) error {
	if initialKm < 0 {
		return fmt.Errorf("initial km reading cannot be negative")
	}
	if servicingInterval <= 0 {
		return fmt.Errorf("servicing interval must be greater than 0")
	}
	return s.servicingRepo.InitializeServicingRecord(vehicleID, initialKm, servicingInterval)
}

// UpdateServicingStatus checks and updates the servicing status based on current km reading
func (s *VehicleServicingService) UpdateServicingStatus(vehicleID int, currentKm float64) error {
	if currentKm < 0 {
		return fmt.Errorf("current km reading cannot be negative")
	}
	return s.servicingRepo.UpdateServicingStatus(vehicleID, currentKm)
}

// MarkAsServiced updates the servicing record after a vehicle has been serviced
func (s *VehicleServicingService) MarkAsServiced(vehicleID int, currentKm float64, servicingType string, cost float64, notes string, servicedBy int) error {
	if currentKm < 0 {
		return fmt.Errorf("current km reading cannot be negative")
	}
	return s.servicingRepo.MarkAsServiced(vehicleID, currentKm, servicingType, cost, notes, servicedBy)
}

// GetCurrentKmAndServicingStatus retrieves the current km reading and servicing status for a vehicle
func (s *VehicleServicingService) GetCurrentKmAndServicingStatus(vehicleID int) (*models.VehicleServicing, error) {
	return s.servicingRepo.GetCurrentKmAndServicingStatus(vehicleID)
}

// GetVehiclesDueForServicing returns a list of vehicles that need servicing
func (s *VehicleServicingService) GetVehiclesDueForServicing() ([]models.VehicleServicing, error) {
	return s.servicingRepo.GetVehiclesDueForServicing()
}

// GetServicingHistory retrieves the servicing history for a vehicle
func (s *VehicleServicingService) GetServicingHistory(vehicleID int) ([]models.VehicleServicingHistory, error) {
	return s.servicingRepo.GetServicingHistory(vehicleID)
}
