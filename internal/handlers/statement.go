package handlers

import (
	"net/http"
	"renting/internal/models"
	"renting/internal/services"
	"renting/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
)

type StatementHandler struct {
	service services.StatementService
}

func NewStatementHandler(service services.StatementService) *StatementHandler {
	return &StatementHandler{
		service: service,
	}
}

// GetStatements godoc
// @Summary Get rental statements
// @Description Retrieves paginated rental statements with filtering capabilities
// @Tags Statements
// @Accept json
// @Produce json
// @Param date_from query string false "Start date filter (format: YYYY-MM-DD)"
// @Param date_to query string false "End date filter (format: YYYY-MM-DD)"
// @Param vehicle_name query string false "Vehicle name filter (partial match)"
// @Param customer_name query string false "Customer name filter (partial match)"
// @Param customer_phone query string false "Customer phone number filter (exact match)"
// @Param status query string false "Status filter (pending/active/completed/cancelled)"
// @Param payment_status query string false "Payment status filter (paid/partial/unpaid)"
// @Param has_damage query bool false "Filter by damage charges"
// @Param has_delay query bool false "Filter by delay charges"
// @Param sort_by query string false "Sort field (booking_date/return_date/total_amount/vehicle_name/customer_name)"
// @Param sort_order query string false "Sort order (asc/desc)"
// @Param limit query int false "Pagination limit (default: 50)"
// @Param offset query int false "Pagination offset (default: 0)"
// @Success 200 {object} utils.StandardResponse
// @Failure 400 {object} utils.StandardResponse
// @Failure 500 {object} utils.StandardResponse
// @Router /statements [get]
func (h *StatementHandler) GetStatements(c *gin.Context) {
	var filter models.StatementFilter

	// Parse query parameters
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(
			http.StatusBadRequest,
			"Invalid query parameters",
			gin.H{"details": err.Error()},
		))
		return
	}

	// Parse and validate dates
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		if parsedDate, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filter.DateFrom = &parsedDate
		} else {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(
				http.StatusBadRequest,
				"Invalid date format",
				gin.H{"expected_format": "YYYY-MM-DD", "field": "date_from"},
			))
			return
		}
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		if parsedDate, err := time.Parse("2006-01-02", dateTo); err == nil {
			filter.DateTo = &parsedDate
		} else {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(
				http.StatusBadRequest,
				"Invalid date format",
				gin.H{"expected_format": "YYYY-MM-DD", "field": "date_to"},
			))
			return
		}
	}

	// Set defaults
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	if filter.SortBy == "" {
		filter.SortBy = "booking_date"
	}
	if filter.SortOrder == "" {
		filter.SortOrder = "desc"
	}

	// Get data from service
	result, err := h.service.GetStatements(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(
			http.StatusInternalServerError,
			"Failed to retrieve statements",
			gin.H{"details": err.Error()},
		))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Statements retrieved successfully",
		result,
	))
}
