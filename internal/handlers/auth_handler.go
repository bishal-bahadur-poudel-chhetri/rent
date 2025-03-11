package handlers

import (
	"net/http"
	"renting/internal/models"
	"renting/internal/services"
	"renting/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		MobileNumber string `json:"mobile_number"`
		Password     string `json:"password"`
		CompanyCode  string `json:"company_code"`
	}

	// Bind the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", nil))
		return
	}

	// Call the AuthService to handle login
	token, user, err := h.authService.Login(c.Request.Context(), req.MobileNumber, req.Password, req.CompanyCode)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, err.Error(), nil))
		return
	}

	// Return success response with token and user data
	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Login successful", map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"username":      user.Username,
			"is_admin":      user.IsAdmin,
			"mobile_number": user.MobileNumber,
		},
	}))
}
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username     string `json:"username"`
		Password     string `json:"password"`
		IsAdmin      bool   `json:"is_admin"`
		CompanyCode  string `json:"company_code"`
		MobileNumber string `json:"mobile_number"`
	}

	// Bind the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", nil))
		return
	}

	// Create a new user
	user := &models.User{
		Username:     req.Username,
		Password:     req.Password,
		IsAdmin:      req.IsAdmin,
		MobileNumber: req.MobileNumber,
	}

	// Call the AuthService to handle registration
	err := h.authService.Register(c.Request.Context(), user, req.CompanyCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, utils.SuccessResponse(http.StatusCreated, "User created successfully", nil))
}
