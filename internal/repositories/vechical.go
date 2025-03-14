package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"renting/internal/models"
	"strconv"
	"time"
)

type VehicleRepository struct {
	db *sql.DB
}

func NewVehicleRepository(db *sql.DB) *VehicleRepository {
	return &VehicleRepository{db: db}
}

func (r *VehicleRepository) GetVehicles(filters models.VehicleFilter, includeBookingDetails bool, includeSaleid bool) ([]models.VehicleResponse, error) {
	// Log filters and includeBookingDetails for debugging
	log.Printf("Fetching vehicles with filters: %+v, includeBookingDetails: %v, includeSaleid: %v", filters, includeBookingDetails, includeSaleid)

	// Base query to fetch vehicles
	baseQuery := `
		SELECT 
			v.vehicle_id, 
			v.vehicle_type_id, 
			v.vehicle_name, 
			v.vehicle_model, 
			v.vehicle_registration_number, 
			v.is_available, 
			v.status
		FROM vehicles v
		WHERE 1=1
	`

	// Add filters dynamically
	if filters.VehicleID != "" {
		baseQuery += " AND v.vehicle_id = " + filters.VehicleID
	}
	if filters.VehicleTypeID != "" {
		baseQuery += " AND v.vehicle_type_id = " + filters.VehicleTypeID
	}
	if filters.VehicleName != "" {
		baseQuery += " AND v.vehicle_name = '" + filters.VehicleName + "'"
	}
	if filters.VehicleModel != "" {
		baseQuery += " AND v.vehicle_model = '" + filters.VehicleModel + "'"
	}
	if filters.VehicleRegistrationNumber != "" {
		baseQuery += " AND v.vehicle_registration_number = '" + filters.VehicleRegistrationNumber + "'"
	}
	if filters.IsAvailable != "" {
		baseQuery += " AND v.is_available = " + filters.IsAvailable
	}
	if filters.Status != "" {
		baseQuery += " AND v.status = '" + filters.Status + "'"
	}

	// Add pagination
	if filters.Limit > 0 {
		baseQuery += " LIMIT " + strconv.Itoa(filters.Limit)
	}
	if filters.Offset > 0 {
		baseQuery += " OFFSET " + strconv.Itoa(filters.Offset)
	}

	// Execute the query
	rows, err := r.db.Query(baseQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vehicles []models.VehicleResponse
	for rows.Next() {
		var v models.VehicleResponse
		err := rows.Scan(
			&v.VehicleID,
			&v.VehicleTypeID,
			&v.VehicleName,
			&v.VehicleModel,
			&v.VehicleRegistrationNumber,
			&v.IsAvailable,
			&v.Status,
		)
		if err != nil {
			return nil, err
		}

		// Fetch future booking details if requested
		if includeBookingDetails {
			log.Printf("Fetching future booking details for vehicle ID: %d", v.VehicleID)
			futureBookings, err := r.getFutureBookingDetails(v.VehicleID)
			if err != nil {
				return nil, err
			}
			v.FutureBookingDetails = futureBookings
		}

		if includeSaleid {
			log.Printf("Fetching sale ID for vehicle ID: %d", v.VehicleID)
			saleID, err := r.getMostRecentSaleID(v.VehicleID)
			if err != nil {
				return nil, err
			}
			v.SaleID = saleID
		}

		vehicles = append(vehicles, v)
	}

	return vehicles, nil
}

// getFutureBookingDetails fetches future booking details for a vehicle
func (r *VehicleRepository) getFutureBookingDetails(vehicleID int) ([]models.FutureBookingDetail, error) {
	rows, err := r.db.Query(`
		SELECT date_of_delivery, number_of_days
		FROM sales
		WHERE vehicle_id = $1 AND date_of_delivery > NOW()
		ORDER BY date_of_delivery ASC
	`, vehicleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var futureBookings []models.FutureBookingDetail
	for rows.Next() {
		var booking models.FutureBookingDetail
		var bookingDate time.Time
		err := rows.Scan(&bookingDate, &booking.NumberOfDays)
		if err != nil {
			return nil, err
		}
		booking.DeliveryDate = bookingDate.Format("2006-01-02") // Format date as YYYY-MM-DD
		futureBookings = append(futureBookings, booking)
	}

	log.Printf("Fetched %d future bookings for vehicle ID: %d", len(futureBookings), vehicleID)
	return futureBookings, nil
}

// getMostRecentSaleID fetches the most recent sale ID for a vehicle
func (r *VehicleRepository) getMostRecentSaleID(vehicleID int) (int, error) {
	fmt.Print(vehicleID)
	var saleID int
	err := r.db.QueryRow(`
		SELECT sale_id
		FROM sales
		WHERE vehicle_id = $1 and status='active'
	`, vehicleID).Scan(&saleID)
	if err != nil {
		if err == sql.ErrNoRows {

			return 0, nil
		}
		return 0, err
	}
	return saleID, nil
}
