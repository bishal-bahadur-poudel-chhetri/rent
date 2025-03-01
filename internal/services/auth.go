package services

import (
	"database/sql"
	"errors"
	"fmt"
	"renting/internal/models"
	"renting/internal/repositories"
	"renting/internal/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists         = errors.New("username or mobile number already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrCompanyNotFound    = errors.New("company code is invalid")
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo    *repositories.UserRepository
	jwtSecret   string
	tokenExpiry time.Duration
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string, tokenExpiry time.Duration) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		tokenExpiry: tokenExpiry,
	}
}

// RegisterUser registers a new user
func (s *AuthService) RegisterUser(req models.RegisterRequest) (int, string, error) {
	// Check if user already exists
	exists, err := s.userRepo.CheckUserExists(req.Username, req.MobileNumber)
	if err != nil {
		return 0, "", err
	}
	if exists {
		return 0, "", ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, "", err
	}

	// Create user
	userID, err := s.userRepo.CreateUser(
		req.Username,
		string(hashedPassword),
		req.MobileNumber,
		req.CompanyID,
		req.IsAdmin,
	)
	if err != nil {
		return 0, "", err
	}

	// Generate token
	token, err := utils.GenerateJWT(userID, s.jwtSecret, s.tokenExpiry)
	if err != nil {
		return 0, "", err
	}

	return userID, token, nil
}

// LoginUser authenticates a user using mobile number, password and company code

func (s *AuthService) LoginUser(req models.LoginRequest) (int, string, *models.User, error) {
	// Get user from database using mobile number and company code
	user, hashedPassword, err := s.userRepo.GetUserByMobileAndCompany(req.MobileNumber, req.CompanyCode, req.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			// Try to check if company exists
			_, companyErr := s.userRepo.GetCompanyByCode(req.CompanyCode)
			if companyErr != nil && companyErr == sql.ErrNoRows {
				fmt.Printf("Error: Company with code '%s' not found\n", req.CompanyCode)
				return 0, "", nil, ErrCompanyNotFound
			}
			fmt.Printf("Error: Invalid credentials for mobile number '%s' and company code '%s'\n", req.MobileNumber, req.CompanyCode)
			return 0, "", nil, ErrInvalidCredentials
		}
		fmt.Printf("Error retrieving user by mobile and company: %v\n", err)
		return 0, "", nil, err
	}

	// Compare password with hash
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		fmt.Printf("Error: Password mismatch for user '%s'\n", req.MobileNumber)
		return 0, "", nil, ErrInvalidCredentials
	}

	// Generate token
	token, err := utils.GenerateJWT(user.ID, s.jwtSecret, s.tokenExpiry)
	if err != nil {
		fmt.Printf("Error generating JWT token for user ID '%d': %v\n", user.ID, err)
		return 0, "", nil, err
	}

	fmt.Printf("Login successful for user '%s' (ID: %d)\n", user.Username, user.ID)

	return user.ID, token, user, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(id int) (*models.User, error) {
	return s.userRepo.GetUserByID(id)
}

// UpdateUserProfile updates a user's profile
func (s *AuthService) UpdateUserProfile(id int, req models.UpdateProfileRequest) error {
	return s.userRepo.UpdateUser(id, req.Username, req.MobileNumber)
}
