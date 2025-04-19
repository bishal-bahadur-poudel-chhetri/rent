package handlers

import (
	"encoding/json"
	"net/http"
	"renting/internal/services"
	"renting/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
)

type VehicleServicingHandler struct {
	servicingService *services.VehicleServicingService
}

func NewVehicleServicingHandler(servicingService *services.VehicleServicingService) *VehicleServicingHandler {
	return &VehicleServicingHandler{servicingService: servicingService}
}

// InitializeServicingRecord handles the creation of a new servicing record
func (h *VehicleServicingHandler) InitializeServicingRecord(w http.ResponseWriter, r *http.Request) {
	var req struct {
		VehicleID         int     `json:"vehicle_id"`
		InitialKm         float64 `json:"initial_km"`
		ServicingInterval float64 `json:"servicing_interval_km"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.servicingService.InitializeServicingRecord(req.VehicleID, req.InitialKm, req.ServicingInterval); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Servicing record initialized successfully"})
}

// UpdateServicingStatus handles updating the servicing status based on current km reading
func (h *VehicleServicingHandler) UpdateServicingStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vehicleID, err := strconv.Atoi(vars["vehicle_id"])
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	var req struct {
		CurrentKm float64 `json:"current_km"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.servicingService.UpdateServicingStatus(vehicleID, req.CurrentKm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Servicing status updated successfully"})
}

// GetVehiclesDueForServicing handles retrieving the list of vehicles that need servicing
func (h *VehicleServicingHandler) GetVehiclesDueForServicing(c *gin.Context) {
	vehicles, err := h.servicingService.GetVehiclesDueForServicing()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Vehicles due for servicing retrieved successfully", vehicles))
}

// GetServicingHistory handles retrieving the servicing history for a vehicle
func (h *VehicleServicingHandler) GetServicingHistory(c *gin.Context) {
	vehicleID, err := strconv.Atoi(c.Param("vehicle_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid vehicle ID", nil))
		return
	}

	history, err := h.servicingService.GetServicingHistory(vehicleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Servicing history retrieved successfully", history))
}

// MarkAsServicedRequest represents the request body for marking a vehicle as serviced
type MarkAsServicedRequest struct {
	CurrentKm float64 `json:"current_km" binding:"required"`
}

// MarkAsServiced marks a vehicle as serviced
func (h *VehicleServicingHandler) MarkAsServiced(c *gin.Context) {
	vehicleID, err := strconv.Atoi(c.Param("vehicle_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid vehicle ID",
			"data":    nil,
		})
		return
	}

	var req MarkAsServicedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request body",
			"data":    nil,
		})
		return
	}

	err = h.servicingService.MarkAsServiced(vehicleID, req.CurrentKm, "", 0, "", 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to mark vehicle as serviced",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Vehicle marked as serviced successfully",
		"data": gin.H{
			"vehicle_id": vehicleID,
			"current_km": req.CurrentKm,
		},
	})
}

// GetServicingStatus handles retrieving the current servicing statusvehicalu
func (h *VehicleServicingHandler) GetServicingStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vehicleID, err := strconv.Atoi(vars["vehicle_id"])
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	status, err := h.servicingService.GetCurrentKmAndServicingStatus(vehicleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// GetCurrentKmAndServicingStatus handles retrieving the current km reading and servicing status
func (h *VehicleServicingHandler) GetCurrentKmAndServicingStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	vehicleID, err := strconv.Atoi(vars["vehicle_id"])
	if err != nil {
		http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
		return
	}

	status, err := h.servicingService.GetCurrentKmAndServicingStatus(vehicleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
