package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	UserID    int    `json:"user_id"`
	CompanyID int    `json:"company_id"`
	Username  string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(userID, companyID int, username, jwtSecret string, tokenExpiry time.Duration) (string, error) {
	expirationTime := time.Now().Add(tokenExpiry)
	claims := &Claims{
		UserID:    userID,
		CompanyID: companyID,
		Username:  username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ExtractUserIDFromToken extracts user ID from JWT token
func ExtractUserIDFromToken(c *gin.Context, jwtSecret string) (int, error) {
	// Extract the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("missing authorization token")
	}

	// Extract token from Bearer token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return 0, fmt.Errorf("missing token in authorization header")
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		// Return the secret key used for signing
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	// Extract user ID from claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		return 0, fmt.Errorf("user_id not found in token")
	}

	// Convert user_id to int
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("user_id is not a valid number")
	}

	return int(userID), nil
}
