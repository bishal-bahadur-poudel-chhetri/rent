package repositories

import (
	"database/sql"
	"fmt"
	"renting/internal/models"
	"time"
)

type DisableDateRepository interface {
	GetDisabledDates(vehicleID int, excludeSaleID *int) (*models.DisableDateResponse, error)
}

type disableDateRepository struct {
	db *sql.DB
}

func NewDisableDateRepository(db *sql.DB) DisableDateRepository {
	return &disableDateRepository{db: db}
}

func (r *disableDateRepository) GetDisabledDates(vehicleID int, excludeSaleID *int) (*models.DisableDateResponse, error) {
	// Query for ALL active rentals
	activeRentalsQuery := `
                SELECT date_of_delivery, return_date
                FROM sales
                WHERE vehicle_id = $1
                        AND status = 'active'
                        AND ($2::integer IS NULL OR sale_id != $2::integer)
        `

	// Query for future pending bookings
	futureBookingsQuery := `
                SELECT date_of_delivery, return_date
                FROM sales
                WHERE vehicle_id = $1
                        AND status = 'pending'
                        AND ($2::integer IS NULL OR sale_id != $2::integer)
                        AND return_date >= CURRENT_DATE
        `

	var excludeParam interface{}
	if excludeSaleID != nil {
		excludeParam = *excludeSaleID
	} else {
		excludeParam = nil
	}

	// Get active rentals
	activeRentals, err := r.fetchDates(activeRentalsQuery, vehicleID, excludeParam)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch active rentals: %v", err)
	}

	// Get future bookings
	futureBookings, err := r.fetchDates(futureBookingsQuery, vehicleID, excludeParam)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch future bookings: %v", err)
	}

	return &models.DisableDateResponse{
		ActiveRentals:  activeRentals,
		FutureBookings: futureBookings,
	}, nil
}

func (r *disableDateRepository) fetchDates(query string, vehicleID int, excludeSaleID interface{}) ([]models.DisabledDateResponse, error) {
	rows, err := r.db.Query(query, vehicleID, excludeSaleID)
	if err != nil {
		return nil, fmt.Errorf("failed to query dates: %v", err)
	}
	defer rows.Close()

	var dates []models.DisabledDateResponse
	for rows.Next() {
		var deliveryDate, returnDate time.Time
		if err := rows.Scan(&deliveryDate, &returnDate); err != nil {
			return nil, fmt.Errorf("failed to scan date: %v", err)
		}

		dates = append(dates, models.DisabledDateResponse{
			DateOfDelivery: deliveryDate.Format("2006-01-02"),
			ReturnDate:     returnDate.Format("2006-01-02"),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return dates, nil
}

