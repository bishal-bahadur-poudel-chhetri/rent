package services

import (
	"renting/internal/repositories"
	"time"
)

type DataAggregateService struct {
	saleRepo *repositories.DataAggregateRepository
}

func NewDataAggregateService(saleRepo *repositories.DataAggregateRepository) *DataAggregateService {
	return &DataAggregateService{saleRepo: saleRepo}
}

type AggregatedData struct {
	Date            string  `json:"date,omitempty"`   // Date in "YYYY-MM-DD" format
	Year            int     `json:"year,omitempty"`   // Year in "YYYY" format
	Month           string  `json:"month,omitempty"`  // Month in "YYYY-MM" format
	PendingRequests int     `json:"pending_requests"` // Count of pending requests
	TotalSales      int     `json:"total_sales"`      // Count of total sales
	FutureBookings  int     `json:"future_bookings"`  // Count of future bookings
	TotalRevenue    float64 `json:"total_revenue"`    // Total revenue
}

// GetAggregatedData returns aggregated data for a specific date, month, or year.
func (s *DataAggregateService) GetAggregatedData(date time.Time, filterType string) (AggregatedData, error) {
	var pendingRequests, totalSales, futureBookings int
	var totalRevenue float64
	var err error

	// Fetch data based on the filter type
	switch filterType {
	case "date":
		pendingRequests, err = s.saleRepo.GetPendingRequests(date, "date")
		if err != nil {
			return AggregatedData{}, err
		}
		totalSales, totalRevenue, err = s.saleRepo.GetTotalSales(date, "date")
		if err != nil {
			return AggregatedData{}, err
		}
		futureBookings, err = s.saleRepo.GetFutureBookings(date, "date")
		if err != nil {
			return AggregatedData{}, err
		}
	case "month":
		pendingRequests, err = s.saleRepo.GetPendingRequests(date, "month")
		if err != nil {
			return AggregatedData{}, err
		}
		totalSales, totalRevenue, err = s.saleRepo.GetTotalSales(date, "month")
		if err != nil {
			return AggregatedData{}, err
		}
		futureBookings, err = s.saleRepo.GetFutureBookings(date, "month")
		if err != nil {
			return AggregatedData{}, err
		}
	case "year":
		pendingRequests, err = s.saleRepo.GetPendingRequests(date, "year")
		if err != nil {
			return AggregatedData{}, err
		}
		totalSales, totalRevenue, err = s.saleRepo.GetTotalSales(date, "year")
		if err != nil {
			return AggregatedData{}, err
		}
		futureBookings, err = s.saleRepo.GetFutureBookings(date, "year")
		if err != nil {
			return AggregatedData{}, err
		}
	default:
		return AggregatedData{}, nil
	}

	// Prepare the response based on the filter type
	switch filterType {
	case "date":
		return AggregatedData{
			Date:            date.Format("2006-01-02"), // Include the date in the response
			PendingRequests: pendingRequests,
			TotalSales:      totalSales,
			FutureBookings:  futureBookings,
			TotalRevenue:    totalRevenue,
		}, nil
	case "month":
		return AggregatedData{
			Month:           date.Format("2006-01"), // Include the month in the response
			PendingRequests: pendingRequests,
			TotalSales:      totalSales,
			FutureBookings:  futureBookings,
			TotalRevenue:    totalRevenue,
		}, nil
	case "year":
		return AggregatedData{
			Year:            date.Year(), // Include the year in the response
			PendingRequests: pendingRequests,
			TotalSales:      totalSales,
			FutureBookings:  futureBookings,
			TotalRevenue:    totalRevenue,
		}, nil
	default:
		return AggregatedData{}, nil
	}
}
