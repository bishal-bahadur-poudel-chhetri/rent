package services

import (
	"fmt"
	"renting/internal/models"
	"renting/internal/repositories"
)

type SaleService struct {
	saleRepo *repositories.SaleRepository
}

func NewSaleService(saleRepo *repositories.SaleRepository) *SaleService {
	return &SaleService{saleRepo: saleRepo}
}

func (s *SaleService) CreateSale(sale models.Sale) (models.SaleSubmitResponse, error) {
	// Call the repository method to create the sale
	response, err := s.saleRepo.CreateSale(sale)
	if err != nil {
		return models.SaleSubmitResponse{}, fmt.Errorf("failed to create sale: %v", err)
	}

	// Return the response from the repository
	return response, nil
}

func (s *SaleService) GetSaleByID(saleID int, include []string) (*models.Sale, error) {
	return s.saleRepo.GetSaleByID(saleID, include)
}
