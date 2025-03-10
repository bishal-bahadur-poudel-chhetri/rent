package models

// VehicleResponse is the response model for a vehicle
type VehicleResponse struct {
	VehicleID                 int    `json:"vehicle_id"`
	VehicleTypeID             int    `json:"vehicle_type_id"`
	VehicleName               string `json:"vehicle_name"`
	VehicleModel              string `json:"vehicle_model"`
	VehicleRegistrationNumber string `json:"vehicle_registration_number"`
	IsAvailable               bool   `json:"is_available"`
	Status                    string `json:"status"`
}
type VehicleRequest struct {
	VehicleTypeID             int    `json:"vehicle_type_id"`
	VehicleName               string `json:"vehicle_name"`
	VehicleModel              string `json:"vehicle_model"`
	VehicleRegistrationNumber string `json:"vehicle_registration_number"`
	IsAvailable               bool   `json:"is_available"`
	Status                    string `json:"status"`
}

type VehicleFilter struct {
	VehicleTypeID             string `json:"vehicle_type_id"`
	VehicleName               string `json:"vehicle_name"`
	VehicleModel              string `json:"vehicle_model"`               // Add this line
	VehicleRegistrationNumber string `json:"vehicle_registration_number"` // Add this line
	IsAvailable               string `json:"is_available"`
	Status                    string `json:"status"`
	Limit                     int    `json:"limit"`
	Offset                    int    `json:"offset"`
}
