package models

import "time"

type Sale struct {
	SaleID             int           `json:"sale_id"`
	VehicleID          int           `json:"vehicle_id"`
	CustomerName       string        `json:"customer_name"`
	TotalAmount        float64       `json:"total_amount"`
	ChargePerDay       float64       `json:"charge_per_day"`
	BookingDate        time.Time     `json:"booking_date"`
	DateOfDelivery     time.Time     `json:"date_of_delivery"`
	ReturnDate         time.Time     `json:"return_date"`
	IsDamaged          bool          `json:"is_damaged"`
	IsWashed           bool          `json:"is_washed"`
	IsDelayed          bool          `json:"is_delayed"`
	NumberOfDays       int           `json:"number_of_days"`
	PaymentID          int           `json:"payment_id"`
	Remark             string        `json:"remark"`
	FuelRangeReceived  float64       `json:"fuel_range_received"`
	FuelRangeDelivered float64       `json:"fuel_range_delivered"`
	KmReceived         float64       `json:"km_received"`
	KmDelivered        float64       `json:"km_delivered"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
	SalesCharges       []SalesCharge `json:"sales_charges"`
	SalesImages        []SalesImage  `json:"sales_images"`
	SalesVideos        []SalesVideo  `json:"sales_videos"`
}

type SalesCharge struct {
	ChargeID   int     `json:"charge_id"`
	SaleID     int     `json:"sale_id"`
	ChargeType string  `json:"charge_type"`
	Amount     float64 `json:"amount"`
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
}
