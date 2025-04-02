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
	UserID int    `json:"user_id"`
	Remark string `json:"remark"`
}

func (h *PaymentVerification) VerifyPayment(c *gin.Context) {
	// Extract payment_id from the URL
	paymentID, err := strconv.Atoi(c.Param("payment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid payment ID", nil))
		return
	}
	saleID, err := strconv.Atoi(c.Param("sale_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid payment ID", nil))
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
	err = h.paymentService.VerifyPayment(paymentID, req.Status, userID, saleID, req.Remark)
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
