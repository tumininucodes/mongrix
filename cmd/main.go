package main

import (
	"context"
	"encoding/json"
	"fmt"
	"mongrix/internal"
	"net/http"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {

	server := gin.Default()

	server.Use(responseModifierMiddleware)

	
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
		// var inputData map[string]interface{}

		var todo internal.Todo

		if err := ctx.ShouldBindJSON(&todo); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
			return
		}

		result, err := internal.InsertObject(todo, coll, &context)
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
		idString := ctx.Param("id")
		id, err := primitive.ObjectIDFromHex(idString)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var data map[string]interface{}
		if err := ctx.ShouldBindJSON(&data); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := internal.ReplaceObject(&id, coll, &context, &data)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, result)
	})

	

	server.Run(":8081")

}



func responseModifierMiddleware(ctx *gin.Context) {

	blw := &interceptResponseWriter{ctx.Writer, []byte{}, ctx.Writer.Status()}
    ctx.Writer = blw
    ctx.Next()

	var objects []map[string]interface{}

	if err := json.Unmarshal(blw.body, &objects); err != nil {
		fmt.Println("Error1:", err.Error())
		return
	}

	fmt.Println("objects ", objects)

	for _, obj := range objects {
		if name, ok := obj["name"].(string); ok {
			obj["name"] = "The " + name
		}
	}

	modifiedJSON, err := json.Marshal(objects)
	if err != nil {
		fmt.Println("Error2:", err)
		return
	}

	blw.body = modifiedJSON

	ctx.Writer.WriteString(string(blw.body))

}

type interceptResponseWriter struct {
    gin.ResponseWriter
    body   []byte
    status int
}


func (w *interceptResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	w.body = append(w.body, b...)
	return len(b), nil
}