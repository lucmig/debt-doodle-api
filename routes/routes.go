package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucmig/debt-doodle-api/api"
)

// Init - routes for debt-doodle api
func Init() *gin.Engine {
	app := gin.New()

	app.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, api.Pong())
	})

	app.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	return app
}
