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

func (s *PaymentVerificationService) VerifyPayment(paymentID int, status string, userID int, saleID int, remark string) error {
	// Validate the status
	if status != "Completed" && status != "Failed" {
		return errors.New("invalid payment status")
	}

	// Call the repository to update the payment status
	err := s.paymentRepo.VerifyPayment(paymentID, status, userID, saleID, remark)
	if err != nil {
		return err
	}

	return nil
}
