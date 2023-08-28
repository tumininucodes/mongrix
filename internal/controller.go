package internal

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongoDB(ctx *context.Context) (*mongo.Client, error) {

	fmt.Println("Connecting to MongoDB....")

	err := godotenv.Load("mongrix.env")
	if err != nil {
		fmt.Println("Error loading .env file")
		return nil, err
	}

	mongoURI := os.Getenv("MOROSHOI_MONGODB_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("MOROSHOI_MONGODB_URI environment variable not set")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)

	client, err := mongo.Connect(*ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the MongoDB server to check the connection
	err = client.Ping(*ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")

	return client, nil
}

func GetObjects(client *mongo.Client, ctx *context.Context) (*bson.A, error) {
	collection := client.Database("todo").Collection("Todo")

	cursor, err := collection.Find(*ctx, bson.M{})
	if err != nil {
		fmt.Println("Error querying collection:", err)
		return nil, err
	}

	var results bson.A
	for cursor.Next(*ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			fmt.Println("Error decoding documents: ", err)
			return nil, err
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		fmt.Println("Cursor error:", err)
		return nil, err
	}

	defer cursor.Close(*ctx)

	return &results, nil
}


func InsertObject(object *bson.M, db *mongo.Database, ctx *context.Context) (bson.M, error) {
	collection := db.Collection("Todo")
	result, err := collection.InsertOne(*ctx, object)
	if err != nil {
		return nil, err
	}
	var inserted bson.M
	collection.FindOne(*ctx, bson.M{"_id": result.InsertedID}).Decode(&inserted)
	println("inserted id -> ", result.InsertedID)
	return inserted, nil
}


func UpdateObject(id primitive.ObjectID, db *mongo.Database, ctx *context.Context, data *map[string]interface{}) (bson.M, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	_, err := db.Collection("Todo").UpdateOne(*ctx, filter, data)
	if err != nil {
		return nil, err
	}
	var upserted bson.M
	println("data -> ", data)
	db.Collection("Todo").FindOne(*ctx, bson.M{"_id": id}).Decode(&upserted)
	return upserted, nil
}

func DeleteObject() {
	
}
