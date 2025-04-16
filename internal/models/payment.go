package models

import (
	"errors"
	"time"
)

type PaymentWithSale struct {
	PaymentID       int          `json:"payment_id"`
	PaymentType     string       `json:"payment_type"`
	AmountPaid      float64      `json:"amount_paid"`
	SaleType        string       `json:"sale_type"`
	PaymentDate     time.Time    `json:"payment_date"`
	PaymentStatus   string       `json:"payment_status"`
	VerifiedByAdmin bool         `json:"verified_by_admin"`
	Remark          string       `json:"remark"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	PaymentUserID   *int         `json:"payment_user_id"`
	PaymentUsername *string      `json:"payment_username"`
	Sale            SalesPayment `json:"sale"`
}

type SalesPayment struct {
	SaleID               int             `json:"sale_id"`
	VehicleID            int             `json:"vehicle_id"`
	UserID               int             `json:"user_id"`
	Username             *string         `json:"username"` // Username of the user who created the sale
	CustomerName         string          `json:"customer_name"`
	Destination          string          `json:"customer_destination"`
	CustomerPhone        string          `json:"customer_phone"`
	TotalAmount          float64         `json:"total_amount"`
	ChargePerDay         float64         `json:"charge_per_day"`
	BookingDate          time.Time       `json:"booking_date"`
	DateOfDelivery       time.Time       `json:"date_of_delivery"`
	ReturnDate           time.Time       `json:"return_date"`
	NumberOfDays         int             `json:"number_of_days"`
	Remark               string          `json:"remark"`
	Status               string          `json:"status"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
	Vehicle              VehicleResponse `json:"vehicle"`
	ActualDateofDelivery *time.Time      `json:"actual_date_of_delivery"`
	ActualReturnDate     *time.Time      `json:"actual_date_of_return"`
	PaymentStatus        string          `json:"payment_status"`
}

type Payment struct {
	PaymentID       int       `json:"payment_id"`
	SaleID          int       `json:"sale_id"`
	AmountPaid      float64   `json:"amount_paid"`
	PaymentDate     time.Time `json:"payment_date"`
	VerifiedByAdmin bool      `json:"verified_by_admin"`
	PaymentType     string    `json:"payment_type"`
	PaymentStatus   string    `json:"payment_status"`
	Remark          string    `json:"remark"`
	UserID          *int      `json:"user_id"`  // Pointer makes it optional
	Username        *string   `json:"username"` // Username associated with UserID
	SaleType        string    `json:"sale_type"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Payment type constants
const (
	TypeBooking  = "booking"  // Initial reservation payment
	TypeDelivery = "delivery" // Payment at vehicle handover
	TypeReturn   = "return"   // Final payment at vehicle return
)

// Validate payment data
func (p *Payment) Validate() error {
	// Validate payment type
	validTypes := map[string]bool{
		TypeBooking:  true,
		TypeDelivery: true,
		TypeReturn:   true,
	}
	if !validTypes[p.SaleType] {
		return errors.New("invalid payment type, must be booking/delivery/return")
	}

	// Validate amount
	if p.AmountPaid <= 0 {
		return errors.New("payment amount must be positive")
	}

	// Validate completed payments are verified
	if p.SaleType == "Completed" && !p.VerifiedByAdmin {
		return errors.New("completed payments must be verified by admin")
	}

	return nil
}
