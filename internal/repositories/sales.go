package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"renting/internal/models"
	"strings"
	"time"
)

type SaleRepository struct {
	db *sql.DB
}

func NewSaleRepository(db *sql.DB) *SaleRepository {
	return &SaleRepository{db: db}
}

func (r *SaleRepository) GetFutureBookings(vehicleID int, dateOfDelivery time.Time) (any, any) {
	panic("unimplemented")
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
	_, err := r.db.Exec(`
		UPDATE sales
		SET status = $2
		WHERE sale_id = $1
	`, saleID, status)
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
	bookingDate := time.Now()

	// Check vehicle availability
	available, err := r.checkVehicleAvailability(sale.VehicleID, sale.DateOfDelivery, sale.ReturnDate)
	if err != nil {
		return models.SaleSubmitResponse{}, fmt.Errorf("failed to check vehicle availability: %v", err)
	}
	if !available {
		return models.SaleSubmitResponse{}, errors.New("vehicle is not available for the selected dates")
	}

	// Set actual delivery date if same day
	var actualDeliveryDate *time.Time
	isFutureBooking := sale.DateOfDelivery.After(bookingDate)
	if !isFutureBooking {
		actualDeliveryDate = &bookingDate
	}

	// Calculate payment status
	paymentStatus := "unpaid"
	totalVerifiedPaid := 0.0
	for _, payment := range sale.Payments {
		if payment.VerifiedByAdmin { // Only count verified payments
			totalVerifiedPaid += payment.AmountPaid
		}
	}

	if totalVerifiedPaid >= sale.TotalAmount {
		paymentStatus = "paid"
	} else if totalVerifiedPaid > 0 {
		paymentStatus = "partial"
	}

	var salesResponse models.SaleSubmitResponse
	tx, err := r.db.Begin()
	if err != nil {
		return salesResponse, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert sale record
	var saleID int
	err = tx.QueryRow(`
		INSERT INTO sales (
			vehicle_id, user_id, customer_name, total_amount, charge_per_day, charge_half_day, booking_date, 
			date_of_delivery, return_date, number_of_days, remark, status, customer_destination, 
			customer_phone, actual_date_of_delivery, payment_status, delivery_time_of_day, return_time_of_day,
			actual_delivery_time_of_day, actual_return_time_of_day, is_short_term_rental, full_days, half_days,
			is_damaged, is_washed, is_delayed, discount, other_charges, is_future_booking
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29)
		RETURNING sale_id
	`,
		sale.VehicleID, sale.UserID, sale.CustomerName, sale.TotalAmount, sale.ChargePerDay, sale.ChargeHalfDay, bookingDate,
		sale.DateOfDelivery, sale.ReturnDate, sale.NumberOfDays, sale.Remark, sale.Status, sale.Destination,
		sale.CustomerPhone, actualDeliveryDate, paymentStatus, sale.DeliveryTimeOfDay, sale.ReturnTimeOfDay,
		sale.ActualDeliveryTimeOfDay.String, sale.ActualReturnTimeOfDay.String,
		sale.IsShortTermRental, sale.FullDays, sale.HalfDays,
		sale.IsDamaged, sale.IsWashed, sale.IsDelayed,
		sale.Discount, sale.OtherCharges, isFutureBooking,
	).Scan(&saleID)
	if err != nil {
		return salesResponse, fmt.Errorf("failed to insert sale: %v", err)
	}

	// Insert related records in a transaction
	if err := r.insertSalesCharges(tx, saleID, sale.SalesCharges); err != nil {
		return salesResponse, fmt.Errorf("failed to insert sales charges: %v", err)
	}

	if err := r.insertSalesImages(tx, saleID, sale.SalesImages); err != nil {
		return salesResponse, fmt.Errorf("failed to insert sales images: %v", err)
	}

	if err := r.insertSalesVideos(tx, saleID, sale.SalesVideos); err != nil {
		return salesResponse, fmt.Errorf("failed to insert sales videos: %v", err)
	}

	if err := r.insertVehicleUsage(tx, saleID, sale.VehicleUsage, sale.UserID); err != nil {
		return salesResponse, fmt.Errorf("failed to insert vehicle usage: %v", err)
	}

	// Insert payments if any
	if len(sale.Payments) > 0 {
		if err := r.insertPayments(tx, saleID, sale.Payments, sale.UserID, sale.VehicleUsage); err != nil {
			return salesResponse, fmt.Errorf("failed to insert payments: %v", err)
		}
	}

	// Only update vehicle status to 'rented' if it's not a future booking
	if !isFutureBooking {
		if err := r.UpdateVehicleStatus(sale.VehicleID, "rented"); err != nil {
			return salesResponse, fmt.Errorf("failed to update vehicle status: %v", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return salesResponse, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Get vehicle name for response
	var vehicleName string
	err = r.db.QueryRow(`
		SELECT vehicle_name 
		FROM vehicles 
		WHERE vehicle_id = $1
	`, sale.VehicleID).Scan(&vehicleName)
	if err != nil {
		return salesResponse, fmt.Errorf("failed to fetch vehicle name: %v", err)
	}

	// Prepare response
	salesResponse = models.SaleSubmitResponse{
		SaleID:      saleID,
		VehicleName: vehicleName,
	}

	return salesResponse, nil
}

// Helper methods for inserting related records
func (r *SaleRepository) insertSalesCharges(tx *sql.Tx, saleID int, charges []models.SalesCharge) error {
	for _, charge := range charges {
		_, err := tx.Exec(`
			INSERT INTO sales_charges (sale_id, charge_type, amount)
			VALUES ($1, $2, $3)
		`, saleID, charge.ChargeType, charge.Amount)
		if err != nil {
			return fmt.Errorf("failed to insert sales charge: %v", err)
		}
	}
	return nil
}

func (r *SaleRepository) insertSalesImages(tx *sql.Tx, saleID int, images []models.SalesImage) error {
	for _, image := range images {
		_, err := tx.Exec(`
			INSERT INTO sales_images (sale_id, image_url)
			VALUES ($1, $2)
		`, saleID, image.ImageURL)
		if err != nil {
			return fmt.Errorf("failed to insert sales image: %v", err)
		}
	}
	return nil
}

func (r *SaleRepository) insertSalesVideos(tx *sql.Tx, saleID int, videos []models.SalesVideo) error {
	for _, video := range videos {
		_, err := tx.Exec(`
			INSERT INTO sales_videos (sale_id, video_url)
			VALUES ($1, $2)
		`, saleID, video.VideoURL)
		if err != nil {
			return fmt.Errorf("failed to insert sales video: %v", err)
		}
	}
	return nil
}

func (r *SaleRepository) insertVehicleUsage(tx *sql.Tx, saleID int, usageRecords []models.VehicleUsage, userID int) error {
	for _, usage := range usageRecords {
		// Insert the usage record
		_, err := tx.Exec(`
			INSERT INTO vehicle_usage (sale_id, vehicle_id, record_type, fuel_range, km_reading, recorded_at, recorded_by)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, saleID, usage.VehicleID, usage.RecordType, usage.FuelRange, usage.KmReading, usage.RecordedAt, userID)
		if err != nil {
			return fmt.Errorf("failed to insert vehicle usage record: %v", err)
		}

	}
	return nil
}

func (r *SaleRepository) insertPayments(tx *sql.Tx, saleID int, payments []models.Payment, userID int, usageRecords []models.VehicleUsage) error {
	for _, payment := range payments {
		if payment.SaleType == "" {
			payment.SaleType = determineSaleType(usageRecords)
		}

		if err := payment.Validate(); err != nil {
			return fmt.Errorf("invalid payment: %v", err)
		}

		_, err := tx.Exec(`
			INSERT INTO payments (
				sale_id, amount_paid, payment_date, verified_by_admin, 
				payment_type, payment_status, remark, sale_type
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, saleID, payment.AmountPaid, payment.PaymentDate, payment.VerifiedByAdmin,
			payment.PaymentType, payment.PaymentStatus, payment.Remark, payment.SaleType)
		if err != nil {
			return fmt.Errorf("failed to insert payment: %v", err)
		}
	}
	return nil
}

func determineSaleType(usageRecords []models.VehicleUsage) string {
	for _, usage := range usageRecords {
		if usage.RecordType == "delivery" {
			return models.TypeDelivery
		} else if usage.RecordType == "return" {
			return models.TypeReturn
		}
	}
	return models.TypeBooking
}
func (r *SaleRepository) GetSaleByID(saleID int, include []string) (*models.Sale, error) {
	sale := &models.Sale{}
	err := r.db.QueryRow(`
        SELECT s.sale_id, s.vehicle_id, s.user_id, s.customer_name, s.customer_phone, s.customer_destination, 
               s.total_amount, s.charge_per_day, s.charge_half_day, s.booking_date, s.date_of_delivery, s.return_date, 
               s.number_of_days, s.actual_date_of_delivery, s.actual_date_of_return, u.username, s.payment_status, 
               s.remark, s.status, s.created_at, s.updated_at, s.delivery_time_of_day, s.return_time_of_day,
               s.actual_delivery_time_of_day, s.actual_return_time_of_day, s.other_charges, s.is_damaged, s.is_washed,
               s.is_delayed, s.is_short_term_rental, s.full_days, s.half_days, s.discount, s.modified_by, s.is_future_booking
        FROM sales s
        LEFT JOIN users u ON s.user_id = u.id
        WHERE s.sale_id = $1
    `, saleID).Scan(
		&sale.SaleID, &sale.VehicleID, &sale.UserID, &sale.CustomerName, &sale.CustomerPhone, &sale.Destination,
		&sale.TotalAmount, &sale.ChargePerDay, &sale.ChargeHalfDay, &sale.BookingDate, &sale.DateOfDelivery, &sale.ReturnDate,
		&sale.NumberOfDays, &sale.ActualDateOfDelivery, &sale.ActualReturnDate, &sale.UserName, &sale.PaymentStatus,
		&sale.Remark, &sale.Status, &sale.CreatedAt, &sale.UpdatedAt, &sale.DeliveryTimeOfDay, &sale.ReturnTimeOfDay,
		&sale.ActualDeliveryTimeOfDay, &sale.ActualReturnTimeOfDay, &sale.OtherCharges, &sale.IsDamaged, &sale.IsWashed,
		&sale.IsDelayed, &sale.IsShortTermRental, &sale.FullDays, &sale.HalfDays, &sale.Discount, &sale.ModifiedBy,
		&sale.IsFutureBooking,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil if sale not found
		}
		return nil, fmt.Errorf("failed to fetch sale: %v", err)
	}

	// Load related data based on include parameter
	for _, inc := range include {
		switch inc {
		case "SalesCharge":
			charges, err := r.getSalesCharges(saleID)
			if err != nil {
				return nil, err
			}
			sale.SalesCharges = charges
		case "SalesImages":
			images, err := r.getSalesImages(saleID)
			if err != nil {
				return nil, err
			}
			sale.SalesImages = images
		case "SalesVideos":
			videos, err := r.getSalesVideos(saleID)
			if err != nil {
				return nil, err
			}
			sale.SalesVideos = videos
		case "VehicleUsage":
			usage, err := r.getVehicleUsage(sale.VehicleID)
			if err != nil {
				return nil, err
			}
			sale.VehicleUsage = usage
		case "Payments":
			payments, err := r.getPayments(saleID)
			if err != nil {
				return nil, err
			}
			sale.Payments = payments
		case "vehicle":
			vehicle, err := r.getVehicle(sale.VehicleID)
			if err != nil {
				return nil, err
			}
			sale.Vehicle = vehicle
		}
	}

	return sale, nil
}
func (r *SaleRepository) getVehicle(vehicleID int) (*models.Vehicle, error) {
	query := `
        SELECT 
            vehicle_id, vehicle_type_id, vehicle_name, vehicle_model, 
            status, vehicle_registration_number, is_available, 
            image_name, created_at, updated_at
        FROM vehicles
        WHERE vehicle_id = $1
    `

	var vehicle models.Vehicle
	err := r.db.QueryRow(query, vehicleID).Scan(
		&vehicle.VehicleID,
		&vehicle.VehicleTypeID,
		&vehicle.VehicleName,
		&vehicle.VehicleModel,
		&vehicle.Status,
		&vehicle.VehicleRegistrationNumber,
		&vehicle.IsAvailable,
		&vehicle.ImageName,
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch vehicle: %v", err)
	}

	return &vehicle, nil
}

// GetAllSales retrieves all sales with optional related data
func (r *SaleRepository) GetAllSales(include []string) ([]models.Sale, error) {
	rows, err := r.db.Query(`
		SELECT s.sale_id, s.vehicle_id, s.user_id, s.customer_name, s.total_amount, s.charge_per_day, 
		       s.booking_date, s.date_of_delivery, s.return_date, s.number_of_days,
		       s.remark, s.status, s.created_at, s.updated_at, u.username, s.payment_status
		FROM sales s
		LEFT JOIN users u ON s.user_id = u.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sales []models.Sale
	for rows.Next() {
		var sale models.Sale
		err := rows.Scan(
			&sale.SaleID, &sale.VehicleID, &sale.UserID, &sale.CustomerName, &sale.TotalAmount, &sale.ChargePerDay,
			&sale.BookingDate, &sale.DateOfDelivery, &sale.ReturnDate, &sale.NumberOfDays, &sale.Remark,
			&sale.Status, &sale.CreatedAt, &sale.UpdatedAt, &sale.UserName, &sale.PaymentStatus,
		)
		if err != nil {
			return nil, err
		}
		sales = append(sales, sale)
	}

	// Load related data for each sale
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

// Helper methods for fetching related records
func (r *SaleRepository) getSalesCharges(saleID int) ([]models.SalesCharge, error) {
	rows, err := r.db.Query(`
		SELECT charge_id, sale_id, charge_type, amount 
		FROM sales_charges 
		WHERE sale_id = $1
	`, saleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var charges []models.SalesCharge
	for rows.Next() {
		var charge models.SalesCharge
		err := rows.Scan(&charge.ChargeID, &charge.SaleID, &charge.ChargeType, &charge.Amount)
		if err != nil {
			return nil, err
		}
		charges = append(charges, charge)
	}
	return charges, nil
}

func (r *SaleRepository) getSalesImages(saleID int) ([]models.SalesImage, error) {
	rows, err := r.db.Query(`
		SELECT image_id, sale_id, image_url, uploaded_at 
		FROM sales_images 
		WHERE sale_id = $1
	`, saleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.SalesImage
	for rows.Next() {
		var image models.SalesImage
		err := rows.Scan(&image.ImageID, &image.SaleID, &image.ImageURL, &image.UploadedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	return images, nil
}

func (r *SaleRepository) getSalesVideos(saleID int) ([]models.SalesVideo, error) {
	rows, err := r.db.Query(`
        SELECT video_id, sale_id, video_url, file_name, uploaded_at 
        FROM sales_videos 
        WHERE sale_id = $1
    `, saleID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sales videos: %v", err)
	}
	defer rows.Close()

	var videos []models.SalesVideo
	for rows.Next() {
		var video models.SalesVideo
		err := rows.Scan(
			&video.VideoID,
			&video.SaleID,
			&video.VideoURL,
			&video.FileName,
			&video.UploadedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sales video: %v", err)
		}
		videos = append(videos, video)
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sales videos: %v", err)
	}

	return videos, nil
}
func (r *SaleRepository) getVehicleUsage(vehicleID int) ([]models.VehicleUsage, error) {
	rows, err := r.db.Query(`
		SELECT usage_id, vehicle_id, record_type, fuel_range, km_reading, recorded_at, recorded_by 
		FROM vehicle_usage 
		WHERE vehicle_id = $1
	`, vehicleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usageRecords []models.VehicleUsage
	for rows.Next() {
		var usage models.VehicleUsage
		err := rows.Scan(&usage.UsageID, &usage.VehicleID, &usage.RecordType, &usage.FuelRange, &usage.KmReading,
			&usage.RecordedAt, &usage.RecordedBy)
		if err != nil {
			return nil, err
		}
		usageRecords = append(usageRecords, usage)
	}
	return usageRecords, nil
}

func (r *SaleRepository) getPayments(saleID int) ([]models.Payment, error) {
	rows, err := r.db.Query(`
        SELECT p.payment_id, p.sale_id, p.amount_paid, p.payment_date, p.verified_by_admin, 
               p.payment_type, p.payment_status, p.remark, p.user_id, p.sale_type, p.created_at, p.updated_at,
               u.username
        FROM payments p
        LEFT JOIN users u ON p.user_id = u.id
        WHERE p.sale_id = $1
    `, saleID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payments: %v", err)
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(
			&payment.PaymentID, &payment.SaleID, &payment.AmountPaid, &payment.PaymentDate, &payment.VerifiedByAdmin,
			&payment.PaymentType, &payment.PaymentStatus, &payment.Remark, &payment.UserID, &payment.SaleType,
			&payment.CreatedAt, &payment.UpdatedAt, &payment.Username,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %v", err)
		}
		payments = append(payments, payment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating payments: %v", err)
	}
	return payments, nil
}
func (r *SaleRepository) GetSales(filters models.SaleFilter) ([]models.Sale, error) {
	baseQuery := `
        SELECT 
			s.sale_id, s.vehicle_id, s.user_id, s.customer_name, s.customer_destination, 
			s.customer_phone, s.total_amount, s.discount, s.other_charges, s.payment_status, 
			s.payment_method, s.booking_date, s.delivery_date, s.return_date, 
			s.delivery_time_of_day, s.return_time_of_day, s.actual_delivery_time_of_day, 
			s.actual_return_time_of_day, s.status, s.created_at, s.updated_at
        FROM sales s
		WHERE 1=1
    `

	// Build WHERE clause from filters
	var whereClauses []string
	var args []interface{}
	argIndex := 1

	if filters.Status != "" {
				whereClauses = append(whereClauses, fmt.Sprintf("s.status = $%d", argIndex))
		args = append(args, filters.Status)
				argIndex++
			}

	if filters.ActualDateOfDelivery != nil {
				whereClauses = append(whereClauses, fmt.Sprintf("s.actual_date_of_delivery = $%d", argIndex))
		args = append(args, filters.ActualDateOfDelivery)
				argIndex++
			}

	if filters.DateOfDeliveryBefore != nil {
			whereClauses = append(whereClauses, fmt.Sprintf("s.date_of_delivery < $%d", argIndex))
		args = append(args, filters.DateOfDeliveryBefore)
			argIndex++
	}

	if filters.CustomerName != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("s.customer_name ILIKE $%d", argIndex))
		args = append(args, "%"+filters.CustomerName+"%")
			argIndex++
	}

	if filters.VehicleID != 0 {
			whereClauses = append(whereClauses, fmt.Sprintf("s.vehicle_id = $%d", argIndex))
		args = append(args, filters.VehicleID)
			argIndex++
	}

	if len(whereClauses) > 0 {
		baseQuery += " AND " + strings.Join(whereClauses, " AND ")
	}

	// Add sorting
	if filters.Sort != "" {
		baseQuery += fmt.Sprintf(" ORDER BY %s", filters.Sort)
	} else {
		baseQuery += " ORDER BY s.created_at DESC" // Default sort
	}

	// Add limit and offset only if limit > 0
	if filters.Limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
		args = append(args, filters.Limit, filters.Offset)
	}

	// Execute the main query
	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sales: %v", err)
	}
	defer rows.Close()

	var sales []models.Sale
	for rows.Next() {
		var sale models.Sale
		err := rows.Scan(
			&sale.SaleID, &sale.VehicleID, &sale.UserID, &sale.CustomerName, &sale.Destination,
			&sale.CustomerPhone, &sale.TotalAmount, &sale.Discount, &sale.OtherCharges, &sale.PaymentStatus,
			&sale.PaymentMethod, &sale.BookingDate, &sale.DateOfDelivery, &sale.ReturnDate,
			&sale.DeliveryTimeOfDay, &sale.ReturnTimeOfDay, &sale.ActualDeliveryTimeOfDay,
			&sale.ActualReturnTimeOfDay, &sale.Status, &sale.CreatedAt, &sale.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sale: %v", err)
		}
		sales = append(sales, sale)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating sales: %v", err)
	}

	return sales, nil
}

// isAdmin checks if the user is an admin
func (r *SaleRepository) isAdmin(userID int) (bool, error) {
	query := `
		SELECT is_admin
		FROM users
		WHERE id = $1
	`
	var isAdmin bool
	err := r.db.QueryRow(query, userID).Scan(&isAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("user not found")
		}
		return false, err
	}
	return isAdmin, nil
}

// isVerified checks if the sale has a verified payment
func (r *SaleRepository) isVerified(saleID int) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM payments
			WHERE sale_id = $1 AND verified_by_admin = true
		)
	`
	var isVerified bool
	err := r.db.QueryRow(query, saleID).Scan(&isVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			// No payment exists, assume not verified
			return false, nil
		}
		return false, fmt.Errorf("failed to check payment verification: %v", err)
	}
	return isVerified, nil
}

func (r *SaleRepository) UpdateSaleByUserID(saleID, userID int, updates map[string]interface{}) error {
	// First check if the user is an admin
	isAdmin, err := r.isAdmin(userID)
	if err != nil {
		return err
	}

	// Check if the sale exists
	var saleUserID int
	err = r.db.QueryRow(`
		SELECT user_id 
		FROM sales 
		WHERE sale_id = $1
	`, saleID).Scan(&saleUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("sale with ID %d not found", saleID)
		}
		return fmt.Errorf("failed to check sale: %v", err)
	}

	// If user is not an admin, they can only update their own sales
	if !isAdmin && saleUserID != userID {
		return fmt.Errorf("cannot update sale with ID %d: user %d is not an admin and does not own this sale", saleID, userID)
	}

	// Define allowed fields for update
	allowedFields := map[string]bool{
		"status":                    true,
		"payment_status":            true,
		"remark":                    true,
		"customer_name":             true,
		"customer_phone":            true,
		"customer_destination":      true,
		"total_amount":              true,
		"charge_per_day":            true,
		"charge_half_day":           true,
		"date_of_delivery":          true,
		"return_date":               true,
		"actual_date_of_delivery":   true,
		"actual_date_of_return":     true,
		"number_of_days":            true,
		"vehicle_id":                true,
		"actual_delivery_time_of_day": true,
		"actual_return_time_of_day":   true,
		"delivery_time_of_day":        true,
		"return_time_of_day":          true,
		"booking_date":                true,
		"is_future_booking":           true,
	}

	// Build the dynamic UPDATE query
	var setClauses []string
	args := []interface{}{saleID}
	argIndex := 2 // Start from 2 since saleID is $1

	// Process updates
	for field, value := range updates {
		if _, ok := allowedFields[field]; !ok {
			return fmt.Errorf("unsupported field for update: %s", field)
		}

		switch field {
		case "status", "payment_status", "remark", "customer_name", "customer_phone", "customer_destination", 
			"actual_delivery_time_of_day", "actual_return_time_of_day", "delivery_time_of_day", "return_time_of_day":
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, argIndex))
			args = append(args, value)
			argIndex++
		case "total_amount", "charge_per_day", "charge_half_day", "number_of_days":
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, argIndex))
			args = append(args, value.(float64))
			argIndex++
		case "date_of_delivery", "return_date", "actual_date_of_delivery", "actual_date_of_return", "booking_date":
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, argIndex))
			args = append(args, value.(time.Time))
			argIndex++
		case "vehicle_id":
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, argIndex))
			args = append(args, value.(int))
			argIndex++
		case "is_future_booking":
			setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, argIndex))
			args = append(args, value.(bool))
			argIndex++
		}
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields provided to update")
	}

	// Construct the UPDATE query
	query := fmt.Sprintf(`
        UPDATE sales 
        SET %s, updated_at = NOW()
        WHERE sale_id = $1
    `, strings.Join(setClauses, ", "))

	// Execute the update
	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update sale: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no sale updated for sale_id %d: sale may not exist or values unchanged", saleID)
	}

	return nil
}

func (r *SaleRepository) AddCharge(saleID int, chargeType string, amount float64, remark string) error {
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
	`, saleID, chargeType, amount)
	if err != nil {
		return fmt.Errorf("failed to insert charge: %v", err)
	}

	// Update the sale's total amount based on charge type
	var updateQuery string
	switch chargeType {
	case "discount":
		updateQuery = `UPDATE sales SET discount = discount + $1 WHERE sale_id = $2`
	case "wash":
		updateQuery = `UPDATE sales SET other_charges = other_charges + $1 WHERE sale_id = $2`
	case "damage":
		updateQuery = `UPDATE sales SET other_charges = other_charges + $1 WHERE sale_id = $2`
	default:
		return fmt.Errorf("invalid charge type: %s", chargeType)
	}

	_, err = tx.Exec(updateQuery, amount, saleID)
	if err != nil {
		return fmt.Errorf("failed to update sale: %v", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
