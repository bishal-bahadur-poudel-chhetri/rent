package middleware

import (
	"net/http"
	"renting/internal/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth middleware for JWT authentication
func JWTAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Use the standardized error response
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Authorization header required", nil))
			return
		}

		// Extract token from "Bearer <token>"
		if !strings.HasPrefix(authHeader, "Bearer ") {
			// Use the standardized error response
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Invalid authorization format", nil))
			return
		}

		tokenString := authHeader[7:]

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			// Use the standardized error response
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Invalid token", nil))
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Check expiration
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				// Use the standardized error response
				c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Token expired", nil))
				return
			}

			// Set user ID in context
			c.Set("userID", int(claims["sub"].(float64)))
			c.Next()
		} else {
			// Use the standardized error response
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "Invalid token claims", nil))
			return
		}
	}
}

// CORS middleware adds CORS headers
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
