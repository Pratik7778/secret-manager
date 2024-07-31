package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"secret-manager/pkg/v1/models"
	"secret-manager/pkg/v1/response"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoginPage(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}
	server.setupRouter()

	t.Run("successful login", func(t *testing.T) {
		user := models.User{
			Username: "testuser",
			Password: "password",
		}
		// token := "a6f7dca5304eb290b826"
		mockService.On("LoginUser", user, mock.AnythingOfType("string")).Return(nil)

		reqBody, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var loginReponse response.LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &loginReponse)
		assert.NoError(t, err)
		// assert.Equal(t, mock.AnythingOfType("string"), loginReponse.Token)
		// assert.Contains(t, w.Body.String(), "Login successful")
		assert.NotEmpty(t, loginReponse.Token)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid JSON request", func(t *testing.T) {
		reqBody := []byte(`{"invalid": "json";,}`)
		req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid JSON")
	})

	t.Run("service error - user doesn't exist", func(t *testing.T) {
		user := models.User{
			Username: "wronguser",
			Password: "wrongpassword",
		}
		mockService.On("LoginUser", user, mock.AnythingOfType("string")).Return(errors.New("user doesn't exist"))

		reqBody, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "user doesn't exist")
		mockService.AssertExpectations(t)
	})

	t.Run("missing username or password", func(t *testing.T) {
		reqBody := []byte(`{"username": ""}`)
		req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Username and password are required")
	})
}
