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

func (r *PaymentVerificationRepository) VerifyPayment(paymentID int, status string, userID int, saleID int, remark string) error {
	// Check if the user is an admin
	isAdmin, err := r.isAdmin(userID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errors.New("only admin users can verify payments")
	}

	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update the specific payment
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

	// Get total sale amount with nullable other_charges
	var totalAmount float64
	var saleCharges sql.NullFloat64
	saleQuery := `
		SELECT total_amount, other_charges
		FROM sales
		WHERE sale_id = $1
	`
	err = tx.QueryRow(saleQuery, saleID).Scan(&totalAmount, &saleCharges)
	if err != nil {
		return err
	}

	// Handle NULL other_charges by defaulting to 0 if invalid/NULL
	totalSaleAmount := totalAmount
	if saleCharges.Valid {
		totalSaleAmount += saleCharges.Float64
	}

	// Get sum of all payments for this sale
	var totalPaid float64
	paymentsQuery := `
		SELECT COALESCE(SUM(amount_paid), 0)
		FROM payments
		WHERE sale_id = $1
	`
	err = tx.QueryRow(paymentsQuery, saleID).Scan(&totalPaid)
	if err != nil {
		return err
	}

	// If total paid meets or exceeds sale amount, update statuses
	const epsilon = 0.01
	if totalPaid+epsilon >= totalSaleAmount {
		// Update sales table
		saleUpdateQuery := `
			UPDATE sales
			SET payment_status = 'paid',
				updated_at = $1
			WHERE sale_id = $2
		`
		_, err = tx.Exec(saleUpdateQuery, time.Now(), saleID)
		if err != nil {
			return err
		}

		// Optionally update all payments for this sale
		paymentsUpdateQuery := `
			UPDATE sales
			SET payment_status = 'paid',
				updated_at = $1,
			WHERE sale_id = $3
		`
		_, err = tx.Exec(paymentsUpdateQuery, time.Now(), saleID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
