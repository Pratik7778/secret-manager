package controller

import (
	"net/http"
	"os"
	"secret-manager/pkg/v1/models"
	"secret-manager/pkg/v1/response"
	"secret-manager/pkg/v1/token"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary      Login existing users
// @Description  Login existing user in the secret manager
// @Tags         Login
// @Accept       json
// @Produce      json
// @Param        detail body models.User true "User details in JSON format"
// @Success      200 {object} response.LoginResponse
// @Failure      400 {object} response.ErrorResponse
// @Failure      400 {object} response.ErrorResponse
// @Router       /login [post]
func (server *Server) loginPage(ctx *gin.Context) {
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

	// token, err := token.CreateJWTToken(detail.Username)
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Error generating token"})
	// 	return
	// }

	token := token.CreateToken(10)
	if token == "" {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Error generating token"})
		return
	}

	// Add prefix to username
	prefix := os.Getenv("SECRET_PREFIX")
	detail.Username = prefix + detail.Username

	err := server.service.LoginUser(detail, token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response.LoginResponse{Token: token})
}
