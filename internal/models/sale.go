package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Sale struct {
	SaleID               int        `json:"sale_id"`
	VehicleID            int        `json:"vehicle_id"`
	UserID               int        `json:"user_id"`
	CustomerName         string     `json:"customer_name"`
	Destination          string     `json:"customer_destination"`
	CustomerPhone        string     `json:"customer_phone"`
	TotalAmount          float64    `json:"total_amount"`
	Discount             float64    `json:"discount"`
	OtherCharges         float64    `json:"other_charges"`
	ChargePerDay         float64    `json:"charge_per_day"`
	ChargeHalfDay        float64    `json:"charge_half_day"`
	BookingDate          time.Time  `json:"booking_date"`
	DateOfDelivery       time.Time  `json:"date_of_delivery"`
	ReturnDate           time.Time  `json:"return_date"`
	IsDamaged            bool       `json:"is_damaged"`
	IsWashed             bool       `json:"is_washed"`
	IsDelayed            bool       `json:"is_delayed"`
	IsShortTermRental    bool       `json:"is_short_term_rental"`
	NumberOfDays         float64    `json:"number_of_days"`
	FullDays             int        `json:"full_days"`
	HalfDays             int        `json:"half_days"`
	DeliveryTimeOfDay    string     `json:"delivery_time_of_day"`
	ReturnTimeOfDay      string     `json:"return_time_of_day"`
	ActualDeliveryTimeOfDay sql.NullString `json:"actual_delivery_time_of_day,omitempty"`
	ActualReturnTimeOfDay sql.NullString `json:"actual_return_time_of_day,omitempty"`
	Remark               string     `json:"remark"`
	Status               string     `json:"status"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	ActualDateOfDelivery *time.Time `json:"actual_date_of_delivery"`
	ActualReturnDate     *time.Time `json:"actual_date_of_return"`
	PaymentStatus        string     `json:"payment_status"`
	PaymentMethod        string     `json:"payment_method"`
	ModifiedBy           sql.NullInt64 `json:"modified_by"`
	// Related fields
	UserName     string         `json:"username"`
	SalesCharges []SalesCharge  `json:"sales_charges"`
	SalesImages  []SalesImage   `json:"sales_images"`
	SalesVideos  []SalesVideo   `json:"sales_videos"`
	VehicleUsage []VehicleUsage `json:"vehicle_usage"`
	Payments     []Payment      `json:"payments"`
	Vehicle      *Vehicle       `json:"vehicle,omitempty"`
	IsFutureBooking bool         `json:"is_future_booking"`
	IsComplete      bool         `json:"is_complete"`
}

// Constants for rental calculations
const (
	MinDaysForFullDayRate = 1
	HalfDayRateMultiplier = 0.5
)

type SalePending struct {
	SaleID               int            `json:"sale_id"`
	VehicleID            int            `json:"vehicle_id"`
	UserID               int            `json:"user_id"` // Change to
	UserName             string         `json:"username"`
	CustomerName         string         `json:"customer_name"`
	Destination          string         `json:"customer_destination"`
	CustomerPhone        string         `json:"customer_phone"`
	TotalAmount          float64        `json:"total_amount"`
	ChargePerDay         float64        `json:"charge_per_day"`
	BookingDate          time.Time      `json:"booking_date"`
	DateOfDelivery       time.Time      `json:"date_of_delivery"`
	ReturnDate           time.Time      `json:"return_date"`
	NumberOfDays         float64        `json:"number_of_days"`
	Remark               string         `json:"remark"`
	PaymentStatus        string         `json:"payment_status"` // Change to PaymentStatus
	Status               string         `json:"status"`
	ActualDateofDelivery *time.Time     `json:"actual_date_of_delivery"`
	ActualReturnDate     *time.Time     `json:"actual_date_of_return"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	SalesCharges         []SalesCharge  `json:"sales_charges"` // Related sales charges
	SalesImages          []SalesImage   `json:"sales_images"`  // Related sales images
	SalesVideos          []SalesVideo   `json:"sales_videos"`  // Related sales videos
	VehicleUsage         []VehicleUsage `json:"vehicle_usage"` // Related vehicle usage records
	Payments             []Payment      `json:"payments"`      // Related payments
	Vehicle              *Vehicle       `json:"vehicle,omitempty"`
}
type SaleSubmitResponse struct {
	SaleID      int    `json:"sale_id"`
	VehicleName string `json:"vehicle_name"`
}

type SalesCharge struct {
	ChargeID   int       `json:"charge_id"`
	SaleID     int       `json:"sale_id"`
	ChargeType string    `json:"charge_type"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type SalesImage struct {
	ImageID    int       `json:"image_id"`
	SaleID     int       `json:"sale_id"`
	ImageURL   string    `json:"image_url"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type VehicleUsage struct {
	UsageID    int       `json:"usage_id"`
	VehicleID  int       `json:"vehicle_id"`
	RecordType string    `json:"record_type"`
	FuelRange  float64   `json:"fuel_range"`
	KmReading  float64   `json:"km_reading"`
	RecordedAt time.Time `json:"recorded_at"`
	RecordedBy int       `json:"recorded_by"`
}

type SaleWithPayment struct {
	SaleID              int       `json:"sale_id"`
	VehicleID           int       `json:"vehicle_id"`
	UserID              int       `json:"user_id"`
	CustomerName        string    `json:"customer_name"`
	CustomerDestination string    `json:"customer_destination"`
	TotalAmount         float64   `json:"total_amount"`
	ChargePerDay        float64   `json:"charge_per_day"`
	BookingDate         time.Time `json:"booking_date"`
	DateOfDelivery      time.Time `json:"date_of_delivery"`
	ReturnDate          time.Time `json:"return_date"`
	IsDamaged           bool      `json:"is_damaged"`
	IsWashed            bool      `json:"is_washed"`
	IsDelayed           bool      `json:"is_delayed"`
	NumberOfDays        float64    `json:"number_of_days"`
	Remark              string    `json:"remark"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	Payment             Payment   `json:"payment"`
	Vehicle             *Vehicle  `json:"vehicle,omitempty"`
}
type Pagination struct {
	CurrentPage int  `json:"current_page"`
	PageSize    int  `json:"page_size"`
	TotalCount  int  `json:"total_count"`
	TotalPages  int  `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrevious bool `json:"has_previous"`
}

// PendingSalesResponse represents the response structure for pending sales
type PendingSalesResponse struct {
	Sales      []Sale     `json:"sales"`
	Pagination Pagination `json:"pagination"`
}

type UpdateSaleRequest struct {
	UserID               *int     `json:"user_id,omitempty"`
	Status               *string  `json:"status,omitempty"`
	PaymentStatus        *string  `json:"payment_status,omitempty"`
	PaymentMethod        *string  `json:"payment_method,omitempty"`
	Remark               *string  `json:"remark,omitempty"`
	CustomerName         *string  `json:"customer_name,omitempty"`
	CustomerPhone        *string  `json:"customer_phone,omitempty"`
	CustomerDestination  *string  `json:"customer_destination,omitempty"`
	TotalAmount          *float64 `json:"total_amount,omitempty"`
	Discount             *float64 `json:"discount,omitempty"`
	OtherCharges         *float64 `json:"other_charges,omitempty"`
	ChargePerDay         *float64 `json:"charge_per_day,omitempty"`
	ChargeHalfDay        *float64 `json:"charge_half_day,omitempty"`
	DateOfDelivery       *string  `json:"date_of_delivery,omitempty"`
	ReturnDate           *string  `json:"return_date,omitempty"`
	ActualDateOfDelivery *string  `json:"actual_date_of_delivery,omitempty"`
	ActualDateOfReturn   *string  `json:"actual_date_of_return,omitempty"`
	NumberOfDays         *float64 `json:"number_of_days,omitempty"`
	FullDays             *int     `json:"full_days,omitempty"`
	HalfDays             *int     `json:"half_days,omitempty"`
	VehicleID            *int     `json:"vehicle_id,omitempty"`
	IsDamaged            *bool    `json:"is_damaged,omitempty"`
	IsWashed             *bool    `json:"is_washed,omitempty"`
	IsDelayed            *bool    `json:"is_delayed,omitempty"`
	IsShortTermRental    *bool    `json:"is_short_term_rental,omitempty"`
	DeliveryTimeOfDay    *string  `json:"delivery_time_of_day,omitempty"`
	ReturnTimeOfDay      *string  `json:"return_time_of_day,omitempty"`
	ActualDeliveryTimeOfDay *string `json:"actual_delivery_time_of_day,omitempty"`
	ActualReturnTimeOfDay   *string `json:"actual_return_time_of_day,omitempty"`
	ModifiedBy           *int     `json:"modified_by,omitempty"`
	IsComplete           *bool    `json:"is_complete,omitempty"`
}

type AddChargeRequest struct {
	SaleID      int     `json:"sale_id"`
	ChargeType  string  `json:"charge_type"` // "discount", "wash", or "damage"
	Amount      float64 `json:"amount"`
	Remark      string  `json:"remark,omitempty"`
}

type SaleFilter struct {
	Status             string
	ActualDateOfDelivery *time.Time
	DateOfDeliveryBefore *time.Time
	CustomerName       string
	VehicleID          int
	Sort               string
	Limit              int
	Offset             int
}

// MarshalJSON customizes JSON output for nullable fields
func (s Sale) MarshalJSON() ([]byte, error) {
	type Alias Sale

	// Handle nullable fields
	var actualDeliveryTimeOfDay *string
	if s.ActualDeliveryTimeOfDay.Valid {
		actualDeliveryTimeOfDay = &s.ActualDeliveryTimeOfDay.String
	}

	var actualReturnTimeOfDay *string
	if s.ActualReturnTimeOfDay.Valid {
		actualReturnTimeOfDay = &s.ActualReturnTimeOfDay.String
	}

	return json.Marshal(&struct {
		ActualDeliveryTimeOfDay *string `json:"actual_delivery_time_of_day"`
		ActualReturnTimeOfDay   *string `json:"actual_return_time_of_day"`
		*Alias
	}{
		ActualDeliveryTimeOfDay: actualDeliveryTimeOfDay,
		ActualReturnTimeOfDay:   actualReturnTimeOfDay,
		Alias:                   (*Alias)(&s),
	})
}
