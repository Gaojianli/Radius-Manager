package main

import (
	"log"

	"github.com/cloudwego/hertz/pkg/app/server"

	"github.com/Gaojianli/raduis_mgnt/config"
	"github.com/Gaojianli/raduis_mgnt/database"
	"github.com/Gaojianli/raduis_mgnt/middleware"
	"github.com/Gaojianli/raduis_mgnt/routes"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := middleware.InitJWT(); err != nil {
		log.Fatal("Failed to initialize JWT middleware:", err)
	}

	h := server.Default(server.WithHostPorts(config.AppConfig.ServerPort))

	routes.SetupRoutes(h)

	log.Printf("Server starting on port %s", config.AppConfig.ServerPort)
	h.Spin()
}
