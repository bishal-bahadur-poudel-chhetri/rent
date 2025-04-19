package services

import (
	"errors"
	"renting/internal/repositories"
)

type PaymentVerificationService struct {
	paymentRepo *repositories.PaymentVerificationRepository
}

func NewPaymentVerificationService(paymentRepo *repositories.PaymentVerificationRepository) *PaymentVerificationService {
	return &PaymentVerificationService{paymentRepo: paymentRepo}
}

func (s *PaymentVerificationService) VerifyPayment(paymentID int, status string, userID int, saleID int) error {
	if status != "Completed" && status != "Failed" {
		return errors.New("invalid payment status")
	}
	return s.paymentRepo.VerifyPayment(paymentID, status, userID, saleID)
}

// GetPaymentDetails retrieves payment details, now accepting userID
func (s *PaymentVerificationService) GetPaymentDetails(paymentID int, userID int) (map[string]interface{}, error) {
	// Optionally check if user is admin or log the request
	// For now, just pass to repository
	return s.paymentRepo.GetPaymentDetails(paymentID)
}

// CancelPayment cancels a payment
func (s *PaymentVerificationService) CancelPayment(paymentID int, userID int) error {
	return s.paymentRepo.CancelPayment(paymentID, userID)
}
