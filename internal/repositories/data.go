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

// GetPendingRequests (unchanged)
func (r *DataAggregateRepository) GetPendingRequests(date time.Time, filterType string) (int, error) {
	var query string
	switch filterType {
	case "date":
		query = `
			SELECT COUNT(*) 
			FROM payments 
			WHERE payment_status = 'Pending' 
			AND verified_by_admin = false 
			AND DATE(payment_date) = $1`
	case "month":
		query = `
			SELECT COUNT(*) 
			FROM payments 
			WHERE payment_status = 'Pending' 
			AND verified_by_admin = false 
			AND EXTRACT(YEAR FROM payment_date) = $1
			AND EXTRACT(MONTH FROM payment_date) = $2`
	case "year":
		query = `
			SELECT COUNT(*) 
			FROM payments 
			WHERE payment_status = 'Pending' 
			AND verified_by_admin = false 
			AND EXTRACT(YEAR FROM payment_date) = $1`
	default:
		return 0, nil
	}

	var count int
	var err error
	switch filterType {
	case "date":
		err = r.db.QueryRow(query, date.Format("2006-01-02")).Scan(&count)
	case "month":
		err = r.db.QueryRow(query, date.Year(), date.Month()).Scan(&count)
	case "year":
		err = r.db.QueryRow(query, date.Year()).Scan(&count)
	}
	return count, err
}

// GetTotalSales returns the count of sales and total verified payments for a given date period
func (r *DataAggregateRepository) GetTotalSales(date time.Time, filterType string) (int, float64, error) {
	var query string
	switch filterType {
	case "date":
		query = `
			SELECT COUNT(DISTINCT p.sale_id), COALESCE(SUM(p.amount_paid), 0)
			FROM payments p
			WHERE p.verified_by_admin = true
			AND DATE(p.payment_date) = $1`
	case "month":
		query = `
			SELECT COUNT(DISTINCT p.sale_id), COALESCE(SUM(p.amount_paid), 0)
			FROM payments p
			WHERE p.verified_by_admin = true
			AND EXTRACT(YEAR FROM p.payment_date) = $1
			AND EXTRACT(MONTH FROM p.payment_date) = $2`
	case "year":
		query = `
			SELECT COUNT(DISTINCT p.sale_id), COALESCE(SUM(p.amount_paid), 0)
			FROM payments p
			WHERE p.verified_by_admin = true
			AND EXTRACT(YEAR FROM p.payment_date) = $1`
	default:
		return 0, 0, nil
	}

	var count int
	var totalRevenue float64
	var err error
	switch filterType {
	case "date":
		err = r.db.QueryRow(query, date.Format("2006-01-02")).Scan(&count, &totalRevenue)
	case "month":
		err = r.db.QueryRow(query, date.Year(), date.Month()).Scan(&count, &totalRevenue)
	case "year":
		err = r.db.QueryRow(query, date.Year()).Scan(&count, &totalRevenue)
	}
	return count, totalRevenue, err
}

// GetFutureBookings (unchanged)
func (r *DataAggregateRepository) GetFutureBookings(date time.Time, filterType string) (int, error) {
	var query string
	switch filterType {
	case "date":
		query = `
			SELECT COUNT(*) 
			FROM sales s
			WHERE s.status = 'pending'
			AND s.booking_date < s.date_of_delivery
			AND DATE(s.booking_date) = $1
			AND EXISTS (
				SELECT 1
				FROM payments p
				WHERE p.sale_id = s.sale_id
				AND p.payment_status = 'Completed'
			)`
	case "month":
		query = `
			SELECT COUNT(*) 
			FROM sales s
			WHERE s.status = 'pending'
			AND s.booking_date < s.date_of_delivery
			AND EXTRACT(YEAR FROM s.booking_date) = $1
			AND EXTRACT(MONTH FROM s.booking_date) = $2
			AND EXISTS (
				SELECT 1
				FROM payments p
				WHERE p.sale_id = s.sale_id
				AND p.payment_status = 'Completed'
			)`
	case "year":
		query = `
			SELECT COUNT(*) 
			FROM sales s
			WHERE s.status = 'pending'
			AND s.booking_date < s.date_of_delivery
			AND EXTRACT(YEAR FROM s.booking_date) = $1
			AND EXISTS (
				SELECT 1
				FROM payments p
				WHERE p.sale_id = s.sale_id
				AND p.payment_status = 'Completed'
			)`
	default:
		return 0, nil
	}

	var count int
	var err error
	switch filterType {
	case "date":
		err = r.db.QueryRow(query, date.Format("2006-01-02")).Scan(&count)
	case "month":
		err = r.db.QueryRow(query, date.Year(), date.Month()).Scan(&count)
	case "year":
		err = r.db.QueryRow(query, date.Year()).Scan(&count)
	}
	return count, err
}

// GetTotalAvailableCars returns the count of available cars
func (r *DataAggregateRepository) GetTotalAvailableCars() (int, error) {
	query := `
		SELECT COUNT(*)
		FROM vehicles
		WHERE is_available = true
		AND status = 'active'`

	var count int
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}
