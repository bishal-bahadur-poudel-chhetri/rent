// models/sales_video.go
package models

import "time"

type SalesVideo struct {
	VideoID    int       `json:"video_id"`    // Primary key, assuming an integer ID
	SaleID     int       `json:"sale_id"`     // Foreign key linking to the sales table
	VideoURL   string    `json:"video_url"`   // URL or path to the video file
	UploadedAt time.Time `json:"uploaded_at"` // Timestamp of when the video was uploaded
	FileName   *string   `json:"file_name"`   // Name of the video file
	MimeType   string    `json:"mime_type"`   // MIME type of the video (e.g., "video/mp4")
	Size       int64     `json:"size"`        // Size of the video file in bytes, using int64 for large files
}
