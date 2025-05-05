package database

import (
	"context"
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoURI = os.Getenv("MONGO_URI")

func New() *mongo.Database {
	clientOpts := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(context.Background(), clientOpts)
	if err != nil {
		panic(fmt.Sprintf("cannot connect to db: %s", err))
	}

	return client.Database("blooprint")
}
