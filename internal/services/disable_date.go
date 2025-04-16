package services

import (
	"fmt"
	"renting/internal/models"
	"renting/internal/repositories"
)

type DisableDateService struct {
	disableDateRepo repositories.DisableDateRepository
}

func NewDisableDateService(disableDateRepo repositories.DisableDateRepository) *DisableDateService {
	return &DisableDateService{disableDateRepo: disableDateRepo}
}

func (s *DisableDateService) GetDisabledDates(vehicleID int, excludeSaleID *int) (*models.DisableDateResponse, error) {
	response, err := s.disableDateRepo.GetDisabledDates(vehicleID, excludeSaleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get disabled dates: %w", err)
	}

	// Ensure we never return nil for the slices
	if response == nil {
		response = &models.DisableDateResponse{
			ActiveRentals:  []models.DisabledDateResponse{},
			FutureBookings: []models.DisabledDateResponse{}, // Changed to FutureBookings
		}
	}

	return response, nil
}

