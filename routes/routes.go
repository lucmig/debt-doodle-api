package routes

import (
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/lucmig/debt-doodle-api/api"
)

// Routes - api end points
func Routes(router *gin.Engine) {
	// swagger:route GET /debts debts GetAllDebts
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

	// swagger:route POST /debts debts AddDebt
	//
	// Adds a debt
	//
	// Add a debt
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
	//	requestBody:
	//		description: a debt
	//		required: true
	//		content:
	//			application/json:
	//				schema:
	//					$ref: #/definitions/Debt
	//
	//	responses:
	//		200:
	//			description: An array of debts
	//			content:
	//				application/json:
	//					schema:
	//						type: Object
	//						items:
	//							$ref: #/definitions/Debt
	//		400:
	//			description: Bad user input
	//			content:
	//				$ref: #/responses/validationError
	router.POST("/debts", api.AddDebt)

	// swagger:route GET /debts/:debtId debts GetDebt
	router.GET("/debts/:debtId", api.GetDebt)

	// swagger:route PUT /debts/:debtId debts EditDebt
	router.PUT("/debts/:debtId", api.EditDebt)

	// swagger:route DELETE /debts/:debtId debts DeleteDebt
	router.DELETE("/debts/:debtId", api.DeleteDebt)

	// swagger:route PUT /debts/:debtId/samples samples UpdateSample
	router.PUT("/debts/:debtId/samples", api.UpdateSample)

	// swagger:route DELETE /debts/:debtId/samples/:date samples DeleteSample
	router.DELETE("/debts/:debtId/samples/:date", api.DeleteSample)

	// swagger:route GET /values/:date values GetValuesDate
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
