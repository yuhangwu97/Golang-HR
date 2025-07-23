package utils

import (
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func IsValidObjectID(id string) bool {
	return primitive.IsValidObjectID(id)
}

func IsValidPassword(password string) bool {
	return len(password) >= 6
}

func SanitizeString(input string) string {
	return strings.TrimSpace(input)
}

func IsValidRole(role string) bool {
	validRoles := []string{"admin", "user", "moderator"}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}