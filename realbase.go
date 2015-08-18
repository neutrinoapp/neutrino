package main

import (
//	"fmt"
//	"gopkg.in/mgo.v2/bson"
//	"realbase/core"
	"github.com/labstack/echo"
	"realbase/core"
	"fmt"
	"realbase/api"
)

func main() {
	e := echo.New()

	realbase.Initialize("localhost:27017")
	api.Initialize(e)

	port := ":1234";

	fmt.Println("Listening on port", port)
	e.Run(port)

//	mongodbHost := "localhost:27017"
//
//	dbService := realbase.NewDbService(mongodbHost, "test", "test")
//
//	socketHandler := realbase.GetMessageService().InitSocketHandler()
//
//	go func() {
//		http.ListenAndServe(":5555", socketHandler)
//	}()
//
//	for {
//		fmt.Print("Enter text: ")
//		var input string
//		fmt.Scanln(&input)
//		dbService.Insert(bson.M{"a": input})
//	}

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
