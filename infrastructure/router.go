package infrastructure

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init() {
	// Echo instance
	e := echo.New()

	api, err := InjectAPIServer("traQToken", "portalToken")
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
			apiUsers.GET("", api.User.GetAll)
			apiUsers.GET("/:userID", api.User.GetByID)
			apiUsers.PATCH("/:userID", api.User.Update)
		}
		apiEvents := v1.Group("/events")
		{
			apiEvents.GET("", api.Event.GetAll)
			apiEvents.GET("/:eventID", api.Event.GetByID)
		}
		apiContests := v1.Group("/contests")
		{
			apiContests.POST("", api.Contest.PostContest)
			apiContests.PATCH("/:contestID", api.Contest.PatchContest)
		}
		apiPing := v1.Group("/ping")
		{
			apiPing.GET("", api.Ping.Ping)
		}
	}

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
