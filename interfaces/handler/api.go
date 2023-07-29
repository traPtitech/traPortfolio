package handler

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type API struct {
	Ping    *PingHandler
	User    *UserHandler
	Project *ProjectHandler
	Event   *EventHandler
	Contest *ContestHandler
	Group   *GroupHandler
}

func NewAPI(ping *PingHandler, user *UserHandler, project *ProjectHandler, event *EventHandler, contest *ContestHandler, group *GroupHandler) API {
	return API{
		Ping:    ping,
		User:    user,
		Project: project,
		Event:   event,
		Contest: contest,
		Group:   group,
	}
}

func setupV1API(g *echo.Group, api API) {
	v1 := g.Group("/v1")

	// ping API
	apiPing := v1.Group("/ping")
	{
		apiPing.GET("", api.Ping.Ping)
	}

	// user API
	userAPI := v1.Group("/users")
	{
		userAPI.GET("", api.User.GetUsers)
		userAPI.GET("/:userID", api.User.GetUser)
		userAPI.PATCH("/:userID", api.User.UpdateUser)
		userAPI.GET("/:userID/accounts", api.User.GetUserAccounts)
		userAPI.POST("/:userID/accounts", api.User.AddUserAccount)
		userAPI.GET("/:userID/accounts/:accountID", api.User.GetUserAccount)
		userAPI.PATCH("/:userID/accounts/:accountID", api.User.EditUserAccount)
		userAPI.DELETE("/:userID/accounts/:accountID", api.User.DeleteUserAccount)
		userAPI.GET("/:userID/projects", api.User.GetUserProjects)
		userAPI.GET("/:userID/contests", api.User.GetUserContests)
		userAPI.GET("/:userID/groups", api.User.GetUserGroups)
		userAPI.GET("/:userID/events", api.User.GetUserEvents)
	}

	// project API
	projectAPI := v1.Group("/projects")
	{
		projectAPI.GET("", api.Project.GetProjects)
		projectAPI.POST("", api.Project.CreateProject)
		projectAPI.GET("/:projectID", api.Project.GetProject)
		projectAPI.PATCH("/:projectID", api.Project.EditProject)
		projectAPI.GET("/:projectID/members", api.Project.GetProjectMembers)
		projectAPI.POST("/:projectID/members", api.Project.AddProjectMembers)
		projectAPI.DELETE("/:projectID/members", api.Project.DeleteProjectMembers)
	}

	// event API
	eventAPI := v1.Group("/events")
	{
		eventAPI.GET("", api.Event.GetEvents)
		eventAPI.GET("/:eventID", api.Event.GetEvent)
		eventAPI.PATCH("/:eventID", api.Event.EditEvent)
	}

	// contest API
	contestAPI := v1.Group("/contests")
	{
		contestAPI.GET("", api.Contest.GetContests)
		contestAPI.POST("", api.Contest.CreateContest)
		contestAPI.GET("/:contestID", api.Contest.GetContest)
		contestAPI.PATCH("/:contestID", api.Contest.EditContest)
		contestAPI.DELETE("/:contestID", api.Contest.DeleteContest)
		contestAPI.GET("/:contestID/teams", api.Contest.GetContestTeams)
		contestAPI.POST("/:contestID/teams", api.Contest.AddContestTeam)
		contestAPI.GET("/:contestID/teams/:teamID", api.Contest.GetContestTeam)
		contestAPI.PATCH("/:contestID/teams/:teamID", api.Contest.EditContestTeam)
		contestAPI.DELETE("/:contestID/teams/:teamID", api.Contest.DeleteContestTeam)
		contestAPI.GET("/:contestID/teams/:teamID/members", api.Contest.GetContestTeamMembers)
		contestAPI.POST("/:contestID/teams/:teamID/members", api.Contest.AddContestTeamMembers)
		contestAPI.PUT("/:contestID/teams/:teamID/members", api.Contest.EditContestTeamMembers)
	}

	// group API
	groupAPI := v1.Group("/groups")
	{
		groupAPI.GET("", api.Group.GetGroups)
		groupAPI.GET("/:groupID", api.Group.GetGroup)
	}
}

type idKey string

const (
	keyUserID        idKey = "userID"
	keyUserAccountID idKey = "accountID"
	keyProject       idKey = "projectID"
	keyEventID       idKey = "eventID"
	keyContestID     idKey = "contestID"
	keyContestTeamID idKey = "teamID"
	keyGroupID       idKey = "groupID"
)

func getID(c echo.Context, key idKey) (uuid.UUID, error) {
	id, err := uuid.FromString(c.Param(string(key)))
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %s", repository.ErrInvalidID, err.Error())
	} else if id.IsNil() {
		return uuid.Nil, repository.ErrNilID
	}

	return id, nil
}
