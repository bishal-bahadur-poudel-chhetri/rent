package models

import "time"

type Sale struct {
	SaleID         int            `json:"sale_id"`
	VehicleID      int            `json:"vehicle_id"`
	UserID         int            `json:"user_id"`
	CustomerName   string         `json:"customer_name"`
<<<<<<< HEAD
	Destination    string         `json:"customer_destination"` // Fixed typo in the tag
=======
	Destination    string         `json:"customer_destination"`
>>>>>>> 6e1f2f9 (update on sales and vehical api)
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
<<<<<<< HEAD
	Payments       []Payment      `json:"payments"`      // List of payments associated with the sale
=======
>>>>>>> 6e1f2f9 (update on sales and vehical api)
	SalesCharges   []SalesCharge  `json:"sales_charges"` // Related sales charges
	SalesImages    []SalesImage   `json:"sales_images"`  // Related sales images
	SalesVideos    []SalesVideo   `json:"sales_videos"`  // Related sales videos
	VehicleUsage   []VehicleUsage `json:"vehicle_usage"` // Related vehicle usage records
<<<<<<< HEAD
}

type Payment struct {
	PaymentID       int       `json:"payment_id"`
	SaleID          int       `json:"sale_id"`           // Foreign key to associate with a sale
	PaymentType     string    `json:"payment_type"`      // e.g., credit_card, cash, etc.
	AmountPaid      float64   `json:"amount_paid"`       // Amount paid in this payment
	PaymentDate     time.Time `json:"payment_date"`      // Date of the payment
	PaymentStatus   string    `json:"payment_status"`    // e.g., Completed, Pending, Failed
	VerifiedByAdmin bool      `json:"verified_by_admin"` // Whether the payment is verified
	Remark          string    `json:"remark"`            // Additional notes about the payment
	CreatedAt       time.Time `json:"created_at"`        // Timestamp when the payment was created
	UpdatedAt       time.Time `json:"updated_at"`        // Timestamp when the payment was last updated
=======
	Payments       []Payment      `json:"payments"`      // Related payments
>>>>>>> 6e1f2f9 (update on sales and vehical api)
}

type SalesCharge struct {
	ChargeID   int       `json:"charge_id"`
	SaleID     int       `json:"sale_id"`
	ChargeType string    `json:"charge_type"` // e.g., damage, wash, delay, discount
	Amount     float64   `json:"amount"`      // Amount of the charge
	CreatedAt  time.Time `json:"created_at"`  // Timestamp when the charge was created
	UpdatedAt  time.Time `json:"updated_at"`  // Timestamp when the charge was last updated
}

type SalesImage struct {
	ImageID    int       `json:"image_id"`
	SaleID     int       `json:"sale_id"`
	ImageURL   string    `json:"image_url"`   // URL of the image
	UploadedAt time.Time `json:"uploaded_at"` // Timestamp when the image was uploaded
}

type SalesVideo struct {
	VideoID    int       `json:"video_id"`
	SaleID     int       `json:"sale_id"`
	VideoURL   string    `json:"video_url"`   // URL of the video
	UploadedAt time.Time `json:"uploaded_at"` // Timestamp when the video was uploaded
}

type VehicleUsage struct {
	UsageID    int       `json:"usage_id"`
	VehicleID  int       `json:"vehicle_id"`
	RecordType string    `json:"record_type"` // e.g., delivery, return
	FuelRange  float64   `json:"fuel_range"`  // Fuel level at the time of recording
	KmReading  float64   `json:"km_reading"`  // Kilometer reading at the time of recording
	RecordedAt time.Time `json:"recorded_at"` // Timestamp when the usage was recorded
	RecordedBy int       `json:"recorded_by"` // ID of the user who recorded the usage
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
