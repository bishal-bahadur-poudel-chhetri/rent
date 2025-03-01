package services

import (
	"log"
	"renting/internal/models"
	"renting/internal/repositories"
)

// SaleService contains the business logic for sales
type SaleService struct {
	SaleRepository repositories.SaleRepository
}

// RegisterSale registers a new sale and returns the sale ID
func (s *SaleService) RegisterSale(sale *models.Sale) (int, error) {
	// Here, we could include business logic like validations, calculations, etc.

	saleID, err := s.SaleRepository.CreateSale(sale)
	if err != nil {
		log.Println("Error in repository while registering sale:", err)
		return 0, err
	}

	return saleID, nil
}
