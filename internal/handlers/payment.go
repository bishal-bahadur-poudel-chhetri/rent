package handlers

import (
	"net/http"
	"renting/internal/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	PaymentService *services.PaymentService
}

func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{PaymentService: paymentService}
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

	payments, err := h.PaymentService.GetPaymentsWithSales(filter, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}
