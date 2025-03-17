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
