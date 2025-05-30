package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Statement struct {
	StatementID          int             `json:"statement_id"`
	VehicleID            int             `json:"vehicle_id"`
	UserID               int             `json:"user_id"`
	CustomerName         string          `json:"customer_name"`
	CustomerDestination  string          `json:"customer_destination"`
	CustomerPhone        string          `json:"customer_phone"`
	TotalAmount          float64         `json:"total_amount"`
	ChargePerDay         float64         `json:"charge_per_day"`
	BookingDate          time.Time       `json:"booking_date"`
	DateOfDelivery       time.Time       `json:"date_of_delivery"`
	ReturnDate           time.Time       `json:"return_date"`
	NumberOfDays         float64         `json:"number_of_days"`
	Remark               string          `json:"remark"`
	Status               string          `json:"status"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
	ActualDateOfDelivery sql.NullTime    `json:"actual_date_of_delivery"`
	ActualDateOfReturn   sql.NullTime    `json:"actual_date_of_return"`
	PaymentStatus        string          `json:"payment_status"`
	OtherCharges         sql.NullFloat64 `json:"other_charges"`
	ModifiedBy           sql.NullString  `json:"modified_by"`
	OutstandingBalance   float64         `json:"outstanding_balance"`
	VehicleName          string          `json:"vehicle_name"`
	VehicleRegistration  string          `json:"vehicle_registration_number"`
	VehicleImage         string          `json:"image_name"`
}

// MarshalJSON customizes JSON output for nullable fields
func (s Statement) MarshalJSON() ([]byte, error) {
	type Alias Statement

	// Handle nullable fields
	var otherCharges *float64
	if s.OtherCharges.Valid {
		otherCharges = &s.OtherCharges.Float64
	}

	var modifiedBy *string
	if s.ModifiedBy.Valid {
		modifiedBy = &s.ModifiedBy.String
	}

	var actualDateOfDelivery *time.Time
	if s.ActualDateOfDelivery.Valid {
		actualDateOfDelivery = &s.ActualDateOfDelivery.Time
	}

	var actualDateOfReturn *time.Time
	if s.ActualDateOfReturn.Valid {
		actualDateOfReturn = &s.ActualDateOfReturn.Time
	}

	return json.Marshal(&struct {
		OtherCharges         *float64   `json:"other_charges"`
		ModifiedBy           *string    `json:"modified_by"`
		ActualDateOfDelivery *time.Time `json:"actual_date_of_delivery"`
		ActualDateOfReturn   *time.Time `json:"actual_date_of_return"`
		*Alias
	}{
		OtherCharges:         otherCharges,
		ModifiedBy:           modifiedBy,
		ActualDateOfDelivery: actualDateOfDelivery,
		ActualDateOfReturn:   actualDateOfReturn,
		Alias:                (*Alias)(&s),
	})
}

