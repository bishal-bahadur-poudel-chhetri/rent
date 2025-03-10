package handlers

import (
	"log"
	"net/http"
	"renting/internal/models"
	"renting/internal/services"
	"renting/internal/utils"

	"github.com/gin-gonic/gin"
)

type VehicleHandler struct {
	vehicleService *services.VehicleService
}

func NewVehicleHandler(vehicleService *services.VehicleService) *VehicleHandler {
	return &VehicleHandler{
		vehicleService: vehicleService,
	}
}

// RegisterVehicleHandler handles vehicle registration requests
func (h *VehicleHandler) RegisterVehicleHandler(c *gin.Context) {
	var vehicle models.VehicleRequest
	if err := c.ShouldBindJSON(&vehicle); err != nil {
		// Use ErrorResponse from utils package
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	vehicleID, err := h.vehicleService.RegisterVehicle(vehicle)
	if err != nil {
		// Use ErrorResponse from utils package
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	// Use SuccessResponse from utils package
	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Vehicle registered successfully", gin.H{"vehicle_id": vehicleID}))
}

// ListVehiclesHandler lists vehicles with filters and pagination
func (h *VehicleHandler) ListVehicles(c *gin.Context) {
	// Parse query parameters
	var filter models.VehicleFilter

	// Bind query parameters to the filter struct
	if err := c.ShouldBindQuery(&filter); err != nil {
		// Use ErrorResponse from utils package
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid query parameters", nil))
		return
	}

	// Set default pagination values if not provided
	if filter.Limit <= 0 {
		filter.Limit = 10 // Default limit
	}
	if filter.Offset < 0 {
		filter.Offset = 0 // Default offset
	}
	filter.IsAvailable = c.DefaultQuery("is_available", "true") // Default value 'true'
	filter.VehicleName = c.DefaultQuery("vehicle_name", "")
	filter.VehicleModel = c.DefaultQuery("vehicle_model", "")
	filter.VehicleRegistrationNumber = c.DefaultQuery("vehicle_registration_number", "")
	filter.Status = c.DefaultQuery("status", "")

	// Debugging logs
	log.Printf("Filter Used: %+v\n", filter)

	// Call the service to fetch vehicles
	vehicles, err := h.vehicleService.ListVehicles(filter)
	if err != nil {
		// Use ErrorResponse from utils package
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	// Use SuccessResponse from utils package
	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Vehicles fetched successfully", gin.H{"vehicles": vehicles}))
}
