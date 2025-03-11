package services

import (
	"renting/internal/models"
	"renting/internal/repositories"
)

type SaleService struct {
	saleRepo *repositories.SaleRepository
}

func NewSaleService(saleRepo *repositories.SaleRepository) *SaleService {
	return &SaleService{saleRepo: saleRepo}
}

func (s *SaleService) CreateSale(sale models.Sale) (int, error) {
	return s.saleRepo.CreateSale(sale)
}

func (s *SaleService) GetSaleByID(saleID int) (*models.Sale, error) {
	return s.saleRepo.GetSaleByID(saleID)
}
