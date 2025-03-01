package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"renting/internal/models"
)

type VehicleRepository struct {
	db *sql.DB
}

func NewVehicleRepository(db *sql.DB) *VehicleRepository {
	return &VehicleRepository{
		db: db,
	}
}

// RegisterVehicle inserts a new vehicle into the database
func (r *VehicleRepository) RegisterVehicle(vehicle models.VehicleRequest) (int, error) {
	query := `
		INSERT INTO vehicles (vehicle_type_id, vehicle_name, vehicle_model, vehicle_registration_number, is_available, status)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING vehicle_id;
	`
	var vehicleID int
	err := r.db.QueryRow(query, vehicle.VehicleTypeID, vehicle.VehicleName, vehicle.VehicleModel, vehicle.VehicleRegistrationNumber, vehicle.IsAvailable, vehicle.Status).Scan(&vehicleID)
	if err != nil {
		return 0, errors.New("failed to insert vehicle")
	}
	return vehicleID, nil
}

func (r *VehicleRepository) ListVehicles(filter models.VehicleFilter) ([]models.VehicleResponse, error) {
	var vehicles []models.VehicleResponse
	var args []interface{}
	query := `
		SELECT vehicle_id, vehicle_type_id, vehicle_name, vehicle_model, vehicle_registration_number, is_available
		FROM vehicles 
		WHERE 1=1
	`

	argIndex := 1
	if filter.VehicleTypeID != "" {
		query += fmt.Sprintf(" AND vehicle_type_id = $%d", argIndex)
		args = append(args, filter.VehicleTypeID)
		argIndex++
	}
	if filter.IsAvailable != "" {
		query += fmt.Sprintf(" AND is_available = $%d", argIndex)
		args = append(args, filter.IsAvailable)
		argIndex++
	}
	if filter.VehicleName != "" {
		query += fmt.Sprintf(" AND vehicle_name ILIKE $%d", argIndex)
		args = append(args, "%"+filter.VehicleName+"%")
		argIndex++
	}

	// Apply pagination
	query += fmt.Sprintf(" ORDER BY vehicle_id LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	// Debugging: Log the query and arguments
	fmt.Printf("Executing Query: %s\n", query)
	fmt.Printf("With Args: %+v\n", args)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var vehicle models.VehicleResponse
		if err := rows.Scan(&vehicle.VehicleID, &vehicle.VehicleTypeID, &vehicle.VehicleName, &vehicle.VehicleModel, &vehicle.VehicleRegistrationNumber, &vehicle.IsAvailable); err != nil {
			return nil, fmt.Errorf("failed to scan rows: %v", err)
		}
		vehicles = append(vehicles, vehicle)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return vehicles, nil
}
