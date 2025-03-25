package repositories

import (
	"database/sql"
	"fmt"
	"renting/internal/models"
	"time"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

type SaleFilter struct {
	SaleID        *int
	PaymentStatus *string
	StartDate     *time.Time
	EndDate       *time.Time
	CustomerName  *string
	SaleStatus    *string
	VerifiedBy    *string
}

func (r *PaymentRepository) GetPaymentsWithSales(filter SaleFilter, limit int, offset int) ([]models.PaymentWithSale, error) {
	query := `
		SELECT 
			p.payment_id,
			p.payment_type,
			p.amount_paid,
			p.payment_date,
			p.payment_status,
			p.verified_by_admin,
			p.remark,
			p.created_at,
			p.updated_at,
			s.sale_id, 
			s.vehicle_id, 
			s.user_id, 
			s.customer_name, 
			s.customer_destination, 
			s.customer_phone, 
			s.total_amount, 
			s.charge_per_day, 
			s.booking_date, 
			s.date_of_delivery, 
			s.return_date, 
			s.is_damaged, 
			s.is_washed, 
			s.is_delayed, 
			s.number_of_days, 
			s.remark, 
			s.status, 
			s.created_at, 
			s.updated_at,
			v.vehicle_id,
			v.vehicle_type_id,
			v.vehicle_name,
			v.vehicle_model,
			v.vehicle_registration_number,
			v.is_available,
			v.image_name,
			v.status
		FROM payments p
		LEFT JOIN sales s ON p.sale_id = s.sale_id
		LEFT JOIN vehicles v ON s.vehicle_id = v.vehicle_id
		WHERE 1=1
	`

	args := []interface{}{}
	argCounter := 1

	if filter.SaleID != nil {
		query += fmt.Sprintf(" AND p.payment_id = $%d", argCounter)
		args = append(args, *filter.SaleID)
		argCounter++
	}

	if filter.PaymentStatus != nil {
		query += fmt.Sprintf(" AND p.payment_status = $%d", argCounter)
		args = append(args, *filter.PaymentStatus)
		argCounter++
	}

	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND p.payment_date >= $%d", argCounter)
		args = append(args, *filter.StartDate)
		argCounter++
	}

	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND p.payment_date <= $%d", argCounter)
		args = append(args, *filter.EndDate)
		argCounter++
	}

	if filter.CustomerName != nil {
		query += fmt.Sprintf(" AND s.customer_name ILIKE $%d", argCounter)
		args = append(args, *filter.CustomerName)
		argCounter++
	}

	if filter.SaleStatus != nil {
		query += fmt.Sprintf(" AND s.status = $%d", argCounter)
		args = append(args, *filter.SaleStatus)
		argCounter++
	}

	if filter.VerifiedBy != nil {
		query += fmt.Sprintf(" AND s.user_id = $%d", argCounter)
		args = append(args, *filter.VerifiedBy)
		argCounter++
	}

	// Add ORDER BY clause to sort by payment_date in descending order
	query += " ORDER BY p.payment_date DESC"

	// Add LIMIT and OFFSET
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.PaymentWithSale
	for rows.Next() {
		var payment models.PaymentWithSale
		err := rows.Scan(
			&payment.PaymentID,
			&payment.PaymentType,
			&payment.AmountPaid,
			&payment.PaymentDate,
			&payment.PaymentStatus,
			&payment.VerifiedByAdmin,
			&payment.Remark,
			&payment.CreatedAt,
			&payment.UpdatedAt,
			&payment.Sale.SaleID,
			&payment.Sale.VehicleID,
			&payment.Sale.UserID,
			&payment.Sale.CustomerName,
			&payment.Sale.Destination,
			&payment.Sale.CustomerPhone,
			&payment.Sale.TotalAmount,
			&payment.Sale.ChargePerDay,
			&payment.Sale.BookingDate,
			&payment.Sale.DateOfDelivery,
			&payment.Sale.ReturnDate,
			&payment.Sale.IsDamaged,
			&payment.Sale.IsWashed,
			&payment.Sale.IsDelayed,
			&payment.Sale.NumberOfDays,
			&payment.Sale.Remark,
			&payment.Sale.Status,
			&payment.Sale.CreatedAt,
			&payment.Sale.UpdatedAt,
			&payment.Sale.Vehicle.VehicleID,
			&payment.Sale.Vehicle.VehicleTypeID,
			&payment.Sale.Vehicle.VehicleName,
			&payment.Sale.Vehicle.VehicleModel,
			&payment.Sale.Vehicle.VehicleRegistrationNumber,
			&payment.Sale.Vehicle.IsAvailable,
			&payment.Sale.Vehicle.SalesImage,
			&payment.Sale.Vehicle.Status,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}
