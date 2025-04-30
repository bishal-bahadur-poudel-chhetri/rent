package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTAuth middleware validates JWT tokens
func JWTAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// The token is expected to be in the format: "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			c.Abort()
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Check if the token is valid and extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Set user information in the context
			if userID, ok := claims["user_id"].(float64); ok {
				c.Set("userID", int(userID))
			}
			if companyID, ok := claims["company_id"].(float64); ok {
				c.Set("companyID", int(companyID))
			}
			if username, ok := claims["username"].(string); ok {
				c.Set("username", username)
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Proceed to the next handler
		c.Next()
	}
}
