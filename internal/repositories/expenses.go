package repositories

import (
	"database/sql"
	"renting/internal/models"
	"strconv"
)

type ExpenseRepository interface {
	Create(expense *models.Expense) error
	Update(expense *models.Expense) error
	Delete(id int) error
	FindByID(id int) (*models.Expense, error)
	FindAll(filter models.ExpenseFilter) ([]models.Expense, error)
}

type expenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) ExpenseRepository {
	return &expenseRepository{db}
}

func (r *expenseRepository) Create(expense *models.Expense) error {
	query := `
		INSERT INTO company_expenses (expense_type, amount, description, expense_date, recorded_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING expense_id, created_at, updated_at`
	err := r.db.QueryRow(query,
		expense.ExpenseType,
		expense.Amount,
		expense.Description,
		expense.ExpenseDate,
		expense.RecordedBy,
	).Scan(&expense.ExpenseID, &expense.CreatedAt, &expense.UpdatedAt)
	return err
}

func (r *expenseRepository) Update(expense *models.Expense) error {
	query := `
		UPDATE company_expenses
		SET expense_type = $1,
			amount = $2,
			description = $3,
			expense_date = $4,
			recorded_by = $5,
			updated_at = NOW()
		WHERE expense_id = $6`
	_, err := r.db.Exec(query,
		expense.ExpenseType,
		expense.Amount,
		expense.Description,
		expense.ExpenseDate,
		expense.RecordedBy,
		expense.ExpenseID,
	)
	return err
}

func (r *expenseRepository) Delete(id int) error {
	query := `DELETE FROM company_expenses WHERE expense_id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *expenseRepository) FindByID(id int) (*models.Expense, error) {
	query := `
		SELECT expense_id, expense_type, amount, description, expense_date, recorded_by, created_at, updated_at
		FROM company_expenses
		WHERE expense_id = $1`
	expense := &models.Expense{}
	err := r.db.QueryRow(query, id).Scan(
		&expense.ExpenseID,
		&expense.ExpenseType,
		&expense.Amount,
		&expense.Description,
		&expense.ExpenseDate,
		&expense.RecordedBy,
		&expense.CreatedAt,
		&expense.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil // Or return a custom error if preferred
	}
	if err != nil {
		return nil, err
	}
	return expense, nil
}

func (r *expenseRepository) FindAll(filter models.ExpenseFilter) ([]models.Expense, error) {
	query := `
		SELECT expense_id, expense_type, amount, description, expense_date, recorded_by, created_at, updated_at
		FROM company_expenses
		WHERE 1=1`
	args := []interface{}{}
	argCount := 1

	if filter.ExpenseType != nil {
		query += " AND expense_type = $" + strconv.Itoa(argCount)
		args = append(args, *filter.ExpenseType)
		argCount++
	}
	if filter.StartDate != nil {
		query += " AND expense_date >= $" + strconv.Itoa(argCount)
		args = append(args, *filter.StartDate)
		argCount++
	}
	if filter.EndDate != nil {
		query += " AND expense_date <= $" + strconv.Itoa(argCount)
		args = append(args, *filter.EndDate)
		argCount++
	}
	query += " ORDER BY expense_date DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var expense models.Expense
		err := rows.Scan(
			&expense.ExpenseID,
			&expense.ExpenseType,
			&expense.Amount,
			&expense.Description,
			&expense.ExpenseDate,
			&expense.RecordedBy,
			&expense.CreatedAt,
			&expense.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return expenses, nil
}
