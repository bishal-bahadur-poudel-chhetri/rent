package services

import (
	"fmt"
	"renting/internal/models"
	"renting/internal/repositories"
	"time"
)

type SaleDetailService struct {
	saleRepo *repositories.SaleDetailRepository
}

func NewSaleDetailService(saleRepo *repositories.SaleDetailRepository) *SaleDetailService {
	return &SaleDetailService{saleRepo: saleRepo}
}

func (s *SaleDetailService) GetSalesWithFilters(filters map[string]string) ([]models.Sale, error) {
	// Validate filters (optional)
	if startDate, ok := filters["start_date"]; ok {
		if _, err := time.Parse("2006-01-02", startDate); err != nil {
			return nil, fmt.Errorf("invalid start_date format: %v", err)
		}
	}
	if endDate, ok := filters["end_date"]; ok {
		if _, err := time.Parse("2006-01-02", endDate); err != nil {
			return nil, fmt.Errorf("invalid end_date format: %v", err)
		}
	}

	// Call the repository
	sales, err := s.saleRepo.GetSalesWithFilters(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales: %v", err)
	}

	return sales, nil
}
