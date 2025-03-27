package models

import "time"

type Expense struct {
	ID                  uint      `json:"id" gorm:"primaryKey"`
	UserID              uint      `json:"user_id"`
	Amount              float64   `json:"amount" gorm:"type:decimal(10,2)"`
	Description         string    `json:"description"`
	Category            string    `json:"category"`
	IsReimbursable      bool      `json:"is_reimbursable"`
	ReimbursementStatus string    `json:"reimbursement_status" gorm:"default:pending"` // pending, reimbursed
	ExpenseDate         time.Time `json:"expense_date" gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

type Reimbursement struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	ExpenseID    uint      `json:"expense_id"`
	Amount       float64   `json:"amount" gorm:"type:decimal(10,2)"`
	ReimbursedAt time.Time `json:"reimbursed_at" gorm:"default:CURRENT_TIMESTAMP"`
}
