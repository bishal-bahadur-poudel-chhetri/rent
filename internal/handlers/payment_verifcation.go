package handlers

import (
	"errors"
	"net/http"
	"renting/internal/services"
	"renting/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentVerification struct {
	paymentService *services.PaymentVerificationService
	jwtSecret      string
}

func NewPaymentVerification(paymentService *services.PaymentVerificationService, jwtSecret string) *PaymentVerification {
	return &PaymentVerification{paymentService: paymentService, jwtSecret: jwtSecret}
}

type VerifyPaymentRequest struct {
	Status string `json:"status"`
}

type CancelPaymentRequest struct {
}

// VerifyPayment handles payment verification (POST)
func (h *PaymentVerification) VerifyPayment(c *gin.Context) {
	// Extract payment_id from the URL
	paymentID, err := strconv.Atoi(c.Param("payment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid payment ID", nil))
		return
	}
	saleID, err := strconv.Atoi(c.Param("sale_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid sale ID", nil))
		return
	}

	// Extract user ID from the JWT token
	userID, err := utils.ExtractUserIDFromToken(c, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Unauthorized", err.Error()))
		return
	}

	// Parse the request body
	var req VerifyPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", nil))
		return
	}

	// Call the service layer to verify the payment
	err = h.paymentService.VerifyPayment(paymentID, req.Status, userID, saleID)
	if err != nil {
		switch {
		case errors.Is(err, errors.New("only admin users can verify payments")):
			c.JSON(http.StatusForbidden, utils.ErrorResponse(http.StatusForbidden, err.Error(), nil))
		case errors.Is(err, errors.New("invalid payment status")):
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		default:
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Internal server error", err.Error()))
		}
		return
	}

	// Return success response
	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Payment verification updated successfully", nil))
}

// GetPaymentDetails handles fetching payment details (GET)
func (h *PaymentVerification) GetPaymentDetails(c *gin.Context) {
	paymentID, err := strconv.Atoi(c.Param("payment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid payment ID", nil))
		return
	}

	userID, err := utils.ExtractUserIDFromToken(c, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Unauthorized", err.Error()))
		return
	}

	// Pass userID to the service layer (e.g., for potential admin check or logging)
	paymentDetails, err := h.paymentService.GetPaymentDetails(paymentID, userID)
	if err != nil {
		if errors.Is(err, errors.New("payment not found")) {
			c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, err.Error(), nil))
		} else {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Internal server error", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Payment details retrieved successfully", paymentDetails))
}

// CancelPayment handles payment cancellation (POST)
func (h *PaymentVerification) CancelPayment(c *gin.Context) {
	paymentID, err := strconv.Atoi(c.Param("payment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid payment ID", nil))
		return
	}

	userID, err := utils.ExtractUserIDFromToken(c, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Unauthorized", err.Error()))
		return
	}

	err = h.paymentService.CancelPayment(paymentID, userID)
	if err != nil {
		switch {
		case errors.Is(err, errors.New("only admin users can cancel payments")):
			c.JSON(http.StatusForbidden, utils.ErrorResponse(http.StatusForbidden, err.Error(), nil))
		case errors.Is(err, errors.New("payment not found or already canceled")):
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		default:
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Internal server error", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Payment canceled successfully", nil))
}
