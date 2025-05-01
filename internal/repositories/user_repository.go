package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"renting/internal/models"
)

type UserRepository interface {
	FindByMobileAndCompanyCode(ctx context.Context, mobileNumber, companyCode string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	GetCompanyIDByCode(ctx context.Context, companyCode string) (int, error)
	GetUserByID(ctx context.Context, userID int) (*models.User, error)
	LockoutUser(ctx context.Context, userID int) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByMobileAndCompanyCode(ctx context.Context, mobileNumber, companyCode string) (*models.User, error) {
	query := `
		SELECT u.id, u.username, u.password, u.is_admin, u.is_locked, u.company_id, u.mobile_number, u.created_at, u.updated_at
		FROM users u
		JOIN companies c ON u.company_id = c.id
		WHERE u.mobile_number = $1 AND c.company_code = $2
	`
	row := r.db.QueryRowContext(ctx, query, mobileNumber, companyCode)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.IsAdmin,
		&user.IsLocked,
		&user.CompanyID,
		&user.MobileNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetCompanyIDByCode(ctx context.Context, companyCode string) (int, error) {
	query := `
		SELECT id FROM companies WHERE company_code = $1
	`
	row := r.db.QueryRowContext(ctx, query, companyCode)

	var companyID int
	err := row.Scan(&companyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("company not found")
		}
		return 0, err
	}

	return companyID, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	fmt.Printf("Executing query: INSERT INTO users (username, password, is_admin, is_locked, company_id, mobile_number) "+
		"VALUES ('%s', '%s', %v, %v, %d, '%s')\n",
		user.Username, user.Password, user.IsAdmin, false, user.CompanyID, user.MobileNumber)
	query := `
		INSERT INTO users (username, password, is_admin, is_locked, company_id, mobile_number)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password,
		user.IsAdmin,
		false,
		user.CompanyID,
		user.MobileNumber,
	).Scan(&user.ID)
	if err != nil {
		log.Printf("Error executing query: %v\n", err)
		return err
	}

	fmt.Printf("User created with ID: %d\n", user.ID)
	return nil
}

func (r *userRepository) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	query := `
		SELECT id, username, password, is_admin, is_locked, company_id, mobile_number, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, userID)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.IsAdmin,
		&user.IsLocked,
		&user.CompanyID,
		&user.MobileNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) LockoutUser(ctx context.Context, userID int) error {
	query := `
		UPDATE users 
		SET is_locked = true 
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to lock out user: %v", err)
	}
	return nil
}
