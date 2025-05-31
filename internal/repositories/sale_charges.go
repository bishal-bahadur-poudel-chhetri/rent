package repositories

import (
	"database/sql"
	"fmt"
	"renting/internal/models"
)

type SaleChargeRepository struct {
	db *sql.DB
}

func NewSaleChargeRepository(db *sql.DB) *SaleChargeRepository {
	return &SaleChargeRepository{db: db}
}

func (r *SaleChargeRepository) AddSalesCharges(saleID int, charges []models.SalesCharge) error {
	for _, charge := range charges {
		_, err := r.db.Exec(`
            INSERT INTO sales_charges (sale_id, charge_type, amount)
            VALUES ($1, $2, $3)
        `, saleID, charge.ChargeType, charge.Amount)
		if err != nil {
			return fmt.Errorf("failed to insert sales charge: %v", err)
		}
	}
	return nil
}

func (r *SaleChargeRepository) GetSalesChargesBySaleID(saleID int) ([]models.SalesCharge, error) {
	rows, err := r.db.Query(`
		SELECT charge_id, sale_id, charge_type, amount
		FROM sales_charges 
		WHERE sale_id = $1
		ORDER BY charge_id DESC
	`, saleID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sales charges: %v", err)
	}
	defer rows.Close()

	var charges []models.SalesCharge
	for rows.Next() {
		var charge models.SalesCharge
		err := rows.Scan(
			&charge.ChargeID,
			&charge.SaleID,
			&charge.ChargeType,
			&charge.Amount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sales charge: %v", err)
		}
		charges = append(charges, charge)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sales charges: %v", err)
	}

	return charges, nil
}

func (r *SaleChargeRepository) UpdateSalesCharge(chargeID int, charge models.SalesCharge) error {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Get the old charge details
	var oldCharge models.SalesCharge
	err = tx.QueryRow(`
		SELECT charge_id, sale_id, charge_type, amount
		FROM sales_charges 
		WHERE charge_id = $1 AND sale_id = $2
	`, chargeID, charge.SaleID).Scan(
		&oldCharge.ChargeID,
		&oldCharge.SaleID,
		&oldCharge.ChargeType,
		&oldCharge.Amount,
	)
	if err != nil {
		return fmt.Errorf("failed to get old charge details: %v", err)
	}

	// Update the charge
	_, err = tx.Exec(`
		UPDATE sales_charges 
		SET charge_type = $1, amount = $2
		WHERE charge_id = $3 AND sale_id = $4
	`, charge.ChargeType, charge.Amount, chargeID, charge.SaleID)
	if err != nil {
		return fmt.Errorf("failed to update sales charge: %v", err)
	}

	// Update the sales table based on the old charge type
	switch oldCharge.ChargeType {
	case "discount":
		_, err = tx.Exec(`
			UPDATE sales 
			SET discount = discount - $1
			WHERE sale_id = $2
		`, oldCharge.Amount, charge.SaleID)
	case "wash", "damage":
		_, err = tx.Exec(`
			UPDATE sales 
			SET other_charges = other_charges - $1
			WHERE sale_id = $2
		`, oldCharge.Amount, charge.SaleID)
	}
	if err != nil {
		return fmt.Errorf("failed to update sales table for old charge: %v", err)
	}

	// Update the sales table based on the new charge type
	switch charge.ChargeType {
	case "discount":
		_, err = tx.Exec(`
			UPDATE sales 
			SET discount = discount + $1
			WHERE sale_id = $2
		`, charge.Amount, charge.SaleID)
	case "wash", "damage":
		_, err = tx.Exec(`
			UPDATE sales 
			SET other_charges = other_charges + $1
			WHERE sale_id = $2
		`, charge.Amount, charge.SaleID)
	}
	if err != nil {
		return fmt.Errorf("failed to update sales table for new charge: %v", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *SaleChargeRepository) DeleteSalesCharge(chargeID int, saleID int) error {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Get the charge details before deleting
	var charge models.SalesCharge
	err = tx.QueryRow(`
		SELECT charge_id, sale_id, charge_type, amount
		FROM sales_charges 
		WHERE charge_id = $1 AND sale_id = $2
	`, chargeID, saleID).Scan(
		&charge.ChargeID,
		&charge.SaleID,
		&charge.ChargeType,
		&charge.Amount,
	)
	if err != nil {
		return fmt.Errorf("failed to get charge details: %v", err)
	}

	// Delete the charge
	result, err := tx.Exec(`
		DELETE FROM sales_charges 
		WHERE charge_id = $1 AND sale_id = $2
	`, chargeID, saleID)
	if err != nil {
		return fmt.Errorf("failed to delete sales charge: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no sales charge found with ID %d for sale %d", chargeID, saleID)
	}

	// Update the sales table based on the charge type
	switch charge.ChargeType {
	case "discount":
		_, err = tx.Exec(`
			UPDATE sales 
			SET discount = COALESCE(discount, 0) - $1
			WHERE sale_id = $2
		`, charge.Amount, saleID)
	case "wash", "damage":
		_, err = tx.Exec(`
			UPDATE sales 
			SET other_charges = COALESCE(other_charges, 0) - $1
			WHERE sale_id = $2
		`, charge.Amount, saleID)
	}
	if err != nil {
		return fmt.Errorf("failed to update sales table: %v", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (r *SaleChargeRepository) AddSalesCharge(charge models.SalesCharge) error {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert the charge
	_, err = tx.Exec(`
		INSERT INTO sales_charges (sale_id, charge_type, amount)
		VALUES ($1, $2, $3)
	`, charge.SaleID, charge.ChargeType, charge.Amount)
	if err != nil {
		return fmt.Errorf("failed to add sales charge: %v", err)
	}

	// Update the sales table based on the charge type
	switch charge.ChargeType {
	case "discount":
		_, err = tx.Exec(`
			UPDATE sales 
			SET discount = COALESCE(discount, 0) + $1
			WHERE sale_id = $2
		`, charge.Amount, charge.SaleID)
	case "wash", "damage":
		_, err = tx.Exec(`
			UPDATE sales 
			SET other_charges = COALESCE(other_charges, 0) + $1
			WHERE sale_id = $2
		`, charge.Amount, charge.SaleID)
	}
	if err != nil {
		return fmt.Errorf("failed to update sales table: %v", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
