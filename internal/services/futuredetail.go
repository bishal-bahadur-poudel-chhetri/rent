package services

import (
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
