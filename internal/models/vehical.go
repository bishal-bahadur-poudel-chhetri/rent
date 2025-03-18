package models

// VehicleResponse is the response model for a vehicle
type VehicleResponse struct {
	VehicleID                 int                   `json:"vehicle_id"`
	VehicleTypeID             int                   `json:"vehicle_type_id"`
	VehicleName               string                `json:"vehicle_name"`
	VehicleModel              string                `json:"vehicle_model"`
	VehicleRegistrationNumber string                `json:"vehicle_registration_number"`
	IsAvailable               bool                  `json:"is_available"`
	SalesImage                string                `json:"image_name"`
	Status                    string                `json:"status"`
	FutureBookingDetails      []FutureBookingDetail `json:"future_booking_details,omitempty"`
	SaleID                    int                   `json:"sale_id"`
}

// FutureBookingDetail represents future booking details
type FutureBookingDetail struct {
	DeliveryDate string `json:"date_of_delivery"`
	NumberOfDays int    `json:"number_of_days"` // Number of days until the booking date
}

// VehicleRequest is the request model for a vehicle
type VehicleRequest struct {
	VehicleTypeID             int    `json:"vehicle_type_id"`
	VehicleName               string `json:"vehicle_name"`
	VehicleModel              string `json:"vehicle_model"`
	VehicleRegistrationNumber string `json:"vehicle_registration_number"`
	IsAvailable               bool   `json:"is_available"`
	Status                    string `json:"status"`
}

// VehicleFilter is the filter model for listing vehicles
type VehicleFilter struct {
	VehicleTypeID             string `json:"vehicle_type_id"`
	VehicleName               string `json:"vehicle_name"`
	VehicleModel              string `json:"vehicle_model"`
	VehicleRegistrationNumber string `json:"vehicle_registration_number"`
	IsAvailable               string `json:"is_available"`
	Status                    string `json:"status"`
	Limit                     int    `json:"limit"`
	Offset                    int    `json:"offset"`
}

// FileMetadata represents metadata for uploaded files
type FileMetadata struct {
	FileName string `json:"file_name"` // Name of the file
	FilePath string `json:"file_path"` // Path where the file is stored
}
