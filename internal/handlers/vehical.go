package handlers

import (
	"net/http"
	"renting/internal/models"
	"renting/internal/services"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vehicleID, err := h.vehicleService.RegisterVehicle(vehicle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"vehicle_id": vehicleID})
}

// ListVehiclesHandler lists vehicles with filters and pagination
func (h *VehicleHandler) ListVehicles(c *gin.Context) {
	// Parse the filter from query parameters
	var filter models.VehicleFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default values if they are zero or not provided
	if filter.Limit == 0 {
		filter.Limit = 10 // Default limit of 10
	}
	if filter.Offset == 0 {
		filter.Offset = 0 // Default offset of 0 (first page)
	}

	// Call the service to list the vehicles
	vehicles, err := h.vehicleService.ListVehicles(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of vehicles
	c.JSON(http.StatusOK, vehicles)
}
