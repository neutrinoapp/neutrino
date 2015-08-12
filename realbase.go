package main

import (
//	"fmt"
//	"gopkg.in/mgo.v2/bson"
//	"realbase/core"
	"net/http"
	"./api"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"realbase/core"
)

func initMiddleware(e *echo.Echo) {
	e.Use(mw.Logger())
	e.Use(mw.Recover())
}

func initRoutes(e *echo.Echo) {
	e.Get("/", func (c *echo.Context) {
		c.String(http.StatusOK, "haha")
	})

	e.Post("/auth", func (c *echo.Context) {
		api.RegisterUserHandler(c)
	})
}

func main() {
	realbase.Initialize("localhost:27017")

	e := echo.New()

	initMiddleware(e)
	initRoutes(e)

	e.Run(":1234")

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
