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

// VerifyPayment updates the payment status in the database only if the user is an admin
func (r *PaymentVerificationRepository) VerifyPayment(paymentID int, status string, userID int, remark string) error {
	// Check if the user is an admin
	isAdmin, err := r.isAdmin(userID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return errors.New("only admin users can verify payments")
	}

	query := `
		UPDATE payments
		SET payment_status = $1, verified_by_admin = $2, remark = $3, updated_at = $4, user_id = $6
		WHERE payment_id = $5
	`
	_, err = r.db.Exec(query, status, true, remark, time.Now(), paymentID, userID)
	if err != nil {
		return err
	}
	return nil
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
