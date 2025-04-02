package handlers

import (
	"net/http"
	"renting/internal/models"
	"renting/internal/services"

	"github.com/gin-gonic/gin"
)

type RevenueHandler struct {
	service *services.RevenueService
}

func NewRevenueHandler(service *services.RevenueService) *RevenueHandler {
	return &RevenueHandler{service: service}
}

// GetRevenue handles GET /api/v1/revenue
func (h *RevenueHandler) GetRevenue(c *gin.Context) {
	var req models.RevenueRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.GetRevenue(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
