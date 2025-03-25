package repositories

import (
	"database/sql"
	"fmt"
	"renting/internal/models"

	"strings"

	"github.com/lib/pq"
)

type StatementRepository struct {
	db *sql.DB
}

func NewStatementRepository(db *sql.DB) *StatementRepository {
	return &StatementRepository{db: db}
}

func (r *StatementRepository) GetStatements(filter models.StatementFilter) ([]models.Statement, error) {
	query := `
		SELECT 
			sale_id, vehicle_id, vehicle_name, vehicle_registration_number,
			vehicle_image, user_id, user_name, customer_name, customer_phone,
			booking_date, date_of_delivery, return_date, total_amount,
			charge_per_day, sale_status, amount_paid, balance,
			has_damage, has_delay, sale_images, sale_videos,
			created_at, updated_at
		FROM statement_view
		WHERE 1=1
	`
	fmt.Print("hi")

	args := []interface{}{}
	argPos := 1

	if filter.DateFrom != nil {
		query += fmt.Sprintf(" AND booking_date >= $%d", argPos)
		args = append(args, *filter.DateFrom)
		argPos++
	}

	if filter.DateTo != nil {
		query += fmt.Sprintf(" AND booking_date <= $%d", argPos)
		args = append(args, *filter.DateTo)
		argPos++
	}

	if filter.VehicleName != "" {
		query += fmt.Sprintf(" AND LOWER(vehicle_name) LIKE LOWER($%d)", argPos)
		args = append(args, "%"+filter.VehicleName+"%")
		argPos++
	}

	if filter.CustomerName != "" {
		query += fmt.Sprintf(" AND LOWER(customer_name) LIKE LOWER($%d)", argPos)
		args = append(args, "%"+filter.CustomerName+"%")
		argPos++
	}

	if filter.CustomerPhone != "" {
		query += fmt.Sprintf(" AND customer_phone = $%d", argPos)
		args = append(args, filter.CustomerPhone)
		argPos++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND sale_status = $%d", argPos)
		args = append(args, filter.Status)
		argPos++
	}

	if filter.PaymentStatus != "" {
		switch filter.PaymentStatus {
		case "paid":
			query += " AND amount_paid >= total_amount"
		case "partial":
			query += " AND amount_paid > 0 AND amount_paid < total_amount"
		case "unpaid":
			query += " AND amount_paid = 0"
		}
	}

	if filter.HasDamage {
		query += " AND has_damage = true"
	}

	if filter.HasDelay {
		query += " AND has_delay = true"
	}

	validSortFields := map[string]bool{
		"booking_date":  true,
		"return_date":   true,
		"total_amount":  true,
		"customer_name": true,
		"vehicle_name":  true,
		"created_at":    true,
	}

	sortBy := "booking_date"
	if validSortFields[filter.SortBy] {
		sortBy = filter.SortBy
	}

	sortOrder := "DESC"
	if strings.ToUpper(filter.SortOrder) == "ASC" {
		sortOrder = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	// Add pagination
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", filter.Limit)
		if filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", filter.Offset)
		}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var statements []models.Statement
	for rows.Next() {
		var s models.Statement
		err := rows.Scan(
			&s.SaleID, &s.VehicleID, &s.VehicleName, &s.VehicleRegistrationNumber,
			&s.VehicleImage, &s.UserID, &s.UserName, &s.CustomerName, &s.CustomerPhone,
			&s.BookingDate, &s.DateOfDelivery, &s.ReturnDate, &s.TotalAmount,
			&s.ChargePerDay, &s.SaleStatus, &s.AmountPaid, &s.Balance,
			&s.HasDamage, &s.HasDelay, pq.Array(&s.SaleImages), pq.Array(&s.SaleVideos),
			&s.CreatedAt, &s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		statements = append(statements, s)
	}

	return statements, nil
}

func (r *StatementRepository) GetTotalCount(filter models.StatementFilter) (int, error) {
	query := `
	SELECT COUNT(*) 
	FROM statement_view
	WHERE 1=1
	`

	args := []interface{}{}
	argPos := 1

	if filter.DateFrom != nil {
		query += fmt.Sprintf(" AND booking_date >= $%d", argPos)
		args = append(args, *filter.DateFrom)
		argPos++
	}

	if filter.DateTo != nil {
		query += fmt.Sprintf(" AND booking_date <= $%d", argPos)
		args = append(args, *filter.DateTo)
		argPos++
	}

	if filter.VehicleName != "" {
		query += fmt.Sprintf(" AND LOWER(vehicle_name) LIKE LOWER($%d)", argPos)
		args = append(args, "%"+filter.VehicleName+"%")
		argPos++
	}

	if filter.CustomerName != "" {
		query += fmt.Sprintf(" AND LOWER(customer_name) LIKE LOWER($%d)", argPos)
		args = append(args, "%"+filter.CustomerName+"%")
		argPos++
	}

	if filter.CustomerPhone != "" {
		query += fmt.Sprintf(" AND customer_phone = $%d", argPos)
		args = append(args, filter.CustomerPhone)
		argPos++
	}

	if filter.Status != "" {
		query += fmt.Sprintf(" AND sale_status = $%d", argPos)
		args = append(args, filter.Status)
		argPos++
	}

	if filter.PaymentStatus != "" {
		switch filter.PaymentStatus {
		case "paid":
			query += " AND amount_paid >= total_amount"
		case "partial":
			query += " AND amount_paid > 0 AND amount_paid < total_amount"
		case "unpaid":
			query += " AND amount_paid = 0"
		}
	}

	if filter.HasDamage {
		query += " AND has_damage = true"
	}

	if filter.HasDelay {
		query += " AND has_delay = true"
	}

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error getting total count: %v", err)
	}

	return count, nil
}
