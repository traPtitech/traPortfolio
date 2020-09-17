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

	echoAPI := e.Group("/api")
	v1 := echoAPI.Group("/v1")
	{
		apiUsers := v1.Group("/users")
		{
			apiUsers.PUT("/:id", api.User.Update)
		}
		apiPing := v1.Group("/ping")
		{
			apiPing.GET("", api.Ping.Ping)
		}
	}

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
