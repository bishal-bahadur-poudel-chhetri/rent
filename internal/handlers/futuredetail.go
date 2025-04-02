package handlers

import (
	"net/http"
	"renting/internal/services"
	"renting/internal/utils" // Import the utils package
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type FuturBookingHandler struct {
	futurBookingService *services.FuturBookingService
}

func NewFuturBookingHandler(futurBookingService *services.FuturBookingService) *FuturBookingHandler {
	return &FuturBookingHandler{futurBookingService: futurBookingService}
}

func (h *FuturBookingHandler) GetFuturBookingsByMonth(c *gin.Context) {
	// Extract query parameters
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	// Validate and convert year
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "invalid year", nil))
		return
	}

	// Validate and convert month
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "invalid month", nil))
		return
	}

	// Extract filters from query parameters
	filters := make(map[string]string)
	for key, value := range c.Request.URL.Query() {
		if key != "year" && key != "month" { // Exclude year and month from filters
			filters[key] = value[0] // Use the first value if there are multiple values for the same key
		}
	}

	// Call the service to get the data
	response, err := h.futurBookingService.GetFuturBookingsByMonth(year, time.Month(month), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	// Return the response
	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "success", response))
}

func (h *FuturBookingHandler) CancelFuturBooking(c *gin.Context) {
	// Extract and validate sale ID
	saleIDStr := c.Param("saleID")
	saleID, err := strconv.Atoi(saleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(
			http.StatusBadRequest,
			"invalid sale ID format",
			nil,
		))
		return
	}

	// Call service to cancel booking
	err = h.futurBookingService.CancelFuturBooking(saleID)
	if err != nil {
		// Handle different error cases
		if strings.Contains(err.Error(), "no booking found") {
			c.JSON(http.StatusNotFound, utils.ErrorResponse(
				http.StatusNotFound,
				err.Error(),
				nil,
			))
		} else {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(
				http.StatusInternalServerError,
				err.Error(),
				nil,
			))
		}
		return
	}

	// Success response
	c.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"booking cancelled successfully",
		nil,
	))
}
