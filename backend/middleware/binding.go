package middleware

import (
	"gin-project/utils"

	"github.com/gin-gonic/gin"
)

func BindJSON(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(obj); err != nil {
			utils.ValidationErrorResponse(c, err)
			c.Abort()
			return
		}
		c.Set("bindedData", obj)
		c.Next()
	}
}

func BindQuery(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindQuery(obj); err != nil {
			utils.ValidationErrorResponse(c, err)
			c.Abort()
			return
		}
		c.Set("bindedQuery", obj)
		c.Next()
	}
}

func BindUri(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.ShouldBindUri(obj); err != nil {
			utils.ValidationErrorResponse(c, err)
			c.Abort()
			return
		}
		c.Set("bindedUri", obj)
		c.Next()
	}
}

func ValidateAndBindJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody interface{}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			utils.ValidationErrorResponse(c, err)
			c.Abort()
			return
		}
		c.Set("requestBody", requestBody)
		c.Next()
	}
}