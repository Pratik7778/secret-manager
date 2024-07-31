package middleware

import (
	"fmt"
	"net/http"
	"secret-manager/pkg/v1/response"
	"secret-manager/pkg/v1/service"
	"secret-manager/pkg/v1/token"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey   = "authorization"
	authorizationTypeBearer  = "bearer"
	authorizationPayloadKey  = "authorization_payload"
	AuthorizationPayloadKey1 = "X-User-Token"
)

func AuthJWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{Error: "authorization isn't provided"})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{Error: "invalid authorization header format"})
			return

		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{Error: err.Error()})
			return
		}

		accessToken := fields[1]
		payload, err := token.VerifyJWTToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{Error: err.Error()})
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		// fmt.Println("______Verified auth____________")
	}
}

func AuthMiddleware(service service.IService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{Error: "authorization isn't provided"})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{Error: "invalid authorization header format"})
			return

		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{Error: err.Error()})
			return
		}

		accessToken := fields[1]
		payload, err := token.VerifyToken(service, accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse{Error: err.Error()})
			return
		}

		ctx.Set(AuthorizationPayloadKey1, payload)
		ctx.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, x-agent-code")
		ctx.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	}
}
