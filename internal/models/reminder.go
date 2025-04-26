package models

import (
	"time"
)

type ReminderType string
type ReminderFrequency string

const (
	ReminderTypeEMI       ReminderType = "emi"
	ReminderTypeInsurance ReminderType = "insurance"
	ReminderTypeBillbook  ReminderType = "billbook"
	ReminderTypeServicing ReminderType = "servicing"

	ReminderFrequencyMonthly ReminderFrequency = "monthly"
	ReminderFrequencyYearly  ReminderFrequency = "yearly"
	ReminderFrequencyCustom  ReminderFrequency = "custom"
)

// Reminder represents a vehicle-related reminder
type Reminder struct {
	ID             int               `json:"id" gorm:"primaryKey;autoIncrement"`
	VehicleID      int               `json:"vehicle_id" gorm:"not null"`
	UserID         int               `json:"user_id" gorm:"not null"`
	Type           ReminderType      `json:"type" gorm:"type:reminder_type;not null"`
	StartDate      time.Time         `json:"start_date" gorm:"not null"`
	Frequency      ReminderFrequency `json:"frequency" gorm:"type:reminder_frequency;not null"`
	CustomInterval *int              `json:"custom_interval,omitempty" gorm:"default:null"` // in days
	NextDueDate    time.Time         `json:"next_due_date" gorm:"not null"`
	CreatedAt      time.Time         `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time         `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// ReminderWithVehicle represents a reminder with associated vehicle details
type ReminderWithVehicle struct {
	Reminder
	VehicleName               string `json:"vehicle_name"`
	VehicleModel              string `json:"vehicle_model"`
	VehicleRegistrationNumber string `json:"vehicle_registration_number"`
}

// ReminderAcknowledgement represents a record of when a reminder was acknowledged
type ReminderAcknowledgement struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ReminderID     int       `json:"reminder_id" gorm:"not null"`
	UserID         int       `json:"user_id" gorm:"not null"`
	AcknowledgedAt time.Time `json:"acknowledged_at" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName specifies the table name for Reminder
func (Reminder) TableName() string {
	return "reminders"
}

// TableName specifies the table name for ReminderAcknowledgement
func (ReminderAcknowledgement) TableName() string {
	return "reminder_acknowledgements"
}
