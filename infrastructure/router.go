package infrastructure

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init() {
	// Echo instance
	e := echo.New()

	api, err := InjectAPIServer()
	if err != nil {
		log.Fatal(err)
	}
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/users", api.User.Users)
	e.GET("/user/:id", api.User.UserByID)
	e.POST("/create", api.User.Add)
	e.PUT("/user/:id", api.User.Update)
	e.DELETE("/user/:id", api.User.DeleteByID)
	e.GET("/ping", api.Ping.Ping)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
