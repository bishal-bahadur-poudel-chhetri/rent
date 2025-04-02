package models

// RevenueRequest defines the request structure for revenue endpoints
type RevenueRequest struct {
	Period      string `form:"period" json:"period"`             // "day", "month", "year", "custom"
	Date        string `form:"date" json:"date"`                 // Reference date (format: 2006-01-02)
	StartDate   string `form:"start_date" json:"start_date"`     // Start date for custom period
	EndDate     string `form:"end_date" json:"end_date"`         // End date for custom period
	RecognizeAt string `form:"recognize_at" json:"recognize_at"` // "start", "end", "prorated"
}

// ErrorResponse defines a standard error response structure
type ErrorResponse struct {
	Error      string `json:"error"`
	Message    string `json:"message"`
	RequestID  string `json:"request_id,omitempty"`
	StatusCode int    `json:"status_code,omitempty"`
}
