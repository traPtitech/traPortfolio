package handler

import (
	"errors"
	"fmt"
	"net/http"

	vd "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

func Setup(e *echo.Echo, api API) error {
	e.HTTPErrorHandler = newHTTPErrorHandler(e)
	e.Binder = &binderWithValidation{}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	apiGroup := e.Group("/api")
	setupV1API(apiGroup, api)

	return nil
}

func newHTTPErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		var (
			code int
			herr *echo.HTTPError
		)

		switch {
		case errors.Is(err, repository.ErrNilID):
			fallthrough
		case errors.Is(err, repository.ErrInvalidID):
			fallthrough
		case errors.Is(err, repository.ErrInvalidArg):
			code = http.StatusBadRequest

		case errors.Is(err, repository.ErrAlreadyExists):
			code = http.StatusConflict

		case errors.Is(err, repository.ErrForbidden):
			code = http.StatusForbidden

		case errors.Is(err, repository.ErrNotFound):
			code = http.StatusNotFound

		default:
			e.Logger.Error(err)
			code = http.StatusInternalServerError
			herr = echo.NewHTTPError(code, http.StatusText(code)).SetInternal(err)
		}

		if herr == nil {
			herr = echo.NewHTTPError(
				code,
				fmt.Sprintf("%s: %s", http.StatusText(code), err.Error()),
			).SetInternal(err)
		}

		e.DefaultHTTPErrorHandler(herr, c)
	}
}

type binderWithValidation struct{}

var _ echo.Binder = (*binderWithValidation)(nil)

func (b *binderWithValidation) Bind(i interface{}, c echo.Context) error {
	if err := (&echo.DefaultBinder{}).Bind(i, c); err != nil {
		return err
	}

	if vld, ok := i.(vd.Validatable); ok {
		if err := vld.Validate(); err != nil {
			if ie, ok := err.(vd.InternalError); ok {
				c.Logger().Fatalf("ozzo-validation internal error: %s", ie.Error())
			}

			return err
		}
	} else {
		c.Logger().Errorf("%T is not validatable", i)
	}

	return nil
}

func setupV1API(g *echo.Group, api API) {
	v1 := g.Group("/v1")
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
	// group API
	groupAPI := v1.Group("/groups")
	{
		groupAPI.GET("", api.Group.GetGroups)
		groupAPI.GET("/:groupID", api.Group.GetGroup)
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
	// ping API
	apiPing := v1.Group("/ping")
	{
		apiPing.GET("", api.Ping.Ping)
	}
}
