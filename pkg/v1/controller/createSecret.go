package controller

import (
	"net/http"
	"os"
	"secret-manager/pkg/v1/middleware"
	"secret-manager/pkg/v1/models"
	"secret-manager/pkg/v1/response"

	"github.com/gin-gonic/gin"
)

// Create Secret godoc
// @Summary      Create secret
// @Description  Create secret in the user's namespace
// @Tags         Create Secret
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        detail body models.Secret true "Secret details in JSON format"
// @Success      200 {object} response.SuccessResponse
// @Failure      401 {object} response.ErrorResponse
// @Failure      400 {object} response.ErrorResponse
// @Router       /secrets/create [post]
func (server *Server) createSecretPage(ctx *gin.Context) {
	var secretval models.Secret
	if err := ctx.BindJSON(&secretval); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid JSON"})
		return
	}

	if secretval.Key == "" || secretval.Value == "" {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "key and value are required"})
		return
	}

	user, exists := ctx.Get(middleware.AuthorizationPayloadKey1)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, "can't get user")
		return
	}

	secretval.Key = os.Getenv("SECRET_PREFIX") + secretval.Key
	err := server.service.CreateSecret(user.(string), secretval)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{Message: "secret created successfully"})
}
