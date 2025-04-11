package repositories

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"renting/internal/config"
	"renting/internal/models"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// VideoRepository defines the interface for video upload operations
type VideoRepository interface {
	UploadFile(filePath string, contentType string, saleID int, videoURL string) (*models.SalesVideo, error)
	UploadFileFromReader(reader io.Reader, size int64, contentType string, saleID int, videoURL string) (*models.SalesVideo, error)
	GetUploadStats() UploadStats
}

type UploadStats struct {
	TotalUploads   int64
	TotalBytes     int64
	AverageSpeedMB float64
	mu             sync.RWMutex
}

type videoRepository struct {
	db       *sql.DB
	config   *config.Config
	stats    UploadStats
	uploader *s3manager.Uploader
	session  *session.Session
}

// NewVideoRepository creates a new instance of videoRepository with optimized settings
func NewVideoRepository(db *sql.DB, config *config.Config) VideoRepository {
	// Initialize AWS session with optimized settings
	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint:         aws.String(config.R2Endpoint),
		Region:           aws.String("auto"),
		Credentials:      credentials.NewStaticCredentials(config.R2AccessKeyID, config.R2SecretAccessKey, ""),
		S3ForcePathStyle: aws.Bool(true),
		HTTPClient: &http.Client{
			Timeout: 1 * time.Minute,
			Transport: &http.Transport{
				MaxIdleConns:          200, // Increased from 100
				MaxIdleConnsPerHost:   200, // Increased from 100
				IdleConnTimeout:       30 * time.Second,
				DisableCompression:    true,             // Videos are already compressed
				MaxConnsPerHost:       200,              // Added to allow more connections per host
				ResponseHeaderTimeout: 30 * time.Second, // Added timeout for response headers
				ExpectContinueTimeout: 1 * time.Second,  // Added timeout for expect continue
				TLSHandshakeTimeout:   10 * time.Second, // Added timeout for TLS handshake
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
			},
		},
	}))

	// Create uploader with optimized settings
	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = 20 * 1024 * 1024 // 20MB per part for faster uploads
		u.Concurrency = 20            // Increased from 10 to 20 concurrent uploads
		u.LeavePartsOnError = false   // Clean up on failure
		u.MaxUploadParts = 10000      // Maximum number of parts for very large files
	})

	return &videoRepository{
		db:       db,
		config:   config,
		uploader: uploader,
		session:  sess,
	}
}

// UploadFileFromReader handles streaming uploads with optimized performance
func (r *videoRepository) UploadFileFromReader(reader io.Reader, size int64, contentType string, saleID int, videoURL string) (*models.SalesVideo, error) {
	startTime := time.Now()
	log.Printf("Starting upload for saleID: %d, size: %.2f MB", saleID, float64(size)/(1024*1024))

	// Generate filename with timestamp and proper extension
	filename := generateFilename(saleID, contentType)

	// Track upload speed
	speedTracker := newSpeedTracker(size)
	wrappedReader := io.TeeReader(reader, speedTracker)

	// Upload to R2 with context for timeout control
	_, err := r.uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(r.config.R2BucketName),
		Key:         aws.String(filename),
		Body:        wrappedReader,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		log.Printf("Upload failed for saleID %d: %v", saleID, err)
		return nil, fmt.Errorf("R2 upload failed: %w", err)
	}

	// Calculate actual upload speed
	uploadDuration := time.Since(startTime)
	uploadSpeedMB := float64(size) / (1024 * 1024) / uploadDuration.Seconds()

	// Update statistics
	r.stats.mu.Lock()
	r.stats.TotalUploads++
	r.stats.TotalBytes += size
	r.stats.AverageSpeedMB = (r.stats.AverageSpeedMB*float64(r.stats.TotalUploads-1) + uploadSpeedMB) / float64(r.stats.TotalUploads)
	r.stats.mu.Unlock()

	// Save to database in a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	fullURL := fmt.Sprintf("%s/%s", videoURL, filename)
	salesVideo := &models.SalesVideo{
		SaleID:     saleID,
		VideoURL:   fullURL,
		FileName:   &filename,
		MimeType:   contentType,
		Size:       size,
		UploadedAt: time.Now(),
	}

	err = tx.QueryRow(
		`INSERT INTO sales_videos (sale_id, video_url, file_name, mime_type, size, uploaded_at)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING video_id`,
		salesVideo.SaleID,
		salesVideo.VideoURL,
		salesVideo.FileName,
		salesVideo.MimeType,
		salesVideo.Size,
		salesVideo.UploadedAt,
	).Scan(&salesVideo.VideoID)

	if err != nil {
		return nil, fmt.Errorf("database insert failed: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	log.Printf("Upload completed in %s (%.2f MB/s) for saleID: %d",
		uploadDuration, uploadSpeedMB, saleID)

	return salesVideo, nil
}

// UploadFile handles file path based uploads
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

// GetUploadStats returns upload statistics
func (r *videoRepository) GetUploadStats() UploadStats {
	r.stats.mu.RLock()
	defer r.stats.mu.RUnlock()
	return r.stats
}

// Helper function to generate consistent filenames
func generateFilename(saleID int, contentType string) string {
	timestamp := time.Now().Format("20060102_150405")
	fileExt := ".mp4"
	if strings.Contains(contentType, "quicktime") {
		fileExt = ".mov"
	}
	return fmt.Sprintf("ios_%d_%s%s", saleID, timestamp, fileExt)
}

// speedTracker tracks upload progress for speed calculation
type speedTracker struct {
	totalBytes int64
	bytesRead  int64
	lastTime   time.Time
	lastBytes  int64
	mu         sync.Mutex
}

func newSpeedTracker(totalBytes int64) *speedTracker {
	return &speedTracker{
		totalBytes: totalBytes,
		lastTime:   time.Now(),
	}
}

func (st *speedTracker) Write(p []byte) (int, error) {
	n := len(p)
	st.mu.Lock()
	st.bytesRead += int64(n)
	now := time.Now()
	elapsed := now.Sub(st.lastTime).Seconds()

	if elapsed > 1.0 { // Update speed every second
		bytesSinceLast := st.bytesRead - st.lastBytes
		speed := float64(bytesSinceLast) / (1024 * 1024) / elapsed
		st.lastTime = now
		st.lastBytes = st.bytesRead
		log.Printf("Upload progress: %.1f%%, Current speed: %.2f MB/s",
			float64(st.bytesRead)/float64(st.totalBytes)*100, speed)
	}
	st.mu.Unlock()
	return n, nil
}
