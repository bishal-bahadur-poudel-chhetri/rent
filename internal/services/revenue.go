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

// MonthlyRevenue represents revenue data for a specific month
type MonthlyRevenue struct {
	Month  time.Time `json:"month"`
	Amount float64   `json:"amount"`
}

// MobileRevenueResponse represents the response format for mobile visualization
type MobileRevenueResponse struct {
	Labels []string  `json:"labels"` // Month names (e.g., "Jan", "Feb")
	Data   []float64 `json:"data"`   // Revenue amounts
	Total  float64   `json:"total"`  // Total revenue for the period
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

// GetMonthlyRevenue returns monthly revenue data within a date range
func (s *RevenueService) GetMonthlyRevenue(startDate, endDate time.Time) ([]repositories.MonthlyRevenue, error) {
	// Ensure dates are at the start of their respective months
	startDate = time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
	endDate = time.Date(endDate.Year(), endDate.Month(), 1, 0, 0, 0, 0, endDate.Location())

	// Call the repository to get monthly revenue data
	monthlyData, err := s.repo.GetMonthlyRevenue(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly revenue: %v", err)
	}

	return monthlyData, nil
}

// GetMobileRevenueVisualization returns formatted data for mobile visualization
func (s *RevenueService) GetMobileRevenueVisualization(startDate, endDate time.Time) (MobileRevenueResponse, error) {
	// Get the monthly revenue data
	monthlyData, err := s.GetMonthlyRevenue(startDate, endDate)
	if err != nil {
		return MobileRevenueResponse{}, err
	}

	// Format the data for mobile consumption
	labels := make([]string, len(monthlyData))
	data := make([]float64, len(monthlyData))
	var total float64

	for i, mr := range monthlyData {
		// Format month as short name (e.g., "Jan", "Feb")
		labels[i] = mr.Month.Format("Jan")
		data[i] = mr.Amount
		total += mr.Amount
	}

	return MobileRevenueResponse{
		Labels: labels,
		Data:   data,
		Total:  total,
	}, nil
}
