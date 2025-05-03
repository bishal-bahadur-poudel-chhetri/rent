package handlers

import (
	"net/http"
	"renting/internal/models"
	"renting/internal/services"
	"renting/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService           services.AuthService
	systemSettingsService *services.SystemSettingsService
}

func NewAuthHandler(authService services.AuthService, systemSettingsService *services.SystemSettingsService) *AuthHandler {
	return &AuthHandler{
		authService:           authService,
		systemSettingsService: systemSettingsService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	// Check if registration is enabled
	settings, err := h.systemSettingsService.GetSystemSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to check system settings", nil))
		return
	}

	if !settings.EnableRegistration {
		c.JSON(http.StatusForbidden, utils.ErrorResponse(http.StatusForbidden, "Registration is currently disabled", nil))
		return
	}

	var request struct {
		Username     string `json:"username"`
		Password     string `json:"password"`
		IsAdmin      bool   `json:"is_admin"`
		CompanyCode  string `json:"company_code"`
		MobileNumber string `json:"mobile_number"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	user := &models.User{
		Username:     request.Username,
		Password:     request.Password,
		IsAdmin:      request.IsAdmin,
		MobileNumber: request.MobileNumber,
	}

	if err := h.authService.Register(c.Request.Context(), user, request.CompanyCode); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse(http.StatusCreated, "User registered successfully", nil))
}

func (h *AuthHandler) Login(c *gin.Context) {
	// Check if login is enabled
	settings, err := h.systemSettingsService.GetSystemSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to check system settings", nil))
		return
	}

	if !settings.EnableLogin {
		c.JSON(http.StatusForbidden, utils.ErrorResponse(http.StatusForbidden, "Login is currently disabled", nil))
		return
	}

	var request struct {
		MobileNumber string `json:"mobile_number"`
		Password     string `json:"password"`
		CompanyCode  string `json:"company_code"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	token, user, err := h.authService.Login(c.Request.Context(), request.MobileNumber, request.Password, request.CompanyCode)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Invalid credentials", nil))
		return
	}

	// Prepare user data for response
	userData := map[string]interface{}{
		"username":      user.Username,
		"is_admin":      user.IsAdmin,
		"mobile_number": user.MobileNumber,
	}
	responseData := map[string]interface{}{
		"token": token,
		"user":  userData,
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Login successful", responseData))
}

// LockoutUser handles locking out a user's login access
func (h *AuthHandler) LockoutUser(c *gin.Context) {
	// Get user ID from path parameter
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid user ID", nil))
		return
	}

	// Get the current user's ID from context (set by JWT middleware)
	currentUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "User not authenticated", nil))
		return
	}

	// Allow self-lockout, but require admin permission for locking out others
	if userID != currentUserID.(int) {
		// Check if current user is admin
		currentUser, err := h.authService.GetUserByID(c.Request.Context(), currentUserID.(int))
		if err != nil || currentUser == nil || !currentUser.IsAdmin {
			c.JSON(http.StatusForbidden, utils.ErrorResponse(http.StatusForbidden, "Admin permission required to lock out other users", nil))
			return
		}
	}

	// Lock out the user
	if err := h.authService.LockoutUser(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "User locked out successfully", nil))
}

// DeleteAccount handles the deletion (soft-delete) of a user's own account
func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	// Get the current user's ID from context (set by JWT middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "User not authenticated", nil))
		return
	}

	// Lock out the user (soft delete)
	if err := h.authService.LockoutUser(c.Request.Context(), userID.(int)); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Account deleted successfully", nil))
}
