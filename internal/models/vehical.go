package models

type VehicleFilter struct {
	VehicleTypeID             string `json:"vehicle_type_id"`
	IsAvailable               string `json:"is_available"`
	VehicleName               string `json:"vehicle_name"`
	VehicleRegistrationNumber string `json:"vehicle_registration_number"`
	Status                    string `json:"status"`
	FromDate                  string `json:"from_date"` // New field for filtering from date
	ToDate                    string `json:"to_date"`   // New field for filtering to date
	Limit                     int    `json:"limit,omitempty"`
	Offset                    int    `json:"offset,omitempty"`
}

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
