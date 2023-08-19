package main

import (
	"context"
	"mongrix/internal"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	server := gin.Default()

	context := context.Background()

	client, err := internal.ConnectToMongoDB(&context)
	if err != nil {
		panic(err.Error())
	}

	server.GET("objects", func(ctx *gin.Context) {
		results, err := internal.GetObjects(client, &context)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, results)
	})

	server.POST("add", func(ctx *gin.Context) {

	})

	server.DELETE(":id", func(ctx *gin.Context) {

	})

	server.PATCH(":id", func(ctx *gin.Context) {

	})

	server.PUT(":id", func(ctx *gin.Context) {

	})

	server.Run(":8080")

}
