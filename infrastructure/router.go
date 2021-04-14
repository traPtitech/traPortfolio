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

	api, err := InjectAPIServer()
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
			apiUsers.PUT("/:userID/accounts", api.User.AddAccount)
			apiUsers.DELETE("/:userID/accounts/:accountID", api.User.DeleteAccount)
		}
		apiProjects := v1.Group("/projects")
		{
			apiProjects.GET("", api.Project.GetAll)
			apiProjects.GET("/:projectID", api.Project.GetByID)
			apiProjects.POST("", api.Project.PostProject)
			apiProjects.PATCH("/:projectID", api.Project.PatchProject)
		}
		apiEvents := v1.Group("/events")
		{
			apiEvents.GET("", api.Event.GetAll)
			apiEvents.GET("/:eventID", api.Event.GetByID)
		}
		apiContests := v1.Group("/contests")
		{
			apiContests.POST("", api.Contest.PostContest)
			apiContestsCID := apiContests.Group("/:contestID")
			{
				apiContestsCID.PATCH("", api.Contest.PatchContest)
				apiContestsCID.POST("", api.Contest.PostContestTeam)
				apiContestsCIDTID := apiContestsCID.Group("/:teamID")
				{
					apiContestsCIDTID.PATCH("", api.Contest.PatchContestTeam)
					apiContestsCIDTID.PUT("", api.Contest.PutContestTeamMember)
					apiContestsCIDTID.DELETE("", api.Contest.DeleteContestTeamMember)
				}
			}
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
