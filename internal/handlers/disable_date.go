package handlers

import (
	"net/http"
	"renting/internal/services"
	"renting/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type DisableDateHandler struct {
	disableDateService *services.DisableDateService
}

func NewDisableDateHandler(disableDateService *services.DisableDateService) *DisableDateHandler {
	return &DisableDateHandler{disableDateService: disableDateService}
}

// GetDisabledDates handles the HTTP request to fetch disabled dates
func (h *DisableDateHandler) GetDisabledDates(c *gin.Context) {
	// Extract query parameters
	vehicleIDStr := c.Query("vehicle_id")
	dateStr := c.Query("date_of_delivery")

	// Convert vehicle_id to int
	vehicleID, err := strconv.Atoi(vehicleIDStr)
	if err != nil {
		response := utils.ErrorResponse(http.StatusBadRequest, "invalid vehicle_id", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Parse date_of_delivery
	dateOfDelivery, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		response := utils.ErrorResponse(http.StatusBadRequest, "invalid date_of_delivery", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Call the service to get disabled dates
	responseData, err := h.disableDateService.GetDisabledDates(vehicleID, dateOfDelivery)
	if err != nil {
		response := utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Return the response as JSON
	response := utils.SuccessResponse(http.StatusOK, "success", responseData)
	c.JSON(http.StatusOK, response)
}
