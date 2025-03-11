package repositories

import (
	"database/sql"
	"renting/internal/models"
)

type SaleRepository struct {
	db *sql.DB
}

func NewSaleRepository(db *sql.DB) *SaleRepository {
	return &SaleRepository{db: db}
}

// CreateSale creates a new sale and its related data
func (r *SaleRepository) CreateSale(sale models.Sale) (int, error) {
	// Insert sale
	var saleID int
	err := r.db.QueryRow(`
		INSERT INTO sales (
			vehicle_id, customer_name, total_amount, charge_per_day, booking_date, 
			date_of_delivery, return_date, is_damaged, is_washed, is_delayed, 
			number_of_days, payment_id, remark, fuel_range_received, fuel_range_delivered, 
			km_received, km_delivered
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		RETURNING sale_id
	`, sale.VehicleID, sale.CustomerName, sale.TotalAmount, sale.ChargePerDay, sale.BookingDate,
		sale.DateOfDelivery, sale.ReturnDate, sale.IsDamaged, sale.IsWashed, sale.IsDelayed,
		sale.NumberOfDays, sale.PaymentID, sale.Remark, sale.FuelRangeReceived, sale.FuelRangeDelivered,
		sale.KmReceived, sale.KmDelivered).Scan(&saleID)
	if err != nil {
		return 0, err
	}

	// Insert sales charges
	for _, charge := range sale.SalesCharges {
		_, err := r.db.Exec(`
			INSERT INTO sales_charges (sale_id, charge_type, amount)
			VALUES ($1, $2, $3)
		`, saleID, charge.ChargeType, charge.Amount)
		if err != nil {
			return 0, err
		}
	}

	// Insert sales images
	for _, image := range sale.SalesImages {
		_, err := r.db.Exec(`
			INSERT INTO sales_images (sale_id, image_url)
			VALUES ($1, $2)
		`, saleID, image.ImageURL)
		if err != nil {
			return 0, err
		}
	}

	// Insert sales videos
	for _, video := range sale.SalesVideos {
		_, err := r.db.Exec(`
			INSERT INTO sales_videos (sale_id, video_url)
			VALUES ($1, $2)
		`, saleID, video.VideoURL)
		if err != nil {
			return 0, err
		}
	}

	return saleID, nil
}

// GetSaleByID retrieves a sale by its ID along with related data
func (r *SaleRepository) GetSaleByID(saleID int) (*models.Sale, error) {
	var sale models.Sale
	err := r.db.QueryRow(`
		SELECT sale_id, vehicle_id, customer_name, total_amount, charge_per_day, booking_date, 
		date_of_delivery, return_date, is_damaged, is_washed, is_delayed, number_of_days, 
		payment_id, remark, fuel_range_received, fuel_range_delivered, km_received, km_delivered, 
		created_at, updated_at
		FROM sales WHERE sale_id = $1
	`, saleID).Scan(
		&sale.SaleID, &sale.VehicleID, &sale.CustomerName, &sale.TotalAmount, &sale.ChargePerDay,
		&sale.BookingDate, &sale.DateOfDelivery, &sale.ReturnDate, &sale.IsDamaged, &sale.IsWashed,
		&sale.IsDelayed, &sale.NumberOfDays, &sale.PaymentID, &sale.Remark, &sale.FuelRangeReceived,
		&sale.FuelRangeDelivered, &sale.KmReceived, &sale.KmDelivered, &sale.CreatedAt, &sale.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Fetch related data (charges, images, videos)
	sale.SalesCharges, _ = r.getSalesCharges(saleID)
	sale.SalesImages, _ = r.getSalesImages(saleID)
	sale.SalesVideos, _ = r.getSalesVideos(saleID)

	return &sale, nil
}

// Helper functions to fetch related data
func (r *SaleRepository) getSalesCharges(saleID int) ([]models.SalesCharge, error) {
	rows, err := r.db.Query("SELECT charge_id, sale_id, charge_type, amount FROM sales_charges WHERE sale_id = $1", saleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var charges []models.SalesCharge
	for rows.Next() {
		var charge models.SalesCharge
		rows.Scan(&charge.ChargeID, &charge.SaleID, &charge.ChargeType, &charge.Amount)
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
