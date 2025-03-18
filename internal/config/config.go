package config

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/lib/pq"
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
		DBConnStr:     getEnv("DATABASE_URL", "postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		TokenExpiry:   time.Hour * 24, // 24 hours
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
func ConnectDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Check database connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Initialize database schema
	err = initDB(db)
	if err != nil {
		return nil, err
	}

	return db, nil
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
		vehicle_id INT,
		user_id INT NOT NULL, 
		customer_name VARCHAR(255),
		customer_destination VARCHAR(255),
		customer_phone VARCHAR(255),
		total_amount DECIMAL(10, 2),
		charge_per_day DECIMAL(10, 2),
		booking_date DATE,
		date_of_delivery DATE,
		return_date DATE,
		is_damaged BOOLEAN,
		is_washed BOOLEAN,
		is_delayed BOOLEAN,
		number_of_days INT,
		remark TEXT,
		status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'active', 'completed', 'cancelled')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (vehicle_id) REFERENCES vehicles(vehicle_id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		 
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
			payment_type VARCHAR(50) NOT NULL,
			amount_paid DECIMAL(10,2) NOT NULL,
			payment_date TIMESTAMP NOT NULL,
			payment_status VARCHAR(50) NOT NULL CHECK (payment_status IN ('Pending', 'Completed', 'Failed')),
			verified_by_admin BOOLEAN DEFAULT FALSE,
			remark TEXT,
			user_id INT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			sale_id INT NOT NULL,
			FOREIGN KEY (sale_id) REFERENCES sales(sale_id) ON DELETE CASCADE
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
