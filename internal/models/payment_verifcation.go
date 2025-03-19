package models

import "time"

type PaymentVerificationRequest struct {
	PaymentID       int       `json:"payment_id"`
	PaymentType     string    `json:"payment_type"`
	AmountPaid      float64   `json:"amount_paid"`
	PaymentDate     time.Time `json:"payment_date"`
	PaymentStatus   string    `json:"payment_status"` // Pending, Completed, Failed
	VerifiedByAdmin bool      `json:"verified_by_admin"`
	Remark          string    `json:"remark"`
	UserID          int       `json:"user_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	SaleID          int       `json:"sale_id"` // Foreign key to Sale
}
