package controller

import (
	"net/http"
	"os"
	"secret-manager/pkg/v1/middleware"
	"secret-manager/pkg/v1/response"

	"github.com/gin-gonic/gin"
)

// Delete Secret godoc
// @Summary      Delete secret
// @Description  Delete secret in the user's namespace
// @Tags         Delete Secret
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        secret_name path string true "Secret name"
// @Success      200 {object} response.SuccessResponse
// @Failure      401 {object} response.ErrorResponse
// @Failure      400 {object} response.ErrorResponse
// @Router       /secrets/{secret_name} [delete]
func (server *Server) deleteSecretPage(ctx *gin.Context) {
	key := ctx.Param("secret_name")
	if key == "" {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "key is required"})
		return
	}

	user, exists := ctx.Get(middleware.AuthorizationPayloadKey1)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, "can't get user")
		return
	}

	prefix := os.Getenv("SECRET_PREFIX")
	key = prefix + key
	err := server.service.DeleteSecret(user.(string), key)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{Message: "secret deleted successfully"})
}
