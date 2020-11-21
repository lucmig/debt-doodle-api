package routes

import (
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/lucmig/debt-doodle-api/api"
)

// Routes - api end points
func Routes(router *gin.Engine) {
	// swagger:route GET /debts debts
	//
	// Lists debts
	//
	// Return all debts.
	//
	// 	Consumes:
	//  	- application/json
	//    - application/x-protobuf
	//
	//	Produces:
	//  	- application/json
	//    - application/x-protobuf
	//
	//	Schemes: http, https, ws, wss
	//
	//	Security:
	//		api_key:
	//    oauth: read, write
	//
	//	responses:
	//		204:
	//			description: No data
	//		200:
	//			description: An array of debts
	//			content:
	//				application/json:
	//					schema:
	//						type: Array
	//						items:
	//							$ref: #/definitions/Debt
	//		400:
	//			description: Bad user input
	//			content:
	//				$ref: #/responses/validationError
	router.GET("/debts", api.GetAllDebts)
	router.POST("/debts", api.AddDebt)
	router.GET("/debts/:debtId", api.GetDebt)
	router.PUT("/debts/:debtId", api.EditDebt)
	router.DELETE("/debts/:debtId", api.DeleteDebt)

	router.PUT("/debts/:debtId/samples", api.UpdateSample)
	router.DELETE("/debts/:debtId/samples/:date", api.DeleteSample)

	router.GET("/values/:date", api.GetValuesDate)

	router.Use(static.Serve("/", static.LocalFile("./swagger-ui", false)))

	router.NoRoute(notFound)
}

func notFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status":  404,
		"message": "Route Not Found",
	})
	return
}
