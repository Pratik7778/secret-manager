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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AuthMiddlewareTest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("X-User-Token", "testuser")
		ctx.Next() //
	}
}

func TestListUserSecrets1(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}
	// server.setupRouter()
	server.router.Use(AuthMiddlewareTest())
	server.router.GET("/api/v1/secrets", server.listPage)

	t.Run("successful listing of secrets", func(t *testing.T) {
		user := "testuser"
		query := ""
		page := "1"
		pageSize := "5"
		totalSecrets := 2

		listOfSecrets := &corev1.SecretList{
			Items: []corev1.Secret{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "secret1",
						Namespace: user,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "secret2",
						Namespace: user,
					},
				},
			},
		}

		mockService.On("ListUserSecrets", user, query, page, pageSize).Return(listOfSecrets, totalSecrets, nil)

		req, _ := http.NewRequest("GET", "/api/v1/secrets?page=1&page_size=5", nil)
		req.Header.Set("authorization", "bearer valid_token")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		response := response.ListResponse{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, []string{"secret1", "secret2"}, response.Secrets)
		assert.Equal(t, totalSecrets, response.Total)
		assert.Equal(t, query, response.Query)
		assert.Equal(t, 5, response.PageSize)
		assert.Equal(t, 1, response.PageNumber)
		mockService.AssertExpectations(t)
	})

	// fakeClient := fake.NewSimpleClientset(
	// 	&corev1.Namespace{
	// 		ObjectMeta: metav1.ObjectMeta{
	// 			Name: "testuser",
	// 		},
	// 	},
	// 	&corev1.Secret{
	// 		ObjectMeta: metav1.ObjectMeta{
	// 			Name:      "secret1",
	// 			Namespace: "testuser",
	// 		},
	// 		StringData: map[string]string{
	// 			"key1": "value1",
	// 		},
	// 	},
	// )
}

func TestListUserSecrets2(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}
	// server.setupRouter()
	server.router.Use(AuthMiddlewareTest())
	server.router.GET("/api/v1/secrets", server.listPage)

	t.Run("failed to list secrets", func(t *testing.T) {
		user := "testuser"
		query := ""
		page := "1"
		pageSize := "5"
		// totalSecrets := 2

		// mockService
		mockService.On("ListUserSecrets", user, query, page, pageSize).Return(&corev1.SecretList{}, 0, errors.New("failed to list secrets"))

		req, _ := http.NewRequest("GET", "/api/v1/secrets?page=1&page_size=5", nil)
		req.Header.Set("authorization", "bearer valid_token")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		errorResp := response.ErrorResponse{}
		_ = json.NewDecoder(w.Body).Decode(&errorResp)
		assert.Equal(t, "failed to list secrets", errorResp.Error)

		mockService.AssertExpectations(t)
	})
}

func TestListUserSecrets3(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	mockService := new(MockService)
	server := &Server{
		router:  gin.Default(),
		service: mockService,
	}
	// server.setupRouter()
	// server.router.Use(AuthMiddlewareTest())
	server.router.GET("/api/v1/secrets", server.listPage)

	t.Run("failed to list secrets", func(t *testing.T) {
		// user := "testuser"
		// query := ""
		// page := "1"
		// pageSize := "5"
		// totalSecrets := 2

		// mockService
		// if not commented, the call to ListUserSecrets is expected but will not occur in the sever.listPage
		// mockService.On("ListUserSecrets", user, query, page, pageSize).Return(&corev1.SecretList{}, 0, nil)

		req, _ := http.NewRequest("GET", "/api/v1/secrets?page=1&page_size=5", nil)
		req.Header.Set("authorization", "bearer valid_token")
		w := httptest.NewRecorder()

		server.router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// errorResp := response.ErrorResponse{}
		// _ = json.NewDecoder(w.Body).Decode(&errorResp)
		var decodedResp string
		_ = json.Unmarshal(w.Body.Bytes(), &decodedResp)
		assert.Equal(t, "can't get user", decodedResp)

		mockService.AssertExpectations(t)
	})
}
