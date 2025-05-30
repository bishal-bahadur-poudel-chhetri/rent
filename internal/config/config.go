package config

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config holds application configuration
type Config struct {
	DBConnStr         string
	JWTSecret         string
	TokenExpiry       time.Duration
	ServerAddress     string
	R2Endpoint        string // R2 endpoint URL
	R2Region          string // R2 region (e.g., "auto")
	R2AccessKeyID     string // R2 access key ID
	R2SecretAccessKey string // R2 secret access key
	R2BucketName      string
}

// LoadConfig loads application configuration from environment variables
func LoadConfig() (*Config, error) {
	return &Config{
		DBConnStr:     getEnv("DATABASE_URL", "postgres://myuser:mypassword@3.7.55.10:5432/mydatabase?sslmode=disable"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		TokenExpiry: time.Hour * 24 * 365, // 24 hours
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),

		// R2 Configuration
		R2Endpoint:        getEnv("R2_ENDPOINT", "https://472fa3e6a29c7eed2cdadcc070098155.r2.cloudflarestorage.com"),
		R2Region:          getEnv("R2_REGION", "auto"),
		R2AccessKeyID:     getEnv("R2_ACCESS_KEY_ID", "9d20e974707b5d0d186f29055252bdb4"),
		R2SecretAccessKey: getEnv("R2_SECRET_ACCESS_KEY", "41731c6315b55507b61757cb6b40ab6b24c04f6edf5d74f64fcbc4cc7bec7171"),
		R2BucketName:      getEnv("R2_BUCKET_NAME", "test"),
	}, nil
}

// ConnectDB establishes a connection to the database
func ConnectDB(connStr string) (*gorm.DB, error) {
	// First connect with standard lib to initialize schema
	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Check database connection
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	// Initialize database schema
	err = initDB(sqlDB)
	if err != nil {
		return nil, err
	}

	// Close the sql.DB connection as we'll use GORM
	sqlDB.Close()

	// Now connect with GORM
	gormDB, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

// initDB creates necessary tables if they don't exist
// initDB creates necessary tables if they don't exist
func initDB(db *sql.DB) error {
	// Create companies table
	companiesQuery := `
	CREATE TABLE IF NOT EXISTS companies (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		company_code VARCHAR(50) NOT NULL UNIQUE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(companiesQuery)
	if err != nil {
		return err
	}

	// Create users table with company_id foreign key
	usersQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(100) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		is_admin BOOLEAN DEFAULT FALSE,
		company_id INTEGER NOT NULL REFERENCES companies(id),
		mobile_number VARCHAR(20) NOT NULL UNIQUE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = db.Exec(usersQuery)
	if err != nil {
		return err
	}

	// Create vehicle_types table
	vehicleTypesQuery := `
	CREATE TABLE IF NOT EXISTS vehicle_types (
		vehicle_type_id SERIAL PRIMARY KEY,
		vehicle_type_name VARCHAR(255) NOT NULL,
		vehicle_model VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = db.Exec(vehicleTypesQuery)
	if err != nil {
		return err
	}

	// Create vehicles table
	vehiclesQuery := `
	CREATE TABLE IF NOT EXISTS vehicles (
		vehicle_id SERIAL PRIMARY KEY,
		vehicle_type_id INT NOT NULL,
		vehicle_name VARCHAR(255) NOT NULL,
		vehicle_model VARCHAR(255) NOT NULL,
		status VARCHAR(50) DEFAULT 'available' CHECK (status IN ('available', 'rented', 'under_maintenance')),
		vehicle_registration_number VARCHAR(50) UNIQUE NOT NULL,
		is_available BOOLEAN DEFAULT TRUE,
		image_name VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (vehicle_type_id) REFERENCES vehicle_types(vehicle_type_id) ON DELETE CASCADE
	);
	`
	_, err = db.Exec(vehiclesQuery)
	if err != nil {
		return err
	}

	salesQuery := `
	CREATE TABLE IF NOT EXISTS sales (
		sale_id SERIAL PRIMARY KEY,
		vehicle_id INTEGER REFERENCES vehicles(vehicle_id),
		user_id INTEGER REFERENCES users(user_id),
		customer_name VARCHAR(255) NOT NULL,
		customer_destination VARCHAR(255) NOT NULL,
		customer_phone VARCHAR(20) NOT NULL,
		total_amount DECIMAL(10,2) NOT NULL,
		discount DECIMAL(10,2) DEFAULT 0,
		other_charges DECIMAL(10,2) DEFAULT 0,
		payment_status VARCHAR(50) NOT NULL,
		payment_method VARCHAR(50) NOT NULL,
		booking_date TIMESTAMP NOT NULL,
		delivery_date TIMESTAMP NOT NULL,
		return_date TIMESTAMP NOT NULL,
		delivery_time_of_day VARCHAR(50) NOT NULL,
		return_time_of_day VARCHAR(50) NOT NULL,
		actual_delivery_time_of_day VARCHAR(50),
		actual_return_time_of_day VARCHAR(50),
		status VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
`
	_, err = db.Exec(salesQuery)
	if err != nil {
		return err
	}
	// Create payments table
	paymentsQuery := `
	CREATE TABLE IF NOT EXISTS payments (
		payment_id SERIAL PRIMARY KEY,
		sale_id INT NOT NULL REFERENCES sales(sale_id) ON DELETE CASCADE,
		payment_stage VARCHAR(20) NOT NULL CHECK (
			payment_stage IN (
				'booking_deposit',    
				'delivery_payment',  
				'return_payment',
				'refund',    
				'other'  
			)
		),
		amount_paid DECIMAL(10,2) NOT NULL,
		payment_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		payment_method VARCHAR(20) NOT NULL CHECK (
			payment_method IN (
				'cash',
				'credit_card',
				'debit_card',
				'bank_transfer',
				'mobile_payment'
			)
		),
		payment_status VARCHAR(20) NOT NULL CHECK (
			payment_status IN (
				'Pending',    -- Payment initiated but not confirmed
				'Completed',   -- Payment successfully received
				'Failed',      -- Payment failed/declined
				'Refunded'     -- Payment was refunded
			)
		),
		collected_by INT REFERENCES users(id),  -- Staff member who collected payment
		verified_by INT REFERENCES users(id),   -- Admin who verified payment
		verified_at TIMESTAMP,                  -- When payment was verified
		receipt_number VARCHAR(50),             -- Official receipt number
		remark TEXT,                            -- Additional notes
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		
		-- Ensure verification is recorded if marked as verified
		CONSTRAINT verification_check CHECK (
			(payment_status = 'Completed' AND verified_by IS NOT NULL AND verified_at IS NOT NULL)
			OR (payment_status != 'Completed')
		)
	);
	`
	_, err = db.Exec(paymentsQuery)
	if err != nil {
		return err
	}

	vehicle_usage := `
CREATE TABLE IF NOT EXISTS vehicle_usage (
    usage_id SERIAL PRIMARY KEY,
    vehicle_id INT NOT NULL,
	sale_id int,
    record_type VARCHAR(50) NOT NULL CHECK (record_type IN ('delivery', 'return')),
    fuel_range DECIMAL(10, 2) NOT NULL, -- Fuel level at delivery or return
    km_reading DECIMAL(10, 2) NOT NULL, -- KM reading at delivery or return
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- When the data was recorded
    recorded_by INT, 
    FOREIGN KEY (vehicle_id) REFERENCES vehicles(vehicle_id) ON DELETE CASCADE,
	FOREIGN KEY (sale_id) REFERENCES sales(sale_id) ON DELETE CASCADE
);
`
	_, err = db.Exec(vehicle_usage)
	if err != nil {
		return err
	}

	// Create vehicle servicing table
	vehicle_servicing := `
CREATE TABLE IF NOT EXISTS vehicle_servicing (
    servicing_id SERIAL PRIMARY KEY,
    vehicle_id INT NOT NULL,
    current_km DECIMAL(10, 2) NOT NULL,
    next_servicing_km DECIMAL(10, 2) NOT NULL,
    servicing_interval_km DECIMAL(10, 2) NOT NULL,
    is_servicing_due BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) DEFAULT 'pending', -- pending, in_progress, completed
    last_serviced_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vehicle_id) REFERENCES vehicles(vehicle_id) ON DELETE CASCADE
);

-- Create vehicle servicing history table
CREATE TABLE IF NOT EXISTS vehicle_servicing_history (
    history_id SERIAL PRIMARY KEY,
    vehicle_id INT NOT NULL,
    servicing_id INT NOT NULL,
    km_reading DECIMAL(10, 2) NOT NULL,
    servicing_type VARCHAR(50) NOT NULL,
    servicing_date TIMESTAMP NOT NULL,
    servicing_cost DECIMAL(10, 2),
    notes TEXT,
    serviced_by INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (vehicle_id) REFERENCES vehicles(vehicle_id) ON DELETE CASCADE,
    FOREIGN KEY (servicing_id) REFERENCES vehicle_servicing(servicing_id) ON DELETE CASCADE,
    FOREIGN KEY (serviced_by) REFERENCES users(id) ON DELETE SET NULL
);
`
	_, err = db.Exec(vehicle_servicing)
	if err != nil {
		return err
	}

	sale_charge := `
	CREATE TABLE IF NOT EXISTS sales_charges (
		charge_id SERIAL PRIMARY KEY,
		sale_id INT NOT NULL,
		charge_type VARCHAR(50) NOT NULL CHECK (charge_type IN ('damage', 'wash', 'delay', 'discount')),
		amount DECIMAL(10,2) NOT NULL DEFAULT 0.00,
		FOREIGN KEY (sale_id) REFERENCES sales(sale_id) ON DELETE CASCADE
	);

	`
	_, err = db.Exec(sale_charge)
	if err != nil {
		return err
	}
	// Create sales table
	saleImageQuery := `
		CREATE TABLE IF NOT EXISTS sales_videos (
		video_id SERIAL PRIMARY KEY,
		sale_id INT NOT NULL,
		video_url VARCHAR(255) NOT NULL,
		uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (sale_id) REFERENCES sales(sale_id) ON DELETE CASCADE
	);`
	_, err = db.Exec(saleImageQuery)
	if err != nil {
		return err
	}
	saleVideoQuery := `
	CREATE TABLE IF NOT EXISTS sales_images (
		image_id SERIAL PRIMARY KEY,
		sale_id INT NOT NULL,
		image_url VARCHAR(255) NOT NULL,
		uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		file_name varchar(255),
		FOREIGN KEY (sale_id) REFERENCES sales(sale_id) ON DELETE CASCADE
	);`
	_, err = db.Exec(saleVideoQuery)
	if err != nil {
		return err
	}

	// Create reminders table
	remindersQuery := `
	CREATE TABLE IF NOT EXISTS reminders (
		reminder_id SERIAL PRIMARY KEY,
		vehicle_id INT NOT NULL,
		reminder_type VARCHAR(255) NOT NULL CHECK (reminder_type IN ('EMI', 'Bill Book Renewal', 'Insurance Renewal')),
		due_date TIMESTAMP NOT NULL,
		is_completed BOOLEAN DEFAULT FALSE,
		remark TEXT,
		next_date TIMESTAMP,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (vehicle_id) REFERENCES vehicles(vehicle_id) ON DELETE CASCADE
	);
	`
	_, err = db.Exec(remindersQuery)
	if err != nil {
		return err
	}

	// Create system settings table
	systemSettingsQuery := `
	CREATE TABLE IF NOT EXISTS system_settings (
		setting_id SERIAL PRIMARY KEY,
		setting_key VARCHAR(50) NOT NULL UNIQUE,
		setting_value BOOLEAN NOT NULL DEFAULT true,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	-- Insert default settings
	INSERT INTO system_settings (setting_key, setting_value, description)
	VALUES 
		('enable_registration', true, 'Controls whether new user registration is allowed'),
		('enable_login', true, 'Controls whether user login is allowed')
	ON CONFLICT (setting_key) DO NOTHING;
	`
	_, err = db.Exec(systemSettingsQuery)
	return err
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
