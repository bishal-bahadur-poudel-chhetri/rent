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
	// Query to fetch future bookings
	futureBookingQuery := `
        SELECT date_of_delivery, return_date
        FROM sales
        WHERE vehicle_id = $1
          AND booking_date != date_of_delivery
          AND status = 'pending'
          AND EXTRACT(YEAR FROM date_of_delivery::DATE) <= EXTRACT(YEAR FROM $2::DATE)
          AND EXTRACT(MONTH FROM date_of_delivery::DATE) <= EXTRACT(MONTH FROM $2::DATE)
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

		// Filter to include only the date 2025-03-22
		if formattedDeliveryDate == "2025-03-22" {
			dates = append(dates, models.DisabledDateResponse{
				DateOfDelivery: formattedDeliveryDate,
				ReturnDate:     formattedReturnDate,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return dates, nil
}
