package handlers

import (
	"net/http"
	"renting/internal/services"
	"renting/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StatementHandler struct {
	svc services.StatementService
}

func NewStatementHandler(svc services.StatementService) *StatementHandler {
	return &StatementHandler{svc: svc}
}

func (h *StatementHandler) GetOutstandingStatements(c *gin.Context) {
	// Parse query parameters
	filters := make(map[string]string)
	
	// Date range filters
	if startDate := c.Query("start_date"); startDate != "" {
		filters["start_date"] = startDate
	}
	if endDate := c.Query("end_date"); endDate != "" {
		filters["end_date"] = endDate
	}
	
	// Sale and customer filters
	if saleID := c.Query("sale_id"); saleID != "" {
		filters["sale_id"] = saleID
	}
	if customerName := c.Query("customer_name"); customerName != "" {
		filters["customer_name"] = customerName
	}
	
	// Other filters
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if paymentStatus := c.Query("payment_status"); paymentStatus != "" {
		filters["payment_status"] = paymentStatus
	}
	if vehicleID := c.Query("vehicle_id"); vehicleID != "" {
		filters["vehicle_id"] = vehicleID
	}
	if vehicleName := c.Query("vehicle_name"); vehicleName != "" {
		filters["vehicle_name"] = vehicleName
	}

	// Parse pagination parameters
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if offset < 0 {
		offset = 0
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit <= 0 || limit > 100 { // Added upper limit for safety
		limit = 10
	}

	// Call service
	statements, err := h.svc.GetOutstandingStatements(c.Request.Context(), filters, offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(
			http.StatusInternalServerError,
			err.Error(),
			nil,
		))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Outstanding statements retrieved successfully",
		statements,
	))
}
