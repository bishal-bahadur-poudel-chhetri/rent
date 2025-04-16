package handlers

import (
	"net/http"
	"renting/internal/services"
	"renting/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DisableDateHandler struct {
	disableDateService *services.DisableDateService
}

func NewDisableDateHandler(disableDateService *services.DisableDateService) *DisableDateHandler {
	return &DisableDateHandler{disableDateService: disableDateService}
}

func (h *DisableDateHandler) GetDisabledDates(c *gin.Context) {
	vehicleIDStr := c.Query("vehicle_id")
	excludeSaleIDStr := c.Query("exclude_sale_id")

	vehicleID, err := strconv.Atoi(vehicleIDStr)
	if err != nil {
		response := utils.ErrorResponse(http.StatusBadRequest, "invalid vehicle_id", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var excludeSaleID *int
	if excludeSaleIDStr != "" {
		id, err := strconv.Atoi(excludeSaleIDStr)
		if err != nil {
			response := utils.ErrorResponse(http.StatusBadRequest, "invalid exclude_sale_id", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}
		excludeSaleID = &id
	}

	responseData, err := h.disableDateService.GetDisabledDates(vehicleID, excludeSaleID)
	if err != nil {
		response := utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.SuccessResponse(http.StatusOK, "success", responseData)
	c.JSON(http.StatusOK, response)
}

