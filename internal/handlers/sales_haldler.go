package handlers

import (
	"net/http"
	"renting/internal/models"
	"renting/internal/services"
	"renting/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SaleHandler handles sales-related requests
type SaleHandler struct {
	saleService *services.SaleService
	jwtSecret   string
}

// NewSaleHandler initializes a SaleHandler
func NewSaleHandler(saleService *services.SaleService, jwtSecret string) *SaleHandler {
	return &SaleHandler{
		saleService: saleService,
		jwtSecret:   jwtSecret,
	}
}

// CreateSale handles the creation of a new sale
func (h *SaleHandler) CreateSale(c *gin.Context) {
	// Extract user ID from the token
	userID, err := utils.ExtractUserIDFromToken(c, h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Unauthorized", err.Error()))
		return
	}

	// Parse request body
	var sale models.Sale
	if err := c.ShouldBindJSON(&sale); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	// Assign extracted user ID
	sale.UserID = userID

	// Create sale
	saleID, err := h.saleService.CreateSale(sale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to create sale", err.Error()))
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, utils.SuccessResponse(http.StatusCreated, "Sale created successfully", gin.H{"sale_id": saleID}))
}

// GetSaleByID retrieves a sale by ID
func (h *SaleHandler) GetSaleByID(c *gin.Context) {
	// Convert sale ID from string to int
	saleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid sale ID", "Sale ID must be a number"))
		return
	}

	// Fetch sale
	sale, err := h.saleService.GetSaleByID(saleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch sale", err.Error()))
		return
	}

	// Check if sale exists
	if sale == nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse(http.StatusNotFound, "Sale not found", nil))
		return
	}

	// Return success response
	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Sale fetched successfully", sale))
}
