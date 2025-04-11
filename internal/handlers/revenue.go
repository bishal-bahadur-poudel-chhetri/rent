package handlers

import (
	"net/http"
	"renting/internal/models"
	"renting/internal/services"
	"time"

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

// GetMonthlyRevenue handles GET /api/v1/revenue/monthly
func (h *RevenueHandler) GetMonthlyRevenue(c *gin.Context) {
	// Parse query parameters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Parse start date
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format. Use YYYY-MM-DD"})
		return
	}

	// Parse end date
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format. Use YYYY-MM-DD"})
		return
	}

	// Validate date range
	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date must be after start_date"})
		return
	}

	// Get monthly revenue data
	monthlyData, err := h.service.GetMonthlyRevenue(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": monthlyData,
	})
}

// GetMobileRevenueVisualization handles GET /api/v1/revenue/mobile-visualization
func (h *RevenueHandler) GetMobileRevenueVisualization(c *gin.Context) {
	// Parse query parameters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Parse start date
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format. Use YYYY-MM-DD"})
		return
	}

	// Parse end date
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format. Use YYYY-MM-DD"})
		return
	}

	// Validate date range
	if endDate.Before(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date must be after start_date"})
		return
	}

	// Get formatted revenue data for mobile visualization
	visualizationData, err := h.service.GetMobileRevenueVisualization(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, visualizationData)
}
