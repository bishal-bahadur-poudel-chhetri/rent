package repositories

import (
	"database/sql"
	"fmt"
	"renting/internal/models"
	"time"
)

type FuturBookingRepository struct {
	db *sql.DB
}

func NewFuturBookingRepository(db *sql.DB) *FuturBookingRepository {
	return &FuturBookingRepository{db: db}
}

// GetFuturBookingsByMonth fetches bookings for a specific month and metadata for future months
func (r *FuturBookingRepository) GetFuturBookingsByMonth(year int, month time.Month, filters map[string]string) (*models.SalesWithMetadataResponse, error) {
	// Fetch bookings for the selected month
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	sales, err := r.getSalesInDateRange(startOfMonth, endOfMonth, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bookings for the month: %v", err)
	}

	// Fetch metadata (count of bookings) for the provided month and future months
	metadata, err := r.getMonthlyMetadata(startOfMonth, filters) // Pass filters to metadata query
	if err != nil {
		return nil, fmt.Errorf("failed to fetch monthly metadata: %v", err)
	}

	// Prepare the response
	response := &models.SalesWithMetadataResponse{
		Sales:    sales,
		Metadata: metadata,
	}

	return response, nil
}

// getSalesInDateRange fetches bookings within a date range and applies filters
func (r *FuturBookingRepository) getSalesInDateRange(start, end time.Time, filters map[string]string) ([]models.Sale_Future, error) {
	// Base query
	query := `
		SELECT sale_id, vehicle_id, user_id, customer_name, customer_destination, customer_phone, 
		total_amount, charge_per_day, booking_date, date_of_delivery, return_date, 
		number_of_days, remark, status, created_at, updated_at
		FROM sales
		WHERE date_of_delivery BETWEEN $1 AND $2
		AND booking_date != date_of_delivery
	`

	// Add filters to the query
	args := []interface{}{start, end}
	argCounter := 3 // Start from $3 because $1 and $2 are already used for start and end dates

	for key, value := range filters {
		switch key {
		case "status":
			query += fmt.Sprintf(" AND status = $%d", argCounter)
			args = append(args, value)
			argCounter++
		case "customer_name":
			query += fmt.Sprintf(" AND customer_name ILIKE $%d", argCounter)
			args = append(args, "%"+value+"%")
			argCounter++
		case "vehicle_id":
			query += fmt.Sprintf(" AND vehicle_id = $%d", argCounter)
			args = append(args, value)
			argCounter++
			// Add more filters as needed
		}
	}

	// Execute the query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query bookings: %v", err)
	}
	defer rows.Close()

	var sales []models.Sale_Future
	for rows.Next() {
		var sale models.Sale_Future
		err := rows.Scan(
			&sale.SaleID, &sale.VehicleID, &sale.UserID, &sale.CustomerName, &sale.Destination, &sale.CustomerPhone,
			&sale.TotalAmount, &sale.ChargePerDay, &sale.BookingDate, &sale.DateOfDelivery, &sale.ReturnDate,
			&sale.NumberOfDays, &sale.Remark, &sale.Status, &sale.CreatedAt, &sale.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sale: %v", err)
		}
		sales = append(sales, sale)
	}

	return sales, nil
}

// getMonthlyMetadata fetches the count of bookings for the provided month and future months, with filters
func (r *FuturBookingRepository) getMonthlyMetadata(startOfMonth time.Time, filters map[string]string) ([]models.MonthlyMetadata, error) {
	// Base query
	query := `
		SELECT TO_CHAR(date_of_delivery, 'YYYY-MM') AS month, COUNT(*) AS count
		FROM sales
		WHERE date_of_delivery >= $1
	`

	// Add filters to the query
	args := []interface{}{startOfMonth}
	argCounter := 2 // Start from $2 because $1 is already used for startOfMonth

	for key, value := range filters {
		switch key {
		case "status":
			query += fmt.Sprintf(" AND status = $%d", argCounter)
			args = append(args, value)
			argCounter++
		case "customer_name":
			query += fmt.Sprintf(" AND customer_name ILIKE $%d", argCounter)
			args = append(args, "%"+value+"%")
			argCounter++
		case "vehicle_id":
			query += fmt.Sprintf(" AND vehicle_id = $%d", argCounter)
			args = append(args, value)
			argCounter++
			// Add more filters as needed
		}
	}

	// Complete the query
	query += `
		GROUP BY TO_CHAR(date_of_delivery, 'YYYY-MM')
		ORDER BY month
	`

	// Execute the query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query monthly metadata: %v", err)
	}
	defer rows.Close()

	var metadata []models.MonthlyMetadata
	for rows.Next() {
		var meta models.MonthlyMetadata
		err := rows.Scan(&meta.Month, &meta.Count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan metadata: %v", err)
		}
		metadata = append(metadata, meta)
	}

	return metadata, nil
}
