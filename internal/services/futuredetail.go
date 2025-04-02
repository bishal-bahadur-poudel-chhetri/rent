package services

import (
	"fmt"
	"renting/internal/models"
	"renting/internal/repositories"
	"time"
)

type FuturBookingService struct {
	repo *repositories.FuturBookingRepository
}

func NewFuturBookingService(repo *repositories.FuturBookingRepository) *FuturBookingService {
	return &FuturBookingService{repo: repo}
}

// GetFuturBookingsByMonth fetches bookings for a specific month and metadata for future months
func (s *FuturBookingService) GetFuturBookingsByMonth(year int, month time.Month, filters map[string]string) (*models.SalesWithMetadataResponse, error) {
	return s.repo.GetFuturBookingsByMonth(year, month, filters)
}

// CancelFuturBooking cancels a future booking by sale ID
func (s *FuturBookingService) CancelFuturBooking(saleID int) error {
	// Input validation
	if saleID <= 0 {
		return fmt.Errorf("invalid sale ID")
	}

	// Call repository method
	err := s.repo.FutureBookingCancellation(saleID)
	if err != nil {
		return fmt.Errorf("failed to cancel booking: %w", err)
	}

	return nil
}
