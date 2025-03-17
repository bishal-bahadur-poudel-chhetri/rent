package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"renting/internal/config"
	"renting/internal/models"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type VideoRepository interface {
	UploadFile(filePath string, contentType string, saleID int, videoURL string) (*models.SalesVideo, error)
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

func (r *videoRepository) UploadFile(filePath string, contentType string, saleID int, videoURL string) (*models.SalesVideo, error) {

	log.Printf("UploadFile called with filePath: %s, contentType: %s, saleID: %d, videoURL: %s", filePath, contentType, saleID, videoURL)

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	log.Printf("File opened successfully: %s", file.Name())

	r2Endpoint := r.config.R2Endpoint
	r2AccessKeyID := r.config.R2AccessKeyID
	r2SecretAccessKey := r.config.R2SecretAccessKey
	r2BucketName := r.config.R2BucketName

	log.Printf("R2 Endpoint: %s", r2Endpoint)
	log.Printf("R2 Access Key ID: %s", r2AccessKeyID)
	log.Printf("R2 Secret Access Key: %s", r2SecretAccessKey)
	log.Printf("R2 Bucket Name: %s", r2BucketName)

	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(r2Endpoint),
		Region:      aws.String("auto"),
		Credentials: credentials.NewStaticCredentials(r2AccessKeyID, r2SecretAccessKey, ""),
	})
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	today := time.Now().Format("2006_01_02")
	modifiedFileName := fmt.Sprintf("sd_%d_%s%s", saleID, today, filepath.Ext(file.Name()))

	uploader := s3manager.NewUploader(sess)

	uploadOutput, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(r2BucketName),
		Key:         aws.String(modifiedFileName),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	log.Printf("File uploaded successfully to R2: %+v", uploadOutput)

	query := `
		INSERT INTO sales_videos (sale_id, video_url, file_name)
		VALUES ($1, $2, $3)
		RETURNING video_id, uploaded_at
	`
	var salesVideo models.SalesVideo
	err = r.db.QueryRow(query, saleID, videoURL, modifiedFileName).Scan(&salesVideo.VideoID, &salesVideo.UploadedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert video metadata into database: %w", err)
	}

	salesVideo.SaleID = saleID
	salesVideo.VideoURL = videoURL
	salesVideo.FileName = modifiedFileName

	log.Printf("SalesVideo object populated: %+v", salesVideo)

	return &salesVideo, nil
}
