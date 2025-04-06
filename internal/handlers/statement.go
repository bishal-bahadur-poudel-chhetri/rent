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
	if bookingDate := c.Query("booking_date"); bookingDate != "" {
		filters["booking_date"] = bookingDate
	}
	if status := c.Query("status"); status != "" {
		filters["status"] = status
	}
	if paymentStatus := c.Query("payment_status"); paymentStatus != "" {
		filters["payment_status"] = paymentStatus
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
	if limit <= 0 {
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
