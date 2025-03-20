package repositories

import (
	"database/sql"
	"fmt"
	"renting/internal/models"
)

type ReturnRepository struct {
	db *sql.DB
}

func NewReturnRepository(db *sql.DB) *ReturnRepository {
	return &ReturnRepository{db: db}
}
func (r *ReturnRepository) UpdateSalesStatus(saleID int) error {
	_, err := r.db.Exec(`
		UPDATE sales
		SET status = 'completed'
		WHERE sale_id = $1
	`, saleID)
	if err != nil {
		return fmt.Errorf("failed to update sales status for saleID %d: %v", saleID, err)
	}
	return nil
}

func (r *ReturnRepository) UpdateVehicleStatus(vehicleID int) error {
	_, err := r.db.Exec(`
		UPDATE vehicles
		SET status = 'available'
		WHERE vehicle_id = $1
	`, vehicleID)
	if err != nil {
		return fmt.Errorf("failed to update vehicle status for vehicleID %d: %v", vehicleID, err)
	}
	return nil
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

	// Insert sales charges
	for _, charge := range sale.SalesCharges {
		_, err := tx.Exec(`
			INSERT INTO sales_charges (sale_id, charge_type, amount)
			VALUES ($1, $2, $3)
		`, sale.SaleID, charge.ChargeType, charge.Amount)
		if err != nil {
			return 0, fmt.Errorf("failed to insert sales charge for saleID %d, chargeType %s: %v", sale.SaleID, charge.ChargeType, err)
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

	// Insert payments
	for _, payment := range sale.Payments {
		_, err := tx.Exec(`
			INSERT INTO payments (
				sale_id, amount_paid, payment_date, verified_by_admin, 
				payment_type, payment_status, remark, user_id
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, sale.SaleID, payment.AmountPaid, payment.PaymentDate, payment.VerifiedByAdmin,
			payment.PaymentType, payment.PaymentStatus, payment.Remark, sale.UserID)
		if err != nil {
			return 0, fmt.Errorf("failed to insert payment for saleID %d, amountPaid %.2f: %v", sale.SaleID, payment.AmountPaid, err)
		}
	}
	// Update sales status to "Completed"
	if err := r.UpdateSalesStatus(sale.SaleID); err != nil {
		return 0, fmt.Errorf("failed to update sales status: %v", err)
	}

	// Update vehicle status to "Available"
	for _, usage := range sale.VehicleUsage {
		if err := r.UpdateVehicleStatus(usage.VehicleID); err != nil {
			return 0, fmt.Errorf("failed to update vehicle status: %v", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction for saleID %d: %v", sale.SaleID, err)
	}

	// Return the sale ID
	return sale.SaleID, nil
}
