package models

import "time"

type Sale struct {
	SaleID         int            `json:"sale_id"`
	VehicleID      int            `json:"vehicle_id"`
	UserID         int            `json:"user_id"`
	CustomerName   string         `json:"customer_name"`
	Destination    string         `json:"customer_destination"`
	CustomerPhone  string         `json:"customer_phone"`
	TotalAmount    float64        `json:"total_amount"`
	ChargePerDay   float64        `json:"charge_per_day"`
	BookingDate    time.Time      `json:"booking_date"`
	DateOfDelivery time.Time      `json:"date_of_delivery"`
	ReturnDate     time.Time      `json:"return_date"`
	IsDamaged      bool           `json:"is_damaged"`
	IsWashed       bool           `json:"is_washed"`
	IsDelayed      bool           `json:"is_delayed"`
	NumberOfDays   int            `json:"number_of_days"`
	Remark         string         `json:"remark"`
	Status         string         `json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	SalesCharges   []SalesCharge  `json:"sales_charges"` // Related sales charges
	SalesImages    []SalesImage   `json:"sales_images"`  // Related sales images
	SalesVideos    []SalesVideo   `json:"sales_videos"`  // Related sales videos
	VehicleUsage   []VehicleUsage `json:"vehicle_usage"` // Related vehicle usage records
	Payments       []Payment      `json:"payments"`      // Related payments
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

type SalesVideo struct {
	VideoID    int       `json:"video_id"`
	SaleID     int       `json:"sale_id"`
	VideoURL   string    `json:"video_url"`
	UploadedAt time.Time `json:"uploaded_at"`
	FileName   string    `json:"file_name"` // Ensure this field exists
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

type Payment struct {
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
}
