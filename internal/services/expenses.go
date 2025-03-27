package services

import (
	"errors"
	"log"
	"renting/internal/models"
	"renting/internal/repositories"

	"time"
)

type ExpenseService interface {
	CreateExpense(expense *models.Expense) error
	ReimburseExpense(expenseID uint) error
	AutoReimburse() error
	GetPendingExpenses() ([]models.Expense, error)
}

type expenseService struct {
	repo repositories.ExpenseRepository
}

func NewExpenseService(repo repositories.ExpenseRepository) ExpenseService {
	return &expenseService{repo: repo}
}

func (s *expenseService) CreateExpense(expense *models.Expense) error {
	if expense.IsReimbursable && expense.ReimbursementStatus == "" {
		expense.ReimbursementStatus = "pending"
	}
	return s.repo.Create(expense)
}

func (s *expenseService) ReimburseExpense(expenseID uint) error {
	expense, err := s.repo.FindByID(expenseID)
	if err != nil {
		return err
	}

	if !expense.IsReimbursable {
		return errors.New("expense is not reimbursable")
	}

	if expense.ReimbursementStatus == "reimbursed" {
		return errors.New("expense already reimbursed")
	}

	if err := s.repo.UpdateStatus(expenseID, "reimbursed"); err != nil {
		return err
	}

	reimbursement := &models.Reimbursement{
		ExpenseID:    expenseID,
		Amount:       expense.Amount,
		ReimbursedAt: time.Now(),
	}

	return s.repo.CreateReimbursement(reimbursement)
}

func (s *expenseService) AutoReimburse() error {
	expenses, err := s.repo.GetPending()
	if err != nil {
		return err
	}

	for _, expense := range expenses {
		if err := s.ReimburseExpense(expense.ID); err != nil {
			log.Printf("Failed to reimburse expense %d: %v", expense.ID, err)
			continue
		}
	}

	return nil
}

func (s *expenseService) GetPendingExpenses() ([]models.Expense, error) {
	return s.repo.GetPending()
}
