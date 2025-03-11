package handlers

import (
	"net/http"
	"strconv"

	"renting/internal/models"
	"renting/internal/services"
	"renting/internal/utils"

	"github.com/gin-gonic/gin"
)

type SaleHandler struct {
	saleService *services.SaleService
}

func NewSaleHandler(saleService *services.SaleService) *SaleHandler {
	return &SaleHandler{saleService: saleService}
}

func (h *SaleHandler) CreateSale(c *gin.Context) {
	var sale models.Sale
	if err := c.ShouldBindJSON(&sale); err != nil {
		// Return error response in standard format
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	saleID, err := h.saleService.CreateSale(sale)
	if err != nil {
		// Return error response in standard format
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to create sale", err.Error()))
		return
	}

	// Return success response in standard format
	c.JSON(http.StatusCreated, utils.SuccessResponse(http.StatusCreated, "Sale created successfully", gin.H{"sale_id": saleID}))
}

func (h *SaleHandler) GetSaleByID(c *gin.Context) {
	saleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// Return error response in standard format
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid sale ID", nil))
		return
	}

	sale, err := h.saleService.GetSaleByID(saleID)
	if err != nil {
		// Return error response in standard format
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch sale", err.Error()))
		return
	}

	// Return success response in standard format
	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Sale fetched successfully", sale))
}
