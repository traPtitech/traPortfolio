package infrastructure

import (
	"log"

	"github.com/traPtitech/traPortfolio/interfaces/handler"

	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init(s *SQLConfig, t *TraQConfig, p *PortalConfig, k *KnoQConfig) {
	// Echo instance
	e := echo.New()
	v := validator.New()
	if err := v.RegisterValidation("is-uuid", handler.IsValidUUID); err != nil {
		log.Fatal(err)
	}
	e.Validator = &Validator{
		validator: v,
	}

	api, err := InjectAPIServer(s, t, p, k)
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

					apiUsersUIDProjects := apiUsersUID.Group("/projects")

					apiUsersUIDProjects.GET("", api.User.GetProjects)

					apiUsersUIDContests := apiUsersUID.Group("/contests")

					apiUsersUIDContests.GET("", api.User.GetContests)

					// apiUsersUIDGroups := apiUsersUID.Group("/groups")

					// apiUsersUIDGroups.GET("", api.User.GetGroups)

					apiUsersUIDEvents := apiUsersUID.Group("/events")

					apiUsersUIDEvents.GET("", api.User.GetEvents)
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

				apiProjectsPIDMembers.GET("", api.Project.GetProjectMembers)
				apiProjectsPIDMembers.POST("", api.Project.AddProjectMembers)
				apiProjectsPIDMembers.DELETE("", api.Project.DeleteProjectMembers)
			}
		}

		{
			apiEvents := v1.Group("/events")

			apiEvents.GET("", api.Event.GetAll)

			apiEventsEID := apiEvents.Group("/:eventID")

			apiEventsEID.GET("", api.Event.GetByID)
			apiEventsEID.PATCH("", api.Event.PatchEvent)
		}

		//{
		// 	apiGroups := v1.Group("/groups")

		//	apiGroups.GET("", api.Group.GetAll)
		// 	apiGroups.GET("/:groupID", api.Group.GetByID)
		//}

		{
			apiContests := v1.Group("/contests")

			apiContests.GET("", api.Contest.GetContests)
			apiContests.POST("", api.Contest.PostContest)
			{
				apiContestsCID := apiContests.Group("/:contestID")

				apiContestsCID.GET("", api.Contest.GetContest)
				apiContestsCID.PATCH("", api.Contest.PatchContest)
				apiContestsCID.DELETE("", api.Contest.DeleteContest)
				{
					apiContestsCIDTeams := apiContestsCID.Group("/teams")

					apiContestsCIDTeams.GET("", api.Contest.GetContestTeams)
					apiContestsCIDTeams.POST("", api.Contest.PostContestTeam)
					{
						apiContestsCIDTeamsTID := apiContestsCIDTeams.Group("/:teamID")

						apiContestsCIDTeamsTID.GET("", api.Contest.GetContestTeam)
						apiContestsCIDTeamsTID.PATCH("", api.Contest.PatchContestTeam)
						// apiContestsCIDTeamsTID.DELETE("", api.Contest.DeleteContestTeam)
						{
							apiContestsCIDTeamsTIDMembers := apiContestsCIDTeamsTID.Group("/members")

							apiContestsCIDTeamsTIDMembers.GET("", api.Contest.GetContestTeamMember)
							apiContestsCIDTeamsTIDMembers.POST("", api.Contest.PostContestTeamMember)
							apiContestsCIDTeamsTIDMembers.DELETE("", api.Contest.DeleteContestTeamMember)
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
