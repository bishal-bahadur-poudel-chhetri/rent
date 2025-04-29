package handlers

import (
	"net/http"
	"renting/internal/services"
	"renting/internal/utils"

	"github.com/gin-gonic/gin"
)

type SystemSettingsHandler struct {
	service *services.SystemSettingsService
}

func NewSystemSettingsHandler(service *services.SystemSettingsService) *SystemSettingsHandler {
	return &SystemSettingsHandler{service: service}
}

func (h *SystemSettingsHandler) GetSystemSettings(c *gin.Context) {
	settings, err := h.service.GetSystemSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "System settings retrieved successfully", settings))
}

func (h *SystemSettingsHandler) UpdateSystemSetting(c *gin.Context) {
	key := c.Param("key")
	if key != "enable_registration" && key != "enable_login" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid setting key", nil))
		return
	}

	var request struct {
		Value bool `json:"value"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid request body", err.Error()))
		return
	}

	if err := h.service.UpdateSystemSetting(key, request.Value); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "System setting updated successfully", nil))
}
