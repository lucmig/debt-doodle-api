package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lucmig/debt-doodle-api/api"
)

// Routes - api end points
func Routes(router *gin.Engine) {
	router.GET("/", welcome)
	router.GET("/debts", api.GetAllDebts)
	router.POST("/debts", api.AddDebt)
	router.GET("/debts/:debtId", api.GetDebt)
	router.PUT("/debts/:debtId", api.EditDebt)
	router.DELETE("/debts/:debtId", api.DeleteDebt)

	router.PUT("/debts/:debtId/samples", api.UpdateSample)
	router.DELETE("/debts/:debtId/samples/:timestamp", api.DeleteSample)

	router.NoRoute(notFound)
}

func welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Welcome To API",
	})
	return
}

func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status":  404,
		"message": "Route Not Found",
	})
	return
}
