package internal

import (
	"context"
	"fmt"
	"os"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongoDB() (*mongo.Client, error) {

	err := godotenv.Load("mongrix.env")
	if err != nil {
		fmt.Println("Error loading .env file")
		return nil, err
	}
	println("got env")


	mongoURI := os.Getenv("MOROSHOI_MONGODB_URI")
	println(mongoURI)
    if mongoURI == "" {
        return nil, fmt.Errorf("MOROSHOI_MONGODB_URI environment variable not set")
    }

    // Set client options
    clientOptions := options.Client().ApplyURI(mongoURI)

    // Connect to MongoDB
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        return nil, err
    }

    // Ping the MongoDB server to check the connection
    err = client.Ping(context.Background(), nil)
    if err != nil {
        return nil, err
    }

    fmt.Println("Connected to MongoDB!")

    return client, nil
}