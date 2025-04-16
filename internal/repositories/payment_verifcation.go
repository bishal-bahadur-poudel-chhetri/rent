package repositories

import (
	"database/sql"
	"errors"
	"time"
)

// PaymentVerificationRepository handles database operations for payment verification
type PaymentVerificationRepository struct {
	db *sql.DB
}

// NewPaymentVerificationRepository creates a new instance of PaymentVerificationRepository
func NewPaymentVerificationRepository(db *sql.DB) *PaymentVerificationRepository {
	return &PaymentVerificationRepository{db: db}
}

// isAdmin checks if the user is an admin
func (r *PaymentVerificationRepository) isAdmin(userID int) (bool, error) {
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

// VerifyPayment updates payment status and related sale records
func (r *PaymentVerificationRepository) VerifyPayment(paymentID int, status string, userID int, saleID int, remark string) error {
	// 1. Admin check
	isAdmin, err := r.isAdmin(userID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errors.New("only admin users can verify payments")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 2. First get the payment amount being verified
	var paymentAmount float64
	err = tx.QueryRow("SELECT amount_paid FROM payments WHERE payment_id = $1", paymentID).Scan(&paymentAmount)
	if err != nil {
		return err
	}

	// 3. Update the payment status
	updateQuery := `
        UPDATE payments
        SET payment_status = $1, 
            verified_by_admin = $2, 
            remark = $3, 
            updated_at = $4, 
            user_id = $5
        WHERE payment_id = $6
    `
	_, err = tx.Exec(updateQuery, status, true, remark, time.Now(), userID, paymentID)
	if err != nil {
		return err
	}

	// 4. Fetch sale details
	var (
		totalAmount    float64
		otherCharges   sql.NullFloat64
		saleStatus     string
		dateOfDelivery time.Time
		returnDate     time.Time
	)
	saleQuery := `
        SELECT total_amount, other_charges, status, date_of_delivery, return_date
        FROM sales
        WHERE sale_id = $1
    `
	err = tx.QueryRow(saleQuery, saleID).Scan(&totalAmount, &otherCharges, &saleStatus, &dateOfDelivery, &returnDate)
	if err != nil {
		return err
	}

	// 5. Calculate total sale amount
	totalSaleAmount := totalAmount
	if otherCharges.Valid {
		totalSaleAmount += otherCharges.Float64
	}

	// 6. Sum completed payments
	var totalPaid float64
	paymentsQuery := `
        SELECT COALESCE(SUM(amount_paid), 0)
        FROM payments
        WHERE sale_id = $1 AND payment_status = 'Completed'
    `
	err = tx.QueryRow(paymentsQuery, saleID).Scan(&totalPaid)
	if err != nil {
		return err
	}

	// 7. Determine sale payment status
	var newSalePaymentStatus string
	switch {
	case totalPaid >= totalSaleAmount && saleStatus == "completed":
		newSalePaymentStatus = "paid"
	case totalPaid > 0:
		newSalePaymentStatus = "partial"
	default:
		newSalePaymentStatus = "unpaid"
	}

	// 8. Update the sale's payment status
	saleUpdateQuery := `
        UPDATE sales
        SET payment_status = $1,
            updated_at = $2
        WHERE sale_id = $3
    `
	_, err = tx.Exec(saleUpdateQuery, newSalePaymentStatus, time.Now(), saleID)
	if err != nil {
		return err
	}

	// 9. Add revenue recognition for THIS payment
	if status == "Completed" {
		revenueAmount := paymentAmount

		// For the final payment that completes the sale, adjust for any rounding differences
		if newSalePaymentStatus == "paid" {
			revenueAmount = totalSaleAmount - (totalPaid - paymentAmount)
		}

		// Calculate recognition period
		recognitionStart := dateOfDelivery
		recognitionEnd := returnDate

		// For immediate recognition (if no return date or same day)
		if returnDate.IsZero() || returnDate.Equal(dateOfDelivery) {
			recognitionEnd = dateOfDelivery
		}

		// Insert into revenue_recognition with all required fields
		_, err = tx.Exec(`
            INSERT INTO revenue_recognition (
                sale_id, 
                total_amount, 
                start_date, 
                end_date,
                updated_at
            ) VALUES (
                $1, $2, $3, $4, $5
            )
        `,
			saleID,
			revenueAmount,
			recognitionStart,
			recognitionEnd,
			time.Now(),
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetPaymentDetails fetches payment details for verification
func (r *PaymentVerificationRepository) GetPaymentDetails(paymentID int) (map[string]interface{}, error) {
	query := `
		SELECT payment_id, sale_id, amount_paid, payment_status, remark, verified_by_admin, user_id, updated_at
		FROM payments
		WHERE payment_id = $1
	`
	var payment struct {
		PaymentID       int
		SaleID          int
		AmountPaid      float64
		PaymentStatus   string
		Remark          sql.NullString
		VerifiedByAdmin bool
		UserID          sql.NullInt64
		UpdatedAt       time.Time
	}

	err := r.db.QueryRow(query, paymentID).Scan(
		&payment.PaymentID,
		&payment.SaleID,
		&payment.AmountPaid,
		&payment.PaymentStatus,
		&payment.Remark,
		&payment.VerifiedByAdmin,
		&payment.UserID,
		&payment.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	result := map[string]interface{}{
		"payment_id":        payment.PaymentID,
		"sale_id":           payment.SaleID,
		"amount_paid":       payment.AmountPaid,
		"payment_status":    payment.PaymentStatus,
		"remark":            payment.Remark.String,
		"verified_by_admin": payment.VerifiedByAdmin,
		"user_id":           payment.UserID.Int64,
		"updated_at":        payment.UpdatedAt,
	}
	if !payment.Remark.Valid {
		result["remark"] = nil
	}
	if !payment.UserID.Valid {
		result["user_id"] = nil
	}

	return result, nil
}

// CancelPayment marks a payment as canceled
func (r *PaymentVerificationRepository) CancelPayment(paymentID int, userID int, remark string) error {
	isAdmin, err := r.isAdmin(userID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errors.New("only admin users can cancel payments")
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		UPDATE payments
		SET payment_status = 'Failed', 
			remark = $1, 
			updated_at = $2, 
			user_id = $3
		WHERE payment_id = $4 AND payment_status != 'Failed'
	`
	result, err := tx.Exec(query, remark, time.Now(), userID, paymentID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("payment not found or already canceled")
	}

	return tx.Commit()
}

