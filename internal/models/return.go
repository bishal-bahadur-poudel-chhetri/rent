package models

type ReturnRequest struct {
	SalesCharges []SalesCharge  `json:"sales_charges"`
	VehicleUsage []VehicleUsage `json:"vehicle_usage"`
	Payments     []Payment      `json:"payments"`
}
