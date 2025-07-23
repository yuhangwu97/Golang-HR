package utils

import (
	"testing"

	"gin-project/utils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co.uk", true},
		{"invalid-email", false},
		{"@example.com", false},
		{"test@", false},
		{"", false},
	}

	for _, test := range tests {
		result := utils.IsValidEmail(test.email)
		assert.Equal(t, test.expected, result, "Email: %s", test.email)
	}
}

func TestIsValidObjectID(t *testing.T) {
	validID := primitive.NewObjectID().Hex()
	invalidID := "invalid-id"

	assert.True(t, utils.IsValidObjectID(validID))
	assert.False(t, utils.IsValidObjectID(invalidID))
	assert.False(t, utils.IsValidObjectID(""))
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		password string
		expected bool
	}{
		{"password123", true},
		{"123456", true},
		{"12345", false},
		{"", false},
	}

	for _, test := range tests {
		result := utils.IsValidPassword(test.password)
		assert.Equal(t, test.expected, result, "Password: %s", test.password)
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  hello world  ", "hello world"},
		{"test", "test"},
		{"  ", ""},
		{"", ""},
	}

	for _, test := range tests {
		result := utils.SanitizeString(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestIsValidRole(t *testing.T) {
	tests := []struct {
		role     string
		expected bool
	}{
		{"admin", true},
		{"user", true},
		{"moderator", true},
		{"invalid", false},
		{"", false},
	}

	for _, test := range tests {
		result := utils.IsValidRole(test.role)
		assert.Equal(t, test.expected, result, "Role: %s", test.role)
	}
}

func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Test@Example.Com", "test@example.com"},
		{"  USER@DOMAIN.COM  ", "user@domain.com"},
		{"normal@email.com", "normal@email.com"},
	}

	for _, test := range tests {
		result := utils.NormalizeEmail(test.input)
		assert.Equal(t, test.expected, result)
	}
}