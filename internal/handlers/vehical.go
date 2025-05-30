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
	saleService    *services.SaleService
}

func NewVehicleHandler(vehicleService *services.VehicleService, saleService *services.SaleService) *VehicleHandler {
	return &VehicleHandler{vehicleService: vehicleService, saleService: saleService}
}

// StandardResponse defines the structure of the API response
type StandardResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // Optional field to send extra data with error response
}

// SuccessResponse creates a dynamic success response
func SuccessResponse(status int, message string, data interface{}) StandardResponse {
	// If no data is passed, set Data to nil
	if data == nil {
		data = struct{}{}
	}
	return StandardResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse creates a dynamic error response
func ErrorResponse(status int, message string, data interface{}) StandardResponse {
	return StandardResponse{
		Status:  status,
		Message: message,
		Data:    data, // Send extra data if provided
	}
}

func (h *VehicleHandler) ListVehicles(c *gin.Context) {
	// Parse query parameters
	queryParams := c.Request.URL.Query()

	// Log query parameters for debugging
	log.Printf("Query Parameters: %+v", queryParams)

	filters := models.VehicleFilter{
		VehicleTypeID:             queryParams.Get("vehicle_type_id"),
		VehicleName:               queryParams.Get("vehicle_name"),
		VehicleModel:              queryParams.Get("vehicle_model"),
		VehicleRegistrationNumber: queryParams.Get("vehicle_registration_number"),
		IsAvailable:               queryParams.Get("is_available"),
		Status:                    queryParams.Get("status"),
	}

	// Parse pagination parameters
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

	// Check if booking details should be included
	includeBookingDetails := queryParams.Get("data") == "true"
	includeSaleid := queryParams.Get("rentedSaleId") == "true"
	log.Printf("Include Booking Details: %v", includeSaleid)

	vehicles, err := h.vehicleService.GetVehicles(filters, includeBookingDetails, includeSaleid)

	if err != nil {
		response := ErrorResponse(http.StatusInternalServerError, err.Error(), nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Prepare success response
	response := SuccessResponse(http.StatusOK, "Vehicles fetched successfully", gin.H{
		"vehicles": vehicles,
	})

	c.JSON(http.StatusOK, response)
}

func (h *VehicleHandler) AddCharge(c *gin.Context) {
	var req models.AddChargeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request body",
			"data":    err.Error(),
		})
		return
	}

	err := h.saleService.AddCharge(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to add charge",
			"data":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Charge added successfully",
		"data":    nil,
	})
}
