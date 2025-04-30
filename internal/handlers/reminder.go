package handlers

import (
	"net/http"
	"renting/internal/models"
	"renting/internal/repositories"
	"renting/internal/services"
	"renting/internal/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ReminderHandler struct {
	reminderService *services.ReminderService
	userRepo        repositories.UserRepository
}

func NewReminderHandler(reminderService *services.ReminderService, userRepo repositories.UserRepository) *ReminderHandler {
	return &ReminderHandler{
		reminderService: reminderService,
		userRepo:        userRepo,
	}
}

// CheckAdminPermission middleware to verify if user is an admin
func (h *ReminderHandler) CheckAdminPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "User not authenticated", nil))
			c.Abort()
			return
		}

		// Get user from database using userID
		user, err := h.userRepo.GetUserByID(c.Request.Context(), userID.(int))
		if err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Error fetching user data", nil))
			c.Abort()
			return
		}

		if !user.IsAdmin {
			c.JSON(http.StatusForbidden, utils.ErrorResponse(http.StatusForbidden, "Admin permission required", nil))
			c.Abort()
			return
		}

		c.Next()
	}
}

// CreateReminder handles the creation of a new reminder
func (h *ReminderHandler) CreateReminder(c *gin.Context) {
	var reminder models.Reminder
	if err := c.ShouldBindJSON(&reminder); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", nil))
		return
	}

	// Get user ID from context (set by JWT middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "User ID not found in context", nil))
		return
	}
	reminder.UserID = userID.(int)

	if err := h.reminderService.CreateReminder(c.Request.Context(), &reminder); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse(http.StatusCreated, "Reminder created successfully", reminder))
}

// AcknowledgeReminder handles the acknowledgment of a reminder
func (h *ReminderHandler) AcknowledgeReminder(c *gin.Context) {
	reminderID, err := strconv.Atoi(c.Param("reminder_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid reminder ID", nil))
		return
	}

	// Get user ID from context (set by JWT middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "User ID not found in context", nil))
		return
	}

	if err := h.reminderService.AcknowledgeReminder(c.Request.Context(), reminderID, userID.(int)); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Reminder acknowledged successfully", nil))
}

// GetRemindersByVehicle handles fetching all reminders for a vehicle
func (h *ReminderHandler) GetRemindersByVehicle(c *gin.Context) {
	vehicleID, err := strconv.Atoi(c.Param("vehicle_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid vehicle ID", nil))
		return
	}

	reminders, err := h.reminderService.GetRemindersByVehicleIDWithVehicleDetails(c.Request.Context(), vehicleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Reminders retrieved successfully", reminders))
}

// GetDueReminders handles fetching all reminders that are due
func (h *ReminderHandler) GetDueReminders(c *gin.Context) {
	// Get pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Get reminder type filter if provided
	reminderType := c.Query("type")

	reminders, err := h.reminderService.GetDueRemindersWithVehicleDetails(c.Request.Context(), limit, offset, reminderType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Due reminders retrieved successfully", reminders))
}

// GetReminderHistory handles fetching acknowledgement history for a reminder
func (h *ReminderHandler) GetReminderHistory(c *gin.Context) {
	reminderID, err := strconv.Atoi(c.Param("reminder_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid reminder ID", nil))
		return
	}

	history, err := h.reminderService.GetReminderHistory(c.Request.Context(), reminderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Reminder history retrieved successfully", history))
}

// UpdateReminder handles updating an existing reminder
func (h *ReminderHandler) UpdateReminder(c *gin.Context) {
	reminderID, err := strconv.Atoi(c.Param("reminder_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid reminder ID", nil))
		return
	}

	var reminder models.Reminder
	if err := c.ShouldBindJSON(&reminder); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", nil))
		return
	}

	// Ensure the ID in the path matches the ID in the body
	reminder.ID = reminderID

	// Get user ID from context (set by JWT middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "User ID not found in context", nil))
		return
	}
	reminder.UserID = userID.(int)

	if err := h.reminderService.UpdateReminder(c.Request.Context(), &reminder); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Reminder updated successfully", reminder))
}

// DeleteReminder handles soft deleting a reminder
func (h *ReminderHandler) DeleteReminder(c *gin.Context) {
	reminderID, err := strconv.Atoi(c.Param("reminder_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid reminder ID", nil))
		return
	}

	if err := h.reminderService.DeleteReminder(c.Request.Context(), reminderID); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Reminder deleted successfully", nil))
}

// GetRemindersByVehicleAndType handles fetching all reminders for a vehicle filtered by multiple types
func (h *ReminderHandler) GetRemindersByVehicleAndType(c *gin.Context) {
	vehicleID, err := strconv.Atoi(c.Param("vehicle_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid vehicle ID", nil))
		return
	}

	// Get reminder types from query parameter
	typesStr := c.Query("types")
	var types []string
	if typesStr != "" {
		types = strings.Split(typesStr, ",")
		// Trim whitespace from each type
		for i := range types {
			types[i] = strings.TrimSpace(types[i])
		}
	}

	reminders, err := h.reminderService.GetRemindersByVehicleIDWithVehicleDetails(c.Request.Context(), vehicleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	// Filter by types if provided
	if len(types) > 0 {
		var filteredReminders []models.ReminderWithVehicle
		for _, reminder := range reminders {
			for _, t := range types {
				if string(reminder.Type) == t {
					filteredReminders = append(filteredReminders, reminder)
					break // Add reminder only once even if it matches multiple types
				}
			}
		}
		reminders = filteredReminders
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Reminders retrieved successfully", reminders))
}
