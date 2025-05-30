package services

import (
	"fmt"
	"renting/internal/models"
	"renting/internal/repositories"
	"strconv"
	"time"
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

func (s *SaleService) GetSales(filters map[string]string, sort string, limit, offset int, include []string) (models.PendingSalesResponse, error) {
	// Convert map filters to SaleFilter struct
	saleFilter := models.SaleFilter{
		Sort:   sort,
		Limit:  limit,
		Offset: offset,
	}

	// Map the string filters to the appropriate fields
	if status, ok := filters["status"]; ok {
		saleFilter.Status = status
	}
	if customerName, ok := filters["customer_name"]; ok {
		saleFilter.CustomerName = customerName
	}
	if vehicleID, ok := filters["vehicle_id"]; ok {
		if id, err := strconv.Atoi(vehicleID); err == nil {
			saleFilter.VehicleID = id
		}
	}
	if deliveryDate, ok := filters["date_of_delivery"]; ok {
		if date, err := time.Parse("2006-01-02", deliveryDate); err == nil {
			saleFilter.DateOfDeliveryBefore = &date
		}
	}
	if actualDeliveryDate, ok := filters["actual_date_of_delivery"]; ok {
		if date, err := time.Parse("2006-01-02", actualDeliveryDate); err == nil {
			saleFilter.ActualDateOfDelivery = &date
		}
	}

	// Get sales from repository
	sales, err := s.saleRepo.GetSales(saleFilter)
	if err != nil {
		return models.PendingSalesResponse{}, fmt.Errorf("failed to fetch sales: %v", err)
	}

	// Calculate total count from the length of sales
	totalCount := len(sales)

	// Handle pagination calculations safely
	currentPage := 1
	totalPages := 1
	hasNext := false
	hasPrevious := offset > 0

	if limit > 0 {
		// Calculate current page (1-based index)
		currentPage = (offset / limit) + 1

		// Calculate total pages with ceiling division
		totalPages = (totalCount + limit - 1) / limit

		// Check if there's a next page
		hasNext = offset+limit < totalCount
	}

	return models.PendingSalesResponse{
		Sales: sales,
		Pagination: models.Pagination{
			CurrentPage: currentPage,
			PageSize:    limit,
			TotalCount:  totalCount,
			TotalPages:  totalPages,
			HasNext:     hasNext,
			HasPrevious: hasPrevious,
		},
	}, nil
}

func (s *SaleService) UpdateSaleByUserID(saleID, userID int, req models.UpdateSaleRequest) error {
	updates := make(map[string]interface{})

	// Populate updates map only with provided (non-nil) fields
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.PaymentStatus != nil {
		updates["payment_status"] = *req.PaymentStatus
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}
	if req.CustomerName != nil {
		updates["customer_name"] = *req.CustomerName
	}
	if req.CustomerPhone != nil {
		updates["customer_phone"] = *req.CustomerPhone
	}
	if req.CustomerDestination != nil {
		updates["customer_destination"] = *req.CustomerDestination
	}
	if req.TotalAmount != nil {
		updates["total_amount"] = *req.TotalAmount
	}
	if req.ChargePerDay != nil {
		updates["charge_per_day"] = *req.ChargePerDay
	}
	if req.VehicleID != nil { // Added this block
		updates["vehicle_id"] = *req.VehicleID
		fmt.Println("Adding vehicle_id to updates:", *req.VehicleID) // Debug
	}
	if req.DateOfDelivery != nil {
		date, err := time.Parse("2006-01-02", *req.DateOfDelivery)
		if err != nil {
			return fmt.Errorf("invalid date_of_delivery format: %v", err)
		}
		updates["date_of_delivery"] = date
	}
	if req.ReturnDate != nil {
		date, err := time.Parse("2006-01-02", *req.ReturnDate)
		if err != nil {
			return fmt.Errorf("invalid return_date format: %v", err)
		}
		updates["return_date"] = date
	}
	if req.ActualDateOfDelivery != nil {
		date, err := time.Parse("2006-01-02", *req.ActualDateOfDelivery)
		if err != nil {
			return fmt.Errorf("invalid actual_date_of_delivery format: %v", err)
		}
		updates["actual_date_of_delivery"] = date
	}
	if req.ActualDateOfReturn != nil {
		date, err := time.Parse("2006-01-02", *req.ActualDateOfReturn)
		if err != nil {
			return fmt.Errorf("invalid actual_date_of_return format: %v", err)
		}
		updates["actual_date_of_return"] = date
	}
	if req.NumberOfDays != nil {
		updates["number_of_days"] = *req.NumberOfDays
	}

	// Check if any fields were provided
	if len(updates) == 0 {
		return fmt.Errorf("no fields provided to update")
	}

	// Call the repository to perform the update
	return s.saleRepo.UpdateSaleByUserID(saleID, userID, updates)
}

