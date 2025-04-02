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
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	sales, err := r.getSalesInDateRange(startOfMonth, endOfMonth, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch bookings for the month: %w", err)
	}

	metadata, err := r.getMonthlyMetadata(startOfMonth, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch monthly metadata: %w", err)
	}

	return &models.SalesWithMetadataResponse{
		Sales:    sales,
		Metadata: metadata,
	}, nil
}

func (r *FuturBookingRepository) getSalesInDateRange(start, end time.Time, filters map[string]string) ([]models.Sale_Future, error) {
	if start.After(end) {
		return nil, fmt.Errorf("start date must be before end date")
	}

	query, args := buildSalesQuery(start, end, filters)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query bookings: %w", err)
	}
	defer rows.Close()

	salesMap := make(map[int]models.Sale_Future)

	for rows.Next() {
		sale, payment, err := scanSaleRow(rows)
		if err != nil {
			return nil, err
		}

		if existingSale, exists := salesMap[sale.SaleID]; exists {
			if payment != nil {
				existingSale.Payment = append(existingSale.Payment, *payment)
			}
			salesMap[sale.SaleID] = existingSale
		} else {
			sale.Payment = []models.Payment_future{}
			if payment != nil {
				sale.Payment = append(sale.Payment, *payment)
			}
			salesMap[sale.SaleID] = sale
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return mapToSalesSlice(salesMap), nil
}

func buildSalesQuery(start, end time.Time, filters map[string]string) (string, []interface{}) {
	baseQuery := `
		SELECT s.sale_id, s.vehicle_id, s.user_id, s.customer_name, s.customer_destination, s.customer_phone, 
		s.total_amount, s.charge_per_day, s.booking_date, s.date_of_delivery, s.return_date, 
		s.number_of_days, s.remark, s.status, s.created_at, s.updated_at,
		p.payment_id, p.amount_paid, p.payment_date, p.payment_type, p.payment_status, v.vehicle_name, v.image_name
		FROM sales s
		LEFT JOIN payments p ON s.sale_id = p.sale_id
		LEFT JOIN vehicles v ON s.vehicle_id = v.vehicle_id
		WHERE s.date_of_delivery BETWEEN $1 AND $2
		AND date(s.booking_date) != date(s.date_of_delivery)
	`

	args := []interface{}{start, end}
	query := baseQuery
	argCounter := 3

	validFilters := map[string]bool{
		"status":        true,
		"customer_name": true,
		"vehicle_id":    true,
	}

	for key, value := range filters {
		if !validFilters[key] {
			continue
		}

		switch key {
		case "status":
			query += fmt.Sprintf(" AND s.status = $%d", argCounter)
			args = append(args, value)
		case "customer_name":
			query += fmt.Sprintf(" AND s.customer_name ILIKE $%d", argCounter)
			args = append(args, "%"+value+"%")
		case "vehicle_id":
			query += fmt.Sprintf(" AND s.vehicle_id = $%d", argCounter)
			args = append(args, value)
		}
		argCounter++
	}

	return query, args
}

func scanSaleRow(rows *sql.Rows) (models.Sale_Future, *models.Payment_future, error) {
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
		return models.Sale_Future{}, nil, fmt.Errorf("failed to scan sale: %w", err)
	}

	if paymentID != nil {
		payment = models.Payment_future{
			PaymentID:     *paymentID,
			AmountPaid:    *amountPaid,
			PaymentDate:   *paymentDate,
			PaymentType:   *paymentType,
			PaymentStatus: *paymentStatus,
		}
		return sale, &payment, nil
	}

	return sale, nil, nil
}

func mapToSalesSlice(salesMap map[int]models.Sale_Future) []models.Sale_Future {
	sales := make([]models.Sale_Future, 0, len(salesMap))
	for _, sale := range salesMap {
		sales = append(sales, sale)
	}
	return sales
}

func (r *FuturBookingRepository) getMonthlyMetadata(startOfMonth time.Time, filters map[string]string) ([]models.MonthlyMetadata, error) {
	query, args := buildMetadataQuery(startOfMonth, filters)
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query monthly metadata: %w", err)
	}
	defer rows.Close()

	var metadata []models.MonthlyMetadata
	for rows.Next() {
		var meta models.MonthlyMetadata
		if err := rows.Scan(&meta.Month, &meta.Count); err != nil {
			return nil, fmt.Errorf("failed to scan metadata: %w", err)
		}
		metadata = append(metadata, meta)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return metadata, nil
}

func buildMetadataQuery(startOfMonth time.Time, filters map[string]string) (string, []interface{}) {
	baseQuery := `
		SELECT TO_CHAR(date_of_delivery, 'YYYY-MM') AS month, COUNT(*) AS count
		FROM sales
		WHERE date_of_delivery >= $1
	`

	args := []interface{}{startOfMonth.Format("2006-01-02")}
	query := baseQuery
	argCounter := 2

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
		}
	}

	query += `
		GROUP BY TO_CHAR(date_of_delivery, 'YYYY-MM')
		ORDER BY month
	`

	return query, args
}

func (r *FuturBookingRepository) FutureBookingCancellation(saleID int) error {
	query := `
		UPDATE sales
		SET status = 'cancelled', updated_at = $1
		WHERE sale_id = $2
		AND status NOT IN ('cancelled', 'completed')
	`

	result, err := r.db.Exec(query, time.Now().UTC(), saleID)
	if err != nil {
		return fmt.Errorf("failed to cancel booking with sale_id %d: %w", saleID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no booking found with sale_id %d or it was already cancelled/completed", saleID)
	}

	return nil
}
