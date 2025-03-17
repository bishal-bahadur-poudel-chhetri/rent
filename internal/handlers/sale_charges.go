package handlers

import (
	"net/http"
	"renting/internal/models"
	"renting/internal/services"
	"renting/internal/utils"

	"github.com/gin-gonic/gin"
)

type SaleChargeHandler struct {
	saleService *services.SaleChargeService
	jwtSecret   string
}

func NewSaleChargeHandler(saleService *services.SaleChargeService, jwtSecret string) *SaleChargeHandler {
	return &SaleChargeHandler{
		saleService: saleService,
		jwtSecret:   jwtSecret,
	}
}

func (h *SaleChargeHandler) AddSalesCharge(c *gin.Context) {
	var req models.SalesCharge
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	if err := h.saleService.AddSalesCharge(req.SaleID, req); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to add sales charge", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusCreated, "Sales charge added successfully", nil))
}
