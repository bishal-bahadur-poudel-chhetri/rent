package services

import (
	"errors"
	"renting/internal/models"
	"renting/internal/repositories"
)

type ExpenseService interface {
	CreateExpense(expense *models.Expense) error
	UpdateExpense(id int, expense *models.Expense) error
	DeleteExpense(id int) error
	GetExpense(id int) (*models.Expense, error)
	GetAllExpenses(filter models.ExpenseFilter) ([]models.Expense, error)
}

type expenseService struct {
	repo repositories.ExpenseRepository
}

func NewExpenseService(repo repositories.ExpenseRepository) ExpenseService {
	return &expenseService{repo}
}

func (s *expenseService) CreateExpense(expense *models.Expense) error {
	if expense.ExpenseType == "" || expense.Amount == 0 || expense.ExpenseDate.IsZero() {
		return errors.New("missing required fields")
	}
	return s.repo.Create(expense)
}

func (s *expenseService) UpdateExpense(id int, expense *models.Expense) error {
	expense.ExpenseID = id
	return s.repo.Update(expense)
}

func (s *expenseService) DeleteExpense(id int) error {
	return s.repo.Delete(id)
}

func (s *expenseService) GetExpense(id int) (*models.Expense, error) {
	return s.repo.FindByID(id)
}

func (s *expenseService) GetAllExpenses(filter models.ExpenseFilter) ([]models.Expense, error) {
	return s.repo.FindAll(filter)
}
