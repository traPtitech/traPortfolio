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
	{
		v1 := echoAPI.Group("/v1")

		{
			apiUsers := v1.Group("/users")

			apiUsers.GET("", api.User.GetAll)
			{
				apiUsersUID := apiUsers.Group("/:userID")

				apiUsersUID.GET("", api.User.GetByID)
				apiUsersUID.PATCH("", api.User.Update)
				{
					apiUsersUIDAccounts := apiUsersUID.Group("/accounts")

					apiUsersUIDAccounts.GET("", api.User.GetAccounts)
					apiUsersUIDAccounts.POST("", api.User.AddAccount)
					{
						apiUsersUIDAccountsAID := apiUsersUIDAccounts.Group("/:accountID")

						apiUsersUIDAccountsAID.GET("", api.User.GetAccount)
						apiUsersUIDAccountsAID.PATCH("", api.User.PatchAccount)
						apiUsersUIDAccountsAID.DELETE("", api.User.DeleteAccount)
					}
				}
			}
		}
		{
			apiProjects := v1.Group("/projects")

			apiProjects.GET("", api.Project.GetAll)
			apiProjects.POST("", api.Project.PostProject)

			{
				apiProjectsPID := apiProjects.Group("/:projectID")

				apiProjectsPID.GET("", api.Project.GetByID)
				apiProjectsPID.PATCH("", api.Project.PatchProject)

				apiProjectsPIDMembers := apiProjectsPID.Group("/members")

				apiProjectsPIDMembers.GET("", api.Project.GetMembers)
				// apiProjectsPIDMembers.POST("", api.Project.PostMembers)
				// apiProjectsPIDMembers.DELETE("", api.Project.DeleteMembers)
			}
		}

		{
			apiEvents := v1.Group("/events")

			apiEvents.GET("", api.Event.GetAll)
			apiEvents.GET("/:eventID", api.Event.GetByID)
		}

		{
			apiContests := v1.Group("/contests")

			apiContests.POST("", api.Contest.PostContest)
			{
				apiContestsCID := apiContests.Group("/:contestID")

				apiContestsCID.PATCH("", api.Contest.PatchContest)
				{
					apiContestsCIDTeams := apiContestsCID.Group("/teams")

					apiContestsCIDTeams.POST("", api.Contest.PostContestTeam)
					{
						apiContestsCIDTeamsTID := apiContestsCIDTeams.Group("/:teamID")

						apiContestsCIDTeamsTID.PATCH("", api.Contest.PatchContestTeam)
						{
							apiContestsCIDTeamsTIDMembers := apiContestsCIDTeamsTID.Group("/members")

							apiContestsCIDTeamsTIDMembers.DELETE("", api.Contest.DeleteContestTeamMember)
							apiContestsCIDTeamsTIDMembers.POST("", api.Contest.PostContestTeamMember)
						}
					}
				}
			}
		}

		{
			apiPing := v1.Group("/ping")

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
