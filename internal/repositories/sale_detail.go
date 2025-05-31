package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"renting/internal/models"
	"time"
)

type SaleDetailRepository struct {
	db *sql.DB
}

func NewSaleDetailRepository(db *sql.DB) *SaleDetailRepository {
	return &SaleDetailRepository{db: db}
}
func (r *SaleDetailRepository) GetSalesWithFilters(filters map[string]string) ([]models.Sale, error) {
	// Base query
	query := `
		SELECT 
			s.sale_id, s.vehicle_id, s.user_id, u.username, s.customer_name, s.customer_destination, s.customer_phone, 
			s.total_amount, s.charge_per_day, s.booking_date, s.date_of_delivery, s.return_date, 
			s.number_of_days, s.remark, s.status, p.sale_type, s.payment_status,
			s.created_at, s.updated_at,
			s.actual_delivery_time_of_day, s.actual_return_time_of_day,
			p.payment_id, p.payment_type, p.amount_paid, p.payment_date, p.payment_status, 
			p.verified_by_admin, p.remark AS payment_remark, p.user_id AS payment_user_id, 
			p.created_at AS payment_created_at, p.updated_at AS payment_updated_at,
			v.vehicle_name,
			(
				SELECT json_agg(json_build_object(
					'charge_id', sc.charge_id,
					'charge_type', sc.charge_type,
					'amount', sc.amount
				))
				FROM sales_charges sc
				WHERE sc.sale_id = s.sale_id
			) as charges,
			vu.usage_id, vu.record_type, vu.fuel_range, vu.km_reading, vu.recorded_at, vu.recorded_by
		FROM sales s
		LEFT JOIN payments p ON s.sale_id = p.sale_id
		LEFT JOIN vehicles v ON s.vehicle_id = v.vehicle_id
		LEFT JOIN vehicle_usage vu ON s.sale_id = vu.sale_id
		LEFT JOIN users u ON s.user_id = u.user_id
		WHERE 1=1
	`

	// Add filters to the query
	args := []interface{}{}
	argCounter := 1

	for key, value := range filters {
		switch key {
		case "start_date":
			query += fmt.Sprintf(" AND s.booking_date >= $%d", argCounter)
			args = append(args, value)
			argCounter++
		case "end_date":
			query += fmt.Sprintf(" AND s.booking_date <= $%d", argCounter)
			args = append(args, value)
			argCounter++
		case "status":
			query += fmt.Sprintf(" AND s.status = $%d", argCounter)
			args = append(args, value)
			argCounter++
		case "verification":
			query += fmt.Sprintf(" AND p.payment_status = $%d", argCounter)
			args = append(args, value)
			argCounter++
		case "car_number":
			query += fmt.Sprintf(" AND v.vehicle_name = $%d", argCounter)
			args = append(args, value)
			argCounter++
		case "is_discount":
			query += " AND EXISTS (SELECT 1 FROM sales_charges sc WHERE sc.sale_id = s.sale_id AND sc.charge_type = 'discount')"
		case "is_washed":
			query += " AND EXISTS (SELECT 1 FROM sales_charges sc WHERE sc.sale_id = s.sale_id AND sc.charge_type = 'wash')"
		case "is_delayed":
			query += " AND EXISTS (SELECT 1 FROM sales_charges sc WHERE sc.sale_id = s.sale_id AND sc.charge_type = 'delay')"
		}
	}

	// Debug: Print the query and arguments
	fmt.Println("Generated Query:", query)
	fmt.Println("Query Arguments:", args)

	// Execute the query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query sales: %v", err)
	}
	defer rows.Close()

	// Map to group sales and related data
	salesMap := make(map[int]models.Sale)

	for rows.Next() {
		var sale models.Sale
		var payment models.Payment
		var paymentID *int
		var amountPaid *float64
		var paymentDate *time.Time
		var paymentType, paymentStatus, paymentRemark *string
		var verifiedByAdmin *bool
		var paymentUserID *int
		var paymentCreatedAt, paymentUpdatedAt *time.Time
		var carNumber *string

		var charges *json.RawMessage

		var usageID *int
		var recordType *string
		var fuelRange *float64
		var kmReading *float64
		var recordedAt *time.Time
		var recordedBy *int

		// Scan the row into variables
		var saleID sql.NullInt32  // Use sql.NullInt32 for nullable integer
		err := rows.Scan(
			&saleID, &sale.VehicleID, &sale.UserID, &sale.UserName, &sale.CustomerName, &sale.Destination, &sale.CustomerPhone,
			&sale.TotalAmount, &sale.ChargePerDay, &sale.BookingDate, &sale.DateOfDelivery, &sale.ReturnDate,
			&sale.NumberOfDays, &sale.Remark, &sale.Status, &payment.SaleType, &sale.PaymentStatus,
			&sale.CreatedAt, &sale.UpdatedAt,
			&sale.ActualDeliveryTimeOfDay, &sale.ActualReturnTimeOfDay,
			&paymentID, &paymentType, &amountPaid, &paymentDate, &paymentStatus,
			&verifiedByAdmin, &paymentRemark, &paymentUserID, &paymentCreatedAt, &paymentUpdatedAt,
			&carNumber,
			&charges,
			&usageID, &recordType, &fuelRange, &kmReading, &recordedAt, &recordedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sale: %v", err)
		}

		// Set the SaleID from the NullInt32
		if saleID.Valid {
			sale.SaleID = int(saleID.Int32)
		}

		// Debug: Print scanned values
		fmt.Printf("Scanned Sale: %+v\n", sale)
		fmt.Printf("Scanned Payment: %+v\n", payment)
		fmt.Printf("Scanned Charges: %v\n", charges)
		fmt.Printf("Scanned VehicleUsage: %+v\n", usageID)

		// If payment data exists, populate the Payment struct
		if paymentID != nil {
			var saleIDPtr *int
			if saleID.Valid {
				saleIDValue := int(saleID.Int32)
				saleIDPtr = &saleIDValue
			}
			payment = models.Payment{
				PaymentID:       *paymentID,
				PaymentType:     *paymentType,
				AmountPaid:      *amountPaid,
				PaymentDate:     *paymentDate,
				PaymentStatus:   *paymentStatus,
				VerifiedByAdmin: *verifiedByAdmin,
				Remark:          *paymentRemark,
				UserID:          paymentUserID,
				CreatedAt:       *paymentCreatedAt,
				UpdatedAt:       *paymentUpdatedAt,
				SaleID:          saleIDPtr,
			}
		}

		// If sales charge data exists, populate the SalesCharge struct
		if charges != nil {
			var chargeArray []struct {
				ChargeID   int     `json:"charge_id"`
				ChargeType string  `json:"charge_type"`
				Amount     float64 `json:"amount"`
			}
			if err := json.Unmarshal([]byte(*charges), &chargeArray); err != nil {
				return nil, fmt.Errorf("failed to unmarshal charges: %v", err)
			}
			for _, c := range chargeArray {
			salesCharge := models.SalesCharge{
					ChargeID:   c.ChargeID,
					ChargeType: c.ChargeType,
					Amount:     c.Amount,
				SaleID:     sale.SaleID,
			}
			sale.SalesCharges = append(sale.SalesCharges, salesCharge)
			}
		}

		// If vehicle usage data exists, populate the VehicleUsage struct
		if usageID != nil {
			vehicleUsage := models.VehicleUsage{
				UsageID:    *usageID,
				RecordType: *recordType,
				FuelRange:  *fuelRange,
				KmReading:  *kmReading,
				RecordedAt: *recordedAt,
				RecordedBy: *recordedBy,
			}
			sale.VehicleUsage = append(sale.VehicleUsage, vehicleUsage)
		}

		// Check if the sale already exists in the map
		if existingSale, ok := salesMap[sale.SaleID]; ok {
			// Append the payment, sales charge, and vehicle usage to the existing sale
			if paymentID != nil {
				existingSale.Payments = append(existingSale.Payments, payment)
			}
			if charges != nil {
				existingSale.SalesCharges = append(existingSale.SalesCharges, sale.SalesCharges...)
			}
			if usageID != nil {
				existingSale.VehicleUsage = append(existingSale.VehicleUsage, sale.VehicleUsage...)
			}
			salesMap[sale.SaleID] = existingSale
		} else {
			// Create a new sale entry in the map
			sale.Payments = []models.Payment{}
			sale.SalesCharges = []models.SalesCharge{}
			sale.VehicleUsage = []models.VehicleUsage{}
			if paymentID != nil {
				sale.Payments = append(sale.Payments, payment)
			}
			if charges != nil {
				sale.SalesCharges = append(sale.SalesCharges, sale.SalesCharges...)
			}
			if usageID != nil {
				sale.VehicleUsage = append(sale.VehicleUsage, models.VehicleUsage{
					UsageID:    *usageID,
					RecordType: *recordType,
					FuelRange:  *fuelRange,
					KmReading:  *kmReading,
					RecordedAt: *recordedAt,
					RecordedBy: *recordedBy,
				})
			}
			salesMap[sale.SaleID] = sale
		}
	}

	// Convert the map to a slice
	var sales []models.Sale
	for _, sale := range salesMap {
		sales = append(sales, sale)
	}

	// Debug: Print final sales data
	fmt.Printf("Final Sales Data: %+v\n", sales)

	return sales, nil
}

