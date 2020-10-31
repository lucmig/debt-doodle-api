package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/lucmig/debt-doodle-api/api"
)

// Connect - Connect to db
func Connect() {
	// Database Config
	clientOptions := options.Client().ApplyURI("mongodb+srv://debt-doodle-rw:debt-doodle-pw@debt-doodle.dxvbs.mongodb.net/debt-doodle?retryWrites=true&w=majority")

	client, err := mongo.NewClient(clientOptions)

	//Set up a context required by mongo.Connect
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	//To close the connection at the end
	defer cancel()

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected!")
	}
	db := client.Database("debt-doodle")
	api.DebtCollection(db)
	return
}
