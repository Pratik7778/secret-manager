package controller

import (
	"fmt"
	"secret-manager/pkg/v1/service"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router  *gin.Engine
	service service.IService
}

func NewServer(service service.IService) (*Server, error) {
	server := &Server{}
	server.router = gin.Default()
	server.service = service
	server.setupRouter()
	return server, nil
}

func (server *Server) Start(port string) error {
	if port == "" {
		port = "8080"
	}

	return server.router.Run(fmt.Sprintf(":%s", port))
}
