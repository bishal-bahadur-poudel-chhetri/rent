package utils

// StandardResponse defines a standard response structure
type StandardResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // Optional field to send extra data with error response
}

// SuccessResponse creates a dynamic success response
func SuccessResponse(status int, message string, data interface{}) StandardResponse {
	// If no data is passed, set Data to nil
	if data == nil {
		data = struct{}{}
	}
	return StandardResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse creates a dynamic error response
func ErrorResponse(status int, message string, data interface{}) StandardResponse {
	return StandardResponse{
		Status:  status,
		Message: message,
		Data:    data, // Send extra data if provided
	}
}
