package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	IsAdmin      bool      `json:"is_admin"`
	CompanyID    int       `json:"company_id"` // Ensure this field exists
	MobileNumber string    `json:"mobile_number"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
type Company struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	CompanyCode string    `json:"company_code"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	IsAdmin      bool   `json:"is_admin"`
	CompanyID    int    `json:"company_id"` // Use company_id instead of company_code
	MobileNumber string `json:"mobile_number"`
}

type Vehicle struct {
	VehicleID                 int       `json:"vehicle_id"`
	VehicleTypeID             int       `json:"vehicle_type_id"`
	VehicleName               string    `json:"vehicle_name"`
	VehicleModel              string    `json:"vehicle_model"`
	Status                    string    `json:"status"`
	VehicleRegistrationNumber string    `json:"vehicle_registration_number"`
	IsAvailable               bool      `json:"is_available"`
	ImageName                 string    `json:"image_name"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}
