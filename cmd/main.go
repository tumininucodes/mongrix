package main

import (
	"fmt"
	"mongrix/internal"
    
)


func main() {


	client, err := internal.ConnectToMongoDB()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("client: ", client)
	
}


