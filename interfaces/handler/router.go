package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Setup(e *echo.Echo, api API) error {
	// Setup validator
	e.Validator = newValidator(e.Logger)

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return h(&Context{Context: c})
		}
	})

	echoAPI := e.Group("/api")
	{
		v1 := echoAPI.Group("/v1")

		{
			apiUsers := v1.Group("/users")

			apiUsers.GET("", api.User.GetUsers)
			{
				apiUsersUID := apiUsers.Group("/:userID")

				apiUsersUID.GET("", api.User.GetUser)
				apiUsersUID.PATCH("", api.User.UpdateUser)
				{
					apiUsersUIDAccounts := apiUsersUID.Group("/accounts")

					apiUsersUIDAccounts.GET("", api.User.GetUserAccounts)
					apiUsersUIDAccounts.POST("", api.User.AddUserAccount)
					{
						apiUsersUIDAccountsAID := apiUsersUIDAccounts.Group("/:accountID")

						apiUsersUIDAccountsAID.GET("", api.User.GetUserAccount)
						apiUsersUIDAccountsAID.PATCH("", api.User.EditUserAccount)
						apiUsersUIDAccountsAID.DELETE("", api.User.DeleteUserAccount)
					}

					apiUsersUIDProjects := apiUsersUID.Group("/projects")

					apiUsersUIDProjects.GET("", api.User.GetUserProjects)

					apiUsersUIDContests := apiUsersUID.Group("/contests")

					apiUsersUIDContests.GET("", api.User.GetUserContests)

					apiUsersUIDGroups := apiUsersUID.Group("/groups")

					apiUsersUIDGroups.GET("", api.User.GetUserGroups)

					apiUsersUIDEvents := apiUsersUID.Group("/events")

					apiUsersUIDEvents.GET("", api.User.GetUserEvents)
				}
			}
		}
		{
			apiProjects := v1.Group("/projects")

			apiProjects.GET("", api.Project.GetProjects)
			apiProjects.POST("", api.Project.CreateProject)

			{
				apiProjectsPID := apiProjects.Group("/:projectID")

				apiProjectsPID.GET("", api.Project.GetProject)
				apiProjectsPID.PATCH("", api.Project.EditProject)

				apiProjectsPIDMembers := apiProjectsPID.Group("/members")

				apiProjectsPIDMembers.GET("", api.Project.GetProjectMembers)
				apiProjectsPIDMembers.POST("", api.Project.AddProjectMembers)
				apiProjectsPIDMembers.DELETE("", api.Project.DeleteProjectMembers)
			}
		}

		{
			apiEvents := v1.Group("/events")

			apiEvents.GET("", api.Event.GetEvents)

			apiEventsEID := apiEvents.Group("/:eventID")

			apiEventsEID.GET("", api.Event.GetEvent)
			apiEventsEID.PATCH("", api.Event.EditEvent)
		}

		{
			apiGroups := v1.Group("/groups")

			apiGroups.GET("", api.Group.GetGroups)
			apiGroups.GET("/:groupID", api.Group.GetGroup)
		}
		{
			apiContests := v1.Group("/contests")

			apiContests.GET("", api.Contest.GetContests)
			apiContests.POST("", api.Contest.CreateContest)
			{
				apiContestsCID := apiContests.Group("/:contestID")

				apiContestsCID.GET("", api.Contest.GetContest)
				apiContestsCID.PATCH("", api.Contest.EditContest)
				apiContestsCID.DELETE("", api.Contest.DeleteContest)
				{
					apiContestsCIDTeams := apiContestsCID.Group("/teams")

					apiContestsCIDTeams.GET("", api.Contest.GetContestTeams)
					apiContestsCIDTeams.POST("", api.Contest.AddContestTeam)
					{
						apiContestsCIDTeamsTID := apiContestsCIDTeams.Group("/:teamID")

						apiContestsCIDTeamsTID.GET("", api.Contest.GetContestTeam)
						apiContestsCIDTeamsTID.PATCH("", api.Contest.EditContestTeam)
						apiContestsCIDTeamsTID.DELETE("", api.Contest.DeleteContestTeam)
						{
							apiContestsCIDTeamsTIDMembers := apiContestsCIDTeamsTID.Group("/members")

							apiContestsCIDTeamsTIDMembers.GET("", api.Contest.GetContestTeamMembers)
							apiContestsCIDTeamsTIDMembers.POST("", api.Contest.AddContestTeamMember)
							apiContestsCIDTeamsTIDMembers.PUT("", api.Contest.EditContestTeamMember)
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

	return nil
}
