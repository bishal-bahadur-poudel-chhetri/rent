package services

import (
	"fmt"
	"mime"
	"mime/multipart"
	"path/filepath"
	"renting/internal/models"
	"renting/internal/repositories"
	"strings"
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

	// Determine content type
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		contentType = mime.TypeByExtension(ext)
		if contentType == "" {
			// Common iOS video types
			switch ext {
			case ".mov":
				contentType = "video/quicktime"
			case ".mp4":
				contentType = "video/mp4"
			default:
				contentType = "application/octet-stream"
			}
		}
	}

	// Upload directly from the reader
	salesVideo, err := s.videoRepo.UploadFileFromReader(src, file.Size, contentType, saleID, videoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return salesVideo, nil
}
