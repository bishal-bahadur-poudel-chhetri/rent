package repositories

import (
	"context"
	"database/sql"
	"renting/internal/models"
	"strconv"
	"strings"
)

type StatementRepository interface {
	GetOutstandingStatements(ctx context.Context, filters map[string]string, offset, limit int) ([]*models.Statement, error)
}

type statementRepository struct {
	db *sql.DB
}

func NewStatementRepository(db *sql.DB) StatementRepository {
	return &statementRepository{db: db}
}

func (r *statementRepository) GetOutstandingStatements(ctx context.Context, filters map[string]string, offset, limit int) ([]*models.Statement, error) {
	baseQuery := `
        SELECT 
            sale_id AS statement_id,
            vehicle_id,
            user_id,
            customer_name,
            customer_destination,
            customer_phone,
            total_amount,
            charge_per_day,
            booking_date,
            date_of_delivery,
            return_date,
            number_of_days,
            remark,
            status,
            created_at,
            updated_at,
            actual_date_of_delivery,
            actual_date_of_return,
            payment_status,
            other_charges,
            modified_by,
            outstanding_balance,
            vehicle_name
        FROM sales_statement_view`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if bookingDate, ok := filters["booking_date"]; ok && bookingDate != "" {
		conditions = append(conditions, "booking_date = $"+strconv.Itoa(argIndex))
		args = append(args, bookingDate)
		argIndex++
	}
	if status, ok := filters["status"]; ok && status != "" {
		conditions = append(conditions, "status = $"+strconv.Itoa(argIndex))
		args = append(args, status)
		argIndex++
	}
	if paymentStatus, ok := filters["payment_status"]; ok && paymentStatus != "" {
		conditions = append(conditions, "payment_status = $"+strconv.Itoa(argIndex))
		args = append(args, paymentStatus)
		argIndex++
	}
	if vehicleName, ok := filters["vehicle_name"]; ok && vehicleName != "" {
		conditions = append(conditions, "vehicle_name ILIKE $"+strconv.Itoa(argIndex))
		args = append(args, "%"+vehicleName+"%")
		argIndex++
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	baseQuery += " ORDER BY booking_date DESC, statement_id DESC"
	baseQuery += " OFFSET $" + strconv.Itoa(argIndex) + " LIMIT $" + strconv.Itoa(argIndex+1)
	args = append(args, offset, limit)

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statements []*models.Statement
	for rows.Next() {
		var s models.Statement
		err := rows.Scan(
			&s.StatementID,
			&s.VehicleID,
			&s.UserID,
			&s.CustomerName,
			&s.CustomerDestination,
			&s.CustomerPhone,
			&s.TotalAmount,
			&s.ChargePerDay,
			&s.BookingDate,
			&s.DateOfDelivery,
			&s.ReturnDate,
			&s.NumberOfDays,
			&s.Remark,
			&s.Status,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.ActualDateOfDelivery,
			&s.ActualDateOfReturn,
			&s.PaymentStatus,
			&s.OtherCharges,
			&s.ModifiedBy,
			&s.OutstandingBalance,
			&s.VehicleName,
		)
		if err != nil {
			return nil, err
		}
		statements = append(statements, &s)
	}

	return statements, nil
}
