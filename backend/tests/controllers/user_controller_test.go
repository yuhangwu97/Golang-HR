package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gin-project/controllers"
	"gin-project/models"
	"gin-project/tests/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestUserController_CreateUser(t *testing.T) {
	mockService := new(mocks.MockUserService)
	controller := controllers.NewUserControllerWithService(mockService)

	router := setupRouter()
	router.POST("/users", controller.CreateUser)

	t.Run("Success", func(t *testing.T) {
		user := &models.User{
			ID:        primitive.NewObjectID(),
			Name:      "John Doe",
			Email:     "john@example.com",
			Role:      "user",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockService.On("CreateUser", &models.User{
			Name:  "John Doe",
			Email: "john@example.com",
			Role:  "user",
		}).Return(user, nil)

		reqBody := map[string]interface{}{
			"name":  "John Doe",
			"email": "john@example.com",
			"role":  "user",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Validation Error", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": "John Doe",
		}
		jsonBody, _ := json.Marshal(reqBody)

		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUserController_GetUser(t *testing.T) {
	mockService := new(mocks.MockUserService)
	controller := controllers.NewUserControllerWithService(mockService)

	router := setupRouter()
	router.GET("/users/:id", controller.GetUser)

	t.Run("Success", func(t *testing.T) {
		userID := primitive.NewObjectID().Hex()
		user := &models.User{
			ID:    primitive.NewObjectID(),
			Name:  "John Doe",
			Email: "john@example.com",
			Role:  "user",
		}

		mockService.On("GetUserByID", userID).Return(user, nil)

		req, _ := http.NewRequest("GET", "/users/"+userID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		userID := primitive.NewObjectID().Hex()

		mockService.On("GetUserByID", userID).Return(nil, mongo.ErrNoDocuments)

		req, _ := http.NewRequest("GET", "/users/"+userID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserController_DeleteUser(t *testing.T) {
	mockService := new(mocks.MockUserService)
	controller := controllers.NewUserControllerWithService(mockService)

	router := setupRouter()
	router.DELETE("/users/:id", controller.DeleteUser)

	t.Run("Success", func(t *testing.T) {
		userID := primitive.NewObjectID().Hex()

		mockService.On("DeleteUser", userID).Return(nil)

		req, _ := http.NewRequest("DELETE", "/users/"+userID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		userID := primitive.NewObjectID().Hex()

		mockService.On("DeleteUser", userID).Return(mongo.ErrNoDocuments)

		req, _ := http.NewRequest("DELETE", "/users/"+userID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		userID := primitive.NewObjectID().Hex()

		mockService.On("DeleteUser", userID).Return(errors.New("database error"))

		req, _ := http.NewRequest("DELETE", "/users/"+userID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}