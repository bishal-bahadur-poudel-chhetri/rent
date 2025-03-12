package repositories

import (
	"database/sql"
	"fmt"
	"renting/internal/models"
)

type SaleRepository struct {
	db *sql.DB
}

func NewSaleRepository(db *sql.DB) *SaleRepository {
	return &SaleRepository{db: db}
}

func (r *SaleRepository) CreateSale(sale models.Sale) (int, error) {
	// Start a database transaction
	tx, err := r.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %v", err)
	}
	// Ensure the transaction is rolled back in case of an error
	defer func() {
		if err != nil {
			tx.Rollback() // Rollback if there's an error
			return
		}
	}()

	// Step 1: Insert payment details into the payments table
	var paymentID int
	err = tx.QueryRow(`
		INSERT INTO payments (
			amount_paid, payment_date, verified_by_admin, 
			payment_type, payment_status, remark
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING payment_id
	`, sale.AmountPaid, sale.PaymentDate, sale.VerifiedByAdmin,
		sale.PaymentType, sale.PaymentStatus, sale.Remark).Scan(&paymentID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert payment details: %v", err)
	}

	// Step 2: Insert sale details into the sales table using the paymentID
	var saleID int
	err = tx.QueryRow(`
		INSERT INTO sales (
			vehicle_id, user_id, customer_name, total_amount, charge_per_day, booking_date, 
			date_of_delivery, return_date, is_damaged, is_washed, is_delayed, 
			number_of_days, payment_id, remark, status,customer_destination
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,$16)
		RETURNING sale_id
	`, sale.VehicleID, sale.UserID, sale.CustomerName, sale.TotalAmount, sale.ChargePerDay, sale.BookingDate,
		sale.DateOfDelivery, sale.ReturnDate, sale.IsDamaged, sale.IsWashed, sale.IsDelayed,
		sale.NumberOfDays, paymentID, sale.Remark, sale.Status, sale.Destination).Scan(&saleID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert sale: %v", err)
	}

	// Step 3: Insert sales charges
	for _, charge := range sale.SalesCharges {
		_, err := tx.Exec(`
			INSERT INTO sales_charges (sale_id, charge_type, amount)
			VALUES ($1, $2, $3)
		`, saleID, charge.ChargeType, charge.Amount)
		if err != nil {
			return 0, fmt.Errorf("failed to insert sales charge: %v", err)
		}
	}

	// Step 4: Insert sales images
	for _, image := range sale.SalesImages {
		_, err := tx.Exec(`
			INSERT INTO sales_images (sale_id, image_url)
			VALUES ($1, $2)
		`, saleID, image.ImageURL)
		if err != nil {
			return 0, fmt.Errorf("failed to insert sales image: %v", err)
		}
	}

	// Step 5: Insert sales videos
	for _, video := range sale.SalesVideos {
		_, err := tx.Exec(`
			INSERT INTO sales_videos (sale_id, video_url)
			VALUES ($1, $2)
		`, saleID, video.VideoURL)
		if err != nil {
			return 0, fmt.Errorf("failed to insert sales video: %v", err)
		}
	}

	// Step 6: Insert vehicle usage records
	for _, usage := range sale.VehicleUsage {
		_, err := tx.Exec(`
			INSERT INTO vehicle_usage (vehicle_id, record_type, fuel_range, km_reading, recorded_by)
			VALUES ($1, $2, $3, $4, $5)
		`, sale.VehicleID, usage.RecordType, usage.FuelRange, usage.KmReading, usage.RecordedBy)
		if err != nil {
			return 0, fmt.Errorf("failed to insert vehicle usage record: %v", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return saleID, nil
}

func (r *SaleRepository) GetSaleByID(saleID int) (*models.Sale, error) {
	var sale models.Sale
	err := r.db.QueryRow(`
		SELECT sale_id, vehicle_id, user_id, customer_name, total_amount, charge_per_day, booking_date, 
		date_of_delivery, return_date, is_damaged, is_washed, is_delayed, number_of_days, 
		payment_id, remark, status, created_at, updated_at
		FROM sales WHERE sale_id = $1
	`, saleID).Scan(
		&sale.SaleID, &sale.VehicleID, &sale.UserID, &sale.CustomerName, &sale.TotalAmount, &sale.ChargePerDay,
		&sale.BookingDate, &sale.DateOfDelivery, &sale.ReturnDate, &sale.IsDamaged, &sale.IsWashed,
		&sale.IsDelayed, &sale.NumberOfDays, &sale.PaymentID, &sale.Remark, &sale.Status, &sale.CreatedAt, &sale.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Fetch related data
	sale.SalesCharges, _ = r.getSalesCharges(saleID)
	sale.SalesImages, _ = r.getSalesImages(saleID)
	sale.SalesVideos, _ = r.getSalesVideos(saleID)
	sale.VehicleUsage, _ = r.getVehicleUsage(sale.VehicleID)

	return &sale, nil
}

// Helper functions to fetch related data
func (r *SaleRepository) getSalesCharges(saleID int) ([]models.SalesCharge, error) {
	rows, err := r.db.Query("SELECT charge_id, sale_id, charge_type, amount, created_at, updated_at FROM sales_charges WHERE sale_id = $1", saleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var charges []models.SalesCharge
	for rows.Next() {
		var charge models.SalesCharge
		rows.Scan(&charge.ChargeID, &charge.SaleID, &charge.ChargeType, &charge.Amount, &charge.CreatedAt, &charge.UpdatedAt)
		charges = append(charges, charge)
	}
	return charges, nil
}

func (r *SaleRepository) getSalesImages(saleID int) ([]models.SalesImage, error) {
	rows, err := r.db.Query("SELECT image_id, sale_id, image_url, uploaded_at FROM sales_images WHERE sale_id = $1", saleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.SalesImage
	for rows.Next() {
		var image models.SalesImage
		rows.Scan(&image.ImageID, &image.SaleID, &image.ImageURL, &image.UploadedAt)
		images = append(images, image)
	}
	return images, nil
}

func (r *SaleRepository) getSalesVideos(saleID int) ([]models.SalesVideo, error) {
	rows, err := r.db.Query("SELECT video_id, sale_id, video_url, uploaded_at FROM sales_videos WHERE sale_id = $1", saleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []models.SalesVideo
	for rows.Next() {
		var video models.SalesVideo
		rows.Scan(&video.VideoID, &video.SaleID, &video.VideoURL, &video.UploadedAt)
		videos = append(videos, video)
	}
	return videos, nil
}

func (r *SaleRepository) getVehicleUsage(vehicleID int) ([]models.VehicleUsage, error) {
	rows, err := r.db.Query("SELECT usage_id, vehicle_id, record_type, fuel_range, km_reading, recorded_at, recorded_by FROM vehicle_usage WHERE vehicle_id = $1", vehicleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usageRecords []models.VehicleUsage
	for rows.Next() {
		var usage models.VehicleUsage
		rows.Scan(&usage.UsageID, &usage.VehicleID, &usage.RecordType, &usage.FuelRange, &usage.KmReading, &usage.RecordedAt, &usage.RecordedBy)
		usageRecords = append(usageRecords, usage)
	}
	return usageRecords, nil
}
