package infrastructure

import (
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

	e.GET("/users", func(c echo.Context) error { return api.User.Index(c) })
	e.GET("/user/:id", func(c echo.Context) error { return api.User.Show(c) })
	e.POST("/create", func(c echo.Context) error { return api.User.Create(c) })
	e.PUT("/user/:id", func(c echo.Context) error { return api.User.Save(c) })
	e.DELETE("/user/:id", func(c echo.Context) error { return api.User.Delete(c) })
	e.GET("/ping", func(c echo.Context) error { return api.Ping.Ping(c) })

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
