package handlers

import (
    "database/sql"
    "log"
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
		Username      string `json:"username"`
		Password      string `json:"password"`
		IsAdmin       bool   `json:"is_admin"`
		HasAccounting bool   `json:"has_accounting"`
		CompanyCode   string `json:"company_code"`
		MobileNumber  string `json:"mobile_number"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	user := &models.User{
		Username:      request.Username,
		Password:      request.Password,
		IsAdmin:       request.IsAdmin,
		HasAccounting: request.HasAccounting,
		MobileNumber:  request.MobileNumber,
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

// UpdateMyProfile allows the authenticated user to update username and mobile
func (h *AuthHandler) UpdateMyProfile(c *gin.Context) {
    userIDAny, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "User not authenticated", nil))
        return
    }
    var req struct {
        Username     string `json:"username"`
        MobileNumber string `json:"mobile_number"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
        return
    }
    if err := h.authService.UpdateProfile(c.Request.Context(), userIDAny.(int), req.Username, req.MobileNumber); err != nil {
        c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
        return
    }
    c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Profile updated", nil))
}

// ChangeMyPassword allows the authenticated user to change password
func (h *AuthHandler) ChangeMyPassword(c *gin.Context) {
    userIDAny, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "User not authenticated", nil))
        return
    }
    var req struct {
        CurrentPassword string `json:"current_password"`
        NewPassword     string `json:"new_password"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
        return
    }
    if err := h.authService.ChangePassword(c.Request.Context(), userIDAny.(int), req.CurrentPassword, req.NewPassword); err != nil {
        c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
        return
    }
    c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Password changed", nil))
}

// CheckAccountingPermission middleware to verify if user has accounting permission
func (h *AuthHandler) CheckAccountingPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[DEBUG] ============ START CheckAccountingPermission ============")
		log.Printf("[DEBUG] Request path: %s %s", c.Request.Method, c.Request.URL.Path)

		// Log all context values for debugging
		for k, v := range c.Keys {
			log.Printf("[DEBUG] Context key: %s, value type: %T", k, v)
		}

		// Get the current user's ID from context (set by JWT middleware)
		userID, exists := c.Get("userID")
		if !exists {
			log.Printf("[ERROR] User ID not found in context")
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "User not authenticated", nil))
			c.Abort()
			return
		}

		log.Printf("[DEBUG] Using user ID: %d from context", userID.(int))

		// Get user data
		user, err := h.authService.GetUserByID(c.Request.Context(), userID.(int))
		if err != nil {
			log.Printf("[ERROR] Error fetching user data: %v", err)
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Error fetching user data", nil))
			c.Abort()
			return
		}

		if user == nil {
			log.Printf("[ERROR] User not found with ID: %d", userID.(int))
			c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "User not found", nil))
			c.Abort()
			return
		}

		log.Printf("[DEBUG] User ID: %d | Username: %s | IsAdmin: %v | HasAccounting: %v",
			user.ID, user.Username, user.IsAdmin, user.HasAccounting)

		// Check if the user has accounting permission
		if !user.HasAccounting {
			log.Printf("[ERROR] ACCESS DENIED - User %d has no accounting permission (IsAdmin: %v, HasAccounting: %v)",
				userID.(int), user.IsAdmin, user.HasAccounting)
			c.JSON(http.StatusForbidden, utils.ErrorResponse(http.StatusForbidden, "User does not have accounting permission", nil))
			c.Abort()
			return
		}

		log.Printf("[DEBUG] ✅ ACCESS GRANTED - User %d has accounting permission (IsAdmin: %v, HasAccounting: %v)",
			userID.(int), user.IsAdmin, user.HasAccounting)
		log.Printf("[DEBUG] ============ END CheckAccountingPermission ============")
		c.Next()
	}
}

// CheckUserPermissions is a debug handler to verify a user's permissions
func (h *AuthHandler) CheckUserPermissions(c *gin.Context) {
	// Get the current user's ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "User not authenticated", nil))
		return
	}

    // Get user data
    user, err := h.authService.GetUserByID(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Error fetching user data", nil))
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "User not found", nil))
		return
	}

    // Optionally enrich with company details
    var companyName string
    var companyCode string
    if dbAny, ok := c.Get("sqlDB"); ok {
        if db, ok := dbAny.(*sql.DB); ok {
            row := db.QueryRowContext(c.Request.Context(),
                "SELECT name, company_code FROM companies WHERE id = $1",
                user.CompanyID,
            )
            _ = row.Scan(&companyName, &companyCode)
        }
    }

    // Return all user permissions for debugging (enriched)
    permissionData := map[string]interface{}{
        "user_id":        user.ID,
        "username":       user.Username,
        "mobile_number":  user.MobileNumber,
        "is_admin":       user.IsAdmin,
        "has_accounting": user.HasAccounting,
        "is_locked":      user.IsLocked,
        "company_id":     user.CompanyID,
        "company_name":   companyName,
        "company_code":   companyCode,
    }

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "User permissions retrieved", permissionData))
}
