package services

import (
	"renting/internal/models"
	"renting/internal/repositories"
)

type ReturnService struct {
	returnRepo *repositories.ReturnRepository
}

func NewReturnService(returnRepo *repositories.ReturnRepository) *ReturnService {
	return &ReturnService{returnRepo: returnRepo}
}

// Update the method signature to accept userID
func (s *ReturnService) CreateReturn(saleID int, userID int, returnRequest models.ReturnRequest) error {
	// Convert the return request to a Sale model
	sale := models.Sale{
		SaleID:       saleID,
		UserID:       userID, // Pass the userID to the Sale model
		SalesCharges: returnRequest.SalesCharges,
		VehicleUsage: returnRequest.VehicleUsage,
		Payments:     returnRequest.Payments,
	}

	// Call the repository to create the return record
	if _, err := s.returnRepo.CreateReturn(sale); err != nil {
		return err
	}

	return nil
}
