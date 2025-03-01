package repositories

import (
	"database/sql"
	"log"
	"renting/internal/models"
)

// SaleRepository interacts with the database to perform CRUD operations related to sales
type SaleRepository struct {
	DB *sql.DB
}

// CreateSale creates a new sale in the database
func (r *SaleRepository) CreateSale(sale *models.Sale) (int, error) {
	query := `
    INSERT INTO sales (vehicle_id, customer_name, total_amount, charge_per_day, booking_date, 
                      date_of_delivery, return_date, is_damaged, is_washed, is_delayed, 
                      number_of_days, payment_id, remark, 
                      fuel_range_received, fuel_range_delivered, km_received, km_delivered, 
                      photo_1, photo_2, photo_3, photo_4)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
    RETURNING sale_id;
    `
	var saleID int
	err := r.DB.QueryRow(query, sale.VehicleID, sale.CustomerName, sale.TotalAmount, sale.ChargePerDay,
		sale.BookingDate, sale.DateOfDelivery, sale.ReturnDate, sale.IsDamaged, sale.IsWashed,
		sale.IsDelayed, sale.NumberOfDays, sale.PaymentID, sale.Remark,
		sale.FuelRangeReceived, sale.FuelRangeDelivered, sale.KmReceived, sale.KmDelivered,
		sale.Photo1, sale.Photo2, sale.Photo3, sale.Photo4).Scan(&saleID)

	if err != nil {
		log.Println("Error inserting sale into database:", err)
		return 0, err
	}

	return saleID, nil
}
