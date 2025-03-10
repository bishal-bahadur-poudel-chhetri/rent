package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"renting/internal/models"
	"strings"
)

type VehicleRepository struct {
	db *sql.DB
}

func NewVehicleRepository(db *sql.DB) *VehicleRepository {
	return &VehicleRepository{
		db: db,
	}
}

func (r *VehicleRepository) RegisterVehicle(vehicle models.VehicleRequest) (int, error) {
	// Debug: Print the input vehicle data
	fmt.Printf("RegisterVehicle - Input: %+v\n", vehicle)

	// Step 1: Check for duplicate vehicle (e.g., based on vehicle_registration_number)
	exists, err := r.checkDuplicateVehicle(vehicle.VehicleRegistrationNumber)
	if err != nil {
		fmt.Printf("RegisterVehicle - Error checking duplicate vehicle: %v\n", err)
		return 0, errors.New("failed to check for duplicate vehicle")
	}
	if exists {
		fmt.Printf("RegisterVehicle - Duplicate vehicle found for registration number: %s\n", vehicle.VehicleRegistrationNumber)
		return 0, errors.New("vehicle with this registration number already exists")
	}

	// Step 2: Validate other conditions (e.g., status, vehicle_type_id, etc.)
	if err := r.validateVehicle(vehicle); err != nil {
		fmt.Printf("RegisterVehicle - Validation error: %v\n", err)
		return 0, err
	}

	// Step 3: Insert the vehicle
	query := `
        INSERT INTO vehicles (vehicle_type_id, vehicle_name, vehicle_model, vehicle_registration_number, is_available, status)
        VALUES ($1, $2, $3, $4, $5, $6) RETURNING vehicle_id;
    `

	// Debug: Print the SQL query and parameters
	fmt.Printf("RegisterVehicle - Query: %s\n", query)
	fmt.Printf("RegisterVehicle - Params: %d, %s, %s, %s, %t, %s\n",
		vehicle.VehicleTypeID, vehicle.VehicleName, vehicle.VehicleModel,
		vehicle.VehicleRegistrationNumber, vehicle.IsAvailable, vehicle.Status)

	var vehicleID int
	err = r.db.QueryRow(query, vehicle.VehicleTypeID, vehicle.VehicleName, vehicle.VehicleModel,
		vehicle.VehicleRegistrationNumber, vehicle.IsAvailable, vehicle.Status).Scan(&vehicleID)

	if err != nil {
		// Debug: Print the error
		fmt.Printf("RegisterVehicle - Error: %v\n", err)
		return 0, errors.New("failed to insert vehicle")
	}

	// Debug: Print the generated vehicle ID
	fmt.Printf("RegisterVehicle - Inserted Vehicle ID: %d\n", vehicleID)

	return vehicleID, nil
}

func (r *VehicleRepository) checkDuplicateVehicle(registrationNumber string) (bool, error) {
	query := `
        SELECT COUNT(*) FROM vehicles WHERE vehicle_registration_number = $1;
    `

	var count int
	err := r.db.QueryRow(query, registrationNumber).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
func (r *VehicleRepository) checkVehicleTypeExists(vehicleTypeID int) (bool, error) {
	query := `
        SELECT COUNT(*) FROM vehicle_types WHERE vehicle_type_id = $1;
    `

	var count int
	err := r.db.QueryRow(query, vehicleTypeID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *VehicleRepository) validateVehicle(vehicle models.VehicleRequest) error {
	// Example: Validate status
	validStatuses := map[string]bool{
		"available":         true,
		"rented":            true,
		"under_maintenance": true,
	}
	if !validStatuses[vehicle.Status] {
		return errors.New("invalid status")
	}

	// Example: Validate vehicle_type_id exists in the vehicle_types table
	exists, err := r.checkVehicleTypeExists(vehicle.VehicleTypeID)
	if err != nil {
		return fmt.Errorf("failed to validate vehicle type: %v", err)
	}
	if !exists {
		return errors.New("invalid vehicle type ID")
	}

	return nil
}

func (r *VehicleRepository) ListVehicles(filter models.VehicleFilter) ([]models.VehicleResponse, error) {
	var vehicles []models.VehicleResponse
	var args []interface{}
	var queryBuilder strings.Builder

	queryBuilder.WriteString(`
		SELECT vehicle_id, vehicle_type_id, vehicle_name, vehicle_model, vehicle_registration_number, is_available, status
		FROM vehicles
		WHERE 1=1
	`)

	fmt.Printf("Filter Used: %+v\n", filter)

	argIndex := 1

	// Add filters dynamically
	if filter.VehicleTypeID != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND vehicle_type_id = $%d", argIndex))
		args = append(args, filter.VehicleTypeID)
		argIndex++
	}
	if filter.IsAvailable != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND is_available = $%d", argIndex))
		args = append(args, filter.IsAvailable)
		argIndex++
	}
	if filter.VehicleName != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND vehicle_name ILIKE $%d", argIndex))
		args = append(args, "%"+filter.VehicleName+"%")
		argIndex++
	}
	if filter.VehicleModel != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND vehicle_model ILIKE $%d", argIndex))
		args = append(args, "%"+filter.VehicleModel+"%")
		argIndex++
	}
	if filter.VehicleRegistrationNumber != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND vehicle_registration_number ILIKE $%d", argIndex))
		args = append(args, "%"+filter.VehicleRegistrationNumber+"%")
		argIndex++
	}
	if filter.Status != "" {
		queryBuilder.WriteString(fmt.Sprintf(" AND status = $%d", argIndex))
		args = append(args, filter.Status)
		argIndex++
	}

	// Apply ordering and pagination
	queryBuilder.WriteString(fmt.Sprintf(" ORDER BY vehicle_id LIMIT $%d OFFSET $%d", argIndex+1, argIndex))
	args = append(args, filter.Offset, filter.Limit) // Swap the order of LIMIT and OFFSET

	query := queryBuilder.String()
	fmt.Printf("Executing Query: %s\n", query)
	fmt.Printf("With Args: %+v\n", args)

	// Execute Query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Parse results
	for rows.Next() {
		var vehicle models.VehicleResponse
		if err := rows.Scan(
			&vehicle.VehicleID, &vehicle.VehicleTypeID, &vehicle.VehicleName,
			&vehicle.VehicleModel, &vehicle.VehicleRegistrationNumber,
			&vehicle.IsAvailable, &vehicle.Status,
		); err != nil {
			return nil, fmt.Errorf("failed to scan rows: %w", err)
		}
		vehicles = append(vehicles, vehicle)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return vehicles, nil
}
