package main

import (
	"mongrix/internal"

	"github.com/gin-gonic/gin"
)


func main() {

    server := gin.Default()

	client, err := internal.ConnectToMongoDB()
	if err != nil {
		panic(err.Error())
	}


    server.GET("reminders", func(ctx *gin.Context) {
        internal.GetReminders(client)
    })


    server.Run(":8080")
	
}


