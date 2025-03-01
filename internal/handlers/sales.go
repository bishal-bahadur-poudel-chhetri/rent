package handlers

import (
	"log"
	"net/http"
	"renting/internal/models"
	"renting/internal/services"

	"github.com/gin-gonic/gin"
)

// SaleHandler handles requests related to sales
type SaleHandler struct {
	SaleService services.SaleService
}

// RegisterSale handles the sale registration request
func (h *SaleHandler) RegisterSale(c *gin.Context) {
	var sale models.Sale

	// Bind the JSON body to the Sale struct
	if err := c.ShouldBindJSON(&sale); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Call the service layer to register the sale
	saleID, err := h.SaleService.RegisterSale(&sale)
	if err != nil {
		log.Println("Error registering sale:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register sale"})
		return
	}

	// Return the sale ID in the response
	c.JSON(http.StatusOK, gin.H{"sale_id": saleID})
}
