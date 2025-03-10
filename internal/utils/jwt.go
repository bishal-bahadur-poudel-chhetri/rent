package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
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
