package services

import (
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"renting/internal/models"
	"renting/internal/repositories"
)

type VideoService interface {
	UploadVideo(file *multipart.FileHeader, saleID int, videoURL string) (*models.SalesVideo, error)
}

type videoService struct {
	videoRepo repositories.VideoRepository
}

func NewVideoService(videoRepo repositories.VideoRepository) VideoService {
	return &videoService{videoRepo: videoRepo}
}

func (s *videoService) UploadVideo(file *multipart.FileHeader, saleID int, videoURL string) (*models.SalesVideo, error) {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Create a temporary file to store the uploaded video
	tmpFile, err := os.CreateTemp("", "video-*.mp4")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tmpFile.Close()

	// Copy the uploaded file content to the temporary file
	if _, err := io.Copy(tmpFile, src); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Determine the MIME type of the file
	contentType := mime.TypeByExtension(".mp4")
	if contentType == "" {
		contentType = "video/mp4" // Fallback to video/mp4 if MIME type detection fails
	}

	// Upload the file to the repository (e.g., Cloudflare R2)
	salesVideo, err := s.videoRepo.UploadFile(tmpFile.Name(), contentType, saleID, videoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return salesVideo, nil
}
