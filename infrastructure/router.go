package infrastructure

import (
	"log"

	"github.com/traPtitech/traPortfolio/interfaces/handler"

	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init() {
	// Echo instance
	e := echo.New()
	e.Validator = &Validator{
		validator: validator.New(),
	}

	api, err := InjectAPIServer("traQToken", "portalToken")
	if err != nil {
		log.Fatal(err)
	}
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return h(&handler.Context{Context: c})
		}
	})

	echoAPI := e.Group("/api")
	v1 := echoAPI.Group("/v1")
	{
		apiUsers := v1.Group("/users")
		{
			apiUsers.GET("", api.User.GetAll)
			apiUsers.GET("/:userID", api.User.GetByID)
			apiUsers.PATCH("/:userID", api.User.Update)
		}
		apiProjects := v1.Group("/projects")
		{
			apiProjects.POST("", api.Project.PostProject)
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

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
