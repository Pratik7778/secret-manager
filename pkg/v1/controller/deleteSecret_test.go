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

func TestDeleteSecretPage1(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}

	server.router.Use(AuthMiddlewareTest())
	server.router.DELETE("/api/v1/secrets/:secret_name", server.deleteSecretPage)

	t.Run("successful deletion of secret", func(t *testing.T) {
		user := "testuser"
		key := "testsecret"
		secretval := models.UpdateSecret{
			Value: "testvalue",
		}

		mockService.On("DeleteSecret", user, key).Return(nil)

		reqBody, _ := json.Marshal(secretval)
		req, _ := http.NewRequest("DELETE", "/api/v1/secrets/testsecret", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("authorization", "bearer valid_token")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		successResp := response.SuccessResponse{}
		err := json.NewDecoder(w.Body).Decode(&successResp)
		assert.NoError(t, err)
		assert.Equal(t, "secret deleted successfully", successResp.Message)

		mockService.AssertExpectations(t)
	})
}

func TestDeleteSecretPage2(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}

	server.router.Use(AuthMiddlewareTest())
	server.router.DELETE("/api/v1/secrets/:secret_name", server.deleteSecretPage)

	t.Run("service error - secret doesn't exist", func(t *testing.T) {
		user := "testuser"
		key := "testsecret"
		secretval := models.UpdateSecret{
			Value: "testvalue",
		}

		mockService.On("DeleteSecret", user, key).Return(errors.New("secret doesn't exist"))

		reqBody, _ := json.Marshal(secretval)
		req, _ := http.NewRequest("DELETE", "/api/v1/secrets/testsecret", bytes.NewBuffer(reqBody))
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

func TestDeleteSecretPage3(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}

	// server.router.Use(AuthMiddlewareTest())
	server.router.DELETE("/api/v1/secrets/:secret_name", server.deleteSecretPage)

	t.Run("status unauthorized", func(t *testing.T) {
		// user := "testuser"
		// key := "testsecret"
		secretval := models.UpdateSecret{
			Value: "testvalue",
		}

		// mockService.On("DeleteSecret", user, key).Return(errors.New("secret doesn't exist"))

		reqBody, _ := json.Marshal(secretval)
		req, _ := http.NewRequest("DELETE", "/api/v1/secrets/testsecret", bytes.NewBuffer(reqBody))
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
