// models/sales_video.go
package models

import "time"

type SalesVideo struct {
	VideoID    int       `json:"video_id"`
	SaleID     int       `json:"sale_id"`
	VideoURL   string    `json:"video_url"`
	FileName   string    `json:"file_name"`
	MimeType   string    `json:"mime_type"` // Add this field
	Size       int64     `json:"size"`      // Add this field
	UploadedAt time.Time `json:"uploaded_at"`
	CreatedAt  time.Time `json:"created_at"`
}
