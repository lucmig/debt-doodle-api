package main

import (
	// "github.com/gin-gonic/gin"
	// graphql "github.com/graph-gophers/graphql-go"
	// "github.com/graph-gophers/graphql-go/relay"
	// "github.com/lucmig/debt-doodle-api/db"

	"github.com/lucmig/debt-doodle-api/routes"
)

// type query struct{}
// func (_ *query) Hello() string { return "Hello, world!" }

func main() {
	app := routes.Init()
	app.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
