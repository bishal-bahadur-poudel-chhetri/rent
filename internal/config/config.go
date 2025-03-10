package config

import (
	"database/sql"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Config holds application configuration
type Config struct {
	DBConnStr     string
	JWTSecret     string
	TokenExpiry   time.Duration
	ServerAddress string
}

// LoadConfig loads application configuration from environment variables
func LoadConfig() (*Config, error) {
	return &Config{
		DBConnStr:     getEnv("DATABASE_URL", "postgres://myuser:mypassword@localhost:5432/mydatabase?sslmode=disable"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		TokenExpiry:   time.Hour * 24, // 24 hours
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
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
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (vehicle_type_id) REFERENCES vehicle_types(vehicle_type_id) ON DELETE CASCADE
	);
	`
	_, err = db.Exec(vehiclesQuery)
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
		user_id INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = db.Exec(paymentsQuery)
	if err != nil {
		return err
	}
	salesQuery := `
CREATE TABLE IF NOT EXISTS sales (
    sale_id SERIAL PRIMARY KEY,  -- Add sale_id as the primary key
    vehicle_id INT,
    customer_name VARCHAR(255),
    total_amount DECIMAL(10, 2),
    charge_per_day DECIMAL(10, 2),
    booking_date DATE,
    date_of_delivery DATE,
    return_date DATE,
    is_damaged BOOLEAN,
    is_washed BOOLEAN,
    is_delayed BOOLEAN,
    number_of_days INT,
    payment_id INT,
    remark TEXT,
    fuel_range_received DECIMAL(10, 2),
    fuel_range_delivered DECIMAL(10, 2),
    km_received DECIMAL(10, 2),
    km_delivered DECIMAL(10, 2),
    photo_1 VARCHAR(255),
    photo_2 VARCHAR(255),
    photo_3 VARCHAR(255),
    photo_4 VARCHAR(255),
    damage_cost DECIMAL(10, 2) DEFAULT 0.00, 
    wash_cost DECIMAL(10, 2) DEFAULT 0.00,   
    delay_cost DECIMAL(10, 2) DEFAULT 0.00,  
    discount_amount DECIMAL(10, 2) DEFAULT 0.00
);
`
	_, err = db.Exec(salesQuery)
	if err != nil {
		return err
	}

	sale_charge := `
	CREATE TABLE IF NOT EXISTS sales_charges (
    charge_id SERIAL PRIMARY KEY,
    sale_id INT NOT NULL,
    charge_type VARCHAR(50) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    FOREIGN KEY (sale_id) REFERENCES sales(sale_id) ON DELETE CASCADE
	);
	`
	_, err = db.Exec(sale_charge)
	if err != nil {
		return err
	}
	// Create sales table

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
