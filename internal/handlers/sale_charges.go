package handlers

import (
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

func (h *SaleChargeHandler) GetSalesChargesBySaleID(c *gin.Context) {
	saleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid sale ID", "Sale ID must be a number"))
		return
	}

	charges, err := h.saleService.GetSalesChargesBySaleID(saleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to fetch sales charges", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Sales charges fetched successfully", charges))
}

func (h *SaleChargeHandler) UpdateSalesCharge(c *gin.Context) {
	chargeID, err := strconv.Atoi(c.Param("chargeID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid charge ID", "Charge ID must be a number"))
		return
	}

	var req models.SalesCharge
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	// Get saleID from the URL parameter
	saleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid sale ID", "Sale ID must be a number"))
		return
	}
	req.SaleID = saleID

	if err := h.saleService.UpdateSalesCharge(chargeID, req); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to update sales charge", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Sales charge updated successfully", nil))
}

func (h *SaleChargeHandler) DeleteSalesCharge(c *gin.Context) {
	chargeID, err := strconv.Atoi(c.Param("chargeID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid charge ID", "Charge ID must be a number"))
		return
	}

	saleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid sale ID", "Sale ID must be a number"))
		return
	}

	if err := h.saleService.DeleteSalesCharge(chargeID, saleID); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to delete sales charge", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Sales charge deleted successfully", nil))
}
