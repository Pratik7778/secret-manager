package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"secret-manager/pkg/v1/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	corev1 "k8s.io/api/core/v1"
)

// MockService is a mock implementation of the service layer
type MockService struct {
	mock.Mock
}

func (m *MockService) CreateUserSecret(user models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockService) GetSecretByLabel(token string) (string, error) {
	args := m.Called(token)
	return args.Error(0).Error(), args.Error(1)
}

func (m *MockService) LoginUser(user models.User, token string) error {
	args := m.Called(user, token)
	return args.Error(0)
}

func (m *MockService) ListUserSecrets(user string, query string, page string, pageSize string) (*corev1.SecretList, int, error) {
	args := m.Called(user, query, page, pageSize)
	return args.Get(0).(*corev1.SecretList), args.Int(1), args.Error(2)
}

func (m *MockService) CreateSecret(user string, secretval models.Secret) error {
	args := m.Called(user, secretval)
	return args.Error(0)
}

func (m *MockService) UpdateSecret(user string, key string, secretval models.UpdateSecret) error {
	args := m.Called(user, key, secretval)
	return args.Error(0)
}

func (m *MockService) DeleteSecret(user string, key string) error {
	args := m.Called(user, key)
	return args.Error(0)
}

func (m *MockService) ViewSecret(user string, key string) (map[string][]byte, error) {
	args := m.Called(user, key)
	return args.Get(0).(map[string][]byte), args.Error(1)
}

func TestRegisterPage(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}
	server.setupRouter()

	t.Run("successful registration", func(t *testing.T) {
		user := models.User{
			Username: "testuser",
			Password: "password",
		}
		mockService.On("CreateUserSecret", user).Return(nil)

		reqBody, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "User created successfully")
		mockService.AssertExpectations(t)
	})

	t.Run("invalid JSON request", func(t *testing.T) {
		reqBody := []byte(`{"invalid": "json",}`)
		req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid JSON")
	})

	t.Run("service error - user already exists", func(t *testing.T) {
		user := models.User{
			Username: "existinguser",
			Password: "password",
		}
		mockService.On("CreateUserSecret", user).Return(errors.New("user already exists"))

		reqBody, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "user already exists")
		mockService.AssertExpectations(t)
	})

	t.Run("service error - can't create user", func(t *testing.T) {
		user := models.User{
			Username: "newuser",
			Password: "password",
		}
		mockService.On("CreateUserSecret", user).Return(errors.New("can't create user"))

		reqBody, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "can't create user")
		mockService.AssertExpectations(t)
	})

	t.Run("missing username or password", func(t *testing.T) {
		reqBody := []byte(`{"username": ""}`)
		req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Username and password are required")
	})
}
