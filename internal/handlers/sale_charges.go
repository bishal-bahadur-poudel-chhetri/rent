package handlers

import (
	"fmt"
	"net/http"
	"renting/internal/models"
	"renting/internal/services"
	"renting/internal/utils"
	"strconv"

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

func (h *SaleChargeHandler) UpdateSalesCharge(c *gin.Context) {
	chargeIDStr := c.Param("chargeId")
	chargeID, err := strconv.Atoi(chargeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid charge ID", "Charge ID must be a number"))
		return
	}

	var req models.SalesCharge
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	// Debug logging
	fmt.Printf("Updating charge ID: %d, Type: %s, Amount: %f\n", chargeID, req.ChargeType, req.Amount)

	if err := h.saleService.UpdateSalesCharge(chargeID, req); err != nil {
		fmt.Printf("Error updating sales charge: %v\n", err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to update sales charge", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Sales charge updated successfully", nil))
}

func (h *SaleChargeHandler) DeleteSalesCharge(c *gin.Context) {
	chargeIDStr := c.Param("chargeId")
	chargeID, err := strconv.Atoi(chargeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid charge ID", "Charge ID must be a number"))
		return
	}

	if err := h.saleService.DeleteSalesCharge(chargeID); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to delete sales charge", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Sales charge deleted successfully", nil))
}
