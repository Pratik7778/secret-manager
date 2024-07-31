package controller

import (
	"net/http"
	"os"
	"secret-manager/pkg/v1/middleware"
	"secret-manager/pkg/v1/response"
	"secret-manager/pkg/v1/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Listpage godoc
// @Summary      List secrets page
// @Tags         List Secrets
// @Security 	 ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        q query string false "Query parameter for search"
// @Param        page query integer false "Page number"
// @Param        page_size query integer false "Page size"
// @Success      200 {object} response.ListResponse
// @Failure      401 {object} response.ErrorResponse
// @Failure      404 {object} response.ErrorResponse
// @Router       /secrets [get]
func (server *Server) listPage(ctx *gin.Context) {
	user, exists := ctx.Get(middleware.AuthorizationPayloadKey1)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, "can't get user")
		return
	}

	// message := "Welcome to list secrets page " + user.(string)
	// ctx.JSON(http.StatusOK, message)

	query := ""
	if ctx.Query("q") != "" {
		query = os.Getenv("SECRET_PREFIX") + ctx.Query("q")
	}
	secretList, totalSecrets, err := server.service.ListUserSecrets(user.(string), query, ctx.Query("page"), ctx.Query("page_size"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Error: err.Error()})
		return
	}

	//Iterate over the list
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = service.DefaultPageNum
	}
	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil || pageSize < 1 || pageSize > service.DefaultPageSize {
		pageSize = service.DefaultPageSize
	}

	response := response.ListResponse{
		Secrets:    make([]string, len(secretList.Items)),
		Total:      totalSecrets,
		PageNumber: page,
		PageSize:   pageSize,
		Query:      ctx.Query("q"),
	}

	for i, s := range secretList.Items {
		// fmt.Println(s.Name)
		response.Secrets[i] = strings.TrimPrefix(s.Name, os.Getenv("SECRET_PREFIX"))
	}

	ctx.JSON(http.StatusOK, response)
}
