package repositories

import (
	"database/sql"
	"fmt"
	"renting/internal/models"
	"time"
)

type VehicleServicingRepository struct {
	db *sql.DB
}

func NewVehicleServicingRepository(db *sql.DB) *VehicleServicingRepository {
	return &VehicleServicingRepository{db: db}
}

// InitializeServicingRecord creates a new servicing record for a vehicle
func (r *VehicleServicingRepository) InitializeServicingRecord(vehicleID int, initialKm float64, servicingInterval float64) error {
	_, err := r.db.Exec(`
		INSERT INTO vehicle_servicing (
			vehicle_id, last_servicing_km, next_servicing_km, 
			servicing_interval_km, is_servicing_due, last_serviced_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`, vehicleID, initialKm, initialKm+servicingInterval, servicingInterval, false, time.Now())
	return err
}

// UpdateServicingStatus checks and updates the servicing status based on current km reading
func (r *VehicleServicingRepository) UpdateServicingStatus(vehicleID int, currentKm float64) error {
	// Get the current servicing record
	var servicing models.VehicleServicing
	err := r.db.QueryRow(`
		SELECT servicing_id, vehicle_id, last_servicing_km, next_servicing_km, 
			servicing_interval_km, is_servicing_due, status, last_serviced_at
		FROM vehicle_servicing
		WHERE vehicle_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, vehicleID).Scan(
		&servicing.ServicingID, &servicing.VehicleID, &servicing.LastServicingKm,
		&servicing.NextServicingKm, &servicing.ServicingInterval,
		&servicing.IsServicingDue, &servicing.Status, &servicing.LastServicedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to get servicing record: %v", err)
	}

	// Check if servicing is due
	isServicingDue := currentKm >= servicing.NextServicingKm
	status := "pending"
	if isServicingDue {
		status = "in_progress"
	}

	// Create a new record if status changes
	if status != servicing.Status {
		_, err = r.db.Exec(`
			INSERT INTO vehicle_servicing (
				vehicle_id, last_servicing_km, next_servicing_km,
				servicing_interval_km, is_servicing_due, status,
				last_serviced_at
			) VALUES (
				$1, $2, $3, $4, $5, $6,
				CASE WHEN $6 = 'completed' THEN CURRENT_TIMESTAMP ELSE NULL END
			)
		`, vehicleID, currentKm, servicing.NextServicingKm,
			servicing.ServicingInterval, isServicingDue, status)
	}
	return err
}

// MarkAsServiced creates a new servicing record after a vehicle has been serviced
func (r *VehicleServicingRepository) MarkAsServiced(vehicleID int, servicedAt time.Time) error {
	// Get the current km reading from vehicle_usage
	var currentKm float64
	err := r.db.QueryRow(`
		SELECT km_reading 
		FROM vehicle_usage 
		WHERE vehicle_id = $1 
		ORDER BY recorded_at DESC 
		LIMIT 1
	`, vehicleID).Scan(&currentKm)
	if err != nil {
		return fmt.Errorf("failed to get current km reading: %v", err)
	}

	// Get the current servicing record to get the interval
	var servicing models.VehicleServicing
	err = r.db.QueryRow(`
		SELECT servicing_interval_km
		FROM vehicle_servicing
		WHERE vehicle_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`, vehicleID).Scan(&servicing.ServicingInterval)

	// If no record exists, use default interval of 5000 km
	if err == sql.ErrNoRows {
		servicing.ServicingInterval = 5000
	} else if err != nil {
		return fmt.Errorf("failed to get servicing interval: %v", err)
	}

	// Create a new servicing record
	_, err = r.db.Exec(`
		INSERT INTO vehicle_servicing (
			vehicle_id, last_servicing_km, next_servicing_km,
			servicing_interval_km, is_servicing_due, status,
			last_serviced_at
		) VALUES (
			$1, $2, $3, $4, false, 'completed',
			$5
		)
	`, vehicleID, currentKm, currentKm+servicing.ServicingInterval,
		servicing.ServicingInterval, servicedAt)
	return err
}

// GetServicingHistory retrieves the servicing history for a vehicle
func (r *VehicleServicingRepository) GetServicingHistory(vehicleID int) ([]models.VehicleServicingHistory, error) {
	rows, err := r.db.Query(`
		SELECT h.history_id, h.vehicle_id, h.servicing_id, h.km_reading,
			h.servicing_type, h.servicing_date, h.servicing_cost,
			h.notes, h.serviced_by, h.created_at, u.username
		FROM vehicle_servicing_history h
		LEFT JOIN users u ON h.serviced_by = u.id
		WHERE h.vehicle_id = $1
		ORDER BY h.servicing_date DESC
	`, vehicleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get servicing history: %v", err)
	}
	defer rows.Close()

	var history []models.VehicleServicingHistory
	for rows.Next() {
		var record models.VehicleServicingHistory
		err := rows.Scan(
			&record.HistoryID, &record.VehicleID, &record.ServicingID,
			&record.KmReading, &record.ServicingType, &record.ServicingDate,
			&record.ServicingCost, &record.Notes, &record.ServicedBy,
			&record.CreatedAt, &record.ServicedByName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan servicing history: %v", err)
		}
		history = append(history, record)
	}

	return history, nil
}

// GetCurrentKmAndServicingStatus retrieves the current km reading and servicing status for a vehicle
func (r *VehicleServicingRepository) GetCurrentKmAndServicingStatus(vehicleID int) (*models.VehicleServicing, error) {
	// Get the latest km reading from vehicle_usage
	var currentKm float64
	err := r.db.QueryRow(`
		SELECT km_reading 
		FROM vehicle_usage 
		WHERE vehicle_id = $1 
		ORDER BY recorded_at DESC 
		LIMIT 1
	`, vehicleID).Scan(&currentKm)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get current km reading: %v", err)
	}

	// Get the servicing record
	var servicing models.VehicleServicing
	err = r.db.QueryRow(`
		SELECT vs.servicing_id, vs.vehicle_id, vs.last_servicing_km, vs.next_servicing_km, 
			vs.servicing_interval_km, vs.is_servicing_due, vs.last_serviced_at,
			vs.created_at, vs.updated_at,
			v.vehicle_name, v.vehicle_registration_number
		FROM vehicle_servicing vs
		JOIN vehicles v ON vs.vehicle_id = v.vehicle_id
		WHERE vs.vehicle_id = $1
	`, vehicleID).Scan(
		&servicing.ServicingID, &servicing.VehicleID, &servicing.LastServicingKm,
		&servicing.NextServicingKm, &servicing.ServicingInterval,
		&servicing.IsServicingDue, &servicing.LastServicedAt,
		&servicing.CreatedAt, &servicing.UpdatedAt,
		&servicing.VehicleName, &servicing.RegistrationNumber,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get servicing status: %v", err)
	}

	// Update the servicing status based on current km
	if currentKm > 0 {
		servicing.CurrentKm = currentKm
		servicing.IsServicingDue = currentKm >= servicing.NextServicingKm
	}

	return &servicing, nil
}

// GetVehiclesDueForServicing returns a list of vehicles that need servicing
func (r *VehicleServicingRepository) GetVehiclesDueForServicing() ([]models.VehicleServicing, error) {
	rows, err := r.db.Query(`
		SELECT DISTINCT ON (vs.vehicle_id)
			vs.servicing_id, vs.vehicle_id, vs.last_servicing_km, vs.next_servicing_km,
			vs.servicing_interval_km, vs.is_servicing_due, vs.last_serviced_at,
			vs.created_at, vs.updated_at,
			v.vehicle_name, v.vehicle_registration_number,
			vu.km_reading as current_km
		FROM vehicle_servicing vs
		JOIN vehicles v ON vs.vehicle_id = v.vehicle_id
		LEFT JOIN (
			SELECT vehicle_id, km_reading
			FROM vehicle_usage
			WHERE (vehicle_id, recorded_at) IN (
				SELECT vehicle_id, MAX(recorded_at)
				FROM vehicle_usage
				GROUP BY vehicle_id
			)
		) vu ON vs.vehicle_id = vu.vehicle_id
		WHERE vs.is_servicing_due = true
		ORDER BY vs.vehicle_id, vs.updated_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicles due for servicing: %v", err)
	}
	defer rows.Close()

	var vehicles []models.VehicleServicing
	for rows.Next() {
		var servicing models.VehicleServicing
		err := rows.Scan(
			&servicing.ServicingID, &servicing.VehicleID, &servicing.LastServicingKm,
			&servicing.NextServicingKm, &servicing.ServicingInterval,
			&servicing.IsServicingDue, &servicing.LastServicedAt,
			&servicing.CreatedAt, &servicing.UpdatedAt,
			&servicing.VehicleName, &servicing.RegistrationNumber,
			&servicing.CurrentKm,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vehicle servicing record: %v", err)
		}
		vehicles = append(vehicles, servicing)
	}

	return vehicles, nil
}
