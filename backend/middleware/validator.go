package middleware

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateObjectID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id != "" {
			if !primitive.IsValidObjectID(id) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid ID format. Must be a valid MongoDB ObjectID",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// ValidateNumericID validates that the ID parameter is a positive integer (for MySQL auto-increment IDs)
func ValidateNumericID() gin.HandlerFunc {
	return ValidateNumericParam("id")
}

// ValidateNumericParam validates that the specified parameter is a positive integer
func ValidateNumericParam(paramName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param(paramName)
		if id != "" {
			numID, err := strconv.ParseUint(id, 10, 32)
			if err != nil || numID == 0 {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid " + paramName + " format. Must be a positive integer",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func ValidateEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Email string `json:"email"`
		}
		
		if err := c.ShouldBindJSON(&body); err == nil {
			emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
			if body.Email != "" && !emailRegex.MatchString(body.Email) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid email format",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func ValidateRequired(fields ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body map[string]interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid JSON format",
			})
			c.Abort()
			return
		}

		for _, field := range fields {
			if value, exists := body[field]; !exists || value == "" || value == nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Field '" + field + "' is required",
				})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}