package models

import "time"

// User represents user data structure
type User struct {
	ID           int       `json:"id"`
	Username     string    `gorm:"unique" json:"username"`
	Password     string    `json:"-"`             // Password is not exposed in JSON
	IsAdmin      bool      `json:"is_admin"`      // Indicates if the user is an admin
	CompanyID    int       `json:"company_id"`    // Associates the user with a company
	MobileNumber string    `json:"mobile_number"` // Mobile number of the user
	CreatedAt    time.Time `json:"created_at"`
	Updated_at   time.Time `json:"updated_at"`
}

// Company represents company data structure
type Company struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	CompanyCode string    `json:"company_code" gorm:"unique"`
	CreatedAt   time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}

// LoginRequest represents login credentials with mobile, password and company code
type LoginRequest struct {
	MobileNumber string `json:"mobile_number" validate:"required"`
	CompanyCode  string `json:"company_code" validate:"required"`
	Password     string `json:"password" validate:"required"`
}

// RegisterRequest represents registration data
type RegisterRequest struct {
	Username     string `json:"username" binding:"required,min=3"`
	Password     string `json:"password" binding:"required,min=6"`
	MobileNumber string `json:"mobile_number" binding:"required"`
	CompanyID    int    `json:"company_id" binding:"required"`
	IsAdmin      bool   `json:"is_admin"`
}

// UpdateProfileRequest represents profile update data
type UpdateProfileRequest struct {
	Username     string `json:"username"`
	MobileNumber string `json:"mobile_number"`
}

// TokenResponse contains JWT token
type TokenResponse struct {
	Token     string `json:"token"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	IsAdmin   bool   `json:"is_admin"`
	CompanyID int    `json:"company_id"`
}

type VehicleRequest struct {
	VehicleTypeID             string `json:"vehicle_type_id"`
	VehicleName               string `json:"vehicle_name"`
	VehicleModel              string `json:"vehicle_model"`
	VehicleRegistrationNumber string `json:"vehicle_registration_number"`
	IsAvailable               bool   `json:"is_available"`
	Status                    string `json:"status"`
}
