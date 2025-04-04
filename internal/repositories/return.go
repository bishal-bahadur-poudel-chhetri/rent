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

	// Update sales table with the additional charge
	if totalCharge > 0 {
		_, err = tx.Exec(`
            UPDATE sales
            SET total_amount = total_amount + $1,
                updated_at = NOW()
            WHERE sale_id = $2
        `, totalCharge, sale.SaleID)
		if err != nil {
			return 0, fmt.Errorf("failed to update total_amount for saleID %d with additional charge %.2f: %v", sale.SaleID, totalCharge, err)
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
