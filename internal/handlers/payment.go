package handlers

import (
	"net/http"
	"renting/internal/services"
	"renting/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	PaymentService *services.PaymentService
	jwtSecret      string // Added jwtSecret field
}

// Updated constructor to accept jwtSecret
func NewPaymentHandler(paymentService *services.PaymentService, jwtSecret string) *PaymentHandler {
	return &PaymentHandler{
		PaymentService: paymentService,
		jwtSecret:      jwtSecret,
	}
}

func (h *PaymentHandler) GetPaymentsWithSales(c *gin.Context) {
	var filter services.SaleFilter

	// Parse filter parameters
	if paymentID := c.Query("payment_id"); paymentID != "" {
		id, err := strconv.Atoi(paymentID)
		if err == nil {
			filter.SaleID = &id
		}
	}

	if paymentStatus := c.Query("payment_status"); paymentStatus != "" {
		filter.PaymentStatus = &paymentStatus
	}

	if startDate := c.Query("start_date"); startDate != "" {
		parsedDate, err := time.Parse("2006-01-02", startDate)
		if err == nil {
			filter.StartDate = &parsedDate
		}
	}

	if endDate := c.Query("end_date"); endDate != "" {
		parsedDate, err := time.Parse("2006-01-02", endDate)
		if err == nil {
			filter.EndDate = &parsedDate
		}
	}

	if customerName := c.Query("customer_name"); customerName != "" {
		filter.CustomerName = &customerName
	}

	if saleStatus := c.Query("status"); saleStatus != "" {
		filter.SaleStatus = &saleStatus
	}

	// Parse pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Validate pagination
	if limit < 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Call the service to get payments
	payments, err := h.PaymentService.GetPaymentsWithSales(filter, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	// Return the response in the StandardResponse format
	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Payments retrieved successfully", payments))
}

func (h *PaymentHandler) UpdatePayment(c *gin.Context) {
	// Extract payment_id from URL parameter
	paymentID, err := strconv.Atoi(c.Param("payment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid payment ID", nil))
		return
	}

	// Extract user_id using utils.ExtractUserIDFromToken
	userID, err := utils.ExtractUserIDFromToken(c, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Unauthorized", err.Error()))
		return
	}

	// Parse request body
	type UpdatePaymentRequest struct {
		PaymentType string  `json:"payment_type" binding:"required"`
		AmountPaid  float64 `json:"amount_paid" binding:"required,gt=0"`
	}
	var req UpdatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	// Call the service to update payment
	err = h.PaymentService.UpdatePayment(paymentID, userID, req.PaymentType, req.AmountPaid)
	if err != nil {
		switch err.Error() {
		case "cannot update a completed payment":
			c.JSON(http.StatusForbidden, utils.ErrorResponse(http.StatusForbidden, err.Error(), nil))
		case "payment not found":
			c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, err.Error(), nil))
		case "only admin users can update payments":
			c.JSON(http.StatusForbidden, utils.ErrorResponse(http.StatusForbidden, err.Error(), nil))
		default:
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to update payment", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Payment updated successfully", nil))
}

func (h *PaymentHandler) InsertPayment(c *gin.Context) {
	// Extract sale_id from URL parameter (using :id as per main.go update)
	saleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid sale ID", nil))
		return
	}

	// Extract user_id using utils.ExtractUserIDFromToken
	userID, err := utils.ExtractUserIDFromToken(c, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Unauthorized", err.Error()))
		return
	}

	// Parse request body
	type InsertPaymentRequest struct {
		PaymentType string  `json:"payment_type" binding:"required"`
		AmountPaid  float64 `json:"amount_paid" binding:"required,gt=0"`
	}
	var req InsertPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	// Call the service to insert payment
	paymentID, err := h.PaymentService.InsertPayment(saleID, userID, req.PaymentType, req.AmountPaid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to create payment", err.Error()))
		return
	}

	// Return the new payment ID in the response
	responseData := map[string]int{"payment_id": paymentID}
	c.JSON(http.StatusCreated, utils.SuccessResponse(http.StatusCreated, "Payment created successfully", responseData))
}
