package models

import "time"

type Sale struct {
	VehicleID          int       `json:"vehicle_id"`
	CustomerName       string    `json:"customer_name"`
	TotalAmount        float64   `json:"total_amount"`
	ChargePerDay       float64   `json:"charge_per_day"`
	BookingDate        time.Time `json:"booking_date"`
	DateOfDelivery     time.Time `json:"date_of_delivery"`
	ReturnDate         time.Time `json:"return_date"`
	IsDamaged          bool      `json:"is_damaged"`
	IsWashed           bool      `json:"is_washed"`
	IsDelayed          bool      `json:"is_delayed"`
	NumberOfDays       int       `json:"number_of_days"`
	PaymentID          int       `json:"payment_id"`
	Remark             string    `json:"remark"`
	FuelRangeReceived  float64   `json:"fuel_range_received"`
	FuelRangeDelivered float64   `json:"fuel_range_delivered"`
	KmReceived         int       `json:"km_received"`
	KmDelivered        int       `json:"km_delivered"`
	Photo1             string    `json:"photo_1"`
	Photo2             string    `json:"photo_2"`
	Photo3             string    `json:"photo_3"`
	Photo4             string    `json:"photo_4"`
}
