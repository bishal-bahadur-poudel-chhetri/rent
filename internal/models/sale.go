package models

import "time"

type Sale struct {
	SaleID               int        `json:"sale_id"`
	VehicleID            int        `json:"vehicle_id"`
	UserID               int        `json:"user_id"`
	CustomerName         string     `json:"customer_name"`
	Destination          string     `json:"customer_destination"`
	CustomerPhone        string     `json:"customer_phone"`
	TotalAmount          float64    `json:"total_amount"`
	ChargePerDay         float64    `json:"charge_per_day"`
	BookingDate          time.Time  `json:"booking_date"`
	DateOfDelivery       time.Time  `json:"date_of_delivery"`
	ReturnDate           time.Time  `json:"return_date"`
	IsDamaged            bool       `json:"is_damaged"`
	IsWashed             bool       `json:"is_washed"`
	IsDelayed            bool       `json:"is_delayed"`
	NumberOfDays         int        `json:"number_of_days"`
	Remark               string     `json:"remark"`
	Status               string     `json:"status"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	ActualDateOfDelivery *time.Time `json:"actual_date_of_delivery"`
	ActualReturnDate     *time.Time `json:"actual_date_of_return"`
	PaymentStatus        string     `json:"payment_status"`
	OtherCharges         *float64   `json:"other_charges"`
	ModifiedBy           int        `json:"modified_by"` // Assuming this is a user ID
	// Related fields from your previous struct
	UserName     string         `json:"username"`
	SalesCharges []SalesCharge  `json:"sales_charges"`
	SalesImages  []SalesImage   `json:"sales_images"`
	SalesVideos  []SalesVideo   `json:"sales_videos"`
	VehicleUsage []VehicleUsage `json:"vehicle_usage"`
	Payments     []Payment      `json:"payments"`
	Vehicle      *Vehicle       `json:"vehicle,omitempty"`
}

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
	NumberOfDays         int            `json:"number_of_days"`
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
	NumberOfDays        int       `json:"number_of_days"`
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
	Remark               *string  `json:"remark,omitempty"`
	CustomerName         *string  `json:"customer_name,omitempty"`
	CustomerPhone        *string  `json:"customer_phone,omitempty"`
	CustomerDestination  *string  `json:"customer_destination,omitempty"`
	TotalAmount          *float64 `json:"total_amount,omitempty"`
	ChargePerDay         *float64 `json:"charge_per_day,omitempty"`
	DateOfDelivery       *string  `json:"date_of_delivery,omitempty"`
	ReturnDate           *string  `json:"return_date,omitempty"`
	ActualDateOfDelivery *string  `json:"actual_date_of_delivery,omitempty"`
	ActualDateOfReturn   *string  `json:"actual_date_of_return,omitempty"`
	NumberOfDays         *int     `json:"number_of_days,omitempty"`
	VehicleID            *int     `json:"vehicle_id"`
}
