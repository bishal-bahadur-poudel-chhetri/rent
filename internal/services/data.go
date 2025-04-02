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
	Date            string  `json:"date,omitempty"`
	Year            int     `json:"year,omitempty"`
	Month           string  `json:"month,omitempty"`
	PendingRequests int     `json:"pending_requests"`
	TotalSales      int     `json:"total_sales"`
	FutureBookings  int     `json:"future_bookings"`
	TotalRevenue    float64 `json:"total_revenue"`
}

func (s *DataAggregateService) GetAggregatedData(date time.Time, filterType string) (AggregatedData, error) {
	var pendingRequests, totalSales, futureBookings int
	var totalRevenue float64
	var err error

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

	switch filterType {
	case "date":
		return AggregatedData{
			Date:            date.Format("2006-01-02"),
			PendingRequests: pendingRequests,
			TotalSales:      totalSales,
			FutureBookings:  futureBookings,
			TotalRevenue:    totalRevenue,
		}, nil
	case "month":
		return AggregatedData{
			Month:           date.Format("2006-01"),
			PendingRequests: pendingRequests,
			TotalSales:      totalSales,
			FutureBookings:  futureBookings,
			TotalRevenue:    totalRevenue,
		}, nil
	case "year":
		return AggregatedData{
			Year:            date.Year(),
			PendingRequests: pendingRequests,
			TotalSales:      totalSales,
			FutureBookings:  futureBookings,
			TotalRevenue:    totalRevenue,
		}, nil
	default:
		return AggregatedData{}, nil
	}
}
