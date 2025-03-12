package utils

type StandardResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(status int, message string, data interface{}) StandardResponse {

	if data == nil {
		data = struct{}{}
	}
	return StandardResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(status int, message string, data interface{}) StandardResponse {
	return StandardResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
