package handlers

import (
	"log"
	"net/http"
	"renting/internal/models"
	"renting/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type VehicleHandler struct {
	vehicleService *services.VehicleService
}

func NewVehicleHandler(vehicleService *services.VehicleService) *VehicleHandler {
	return &VehicleHandler{vehicleService: vehicleService}
}

type StandardResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(status int, message string, data interface{}) StandardResponse {

	if data == nil {
		data = struct{}{}
	}
	return StandardResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(status int, message string, data interface{}) StandardResponse {
	return StandardResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func (h *VehicleHandler) ListVehicles(c *gin.Context) {

	queryParams := c.Request.URL.Query()

	log.Printf("Query Parameters: %+v", queryParams)

	filters := models.VehicleFilter{
		VehicleID:                 queryParams.Get("vehicle_id"),
		VehicleTypeID:             queryParams.Get("vehicle_type_id"),
		VehicleName:               queryParams.Get("vehicle_name"),
		VehicleModel:              queryParams.Get("vehicle_model"),
		VehicleRegistrationNumber: queryParams.Get("vehicle_registration_number"),
		IsAvailable:               queryParams.Get("is_available"),
		Status:                    queryParams.Get("status"),
	}

	if limit := queryParams.Get("limit"); limit != "" {
		limitValue, err := strconv.Atoi(limit)
		if err != nil {
			response := ErrorResponse(http.StatusBadRequest, "Invalid limit value", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}
		filters.Limit = limitValue
	}
	if offset := queryParams.Get("offset"); offset != "" {
		offsetValue, err := strconv.Atoi(offset)
		if err != nil {
			response := ErrorResponse(http.StatusBadRequest, "Invalid offset value", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}
		filters.Offset = offsetValue
	}

	includeBookingDetails := queryParams.Get("data") == "true"
	includeSaleid := queryParams.Get("rentedSaleId") == "true"
	log.Printf("Include Booking Details: %v", includeSaleid)

	vehicles, err := h.vehicleService.GetVehicles(filters, includeBookingDetails, includeSaleid)
	if err != nil {
		response := ErrorResponse(http.StatusInternalServerError, err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := SuccessResponse(http.StatusOK, "Vehicles fetched successfully", gin.H{
		"vehicles": vehicles,
	})

	c.JSON(http.StatusOK, response)
}
