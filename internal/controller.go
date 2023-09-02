package internal

import (
	"context"
	"errors"
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

func GetObjects(coll *mongo.Collection, ctx *context.Context) (*bson.A, error) {
	cursor, err := coll.Find(*ctx, bson.M{})
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


func InsertObject(object *bson.M, coll *mongo.Collection, ctx *context.Context) (*bson.M, error) {
	result, err := coll.InsertOne(*ctx, object)
	if err != nil {
		return nil, err
	}
	var inserted bson.M
	coll.FindOne(*ctx, bson.M{"_id": result.InsertedID}).Decode(&inserted)
	println("inserted id -> ", result.InsertedID)
	return &inserted, nil
}


func UpdateObject(id *primitive.ObjectID, coll *mongo.Collection, ctx *context.Context, data *map[string]interface{}) (*bson.M, error) {
	filter := bson.M{
		"_id": id,
	}
	update := bson.M{
		"$set": data,
	}
	_, err := coll.UpdateOne(*ctx, filter, update)
	if err != nil {
		return nil, err
	}
	var upserted bson.M
	coll.FindOne(*ctx, bson.M{"_id": id}).Decode(&upserted)
	return &upserted, nil
}

func ReplaceObject(id *primitive.ObjectID, coll *mongo.Collection, ctx *context.Context, data *map[string]interface{}) (*bson.M, error) {
	filter := bson.M{
		"_id": id,
	}
	replacement := bson.M{
		"_id": id,
		"data": data,
	}
	
	result, err := coll.ReplaceOne(*ctx, filter, replacement)
	if err != nil {
		return nil, err
	}

	if result.ModifiedCount > 0 {
		return &replacement, nil
	} else {
		return nil, errors.New("error replacing object")
	}
}

func DeleteObject(id *primitive.ObjectID, coll *mongo.Collection, ctx *context.Context) (bool, error) {
	filter := bson.M{
		"_id": id,
	}
	result, err := coll.DeleteOne(*ctx, filter)
	if err != nil {
		println("error deleting ->", err.Error())
		return false, err
	}
	if result.DeletedCount > 0 {
		return true, nil
	} else {
		return false, nil
	}
	
}
