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

func (r *SaleChargeRepository) UpdateSalesCharge(chargeID int, charge models.SalesCharge) error {
	fmt.Printf("Repository: Updating charge ID %d with type '%s' and amount %f\n", chargeID, charge.ChargeType, charge.Amount)
	
	_, err := r.db.Exec(`
		UPDATE sales_charges 
		SET charge_type = $1, amount = $2
		WHERE charge_id = $3
	`, charge.ChargeType, charge.Amount, chargeID)
	if err != nil {
		fmt.Printf("Repository: Database error: %v\n", err)
		return fmt.Errorf("failed to update sales charge: %v", err)
	}
	fmt.Printf("Repository: Update successful\n")
	return nil
}

func (r *SaleChargeRepository) DeleteSalesCharge(chargeID int) error {
	_, err := r.db.Exec(`
		DELETE FROM sales_charges 
		WHERE charge_id = $1
	`, chargeID)
	if err != nil {
		return fmt.Errorf("failed to delete sales charge: %v", err)
	}
	return nil
}
