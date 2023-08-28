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

	db := client.Database("todo")

	server.GET("objects", func(ctx *gin.Context) {
		results, err := internal.GetObjects(client, &context)
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

		result, err := internal.InsertObject(&bsonObject, db, &context)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusCreated, result)
	})

	server.PATCH(":id", func(ctx *gin.Context) {
		idString := ctx.Param("id")
		mongoId, err := primitive.ObjectIDFromHex(idString)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := internal.UpdateObject(mongoId, db, &context)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	})

	server.DELETE(":id", func(ctx *gin.Context) {

	})


	server.PUT(":id", func(ctx *gin.Context) {

	})

	server.Run(":8080")

}
