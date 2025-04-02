package handlers

import (
	"net/http"
	"renting/internal/services"
	"renting/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type DataAggregateHandler struct {
	DataAggregateService *services.DataAggregateService
}

func NewDataAggregateHandler(dataAggregateService *services.DataAggregateService) *DataAggregateHandler {
	return &DataAggregateHandler{DataAggregateService: dataAggregateService}
}

func (h *DataAggregateHandler) GetAggregatedData(c *gin.Context) {
	date := c.Query("date")
	year := c.Query("year")
	month := c.Query("month")

	var aggregatedData services.AggregatedData
	var err error

	if date != "" {
		parsedTime, parseErr := time.Parse("2006-01-02", date)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid date format", nil))
			return
		}
		aggregatedData, err = h.DataAggregateService.GetAggregatedData(parsedTime, "date")
	} else if year != "" {
		yearInt, parseErr := strconv.Atoi(year)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid year format", nil))
			return
		}
		parsedTime := time.Date(yearInt, 1, 1, 0, 0, 0, 0, time.UTC)
		aggregatedData, err = h.DataAggregateService.GetAggregatedData(parsedTime, "year")
	} else if month != "" {
		parsedTime, parseErr := time.Parse("2006-01", month)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid month format", nil))
			return
		}
		aggregatedData, err = h.DataAggregateService.GetAggregatedData(parsedTime, "month")
	} else {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Please provide a date, year, or month", nil))
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "Data retrieved successfully", aggregatedData))
}
