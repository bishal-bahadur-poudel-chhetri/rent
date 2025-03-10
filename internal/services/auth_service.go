package services

import (
	"context"
	"errors"
	"fmt"
	"renting/internal/models"
	"renting/internal/repositories"
	"renting/internal/utils"
	"time"
)

type AuthService interface {
	Login(ctx context.Context, mobileNumber, password, companyCode string) (string, *models.User, error)
	Register(ctx context.Context, user *models.User, companyCode string) error
}
type authService struct {
	userRepo    repositories.UserRepository
	jwtSecret   string
	tokenExpiry time.Duration
}

func NewAuthService(userRepo repositories.UserRepository, jwtSecret string, tokenExpiry time.Duration) AuthService {
	return &authService{
		userRepo:    userRepo,
		jwtSecret:   jwtSecret,
		tokenExpiry: tokenExpiry,
	}
}

func (s *authService) Login(ctx context.Context, mobileNumber, password, companyCode string) (string, *models.User, error) {
	// Find the user by mobile number and company code
	user, err := s.userRepo.FindByMobileAndCompanyCode(ctx, mobileNumber, companyCode)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	if user == nil {
		return "", nil, errors.New("invalid credentials")
	}

	// Verify the password
	if !utils.CheckPasswordHash(password, user.Password) {
		return "", nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.CompanyID, user.Username, s.jwtSecret, s.tokenExpiry)
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	return token, user, nil
}

func (s *authService) Register(ctx context.Context, user *models.User, companyCode string) error {
	// Fetch the company ID using the company code
	companyID, err := s.userRepo.GetCompanyIDByCode(ctx, companyCode)
	if err != nil {
		return fmt.Errorf("failed to fetch company ID: %v", err)
	}

	// Set the company ID in the user model
	user.CompanyID = companyID

	// Check if the mobile number already exists
	existingUser, err := s.userRepo.FindByMobileAndCompanyCode(ctx, user.MobileNumber, companyCode)
	if err != nil {
		return errors.New("failed to check mobile number")
	}
	if existingUser != nil {
		return errors.New("mobile number already exists")
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = hashedPassword

	// Create the user
	err = s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return errors.New("failed to create user")
	}

	return nil
}
