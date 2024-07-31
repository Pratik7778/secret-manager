package controller

import (
	"net/http"
	"os"
	"secret-manager/pkg/v1/middleware"
	"secret-manager/pkg/v1/models"
	"secret-manager/pkg/v1/response"

	"github.com/gin-gonic/gin"
)

// Update Secret godoc
// @Summary      Update secret
// @Description  Update secret in the user's namespace
// @Tags         Update Secret
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        secret_name path string true "Secret name"
// @Param        detail body models.UpdateSecret true "Secret value in JSON format"
// @Success      200 {object} response.SuccessResponse
// @Failure      401 {object} response.ErrorResponse
// @Failure      400 {object} response.ErrorResponse
// @Router       /secrets/{secret_name} [put]
func (server *Server) updateSecretPage(ctx *gin.Context) {
	var secretval models.UpdateSecret
	if err := ctx.BindJSON(&secretval); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid JSON"})
		return
	}

	key := ctx.Param("secret_name")
	if key == "" || secretval.Value == "" {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "key and value are required"})
		return
	}

	user, exists := ctx.Get(middleware.AuthorizationPayloadKey1)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, "can't get user")
		return
	}

	key = os.Getenv("SECRET_PREFIX") + key
	err := server.service.UpdateSecret(user.(string), key, secretval)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{Message: "secret updated successfully"})
}
