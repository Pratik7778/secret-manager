package controller

import (
	"net/http"
	"os"
	"secret-manager/pkg/v1/models"
	"secret-manager/pkg/v1/response"

	"github.com/gin-gonic/gin"
)

// RegisterUser godoc
// @Summary      Register a new user
// @Description  Registers a new user in the secret manager
// @Tags         Registration
// @Accept       json
// @Produce      json
// @Param        detail body models.User true "User details in JSON format"
// @Success      200 {object} response.SuccessResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      400 {object} response.ErrorResponse
// @Router       /register [post]
func (server *Server) registerPage(ctx *gin.Context) {
	var detail models.User
	if err := ctx.BindJSON(&detail); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Invalid JSON"})
		return
	}

	if detail.Username == "" || detail.Password == "" {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Username and password are required"})
		return
	}

	// err := client.CreateUserSecret(detail)
	//Add prefix to username
	prefix := os.Getenv("SECRET_PREFIX")
	detail.Username = prefix + detail.Username

	err := server.service.CreateUserSecret(detail)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{Message: "User created successfully"})
}
