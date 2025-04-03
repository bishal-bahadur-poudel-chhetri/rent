package repositories

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"renting/internal/config"
	"renting/internal/models"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type VideoRepository interface {
	UploadFile(filePath string, contentType string, saleID int, videoURL string) (*models.SalesVideo, error)
	UploadFileFromReader(reader io.Reader, size int64, contentType string, saleID int, videoURL string) (*models.SalesVideo, error)
}

type videoRepository struct {
	db     *sql.DB
	config *config.Config
}

func NewVideoRepository(db *sql.DB, config *config.Config) VideoRepository {
	return &videoRepository{
		db:     db,
		config: config,
	}
}

func (r *videoRepository) UploadFileFromReader(reader io.Reader, size int64, contentType string, saleID int, videoURL string) (*models.SalesVideo, error) {
	startTime := time.Now()
	log.Printf("Starting upload for saleID: %d, size: %d bytes", saleID, size)

	// Configure AWS session for R2
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(r.config.R2Endpoint),
		Region:      aws.String("auto"),
		Credentials: credentials.NewStaticCredentials(r.config.R2AccessKeyID, r.config.R2SecretAccessKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	// Generate filename
	timestamp := time.Now().Format("20060102_150405")
	fileExt := ".mp4"
	if strings.Contains(contentType, "quicktime") {
		fileExt = ".mov"
	}
	filename := fmt.Sprintf("ios_%d_%s%s", saleID, timestamp, fileExt)

	// Upload to R2
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(r.config.R2BucketName),
		Key:         aws.String(filename),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to R2: %w", err)
	}

	// Save to database
	query := `
        INSERT INTO sales_videos (sale_id, video_url, file_name, mime_type, size)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING video_id, uploaded_at
    `
	fullURL := fmt.Sprintf("%s/%s", videoURL, filename)
	var salesVideo models.SalesVideo
	var fileNamePtr *string = &filename // Convert string to *string for nullable field
	err = r.db.QueryRow(
		query,
		saleID,
		fullURL,
		fileNamePtr, // Use pointer for nullable field
		contentType,
		size,
	).Scan(&salesVideo.VideoID, &salesVideo.UploadedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert video metadata into database: %w", err)
	}

	// Set all fields
	salesVideo.SaleID = saleID
	salesVideo.VideoURL = fullURL
	salesVideo.FileName = fileNamePtr // Already a *string
	salesVideo.MimeType = contentType
	salesVideo.Size = size

	log.Printf("Upload completed in %s", time.Since(startTime))
	return &salesVideo, nil
}

// Keep the original implementation for backward compatibility
func (r *videoRepository) UploadFile(filePath string, contentType string, saleID int, videoURL string) (*models.SalesVideo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return r.UploadFileFromReader(file, fileInfo.Size(), contentType, saleID, videoURL)
}
