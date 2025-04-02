package services

import (
	"fmt"
	"renting/internal/models"
	"renting/internal/repositories"
	"time"
)

type RevenueService struct {
	repo *repositories.RevenueRepository
}

func NewRevenueService(repo *repositories.RevenueRepository) *RevenueService {
	return &RevenueService{repo: repo}
}

type RevenueResponse struct {
	Period       string  `json:"period"`
	TotalRevenue float64 `json:"total_revenue"`
}

// GetRevenue fetches revenue based on the request
func (s *RevenueService) GetRevenue(req models.RevenueRequest) (RevenueResponse, error) {
	filter := repositories.RevenueFilter{Period: req.Period}

	// Parse date if provided
	if req.Date != "" {
		date, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return RevenueResponse{}, fmt.Errorf("invalid date format: %v", err)
		}
		filter.Date = date
	}

	total, err := s.repo.GetTotalRevenue(filter)
	if err != nil {
		return RevenueResponse{}, err
	}

	return RevenueResponse{
		Period:       req.Period,
		TotalRevenue: total,
	}, nil
}
