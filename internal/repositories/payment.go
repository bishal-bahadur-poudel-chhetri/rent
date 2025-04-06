package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"renting/internal/models"
	"time"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

type SaleFilter struct {
	SaleID        *int
	PaymentStatus *string
	StartDate     *time.Time
	EndDate       *time.Time
	CustomerName  *string
	SaleStatus    *string
	VerifiedBy    *string
}

// GetPaymentsWithSales (unchanged)
func (r *PaymentRepository) GetPaymentsWithSales(filter SaleFilter, limit int, offset int) ([]models.PaymentWithSale, error) {
	query := `
        SELECT 
            p.payment_id,
            p.payment_type,
            p.amount_paid,
            p.payment_date,
            p.payment_status,
            p.verified_by_admin,
            p.remark,
            p.sale_type,
            p.created_at,
            p.updated_at,
            p.user_id AS payment_user_id,
            pu.username AS payment_username,
            s.sale_id, 
            s.vehicle_id, 
            s.user_id AS sale_user_id,
            s.customer_name, 
            s.customer_destination, 
            s.customer_phone,
            s.total_amount, 
            s.charge_per_day, 
            s.booking_date, 
            s.date_of_delivery, 
            s.return_date, 
            s.number_of_days, 
            s.remark, 
			s.actual_date_of_delivery,
			s.actual_date_of_return,
			s.payment_status,
            s.status, 
            s.created_at AS sale_created_at, 
            s.updated_at AS sale_updated_at,
            v.vehicle_id,
            v.vehicle_type_id,
            v.vehicle_name,
            v.vehicle_model,
            v.vehicle_registration_number,
            v.is_available,
            v.image_name,
            v.status
        FROM payments p
        LEFT JOIN sales s ON p.sale_id = s.sale_id
        LEFT JOIN vehicles v ON s.vehicle_id = v.vehicle_id
        LEFT JOIN users pu ON p.user_id = pu.ID
        LEFT JOIN users su ON s.user_id = su.ID
        WHERE 1=1
    `

	args := []interface{}{}
	argCounter := 1

	if filter.SaleID != nil {
		query += fmt.Sprintf(" AND p.payment_id = $%d", argCounter)
		args = append(args, *filter.SaleID)
		argCounter++
	}

	if filter.PaymentStatus != nil {
		query += fmt.Sprintf(" AND p.payment_status = $%d", argCounter)
		args = append(args, *filter.PaymentStatus)
		argCounter++
	}

	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND p.payment_date >= $%d", argCounter)
		args = append(args, *filter.StartDate)
		argCounter++
	}

	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND p.payment_date <= $%d", argCounter)
		args = append(args, *filter.EndDate)
		argCounter++
	}

	if filter.CustomerName != nil {
		query += fmt.Sprintf(" AND s.customer_name ILIKE $%d", argCounter)
		args = append(args, "%"+*filter.CustomerName+"%")
		argCounter++
	}

	if filter.SaleStatus != nil {
		query += fmt.Sprintf(" AND s.status = $%d", argCounter)
		args = append(args, *filter.SaleStatus)
		argCounter++
	}

	if filter.VerifiedBy != nil {
		query += fmt.Sprintf(" AND s.user_id = $%d", argCounter)
		args = append(args, *filter.VerifiedBy)
		argCounter++
	}

	query += " ORDER BY p.payment_date DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.PaymentWithSale
	for rows.Next() {
		var payment models.PaymentWithSale
		err := rows.Scan(
			&payment.PaymentID,
			&payment.PaymentType,
			&payment.AmountPaid,
			&payment.PaymentDate,
			&payment.PaymentStatus,
			&payment.VerifiedByAdmin,
			&payment.Remark,
			&payment.SaleType,
			&payment.CreatedAt,
			&payment.UpdatedAt,
			&payment.PaymentUserID,
			&payment.PaymentUsername,
			&payment.Sale.SaleID,
			&payment.Sale.VehicleID,
			&payment.Sale.UserID,
			&payment.Sale.CustomerName,
			&payment.Sale.Destination,
			&payment.Sale.CustomerPhone,
			&payment.Sale.TotalAmount,
			&payment.Sale.ChargePerDay,
			&payment.Sale.BookingDate,
			&payment.Sale.DateOfDelivery,
			&payment.Sale.ReturnDate,
			&payment.Sale.NumberOfDays,
			&payment.Sale.Remark,
			&payment.Sale.ActualDateofDelivery,
			&payment.Sale.ActualReturnDate,
			&payment.Sale.PaymentStatus,
			&payment.Sale.Status,
			&payment.Sale.CreatedAt,
			&payment.Sale.UpdatedAt,
			&payment.Sale.Vehicle.VehicleID,
			&payment.Sale.Vehicle.VehicleTypeID,
			&payment.Sale.Vehicle.VehicleName,
			&payment.Sale.Vehicle.VehicleModel,
			&payment.Sale.Vehicle.VehicleRegistrationNumber,
			&payment.Sale.Vehicle.IsAvailable,
			&payment.Sale.Vehicle.SalesImage,
			&payment.Sale.Vehicle.Status,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

// UpdatePayment updates payment_type and amount_paid for a payment
func (r *PaymentRepository) UpdatePayment(paymentID int, userID int, paymentType string, amountPaid float64) error {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check current payment status
	var currentStatus string
	statusQuery := `
		SELECT payment_status
		FROM payments
		WHERE payment_id = $1
	`
	err = tx.QueryRow(statusQuery, paymentID).Scan(&currentStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("payment not found")
		}
		return err
	}

	// Prevent update if payment is already Completed
	if currentStatus == "Completed" {
		return errors.New("cannot update a completed payment")
	}

	// Determine new status based on user role
	isAdmin, err := r.isAdmin(userID)
	if err != nil {
		return err
	}
	newStatus := "Pending"
	if isAdmin {
		newStatus = "Completed"
	}

	// Update payment
	query := `
		UPDATE payments
		SET 
			payment_type = $1,
			amount_paid = $2,
			payment_status = $3,
			updated_at = $4,
			user_id = $5
		WHERE payment_id = $6
	`
	result, err := tx.Exec(query,
		paymentType,
		amountPaid,
		newStatus,
		time.Now(),
		userID,
		paymentID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("payment not found")
	}

	return tx.Commit()
}

// InsertPayment adds a new payment
func (r *PaymentRepository) InsertPayment(saleID int, paymentType string, amountPaid float64, remark string) (int, error) {
	paymentStatus := "Pending"

	query := `
        INSERT INTO payments (
            sale_id,
            payment_type,
            amount_paid,
            payment_date,
            payment_status,
            created_at,
            updated_at,
            remark
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING payment_id
    `

	var paymentID int
	err := r.db.QueryRow(query,
		saleID,
		paymentType,
		amountPaid,
		time.Now(),
		paymentStatus,
		time.Now(),
		time.Now(),
		remark,
	).Scan(&paymentID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert payment: %w", err)
	}

	return paymentID, nil
}

// isAdmin checks if the user is an admin
func (r *PaymentRepository) isAdmin(userID int) (bool, error) {
	query := `
		SELECT is_admin
		FROM users
		WHERE id = $1
	`
	var isAdmin bool
	err := r.db.QueryRow(query, userID).Scan(&isAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("user not found")
		}
		return false, err
	}
	return isAdmin, nil
}

