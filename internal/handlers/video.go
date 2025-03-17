package handlers

import (
	"net/http"
	"strconv"

	"renting/internal/services"
	"renting/internal/utils" // Import the utils package

	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	videoService services.VideoService
}

func NewVideoHandler(videoService services.VideoService) *VideoHandler {
	return &VideoHandler{videoService: videoService}
}

func (h *VideoHandler) UploadVideo(c *gin.Context) {
	// Parse the form data to get the file
	file, err := c.FormFile("file")
	if err != nil {
		// Use ErrorResponse for error handling
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Failed to read file", nil))
		return
	}

	// Parse additional fields from the form data
	saleIDStr := c.PostForm("sale_id")
	saleID, err := strconv.Atoi(saleIDStr)
	if err != nil {
		// Use ErrorResponse for error handling
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid sale_id", nil))
		return
	}

	videoURL := "https://pub-8da91f66939f4cdc9e4206024a0e68e9.r2.dev"

	// Upload the file using the service layer
	salesVideo, err := h.videoService.UploadVideo(file, saleID, videoURL)
	if err != nil {
		// Use ErrorResponse for error handling
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	// Use SuccessResponse for success handling
	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "File uploaded successfully", salesVideo))
}
