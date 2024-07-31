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

func TestCreateSecretPage1(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}

	server.router.Use(AuthMiddlewareTest())
	server.router.POST("/api/v1/secrets/create", server.createSecretPage)

	t.Run("successful creation of secret", func(t *testing.T) {
		user := "testuser"
		secretval := models.Secret{
			Key:   "testsecret",
			Value: "testvalue",
		}

		mockService.On("CreateSecret", user, secretval).Return(nil)

		reqBody, _ := json.Marshal(secretval)
		req, _ := http.NewRequest("POST", "/api/v1/secrets/create", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("authorization", "bearer valid_token")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		successResp := response.SuccessResponse{}
		err := json.NewDecoder(w.Body).Decode(&successResp)
		assert.NoError(t, err)
		assert.Equal(t, "secret created successfully", successResp.Message)

		mockService.AssertExpectations(t)
	})
}

func TestCreateSecretPage2(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}

	server.router.Use(AuthMiddlewareTest())
	server.router.POST("/api/v1/secrets/create", server.createSecretPage)

	t.Run("service error - secret already exists", func(t *testing.T) {
		user := "testuser"
		secretval := models.Secret{
			Key:   "testsecret",
			Value: "testvalue",
		}

		mockService.On("CreateSecret", user, secretval).Return(errors.New("secret already exists"))

		reqBody, _ := json.Marshal(secretval)
		req, _ := http.NewRequest("POST", "/api/v1/secrets/create", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("authorization", "bearer valid_token")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		errorResp := response.ErrorResponse{}
		err := json.NewDecoder(w.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Equal(t, "secret already exists", errorResp.Error)

		mockService.AssertExpectations(t)
	})
}

func TestCreateSecretPage3(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}

	server.router.Use(AuthMiddlewareTest())
	server.router.POST("/api/v1/secrets/create", server.createSecretPage)

	t.Run("invalid json", func(t *testing.T) {

		reqBody := []byte(`{"invalid": "json",}`)
		req, _ := http.NewRequest("POST", "/api/v1/secrets/create", bytes.NewBuffer(reqBody))
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

func TestCreateSecretPage4(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}

	server.router.Use(AuthMiddlewareTest())
	server.router.POST("/api/v1/secrets/create", server.createSecretPage)

	t.Run("key value is empty", func(t *testing.T) {

		// also works, because gin will ignore the extra fields and will not return error
		// reqBody := []byte(`{"invalid": "json"}`)

		reqBody := []byte(`{"key": "", "value": "json"}`)
		req, _ := http.NewRequest("POST", "/api/v1/secrets/create", bytes.NewBuffer(reqBody))
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

func TestCreateSecretPage5(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}

	// server.router.Use(AuthMiddlewareTest())
	server.router.POST("/api/v1/secrets/create", server.createSecretPage)

	t.Run("status unauthorized", func(t *testing.T) {
		// user := "testuser"
		secretval := models.Secret{
			Key:   "testsecret",
			Value: "testvalue",
		}

		// mockService.On("CreateSecret", user, secretval).Return(nil)

		reqBody, _ := json.Marshal(secretval)
		req, _ := http.NewRequest("POST", "/api/v1/secrets/create", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("authorization", "bearer valid_token")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var decodedResp string
		_ = json.Unmarshal(w.Body.Bytes(), &decodedResp)
		assert.Equal(t, "can't get user", decodedResp)

		mockService.AssertExpectations(t)
	})
}
