package models

import "time"

type PaymentWithSale struct {
	PaymentID       int           `json:"payment_id"`
	PaymentType     string        `json:"payment_type"`
	AmountPaid      float64       `json:"amount_paid"`
	PaymentDate     time.Time     `json:"payment_date"`
	PaymentStatus   string        `json:"payment_status"`
	VerifiedByAdmin bool          `json:"verified_by_admin"`
	Remark          string        `json:"remark"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	Sale            Sales_payment `json:"sale"` // Use Sales_payment instead of Sale
}

type Sales_payment struct {
	SaleID         int             `json:"sale_id"`
	VehicleID      int             `json:"vehicle_id"`
	UserID         int             `json:"user_id"`
	CustomerName   string          `json:"customer_name"`
	Destination    string          `json:"customer_destination"`
	CustomerPhone  string          `json:"customer_phone"`
	TotalAmount    float64         `json:"total_amount"`
	ChargePerDay   float64         `json:"charge_per_day"`
	BookingDate    time.Time       `json:"booking_date"`
	DateOfDelivery time.Time       `json:"date_of_delivery"`
	ReturnDate     time.Time       `json:"return_date"`
	IsDamaged      bool            `json:"is_damaged"`
	IsWashed       bool            `json:"is_washed"`
	IsDelayed      bool            `json:"is_delayed"`
	NumberOfDays   int             `json:"number_of_days"`
	Remark         string          `json:"remark"`
	Status         string          `json:"status"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Vehicle        VehicleResponse `json:"vehicle"`
}
