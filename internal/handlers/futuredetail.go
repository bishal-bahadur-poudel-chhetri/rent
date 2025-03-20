package handlers

import (
	"net/http"
	"renting/internal/services"
	"renting/internal/utils" // Import the utils package
	"strconv"
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
