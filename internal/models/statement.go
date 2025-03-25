package models

import "time"

type Statement struct {
	SaleID                    int       `json:"sale_id"`
	VehicleID                 int       `json:"vehicle_id"`
	VehicleName               string    `json:"vehicle_name"`
	VehicleRegistrationNumber string    `json:"vehicle_registration_number"`
	VehicleImage              string    `json:"vehicle_image"`
	UserID                    int       `json:"user_id"`
	UserName                  string    `json:"user_name"`
	CustomerName              string    `json:"customer_name"`
	CustomerPhone             string    `json:"customer_phone"`
	BookingDate               time.Time `json:"booking_date"`
	DateOfDelivery            time.Time `json:"date_of_delivery"`
	ReturnDate                time.Time `json:"return_date"`
	TotalAmount               float64   `json:"total_amount"`
	ChargePerDay              float64   `json:"charge_per_day"`
	SaleStatus                string    `json:"sale_status"`
	AmountPaid                float64   `json:"amount_paid"`
	Balance                   float64   `json:"balance"`
	HasDamage                 bool      `json:"has_damage"`
	HasDelay                  bool      `json:"has_delay"`
	SaleImages                []string  `json:"sale_images"`
	SaleVideos                []string  `json:"sale_videos"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

type StatementFilter struct {
	DateFrom      *time.Time `form:"date_from" time_format:"2006-01-02"`
	DateTo        *time.Time `form:"date_to" time_format:"2006-01-02"`
	VehicleName   string     `form:"vehicle_name"`
	CustomerName  string     `form:"customer_name"`
	CustomerPhone string     `form:"customer_phone"`
	Status        string     `form:"status"`
	PaymentStatus string     `form:"payment_status"`
	HasDamage     bool       `form:"has_damage"`
	HasDelay      bool       `form:"has_delay"`
	SortBy        string     `form:"sort_by"`
	SortOrder     string     `form:"sort_order"`
	Limit         int        `form:"limit"`
	Offset        int        `form:"offset"`
}

type PaginatedStatements struct {
	Data       []Statement `json:"data"`
	TotalCount int         `json:"total_count"`
	Limit      int         `json:"limit"`
	Offset     int         `json:"offset"`
	HasMore    bool        `json:"has_more"`
}
