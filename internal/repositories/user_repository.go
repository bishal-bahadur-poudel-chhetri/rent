package repositories

import (
	"context"
	"database/sql"
	"renting/internal/models"
)

type UserRepository interface {
	FindByMobileAndCompanyCode(ctx context.Context, mobileNumber, companyCode string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByMobileAndCompanyCode(ctx context.Context, mobileNumber, companyCode string) (*models.User, error) {
	query := `
		SELECT u.id, u.username, u.password, u.is_admin, u.company_id, u.mobile_number, u.created_at, u.updated_at
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
		&user.CompanyID,
		&user.MobileNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (username, password, is_admin, company_id, mobile_number)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password,
		user.IsAdmin,
		user.CompanyID,
		user.MobileNumber,
	).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}
