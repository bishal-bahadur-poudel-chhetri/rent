package handlers

import (
	"net/http"
	"renting/internal/services"
	"renting/internal/utils" // Import the utils package

	"github.com/gin-gonic/gin"
)

type SaleDetailHandler struct {
	saleService *services.SaleDetailService
}

func NewSaleDetailHandler(saleService *services.SaleDetailService) *SaleDetailHandler {
	return &SaleDetailHandler{saleService: saleService}
}

func (h *SaleDetailHandler) GetSalesWithFilters(c *gin.Context) {
	// Parse query parameters
	queryParams := c.Request.URL.Query()
	filters := make(map[string]string)

	for key, values := range queryParams {
		if len(values) > 0 {
			filters[key] = values[0]
		}
	}

	// Call the service
	sales, err := h.saleService.GetSalesWithFilters(filters)
	if err != nil {
		// Return error response in StandardResponse format
		response := utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch sales", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	// Return success response in StandardResponse format
	response := utils.SuccessResponse(http.StatusOK, "Sales fetched successfully", sales)
	c.JSON(http.StatusOK, response)
}
