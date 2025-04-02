package repositories

import (
	"database/sql"
	"fmt"
	"time"
)

type RevenueRepository struct {
	db *sql.DB
}

func NewRevenueRepository(db *sql.DB) *RevenueRepository {
	return &RevenueRepository{db: db}
}

// RevenueFilter defines filter options
type RevenueFilter struct {
	Period string    // "daily", "monthly", "yearly"
	Date   time.Time // Specific date to filter (optional)
}

// GetTotalRevenue calculates total revenue based on the filter
func (r *RevenueRepository) GetTotalRevenue(filter RevenueFilter) (float64, error) {
	var query string
	var args []interface{}
	today := time.Now().Truncate(24 * time.Hour)

	switch filter.Period {
	case "daily":
		// Default to today if no date provided
		date := filter.Date
		if filter.Date.IsZero() {
			date = today
		}
		// Sum daily_amount for sales active on the given date
		query = `
            SELECT COALESCE(SUM(daily_amount), 0)
            FROM revenue_recognition
            WHERE $1 BETWEEN start_date AND end_date`
		args = []interface{}{date}
	case "monthly":
		// Default to current month
		date := filter.Date
		if filter.Date.IsZero() {
			date = today
		}
		monthStart := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
		monthEnd := monthStart.AddDate(0, 1, -1)
		// Sum prorated daily_amount for sales active within the month
		query = `
            SELECT COALESCE(SUM(daily_amount * (
                LEAST(end_date, $2) - GREATEST(start_date, $1) + 1
            )), 0)
            FROM revenue_recognition
            WHERE start_date <= $2 AND end_date >= $1`
		args = []interface{}{monthStart, monthEnd}
	case "yearly":
		date := filter.Date
		if filter.Date.IsZero() {
			date = today
		}
		yearStart := time.Date(date.Year(), 1, 1, 0, 0, 0, 0, date.Location())
		yearEnd := time.Date(date.Year(), 12, 31, 0, 0, 0, 0, date.Location())
		query = `
			SELECT COALESCE(SUM(daily_amount * (
				LEAST(end_date, $2) - GREATEST(start_date, $1) + 1
			)), 0)
			FROM revenue_recognition
			WHERE start_date <= $2 AND end_date >= $1`
		args = []interface{}{yearStart, yearEnd}
	default:
		return 0, fmt.Errorf("invalid period: %s", filter.Period)
	}

	var total float64
	err := r.db.QueryRow(query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to query revenue: %v", err)
	}
	return total, nil
}
