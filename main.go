package main

import (
	"os"
	"secret-manager/pkg/v1/controller"
	"secret-manager/pkg/v1/service"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// @title Secret Manager API
// @Secret-Manager-API
// @version 2.0
// @description This is a Secret-Manager-API server.
// @termsOfService http://swagger.io/terms/

// @contact.email info@sercretmanager.com
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// fmt.Println("Hello world")

	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}

	client, err := service.CreateClient()
	if err != nil {
		logrus.Fatal("Can't create client: ", err)
	}

	server, err := controller.NewServer(service.NewService(client))
	if err != nil {
		logrus.Fatal("Can't create server: ", err)
	}

	err = server.Start(os.Getenv("PORT"))
	if err != nil {
		logrus.Fatal("Unable to start server ", err)
	}
}
