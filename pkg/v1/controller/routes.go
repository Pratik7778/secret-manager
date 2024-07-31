package controller

import (
	"net/http"
	"os"
	"secret-manager/docs"
	"secret-manager/pkg/v1/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (server *Server) setupRouter() {
	server.router.Use(middleware.CORSMiddleware())
	server.setupSwagger()
	apiRoutes := server.router.Group("/api/v1")
	server.setupBasicRoutes(apiRoutes)

	secretsRoute := apiRoutes.Group("/secrets")
	_ = secretsRoute.Use(middleware.AuthMiddleware(server.service))
	server.setupuserRoutes(secretsRoute)
}

func (server *Server) setupBasicRoutes(routes *gin.RouterGroup) {
	routes.GET("", server.homePage)
	routes.POST("/register", server.registerPage)
	routes.POST("/login", server.loginPage)
}

func (server *Server) setupuserRoutes(routes *gin.RouterGroup) {
	routes.GET("", server.listPage)
	routes.POST("/create", server.createSecretPage)
	routes.PUT("/:secret_name", server.updateSecretPage)
	routes.DELETE("/:secret_name", server.deleteSecretPage)
	routes.GET("/:secret_name", server.viewSecretPage)
}

// Homepage godoc
// @Summary      Show Secret Manager home page message
// @Tags         Home
// @Accept       json
// @Produce      json
// @Success      200
// @Router       / [get]
func (server *Server) homePage(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "Welcome to the secret manager")
}

func (server *Server) setupSwagger() {
	hostPath := os.Getenv("API_URL")
	if hostPath == "" {
		hostPath = "localhost:8080"
	}
	basePath := os.Getenv("BASEPATH")
	if basePath == "" {
		basePath = "/api/v1"
	}
	logrus.Info("hostapth, basePath: ", hostPath, basePath)
	docs.SwaggerInfo.Host = hostPath
	docs.SwaggerInfo.BasePath = basePath
	server.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
