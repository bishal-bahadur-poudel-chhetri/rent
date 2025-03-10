package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	IsAdmin      bool      `json:"is_admin"`
	CompanyID    int       `json:"company_id"`
	MobileNumber string    `json:"mobile_number"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
