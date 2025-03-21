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
func (r *FuturBookingRepository) getSalesInDateRange(start, end time.Time, filters map[string]string) ([]models.Sale_Future, error) {
	// Validate date range
	if start.After(end) {
		return nil, fmt.Errorf("start date must be before end date")
	}

	// Base query
	query := `
		SELECT s.sale_id, s.vehicle_id, s.user_id, s.customer_name, s.customer_destination, s.customer_phone, 
		s.total_amount, s.charge_per_day, s.booking_date, s.date_of_delivery, s.return_date, 
		s.number_of_days, s.remark, s.status, s.created_at, s.updated_at,
		p.payment_id, p.amount_paid, p.payment_date, p.payment_type, p.payment_status,v.vehicle_name,v.image_name
		FROM sales s
		LEFT JOIN payments p ON s.sale_id = p.sale_id
		LEFT JOIN vehicles v on s.vehicle_id = v.vehicle_id
		WHERE s.booking_date BETWEEN $1 AND $2
		AND date(s.booking_date) != date(s.date_of_delivery)
	`

	// Add filters to the query
	args := []interface{}{start, end}
	argCounter := 3 // Start from $3 because $1 and $2 are already used for start and end dates

	// Validate filter keys
	validFilters := map[string]bool{
		"status":        true,
		"customer_name": true,
		"vehicle_id":    true,
	}

	for key, value := range filters {
		if !validFilters[key] {
			return nil, fmt.Errorf("invalid filter key: %s", key)
		}

		switch key {
		case "status":
			query += fmt.Sprintf(" AND s.status = $%d", argCounter)
			args = append(args, value)
			argCounter++
		case "customer_name":
			query += fmt.Sprintf(" AND s.customer_name ILIKE $%d", argCounter)
			args = append(args, "%"+value+"%")
			argCounter++
		case "vehicle_id":
			query += fmt.Sprintf(" AND s.vehicle_id = $%d", argCounter)
			args = append(args, value)
			argCounter++
		}
	}

	// Execute the query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query bookings: %v", err)
	}
	defer rows.Close()

	// Map to group payments by sale_id
	salesMap := make(map[int]models.Sale_Future)

	for rows.Next() {
		var sale models.Sale_Future
		var payment models.Payment_future
		var paymentID *int
		var amountPaid *float64
		var paymentDate *time.Time
		var paymentType, paymentStatus *string

		err := rows.Scan(
			&sale.SaleID, &sale.VehicleID, &sale.UserID, &sale.CustomerName, &sale.Destination, &sale.CustomerPhone,
			&sale.TotalAmount, &sale.ChargePerDay, &sale.BookingDate, &sale.DateOfDelivery, &sale.ReturnDate,
			&sale.NumberOfDays, &sale.Remark, &sale.Status, &sale.CreatedAt, &sale.UpdatedAt,
			&paymentID, &amountPaid, &paymentDate, &paymentType, &paymentStatus, &sale.VehicleName, &sale.VehicleImage,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sale (SaleID: %d): %v", sale.SaleID, err)
		}

		// If payment data exists, populate the Payment_future struct
		if paymentID != nil {
			payment = models.Payment_future{
				PaymentID:     *paymentID,
				AmountPaid:    *amountPaid,
				PaymentDate:   *paymentDate,
				PaymentType:   *paymentType,
				PaymentStatus: *paymentStatus,
			}
		}

		// Check if the sale already exists in the map
		if existingSale, ok := salesMap[sale.SaleID]; ok {
			// Append the payment to the existing sale
			existingSale.Payment = append(existingSale.Payment, payment)
			salesMap[sale.SaleID] = existingSale
		} else {
			// Create a new sale entry in the map
			sale.Payment = []models.Payment_future{}
			if paymentID != nil {
				sale.Payment = append(sale.Payment, payment)
			}
			salesMap[sale.SaleID] = sale
		}
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	// Convert the map to a slice
	var sales []models.Sale_Future
	for _, sale := range salesMap {
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
