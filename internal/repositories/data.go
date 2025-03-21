package repositories

import (
	"database/sql"
	"time"
)

type DataAggregateRepository struct {
	db *sql.DB
}

func NewDataAggregateRepository(db *sql.DB) *DataAggregateRepository {
	return &DataAggregateRepository{db: db}
}

// GetPendingRequests returns the count of pending payment requests for a specific date, year, or month.
func (r *DataAggregateRepository) GetPendingRequests(date time.Time) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM payments 
		WHERE payment_status = 'Pending' 
		AND verified_by_admin = false 
		AND DATE(payment_date) = $1`
	var count int
	err := r.db.QueryRow(query, date.Format("2006-01-02")).Scan(&count)
	return count, err
}

// GetTotalSales returns the total number of completed sales and the total revenue for a specific date, year, or month.
func (r *DataAggregateRepository) GetTotalSales(date time.Time) (int, float64, error) {
	query := `
		SELECT COUNT(*), COALESCE(SUM(total_amount), 0) 
		FROM sales 
		WHERE status = 'completed' 
		AND DATE(booking_date) = $1`
	var count int
	var totalRevenue float64
	err := r.db.QueryRow(query, date.Format("2006-01-02")).Scan(&count, &totalRevenue)
	return count, totalRevenue, err
}

// GetFutureBookings returns the count of future bookings for a specific date, year, or month.
func (r *DataAggregateRepository) GetFutureBookings(date time.Time) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM sales 
		WHERE booking_date < date_of_delivery 
		AND DATE(booking_date) = $1`
	var count int
	err := r.db.QueryRow(query, date.Format("2006-01-02")).Scan(&count)
	return count, err
}

// GetPendingRequestsByYear returns the count of pending payment requests for a specific year.
func (r *DataAggregateRepository) GetPendingRequestsByYear(year int) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM payments 
		WHERE payment_status = 'Pending' 
		AND verified_by_admin = false 
		AND EXTRACT(YEAR FROM payment_date) = $1`
	var count int
	err := r.db.QueryRow(query, year).Scan(&count)
	return count, err
}

// GetTotalSalesByYear returns the total number of completed sales and the total revenue for a specific year.
func (r *DataAggregateRepository) GetTotalSalesByYear(year int) (int, float64, error) {
	query := `
		SELECT COUNT(*), COALESCE(SUM(total_amount), 0) 
		FROM sales 
		WHERE status = 'completed' 
		AND EXTRACT(YEAR FROM booking_date) = $1`
	var count int
	var totalRevenue float64
	err := r.db.QueryRow(query, year).Scan(&count, &totalRevenue)
	return count, totalRevenue, err
}

// GetFutureBookingsByYear returns the count of future bookings for a specific year.
func (r *DataAggregateRepository) GetFutureBookingsByYear(year int) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM sales 
		WHERE booking_date > date_of_delivery 
		AND EXTRACT(YEAR FROM booking_date) = $1`
	var count int
	err := r.db.QueryRow(query, year).Scan(&count)
	return count, err
}
