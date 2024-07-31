package controller

import (
	"net/http"
	"os"
	"secret-manager/pkg/v1/middleware"
	"secret-manager/pkg/v1/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// View Secret godoc
// @Summary      View secret
// @Description  View secret value in the user's namespace
// @Tags         View Secret
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        secret_name path string true "Secret name"
// @Success      200 {object} response.SuccessResponse
// @Failure      401 {object} response.ErrorResponse
// @Failure      400 {object} response.ErrorResponse
// @Router       /secrets/{secret_name} [get]
func (server *Server) viewSecretPage(ctx *gin.Context) {
	key := ctx.Param("secret_name")
	if key == "" {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "secret_name is required"})
		return
	}

	user, exists := ctx.Get(middleware.AuthorizationPayloadKey1)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, "can't get user")
		return
	}

	key = os.Getenv("SECRET_PREFIX") + key
	secretData, err := server.service.ViewSecret(user.(string), key)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	stringMap := make(map[string]string, len(secretData))
	for k, v := range secretData {
		k = strings.TrimPrefix(k, os.Getenv("SECRET_PREFIX"))
		stringMap[k] = string(v)
	}

	ctx.JSON(http.StatusOK, stringMap)
}
