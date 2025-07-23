package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func HandleServiceError(c *gin.Context, err error, resourceName string) {
	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		NotFoundResponse(c, resourceName)
	case mongo.IsDuplicateKeyError(err):
		ErrorResponse(c, 409, resourceName+" already exists")
	case IsValidationError(err):
		ValidationErrorResponse(c, err)
	default:
		InternalServerErrorResponse(c, err.Error())
	}
}

func IsValidationError(err error) bool {
	return err != nil && (err.Error() == "invalid email format" ||
		err.Error() == "invalid password length" ||
		err.Error() == "invalid role")
}

func BindAndValidate(c *gin.Context, obj interface{}) bool {

	if err := c.ShouldBindJSON(obj); err != nil {
		ValidationErrorResponse(c, err)
		return false
	}
	return true
}

func GetValidatedParamID(c *gin.Context) (string, bool) {
	id := c.Param("id")
	if !IsValidObjectID(id) {
		ErrorResponse(c, 400, "Invalid ID format")
		return "", false
	}
	return id, true
}

func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		UnauthorizedResponse(c, "User ID not found in token")
		return "", false
	}
	return userID.(string), true
}

func GetUserRoleFromContext(c *gin.Context) (string, bool) {
	userRole, exists := c.Get("user_role")
	if !exists {
		ForbiddenResponse(c, "User role not found in token")
		return "", false
	}
	return userRole.(string), true
}

func RequireOwnershipOrAdmin(c *gin.Context, resourceOwnerID string) bool {
	userID, ok := GetUserIDFromContext(c)
	if !ok {
		return false
	}

	userRole, ok := GetUserRoleFromContext(c)
	if !ok {
		return false
	}

	if userRole == "admin" || userID == resourceOwnerID {
		return true
	}

	ForbiddenResponse(c, "Access denied: insufficient permissions")
	return false
}
