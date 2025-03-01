package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetQueryInt retrieves an integer query parameter or returns a default value
func GetQueryInt(c *gin.Context, key string, defaultValue int) int {
	valueStr := c.DefaultQuery(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
