package services

import (
	"fmt"
	"renting/internal/models"
	"renting/internal/repositories"
	"time"
)

type DisableDateService struct {
	disableDateRepo *repositories.DisableDateRepository
}

func NewDisableDateService(disableDateRepo *repositories.DisableDateRepository) *DisableDateService {
	return &DisableDateService{disableDateRepo: disableDateRepo}
}

// GetDisabledDates fetches disabled dates for a specific vehicle and date range
func (s *DisableDateService) GetDisabledDates(vehicleID int, dateOfDelivery time.Time) (*models.DisableDateResponse, error) {
	// Call the repository to fetch disabled dates
	response, err := s.disableDateRepo.GetDisabledDates(vehicleID, dateOfDelivery)
	if err != nil {
		return nil, fmt.Errorf("failed to get disabled dates: %v", err)
	}

	return response, nil
}
