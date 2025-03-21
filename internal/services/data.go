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

func (s *DataAggregateService) GetAggregatedData(date time.Time) (AggregatedData, error) {
	pendingRequests, err := s.saleRepo.GetPendingRequests(date)
	if err != nil {
		return AggregatedData{}, err
	}

	totalSales, totalRevenue, err := s.saleRepo.GetTotalSales(date)
	if err != nil {
		return AggregatedData{}, err
	}

	futureBookings, err := s.saleRepo.GetFutureBookings(date)
	if err != nil {
		return AggregatedData{}, err
	}

	return AggregatedData{
		Date:            date.Format("2006-01-02"), // Include the date in the response
		PendingRequests: pendingRequests,
		TotalSales:      totalSales,
		FutureBookings:  futureBookings,
		TotalRevenue:    totalRevenue,
	}, nil
}

func (s *DataAggregateService) GetAggregatedDataByYear(year int) (AggregatedData, error) {
	pendingRequests, err := s.saleRepo.GetPendingRequestsByYear(year)
	if err != nil {
		return AggregatedData{}, err
	}

	totalSales, totalRevenue, err := s.saleRepo.GetTotalSalesByYear(year)
	if err != nil {
		return AggregatedData{}, err
	}

	futureBookings, err := s.saleRepo.GetFutureBookingsByYear(year)
	if err != nil {
		return AggregatedData{}, err
	}

	return AggregatedData{
		Year:            year, // Include the year in the response
		PendingRequests: pendingRequests,
		TotalSales:      totalSales,
		FutureBookings:  futureBookings,
		TotalRevenue:    totalRevenue,
	}, nil
}

func (s *DataAggregateService) GetAggregatedDataByMonth(month time.Time) (AggregatedData, error) {
	pendingRequests, err := s.saleRepo.GetPendingRequests(month)
	if err != nil {
		return AggregatedData{}, err
	}

	totalSales, totalRevenue, err := s.saleRepo.GetTotalSales(month)
	if err != nil {
		return AggregatedData{}, err
	}

	futureBookings, err := s.saleRepo.GetFutureBookings(month)
	if err != nil {
		return AggregatedData{}, err
	}

	return AggregatedData{
		Month:           month.Format("2006-01"), // Include the month in the response
		PendingRequests: pendingRequests,
		TotalSales:      totalSales,
		FutureBookings:  futureBookings,
		TotalRevenue:    totalRevenue,
	}, nil
}
