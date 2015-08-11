package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"realbase/core"

	//"log"
	"net/http"
)

func main() {
	mongodbHost := "localhost:27017"

	dbService := realbase.GetDbService(mongodbHost, "test", "test")

	socketHandler := realbase.GetMessageService().InitSocketHandler()

	go func() {
		http.ListenAndServe(":5555", socketHandler)
	}()

	for {
		fmt.Print("Enter text: ")
		var input string
		fmt.Scanln(&input)
		dbService.Insert(bson.M{"a": input})
	}

	// service :=
	// result := bson.M{"test11": "test11"}
	// err := service.Insert(result)
	//
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Println(result)
}
