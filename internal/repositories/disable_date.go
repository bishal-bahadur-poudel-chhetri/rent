package repositories

import (
	"database/sql"
	"fmt"
	"renting/internal/models"
	"time"
)

type DisableDateRepository struct {
	db *sql.DB
}

func NewDisableDateRepository(db *sql.DB) *DisableDateRepository {
	return &DisableDateRepository{db: db}
}

// GetDisabledDates fetches disabled dates for a specific vehicle and date range
func (r *DisableDateRepository) GetDisabledDates(vehicleID int, dateOfDelivery time.Time) (*models.DisableDateResponse, error) {
	// Query to fetch future bookings that overlap with the target month
	futureBookingQuery := `
        WITH target_month AS (
            SELECT 
                DATE_TRUNC('MONTH', $2::DATE) AS month_start,
                DATE_TRUNC('MONTH', $2::DATE) + INTERVAL '1 MONTH' - INTERVAL '1 DAY' AS month_end
        )
        SELECT date_of_delivery, return_date
        FROM sales, target_month
        WHERE vehicle_id = $1
          AND booking_date != date_of_delivery
          AND status = 'pending'
          AND (
              -- Case 1: date_of_delivery is in the target month
              date_of_delivery BETWEEN target_month.month_start AND target_month.month_end
              -- Case 2: return_date is in the target month
              OR return_date BETWEEN target_month.month_start AND target_month.month_end
              -- Case 3: Booking spans the target month
              OR (
                  date_of_delivery < target_month.month_start 
                  AND return_date > target_month.month_end
              )
          )
    `

	// Query to fetch today's sales
	todaySalesQuery := `
        SELECT date_of_delivery, return_date
        FROM sales
        WHERE vehicle_id = $1
          AND booking_date = date_of_delivery
          AND status = 'active'
          AND date_of_delivery = $2
    `

	// Fetch future bookings
	futureBookings, err := r.fetchDates(futureBookingQuery, vehicleID, dateOfDelivery)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch future bookings: %v", err)
	}

	// Fetch today's sales
	todaySales, err := r.fetchDates(todaySalesQuery, vehicleID, dateOfDelivery)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch today's sales: %v", err)
	}

	// Combine the results
	response := &models.DisableDateResponse{
		FutureBooking: futureBookings,
		TodaySales:    todaySales,
	}

	return response, nil
}

// fetchDates is a helper function to fetch dates based on a query
func (r *DisableDateRepository) fetchDates(query string, vehicleID int, dateOfDelivery time.Time) ([]models.DisabledDateResponse, error) {
	rows, err := r.db.Query(query, vehicleID, dateOfDelivery)
	if err != nil {
		return nil, fmt.Errorf("failed to query dates: %v", err)
	}
	defer rows.Close()

	var dates []models.DisabledDateResponse
	for rows.Next() {
		var deliveryDate time.Time
		var returnDate time.Time
		err := rows.Scan(&deliveryDate, &returnDate)
		if err != nil {
			return nil, fmt.Errorf("failed to scan date: %v", err)
		}

		// Format the dates to YYYY-MM-DD
		formattedDeliveryDate := deliveryDate.Format("2006-01-02")
		formattedReturnDate := returnDate.Format("2006-01-02")

		// Append the result without filtering for a specific date
		dates = append(dates, models.DisabledDateResponse{
			DateOfDelivery: formattedDeliveryDate,
			ReturnDate:     formattedReturnDate,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return dates, nil
}
