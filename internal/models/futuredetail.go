package models

import "time"

type Payment_future struct {
	PaymentID     int
	AmountPaid    float64
	PaymentDate   time.Time
	PaymentType   string
	PaymentStatus string
}
type Sale_Future struct {
	SaleID         int       `json:"sale_id"`
	VehicleID      int       `json:"vehicle_id"`
	UserID         int       `json:"user_id"`
	CustomerName   string    `json:"customer_name"`
	Destination    string    `json:"customer_destination"`
	CustomerPhone  string    `json:"customer_phone"`
	TotalAmount    float64   `json:"total_amount"`
	ChargePerDay   float64   `json:"charge_per_day"`
	BookingDate    time.Time `json:"booking_date"`
	DateOfDelivery time.Time `json:"date_of_delivery"`
	ReturnDate     time.Time `json:"return_date"`
	NumberOfDays   int       `json:"number_of_days"`
	Remark         string    `json:"remark"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	VehicleName    string    `json:"vehicle_name"`
	VehicleImage   string    `json:"image_name"`

	Payment []Payment_future `json:"payments"`
}

type MonthlyMetadata struct {
	Month string `json:"month"` // Format: "YYYY-MM"
	Count int    `json:"count"` // Number of bookings in the month
}

type SalesWithMetadataResponse struct {
	Sales    []Sale_Future     `json:"sales"`    // Sales for the selected month
	Metadata []MonthlyMetadata `json:"metadata"` // Metadata for other months
}
