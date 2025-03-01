package repositories

import (
	"database/sql"
	"fmt"
	"renting/internal/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserRepository handles database operations related to users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser adds a new user to the database with hashed password
func (r *UserRepository) CreateUser(username, password, mobileNumber string, companyID int, isAdmin bool) (int, error) {
	var userID int
	currentTime := time.Now()

	// Hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	query := `INSERT INTO users (username, password, mobile_number, company_id, is_admin, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err = r.db.QueryRow(
		query,
		username,
		hashedPassword,
		mobileNumber,
		companyID,
		isAdmin,
		currentTime,
		currentTime,
	).Scan(&userID)

	return userID, err
}

// GetUserByMobileAndCompany retrieves a user by mobile number and company code, and compares password
func (r *UserRepository) GetUserByMobileAndCompany(mobileNumber, companyCode, enteredPassword string) (*models.User, string, error) {
	var user models.User
	var hashedPassword string

	// Use the exact values as in the working query for debugging
	query := `SELECT u.id, u.username, u.password, u.is_admin, u.company_id, u.mobile_number, u.created_at, u.updated_at 
              FROM users u
              JOIN companies c ON u.company_id = c.id
              WHERE u.mobile_number = $1 AND c.company_code = $2`

	// Check if the parameters are being passed correctly
	fmt.Printf("Query: %s\n", query)
	fmt.Printf("Parameters - mobileNumber: %s, companyCode: %s\n", mobileNumber, companyCode)

	// Run the query with parameters
	err := r.db.QueryRow(query, mobileNumber, companyCode).Scan(
		&user.ID,
		&user.Username,
		&hashedPassword,
		&user.IsAdmin,
		&user.CompanyID,
		&user.MobileNumber,
		&user.CreatedAt,
		&user.Updated_at,
	)

	// Debugging the error
	if err != nil {
		fmt.Println("Error:", err)
		return nil, "", err
	}

	return &user, hashedPassword, nil
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, is_admin, company_id, mobile_number, created_at, updated_at 
			  FROM users 
			  WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.IsAdmin,
		&user.CompanyID,
		&user.MobileNumber,
		&user.CreatedAt,
		&user.Updated_at,
	)

	return &user, err
}

// UpdateUser updates a user's profile
func (r *UserRepository) UpdateUser(id int, username, mobileNumber string) error {
	currentTime := time.Now()
	query := `UPDATE users 
			  SET username = $1, mobile_number = $2, updated_at = $3 
			  WHERE id = $4`

	_, err := r.db.Exec(query, username, mobileNumber, currentTime, id)
	return err
}

// CheckUserExists checks if a user with the given username or mobile number already exists
func (r *UserRepository) CheckUserExists(username, mobileNumber string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE username = $1 OR mobile_number = $2"
	err := r.db.QueryRow(query, username, mobileNumber).Scan(&count)
	return count > 0, err
}

// GetCompanyByCode retrieves a company by its code
func (r *UserRepository) GetCompanyByCode(companyCode string) (*models.Company, error) {
	var company models.Company
	query := "SELECT id, name, company_code, created_at, updated_at FROM companies WHERE company_code = $1"

	err := r.db.QueryRow(query, companyCode).Scan(
		&company.ID,
		&company.Name,
		&company.CompanyCode,
		&company.CreatedAt,
		&company.Updated_at,
	)

	return &company, err
}
