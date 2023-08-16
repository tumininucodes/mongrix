package main

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func main() {


	client, err := connectToMongoDB()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("client: ", client)
	
}


func connectToMongoDB() (*mongo.Client, error) {

	mongoURI := os.Getenv("MOROSHOI_MONGODB_URI")
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