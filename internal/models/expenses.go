package models

import "time"

type Expense struct {
	ExpenseID   int       `json:"expense_id,omitempty" gorm:"primaryKey;autoIncrement"`
	ExpenseType string    `json:"expense_type" gorm:"type:varchar(50);not null"`
	Amount      float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	Description *string   `json:"description,omitempty" gorm:"type:text"`
	ExpenseDate time.Time `json:"expense_date" gorm:"not null"`
	RecordedBy  *int      `json:"recorded_by,omitempty"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

type ExpenseFilter struct {
	ExpenseType *string    `form:"expense_type"`
	StartDate   *time.Time `form:"start_date"`
	EndDate     *time.Time `form:"end_date"`
}
