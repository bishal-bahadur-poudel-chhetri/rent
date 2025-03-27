package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"renting/internal/models"
	"time"
)

type SaleRepository struct {
	db *sql.DB
}

func (r *SaleRepository) GetFutureBookings(vehicleID int, dateOfDelivery time.Time) (any, any) {
	panic("unimplemented")
}

func NewSaleRepository(db *sql.DB) *SaleRepository {
	return &SaleRepository{db: db}
}
func (r *SaleRepository) UpdateVehicleStatus(vehicleID int, status string) error {
	_, err := r.db.Exec(`
		UPDATE vehicles
		SET status = $1
		WHERE vehicle_id = $2
	`, status, vehicleID)
	if err != nil {
		return fmt.Errorf("failed to update vehicle status: %v", err)
	}
	return nil
}
func (r *SaleRepository) UpdateSaleStatus(saleID int, status string) error {
	fmt.Print(saleID, status)
	_, err := r.db.Exec(`
		UPDATE sales
		SET status = $2
		WHERE sale_id = $1
	`, saleID, status) // Correct order of parameters: saleID ($1), status ($2)
	if err != nil {
		return fmt.Errorf("failed to update sale status: %v", err)
	}
	return nil
}

func (r *SaleRepository) checkVehicleAvailability(vehicleID int, startDate, endDate time.Time) (bool, error) {
	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) 
		FROM sales 
		WHERE vehicle_id = $1 
		AND (
			(date_of_delivery <= $2 AND return_date >= $3) OR
			(date_of_delivery BETWEEN $2 AND $3) OR
			(return_date BETWEEN $2 AND $3)
		) AND status IN ('active', 'pending')
	`, vehicleID, endDate, startDate).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to check vehicle availability: %v", err)
	}

	return count == 0, nil
}

func (r *SaleRepository) CreateSale(sale models.Sale) (models.SaleSubmitResponse, error) {
	BookingDate := time.Now()
	// Check vehicle availability first
	available, err := r.checkVehicleAvailability(sale.VehicleID, sale.DateOfDelivery, sale.ReturnDate)
	if err != nil {
		return models.SaleSubmitResponse{}, err
	}
	if !available {
		return models.SaleSubmitResponse{}, errors.New("vehicle is not available for the selected dates")
	}
	actualDeliveryDate := sale.DateOfDelivery

	if sale.BookingDate.Before(sale.DateOfDelivery) {
		actualDeliveryDate = time.Now()
	}

	// Initialize the response object
	var salesResponse models.SaleSubmitResponse

	// Begin the transaction
	tx, err := r.db.Begin()
	if err != nil {
		return salesResponse, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert sale
	var saleID int
	err = tx.QueryRow(`
        INSERT INTO sales (
            vehicle_id, user_id, customer_name, total_amount, charge_per_day, booking_date, 
            date_of_delivery, return_date,
            number_of_days, remark, status, customer_destination, customer_phone,actual_date_of_delivery
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,$17)

        RETURNING sale_id
    `, sale.VehicleID, sale.UserID, sale.CustomerName, sale.TotalAmount, sale.ChargePerDay, BookingDate,
		sale.DateOfDelivery, sale.ReturnDate,
		sale.NumberOfDays, sale.Remark, sale.Status, sale.Destination, sale.CustomerPhone, actualDeliveryDate).Scan(&saleID)

	if err != nil {
		return salesResponse, fmt.Errorf("failed to insert sale: %v", err)
	}

	// Insert sales charges
	for _, charge := range sale.SalesCharges {
		_, err := tx.Exec(`
            INSERT INTO sales_charges (sale_id, charge_type, amount)
            VALUES ($1, $2, $3)
        `, saleID, charge.ChargeType, charge.Amount)
		if err != nil {
			return salesResponse, fmt.Errorf("failed to insert sales charge: %v", err)
		}
	}

	// Insert sales images
	for _, image := range sale.SalesImages {
		_, err := tx.Exec(`
            INSERT INTO sales_images (sale_id, image_url)
            VALUES ($1, $2)
        `, saleID, image.ImageURL)
		if err != nil {
			return salesResponse, fmt.Errorf("failed to insert sales image: %v", err)
		}
	}

	// Insert sales videos
	for _, video := range sale.SalesVideos {
		_, err := tx.Exec(`
            INSERT INTO sales_videos (sale_id, video_url)
            VALUES ($1, $2)
        `, saleID, video.VideoURL)
		if err != nil {
			return salesResponse, fmt.Errorf("failed to insert sales video: %v", err)
		}
	}

	// Insert vehicle usage records and set payment sale_type accordingly
	for _, usage := range sale.VehicleUsage {
		_, err := tx.Exec(`
            INSERT INTO vehicle_usage (sale_id, vehicle_id, record_type, fuel_range, km_reading, recorded_at, recorded_by)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
        `, saleID, usage.VehicleID, usage.RecordType, usage.FuelRange, usage.KmReading, usage.RecordedAt, sale.UserID)
		if err != nil {
			return salesResponse, fmt.Errorf("failed to insert vehicle usage record: %v", err)
		}
	}

	for _, payment := range sale.Payments {
		if payment.SaleType == "" {
			for _, usage := range sale.VehicleUsage {
				if usage.RecordType == "delivery" {
					payment.SaleType = models.TypeDelivery
					break
				} else if usage.RecordType == "return" {
					payment.SaleType = models.TypeReturn
					break
				}
			}

			if payment.SaleType == "" {
				payment.SaleType = models.TypeBooking
			}
		}

		// Validate payment
		if err := payment.Validate(); err != nil {
			return salesResponse, fmt.Errorf("invalid payment: %v", err)
		}

		_, err := tx.Exec(`
            INSERT INTO payments (
                sale_id, amount_paid, payment_date, verified_by_admin, 
                payment_type, payment_status, remark, user_id, sale_type
            ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        `, saleID, payment.AmountPaid, payment.PaymentDate, payment.VerifiedByAdmin,
			payment.PaymentType, payment.PaymentStatus, payment.Remark, sale.UserID, payment.SaleType)
		if err != nil {
			return salesResponse, fmt.Errorf("failed to insert payment: %v", err)
		}
	}

	// Fetch vehicle name
	var vehicleName string
	err = tx.QueryRow(`
        SELECT vehicle_name FROM vehicles WHERE vehicle_id = $1
    `, sale.VehicleID).Scan(&vehicleName)
	if err != nil {
		return salesResponse, fmt.Errorf("failed to fetch vehicle name: %v", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return salesResponse, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Update vehicle status after the transaction is committed
	if sale.BookingDate.Year() == sale.DateOfDelivery.Year() &&
		sale.BookingDate.Month() == sale.DateOfDelivery.Month() &&
		sale.BookingDate.Day() == sale.DateOfDelivery.Day() {

		if err := r.UpdateVehicleStatus(sale.VehicleID, "rented"); err != nil {
			return salesResponse, fmt.Errorf("failed to update vehicle status: %v", err)
		}
		if err := r.UpdateSaleStatus(saleID, "active"); err != nil {
			return salesResponse, fmt.Errorf("failed to update sale status: %v", err)
		}
	}

	// Populate the response object
	salesResponse.SaleID = saleID
	salesResponse.VehicleName = vehicleName

	return salesResponse, nil
}

func (r *SaleRepository) GetSaleByID(saleID int, include []string) (*models.Sale, error) {

	sale := &models.Sale{}
	err := r.db.QueryRow(`
		SELECT s.sale_id, s.vehicle_id, s.user_id, s.customer_name,s.customer_phone,s.customer_destination,s.total_amount, s.charge_per_day, s.booking_date, 
		s.date_of_delivery, s.return_date, s.number_of_days, s.actual_date_of_delivery,s.actual_date_of_return,u.username,
		s.remark, s.status, s.created_at, s.updated_at
		FROM sales s
		LEFT JOIN users u ON s.user_id = u.id
		WHERE s.sale_id = $1 
	`, saleID).Scan(
		&sale.SaleID, &sale.VehicleID, &sale.UserID, &sale.CustomerName, &sale.CustomerPhone, &sale.Destination, &sale.TotalAmount, &sale.ChargePerDay,
		&sale.BookingDate, &sale.DateOfDelivery, &sale.ReturnDate, &sale.NumberOfDays, &sale.ActualDateofDelivery, &sale.ActualReturnDate, &sale.UserName, &sale.Remark, &sale.Status, &sale.CreatedAt, &sale.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Fetch related data based on the include parameter
	for _, inc := range include {
		switch inc {
		case "SalesCharge":
			fmt.Println("Fetching sales charges...") // Debug log
			charges, err := r.getSalesCharges(saleID)
			if err != nil {
				fmt.Printf("Error fetching sales charges: %v\n", err) // Debug log
				return nil, err
			}
			fmt.Printf("Fetched sales charges: %+v\n", charges) // Debug log
			sale.SalesCharges = charges

		case "SalesImages":
			fmt.Println("Fetching sales images...") // Debug log
			images, err := r.getSalesImages(saleID)
			if err != nil {
				fmt.Printf("Error fetching sales images: %v\n", err) // Debug log
				return nil, err
			}
			fmt.Printf("Fetched sales images: %+v\n", images) // Debug log
			sale.SalesImages = images

		case "SalesVideos":
			fmt.Println("Fetching sales videos...") // Debug log
			videos, err := r.getSalesVideos(saleID)
			if err != nil {
				fmt.Printf("Error fetching sales videos: %v\n", err) // Debug log
				return nil, err
			}
			fmt.Printf("Fetched sales videos: %+v\n", videos) // Debug log
			sale.SalesVideos = videos

		case "VehicleUsage":
			fmt.Println("Fetching vehicle usage...") // Debug log
			usage, err := r.getVehicleUsage(sale.VehicleID)
			if err != nil {
				fmt.Printf("Error fetching vehicle usage: %v\n", err) // Debug log
				return nil, err
			}
			fmt.Printf("Fetched vehicle usage: %+v\n", usage) // Debug log
			sale.VehicleUsage = usage

		case "Payments":
			fmt.Println("Fetching payments...") // Debug log
			payments, err := r.getPayments(saleID)
			if err != nil {
				fmt.Printf("Error fetching payments: %v\n", err) // Debug log
				return nil, err
			}
			fmt.Printf("Fetched payments: %+v\n", payments) // Debug log
			sale.Payments = payments
		}
	}

	return sale, nil
}
func (r *SaleRepository) GetAllSales(include []string) ([]models.Sale, error) {

	rows, err := r.db.Query(`
	SELECT s.sale_id, s.vehicle_id, s.user_id, s.customer_name, s.total_amount, s.charge_per_day, 
       s.booking_date, s.date_of_delivery, s.return_date, s.number_of_days, 
       s.remark, s.status, s.created_at, s.updated_at, u.username
	FROM sales s
	LEFT JOIN users u ON s.user_id = u.id

	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []models.Sale
	for rows.Next() {
		sale := models.Sale{}
		err := rows.Scan(
			&sale.SaleID, &sale.VehicleID, &sale.UserID, &sale.CustomerName, &sale.TotalAmount, &sale.ChargePerDay,
			&sale.BookingDate, &sale.DateOfDelivery, &sale.ReturnDate, &sale.NumberOfDays, &sale.Remark, &sale.Status, &sale.CreatedAt, &sale.UpdatedAt, &sale.UserName,
		)
		if err != nil {
			return nil, err
		}
		sales = append(sales, sale)
	}

	// Fetch related data based on the include parameter
	for i, sale := range sales {
		for _, inc := range include {
			switch inc {
			case "SalesCharge":
				charges, err := r.getSalesCharges(sale.SaleID)
				if err != nil {
					return nil, err
				}
				sales[i].SalesCharges = charges

			case "SalesImages":
				images, err := r.getSalesImages(sale.SaleID)
				if err != nil {
					return nil, err
				}
				sales[i].SalesImages = images

			case "SalesVideos":
				videos, err := r.getSalesVideos(sale.SaleID)
				if err != nil {
					return nil, err
				}
				sales[i].SalesVideos = videos

			case "VehicleUsage":
				usage, err := r.getVehicleUsage(sale.VehicleID)
				if err != nil {
					return nil, err
				}
				sales[i].VehicleUsage = usage

			case "Payments":
				payments, err := r.getPayments(sale.SaleID)
				if err != nil {
					return nil, err
				}
				sales[i].Payments = payments
			}
		}
	}

	return sales, nil
}

func (r *SaleRepository) getPayments(saleID int) ([]models.Payment, error) {
	rows, err := r.db.Query(`
		SELECT payment_id, sale_id, amount_paid, payment_date, verified_by_admin, 
		payment_type, payment_status, remark, user_id,sale_type,created_at, updated_at
		FROM payments WHERE sale_id = $1
	`, saleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(
			&payment.PaymentID, &payment.SaleID, &payment.AmountPaid, &payment.PaymentDate, &payment.VerifiedByAdmin,
			&payment.PaymentType, &payment.PaymentStatus, &payment.Remark, &payment.UserID, &payment.SaleType, &payment.CreatedAt, &payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *SaleRepository) getSalesCharges(saleID int) ([]models.SalesCharge, error) {
	rows, err := r.db.Query("SELECT charge_id, sale_id, charge_type, amount FROM sales_charges WHERE sale_id = $1", saleID)
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
	rows, err := r.db.Query("SELECT video_id, sale_id, video_url,file_name, uploaded_at FROM sales_videos WHERE sale_id = $1", saleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videos []models.SalesVideo
	for rows.Next() {
		var video models.SalesVideo
		rows.Scan(&video.VideoID, &video.SaleID, &video.VideoURL, &video.FileName, &video.UploadedAt)
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
