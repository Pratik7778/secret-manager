package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"secret-manager/pkg/v1/response"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestViewSecretPage1(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}

	server.router.Use(AuthMiddlewareTest())
	server.router.GET("/api/v1/secrets/:secret_name", server.viewSecretPage)

	t.Run("successful view of secret", func(t *testing.T) {
		user := "testuser"
		key := "testsecret"

		secretMap := map[string][]byte{
			"testsecret": []byte("testvalue"),
		}
		mockService.On("ViewSecret", user, key).Return(secretMap, nil)

		req, _ := http.NewRequest("GET", "/api/v1/secrets/testsecret", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("authorization", "bearer valid_token")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		stringMap := make(map[string]string, 1)
		for k, v := range secretMap {
			stringMap[k] = string(v)
		}

		successResp := make(map[string]string)
		err := json.NewDecoder(w.Body).Decode(&successResp)
		assert.NoError(t, err)
		assert.Equal(t, stringMap, successResp)

		mockService.AssertExpectations(t)
	})
}

func TestViewSecretPage2(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}

	server.router.Use(AuthMiddlewareTest())
	server.router.GET("/api/v1/secrets/:secret_name", server.viewSecretPage)

	t.Run("service error - secret doesn't exist", func(t *testing.T) {
		user := "testuser"
		key := "testsecret"

		mockService.On("ViewSecret", user, key).Return(map[string][]byte{}, errors.New("secret doesn't exist"))

		req, _ := http.NewRequest("GET", "/api/v1/secrets/testsecret", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("authorization", "bearer valid_token")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		errorResp := response.ErrorResponse{}
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, "secret doesn't exist", errorResp.Error)

		mockService.AssertExpectations(t)
	})
}

func TestViewSecretPage3(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}

	// server.router.Use(AuthMiddlewareTest())
	server.router.GET("/api/v1/secrets/:secret_name", server.viewSecretPage)

	t.Run("status unauthorized", func(t *testing.T) {
		// user := "testuser"
		// key := "testsecret"

		// mockService.On("ViewSecret", user, key).Return(map[string][]byte{}, errors.New("secret doesn't exist"))

		req, _ := http.NewRequest("GET", "/api/v1/secrets/testsecret", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("authorization", "bearer valid_token")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var decodedResp string
		err := json.Unmarshal(w.Body.Bytes(), &decodedResp)
		assert.NoError(t, err)
		assert.Equal(t, "can't get user", decodedResp)

		mockService.AssertExpectations(t)
	})
}
