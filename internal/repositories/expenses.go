package repositories

import (
	"database/sql"
	"errors"
	"log"
	"renting/internal/models"

	"time"
)

type ExpenseRepository interface {
	Create(expense *models.Expense) error
	FindByID(id uint) (*models.Expense, error)
	UpdateStatus(id uint, status string) error
	GetPending() ([]models.Expense, error)
	CreateReimbursement(reimbursement *models.Reimbursement) error
}

type expenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) ExpenseRepository {
	return &expenseRepository{db: db}
}

// Create a new expense
func (r *expenseRepository) Create(expense *models.Expense) error {
	query := `
		INSERT INTO expenses (
			user_id, amount, description, category, 
			is_reimbursable, reimbursement_status, expense_date
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(
		query,
		expense.UserID,
		expense.Amount,
		expense.Description,
		expense.Category,
		expense.IsReimbursable,
		expense.ReimbursementStatus,
		expense.ExpenseDate,
	).Scan(&expense.ID, &expense.CreatedAt, &expense.UpdatedAt)

	return err
}

// Find expense by ID
func (r *expenseRepository) FindByID(id uint) (*models.Expense, error) {
	query := `
		SELECT id, user_id, amount, description, category, 
			   is_reimbursable, reimbursement_status, expense_date,
			   created_at, updated_at
		FROM expenses
		WHERE id = $1
	`
	var expense models.Expense
	err := r.db.QueryRow(query, id).Scan(
		&expense.ID,
		&expense.UserID,
		&expense.Amount,
		&expense.Description,
		&expense.Category,
		&expense.IsReimbursable,
		&expense.ReimbursementStatus,
		&expense.ExpenseDate,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("expense not found")
		}
		return nil, err
	}

	return &expense, nil
}

// Update reimbursement status
func (r *expenseRepository) UpdateStatus(id uint, status string) error {
	query := `
		UPDATE expenses
		SET reimbursement_status = $1,
			updated_at = $2
		WHERE id = $3
	`
	_, err := r.db.Exec(query, status, time.Now(), id)
	return err
}

// Get all pending expenses
func (r *expenseRepository) GetPending() ([]models.Expense, error) {
	query := `
		SELECT id, user_id, amount, description, category, 
			   is_reimbursable, reimbursement_status, expense_date,
			   created_at, updated_at
		FROM expenses
		WHERE is_reimbursable = true AND reimbursement_status = 'pending'
		ORDER BY expense_date ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		err := rows.Scan(
			&expense.ID,
			&expense.UserID,
			&expense.Amount,
			&expense.Description,
			&expense.Category,
			&expense.IsReimbursable,
			&expense.ReimbursementStatus,
			&expense.ExpenseDate,
			&expense.CreatedAt,
			&expense.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning expense row: %v", err)
			continue
		}
		expenses = append(expenses, expense)
	}

	return expenses, nil
}

// Create reimbursement record
func (r *expenseRepository) CreateReimbursement(reimbursement *models.Reimbursement) error {
	query := `
		INSERT INTO reimbursements (
			expense_id, amount, reimbursed_at
		) VALUES ($1, $2, $3)
		RETURNING id
	`
	err := r.db.QueryRow(
		query,
		reimbursement.ExpenseID,
		reimbursement.Amount,
		reimbursement.ReimbursedAt,
	).Scan(&reimbursement.ID)

	return err
}
