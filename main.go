package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/lucmig/debt-doodle-api/config"
	"github.com/lucmig/debt-doodle-api/routes"
)

func main() {
	// Database
	config.Connect()

	// Init Router
	router := gin.Default()

	// Route Handlers / Endpoints
	routes.Routes(router)

	log.Fatal(router.Run(":9000"))
}
