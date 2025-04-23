package repositories

import (
	"database/sql"
	"fmt"
	"renting/internal/models"
	"time"
)

type ReturnRepository struct {
	db *sql.DB
}

func NewReturnRepository(db *sql.DB) *ReturnRepository {
	return &ReturnRepository{db: db}
}
func (r *ReturnRepository) CreateReturn(sale models.Sale) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// Insert sales charges and calculate total additional charge
	var totalCharge float64
	for _, charge := range sale.SalesCharges {
		_, err := tx.Exec(`
            INSERT INTO sales_charges (sale_id, charge_type, amount)
            VALUES ($1, $2, $3)
        `, sale.SaleID, charge.ChargeType, charge.Amount)
		if err != nil {
			return 0, fmt.Errorf("failed to insert sales charge for saleID %d, chargeType %s: %v", sale.SaleID, charge.ChargeType, err)
		}
		totalCharge += charge.Amount // Accumulate total charge
	}

	// Update sales table with the additional charge in other_charges
	if totalCharge > 0 {
		_, err = tx.Exec(`
            UPDATE sales
            SET other_charges = COALESCE(other_charges, 0) + $1,
                updated_at = NOW()
            WHERE sale_id = $2
        `, totalCharge, sale.SaleID)
		if err != nil {
			return 0, fmt.Errorf("failed to update other_charges for saleID %d with additional charge %.2f: %v", sale.SaleID, totalCharge, err)
		}
	}

	// Insert vehicle usage records
	for _, usage := range sale.VehicleUsage {
		_, err := tx.Exec(`
            INSERT INTO vehicle_usage (sale_id, vehicle_id, record_type, fuel_range, km_reading, recorded_at, recorded_by)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
        `, sale.SaleID, usage.VehicleID, usage.RecordType, usage.FuelRange, usage.KmReading, usage.RecordedAt, sale.UserID)
		if err != nil {
			return 0, fmt.Errorf("failed to insert vehicle usage record for saleID %d, vehicleID %d: %v", sale.SaleID, usage.VehicleID, err)
		}
		fmt.Println(usage.RecordType)
		// Only check servicing status for return records
		if usage.RecordType == "return" {
			// Check if servicing record exists for this vehicle
			var servicingExists bool
			err = tx.QueryRow(`
				SELECT EXISTS (
					SELECT 1 FROM vehicle_servicing WHERE vehicle_id = $1
				)
			`, usage.VehicleID).Scan(&servicingExists)
			if err != nil {
				return 0, fmt.Errorf("failed to check servicing record existence: %v", err)
			}

			if !servicingExists {
				// Initialize servicing record if it doesn't exist
				_, err = tx.Exec(`
					INSERT INTO vehicle_servicing (
						vehicle_id, current_km, next_servicing_km, 
						servicing_interval_km, is_servicing_due, last_serviced_at,
						status
					) VALUES ($1, $2, $3, $4, $5, $6, $7)
				`, usage.VehicleID, usage.KmReading, usage.KmReading+5000, 5000, false, usage.RecordedAt, "pending")
				if err != nil {
					return 0, fmt.Errorf("failed to initialize servicing record: %v", err)
				}
			} else {
				// Update servicing status based on new km reading
				_, err = tx.Exec(`
					UPDATE vehicle_servicing
					SET current_km = $1,
						is_servicing_due = CASE 
							WHEN $1 >= next_servicing_km THEN true 
							ELSE is_servicing_due 
						END,
						status = CASE
							WHEN $1 >= next_servicing_km AND status = 'pending' THEN 'in_progress'
							ELSE status
						END,
						updated_at = CURRENT_TIMESTAMP
					WHERE vehicle_id = $2
					AND status IN ('pending', 'in_progress')
				`, usage.KmReading, usage.VehicleID)
				if err != nil {
					return 0, fmt.Errorf("failed to update servicing status: %v", err)
				}
			}
		}
	}

	currentTime := time.Now()

	// Determine sale type based on vehicle usage records
	saleType := "delivery" // default to delivery
	for _, usage := range sale.VehicleUsage {
		if usage.RecordType == "return" {
			saleType = "return"
			break
		}
	}

	// Insert payments with sale_type
	for _, payment := range sale.Payments {
		_, err := tx.Exec(`
            INSERT INTO payments (
                sale_id, amount_paid, payment_date, verified_by_admin, 
                payment_type, payment_status, remark, user_id, sale_type
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        `, sale.SaleID, payment.AmountPaid, payment.PaymentDate, payment.VerifiedByAdmin,
			payment.PaymentType, payment.PaymentStatus, payment.Remark, sale.UserID, saleType)
		if err != nil {
			return 0, fmt.Errorf("failed to insert payment for saleID %d, amountPaid %.2f: %v", sale.SaleID, payment.AmountPaid, err)
		}
	}

	// Update sales status, vehicle status, and delivery/return dates based on record_type
	for _, usage := range sale.VehicleUsage {
		var saleStatus, vehicleStatus string
		var updateQuery string
		var queryArgs []interface{}

		if usage.RecordType == "delivery" {
			saleStatus = "active"
			vehicleStatus = "rented"
			updateQuery = `
                UPDATE sales
                SET status = $1,
                    actual_date_of_delivery = $2
                WHERE sale_id = $3
            `
			queryArgs = []interface{}{saleStatus, currentTime, sale.SaleID}
		} else if usage.RecordType == "return" {
			saleStatus = "completed"
			vehicleStatus = "available"
			updateQuery = `
                UPDATE sales
                SET status = $1,
                    actual_date_of_return = $2
                WHERE sale_id = $3
            `
			queryArgs = []interface{}{saleStatus, currentTime, sale.SaleID}
		} else {
			saleStatus = "active"
			vehicleStatus = "available"
			updateQuery = `
                UPDATE sales
                SET status = $1
                WHERE sale_id = $2
            `
			queryArgs = []interface{}{saleStatus, sale.SaleID}
		}

		// Update sale status and dates
		_, err := tx.Exec(updateQuery, queryArgs...)
		if err != nil {
			return 0, fmt.Errorf("failed to update sales status for saleID %d: %v", sale.SaleID, err)
		}

		// Update vehicle status
		_, err = tx.Exec(`
            UPDATE vehicles
            SET status = $1
            WHERE vehicle_id = $2
        `, vehicleStatus, usage.VehicleID)
		if err != nil {
			return 0, fmt.Errorf("failed to update vehicle status for vehicleID %d: %v", usage.VehicleID, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction for saleID %d: %v", sale.SaleID, err)
	}

	return sale.SaleID, nil
}
