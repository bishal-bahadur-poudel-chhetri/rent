package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"renting/internal/services"
	"renting/internal/utils"

	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	videoService services.VideoService
	videoURL     string
}

func NewVideoHandler(videoService services.VideoService, videoURL string) *VideoHandler {
	return &VideoHandler{
		videoService: videoService,
		videoURL:     videoURL,
	}
}

func (h *VideoHandler) UploadVideo(c *gin.Context) {
	startTime := time.Now()

	// Limit form size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 32<<20) // 32MB
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Form data too large or invalid", nil))
		return
	}

	// Get file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Failed to read file", nil))
		return
	}

	// Validate sale_id
	saleIDStr := c.PostForm("sale_id")
	if saleIDStr == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "sale_id is required", nil))
		return
	}
	saleID, err := strconv.Atoi(saleIDStr)
	if err != nil || saleID <= 0 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Invalid sale_id", nil))
		return
	}

	// Upload via service
	salesVideo, err := h.videoService.UploadVideo(file, saleID, h.videoURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	// Log and respond
	go log.Printf("VideoHandler: Upload for saleID %d took %s", saleID, time.Since(startTime))
	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "File uploaded successfully", salesVideo))
}
