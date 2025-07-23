package utils

import (
	"github.com/gin-gonic/gin"
)

func BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		ValidationErrorResponse(c, err)
		return false
	}
	return true
}

func BindQuery(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		ValidationErrorResponse(c, err)
		return false
	}
	return true
}

func BindUri(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindUri(obj); err != nil {
		ValidationErrorResponse(c, err)
		return false
	}
	return true
}

func BindHeader(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindHeader(obj); err != nil {
		ValidationErrorResponse(c, err)
		return false
	}
	return true
}

func GetParamID(c *gin.Context) (string, bool) {
	id := c.Param("id")
	if !IsValidObjectID(id) {
		ErrorResponse(c, 400, "Invalid ID format")
		return "", false
	}
	return id, true
}

func GetQueryParam(c *gin.Context, key, defaultValue string) string {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetRequiredQueryParam(c *gin.Context, key string) (string, bool) {
	value := c.Query(key)
	if value == "" {
		ErrorResponse(c, 400, "Missing required query parameter: "+key)
		return "", false
	}
	return value, true
}