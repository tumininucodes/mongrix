package main

import (
	"context"
	"mongrix/internal"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {

	server := gin.Default()

	context := context.Background()

	client, err := internal.ConnectToMongoDB(&context)
	if err != nil {
		panic(err.Error())
	}
	coll := client.Database("todo").Collection("Todo")


	server.GET("objects", func(ctx *gin.Context) {
		results, err := internal.GetObjects(coll, &context)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, results)
	})

	server.POST("add", func(ctx *gin.Context) {
		var inputData map[string]interface{}
		if err := ctx.ShouldBindJSON(&inputData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
			return
		}
		bsonObject := bson.M(inputData)

		result, err := internal.InsertObject(&bsonObject, coll, &context)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, result)
	})

	server.PATCH(":id", func(ctx *gin.Context) {
		idString := ctx.Param("id")
		var updateData map[string]interface{}

		if err := ctx.ShouldBindJSON(&updateData); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		mongoId, err := primitive.ObjectIDFromHex(idString)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := internal.UpdateObject(&mongoId, coll, &context, &updateData)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	})

	server.DELETE(":id", func(ctx *gin.Context) {
		idString := ctx.Param("id")
		id, err := primitive.ObjectIDFromHex(idString)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		deleted, err := internal.DeleteObject(&id, coll, &context)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"deleted": deleted})
	})


	server.PUT(":id", func(ctx *gin.Context) {

	})

	server.Run(":8080")

}
