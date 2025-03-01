package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT creates a JWT token for a user
func GenerateJWT(userID int, jwtSecret string, expiry time.Duration) (string, error) {
	// Create claims
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(expiry).Unix(),
		"iat": time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and return token
	return token.SignedString([]byte(jwtSecret))
}
