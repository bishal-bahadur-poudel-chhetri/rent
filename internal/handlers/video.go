package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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

	// Increased form size limit to 100MB
	maxSize := int64(100 << 20) // 100MB
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
	if err := c.Request.ParseMultipartForm(maxSize); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "File too large. Maximum size is 100MB", nil))
		return
	}

	// Get file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "Failed to read file: "+err.Error(), nil))
		return
	}

	// Check file size
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest, "File too large. Maximum size is 100MB", nil))
		return
	}
	contentType := file.Header.Get("Content-Type")
	log.Printf("Received Content-Type: %s, Filename: %s", contentType, file.Filename)

	// Allow only video/* or exactly application/octet-stream
	if !(strings.HasPrefix(contentType, "video/") || contentType == "application/octet-stream") {
		log.Printf("Rejected file: Filename: %s, Content-Type: %s", file.Filename, contentType)
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(http.StatusBadRequest,
			"Invalid file type. Only video files (video/* or application/octet-stream) are allowed", nil))
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
		log.Printf("Error uploading video for saleID %d: %v", saleID, err)
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse(http.StatusInternalServerError, "Failed to upload video: "+err.Error(), nil))
		return
	}

	// Log and respond
	duration := time.Since(startTime)
	uploadSpeedMBps := float64(file.Size) / (1024 * 1024) / duration.Seconds()
	log.Printf("VideoHandler: Upload for saleID %d completed. Size: %.2fMB, Duration: %s, Speed: %.2f MB/s",
		saleID, float64(file.Size)/(1024*1024), duration, uploadSpeedMBps)

	c.JSON(http.StatusOK, utils.SuccessResponse(http.StatusOK, "File uploaded successfully", gin.H{
		"video": salesVideo,
		"stats": gin.H{
			"size":     file.Size,
			"duration": duration.String(),
			"speed":    fmt.Sprintf("%.2f MB/s", uploadSpeedMBps),
		},
	}))
}

