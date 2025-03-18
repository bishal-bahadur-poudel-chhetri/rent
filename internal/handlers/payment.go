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

func NewPaymentHandler(PaymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{PaymentService: PaymentService}
}

func (h *PaymentHandler) GetPaymentsWithSales(c *gin.Context) {
	var filter services.SaleFilter

	// Parse query parameters
	if paymentID := c.Query("payment_id"); paymentID != "" {
		id, err := strconv.Atoi(paymentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payment_id"})
			return
		}
		filter.SaleID = &id
	}

	if paymentStatus := c.Query("payment_status"); paymentStatus != "" {
		filter.PaymentStatus = &paymentStatus
	}

	if startDate := c.Query("start_date"); startDate != "" {
		parsedStartDate, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date"})
			return
		}
		filter.StartDate = &parsedStartDate
	}

	if endDate := c.Query("end_date"); endDate != "" {
		parsedEndDate, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date"})
			return
		}
		filter.EndDate = &parsedEndDate
	}

	payments, err := h.PaymentService.GetPaymentsWithSales(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}
