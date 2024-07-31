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
)

func TestUpdateSecretPage1(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}

	server.router.Use(AuthMiddlewareTest())
	server.router.PUT("/api/v1/secrets/:secret_name", server.updateSecretPage)

	t.Run("successful update of secret", func(t *testing.T) {
		user := "testuser"
		key := "testsecret"
		secretval := models.UpdateSecret{
			Value: "testvalue",
		}

		mockService.On("UpdateSecret", user, key, secretval).Return(nil)

		reqBody, _ := json.Marshal(secretval)
		req, _ := http.NewRequest("PUT", "/api/v1/secrets/testsecret", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("authorization", "bearer valid_token")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		successResp := response.SuccessResponse{}
		err := json.NewDecoder(w.Body).Decode(&successResp)
		assert.NoError(t, err)
		assert.Equal(t, "secret updated successfully", successResp.Message)

		mockService.AssertExpectations(t)
	})
}

func TestUpdateSecretPage2(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}

	server.router.Use(AuthMiddlewareTest())
	server.router.PUT("/api/v1/secrets/:secret_name", server.updateSecretPage)

	t.Run("service error - secret doesn't exist", func(t *testing.T) {
		user := "testuser"
		key := "testsecret"
		secretval := models.UpdateSecret{
			Value: "testvalue",
		}

		mockService.On("UpdateSecret", user, key, secretval).Return(errors.New("secret doesn't exist"))

		reqBody, _ := json.Marshal(secretval)
		req, _ := http.NewRequest("PUT", "/api/v1/secrets/testsecret", bytes.NewBuffer(reqBody))
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

func TestUpdateSecretPage3(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}

	server.router.Use(AuthMiddlewareTest())
	server.router.POST("/api/v1/secrets/:secret_name", server.updateSecretPage)

	t.Run("invalid json", func(t *testing.T) {

		reqBody := []byte(`{"invalid": "json",}`)
		req, _ := http.NewRequest("POST", "/api/v1/secrets/testsecret", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("authorization", "bearer valid_token")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		errorResp := response.ErrorResponse{}
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid JSON", errorResp.Error)

		mockService.AssertExpectations(t)
	})
}

func TestUpdateSecretPage4(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}

	server.router.Use(AuthMiddlewareTest())
	server.router.POST("/api/v1/secrets/:secret_name", server.updateSecretPage)

	t.Run("key and value are required", func(t *testing.T) {

		reqBody := []byte(`{"data": "data"}`)
		req, _ := http.NewRequest("POST", "/api/v1/secrets/testsecret", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("authorization", "bearer valid_token")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		errorResp := response.ErrorResponse{}
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, "key and value are required", errorResp.Error)

		mockService.AssertExpectations(t)
	})
}

func TestUpdateSecretPage5(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}

	// server.router.Use(AuthMiddlewareTest())
	server.router.PUT("/api/v1/secrets/:secret_name", server.updateSecretPage)

	t.Run("status unauthorized", func(t *testing.T) {
		// user := "testuser"
		// key := "testsecret"
		secretval := models.UpdateSecret{
			Value: "testvalue",
		}

		// mockService.On("UpdateSecret", user, key, secretval).Return(errors.New("secret doesn't exist"))

		reqBody, _ := json.Marshal(secretval)
		req, _ := http.NewRequest("PUT", "/api/v1/secrets/testsecret", bytes.NewBuffer(reqBody))
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
