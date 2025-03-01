package handlers

import (
	"net/http"
	"renting/internal/models"
	"renting/internal/services"
	"renting/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request payload", err.Error()))
		return
	}

	userID, token, err := h.authService.RegisterUser(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrUserExists {
			status = http.StatusConflict
		}
		c.JSON(status, utils.ErrorResponse(status, err.Error(), nil))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse(http.StatusCreated, "Registration successful", models.TokenResponse{
		Token:     token,
		UserID:    userID,
		Username:  req.Username,
		IsAdmin:   req.IsAdmin,
		CompanyID: req.CompanyID,
	}))
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request payload", err.Error()))
		return
	}

	userID, token, user, err := h.authService.LoginUser(req)
	if err != nil {
		status := http.StatusInternalServerError
		var errorMessage string
		var errorData interface{}

		// Switch case for handling different errors
		switch err {
		case services.ErrInvalidCredentials:
			status = http.StatusUnauthorized
			errorMessage = "Invalid credentials. Please check your mobile number and company code."
			errorData = gin.H{"help": "Ensure the mobile number and company code match an existing user"}
		case services.ErrCompanyNotFound:
			status = http.StatusNotFound
			errorMessage = "Company not found. Please check the company code."
			errorData = gin.H{"help": "Verify the company code for the user"}
		default:
			errorMessage = "An unexpected error occurred."
			errorData = gin.H{"help": "Please try again later"}
		}

		// Responding with error response in standardized format
		c.JSON(status, utils.ErrorResponse(status, errorMessage, errorData))
		return
	}

	// Standardized success response format
	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Login successful", models.TokenResponse{
		Token:     token,
		UserID:    userID,
		Username:  user.Username,
		IsAdmin:   user.IsAdmin,
		CompanyID: user.CompanyID,
	}))
}

// GetProfile returns user profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	user, err := h.authService.GetUserByID(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to retrieve profile", nil))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Profile retrieved successfully", user))
}

// UpdateProfile updates user profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request data", err.Error()))
		return
	}

	err := h.authService.UpdateUserProfile(userID.(int), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to update profile", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Profile updated successfully", nil))
}
