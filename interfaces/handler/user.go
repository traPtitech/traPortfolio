package handler

import (
	"fmt"
	"net/http"

	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/handler/schema"
	"github.com/traPtitech/traPortfolio/util/optional"

	"github.com/traPtitech/traPortfolio/usecases/service"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type UserHandler struct {
	s service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{s}
}

// GetUsers GET /users
func (h *UserHandler) GetUsers(c echo.Context) error {
	req := schema.GetUsersParams{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	if req.Me != nil && *req.Me {
		if req.Name != nil {
			return fmt.Errorf("%w: cannot specify both 'me' and 'name'", repository.ErrInvalidArg)
		}

		name, ok := getMyName(c)
		if !ok {
			return fmt.Errorf("%w: user name not found", repository.ErrForbidden)
		}

		req.Name = &name
	}

	ctx := c.Request().Context()
	args := repository.GetUsersArgs{
		IncludeSuspended: optional.FromPtr((*bool)(req.IncludeSuspended)),
		Name:             optional.FromPtr((*string)(req.Name)),
		Limit:            optional.FromPtr((*int)(req.Limit)),
	}

	users, err := h.s.GetUsers(ctx, &args)
	if err != nil {
		return err
	}

	res := make([]schema.User, len(users))
	for i, v := range users {
		res[i] = newUser(v.ID, v.Name, v.RealName())
	}

	return c.JSON(http.StatusOK, res)
}

// SyncUsers POST /users/sync
func (h *UserHandler) SyncUsers(c echo.Context) error {
	ctx := c.Request().Context()
	if err := h.s.SyncUsers(ctx); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// GetUser GET /users/:userID
func (h *UserHandler) GetUser(c echo.Context) error {
	userID, err := getID(c, keyUserID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	user, err := h.s.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	accounts := make([]schema.Account, len(user.Accounts))
	for i, v := range user.Accounts {
		accounts[i] = newAccount(v.ID, v.DisplayName, schema.AccountType(v.Type), v.URL, v.PrPermitted)
	}

	return c.JSON(http.StatusOK, newUserDetail(
		newUser(user.ID, user.Name, user.RealName()),
		accounts,
		user.Bio,
		user.State,
	))
}

// UpdateUser PATCH /users/:userID
func (h *UserHandler) UpdateUser(c echo.Context) error {
	userID, err := getID(c, keyUserID)
	if err != nil {
		return err
	}

	req := schema.EditUserRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	u := repository.UpdateUserArgs{
		Description: optional.FromPtr(req.Bio),
		Check:       optional.FromPtr(req.Check),
	}

	if err := h.s.Update(ctx, userID, &u); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// GetUserAccounts GET /users/:userID/accounts
func (h *UserHandler) GetUserAccounts(c echo.Context) error {
	userID, err := getID(c, keyUserID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	accounts, err := h.s.GetAccounts(ctx, userID)
	if err != nil {
		return err
	}

	res := make([]schema.Account, len(accounts))
	for i, v := range accounts {
		res[i] = newAccount(v.ID, v.DisplayName, schema.AccountType(v.Type), v.URL, v.PrPermitted)
	}

	return c.JSON(http.StatusOK, res)
}

// GetUserAccount GET /users/:userID/accounts/:accountID
func (h *UserHandler) GetUserAccount(c echo.Context) error {
	userID, err := getID(c, keyUserID)
	if err != nil {
		return err
	}

	accountID, err := getID(c, keyUserAccountID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	account, err := h.s.GetAccount(ctx, userID, accountID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, newAccount(account.ID, account.DisplayName, schema.AccountType(account.Type), account.URL, account.PrPermitted))
}

// AddUserAccount POST /users/:userID/accounts
func (h *UserHandler) AddUserAccount(c echo.Context) error {
	userID, err := getID(c, keyUserID)
	if err != nil {
		return err
	}

	req := schema.AddAccountRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	args := repository.CreateAccountArgs{
		DisplayName: req.DisplayName,
		Type:        domain.AccountType(req.Type),
		PrPermitted: bool(req.PrPermitted),
		URL:         req.Url,
	}
	account, err := h.s.CreateAccount(ctx, userID, &args)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, newAccount(account.ID, account.DisplayName, schema.AccountType(account.Type), account.URL, account.PrPermitted))
}

// EditUserAccount PATCH /users/:userID/accounts/:accountID
func (h *UserHandler) EditUserAccount(c echo.Context) error {
	userID, err := getID(c, keyUserID)
	if err != nil {
		return err
	}

	accountID, err := getID(c, keyUserAccountID)
	if err != nil {
		return err
	}

	req := schema.EditUserAccountRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()

	args := repository.UpdateAccountArgs{
		DisplayName: optional.FromPtr(req.DisplayName),
		Type:        optional.FromPtr((*domain.AccountType)(req.Type)),
		URL:         optional.FromPtr(req.Url),
		PrPermitted: optional.FromPtr(req.PrPermitted),
	}

	err = h.s.EditAccount(ctx, userID, accountID, &args)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteUserAccount DELETE /users/:userID/accounts/:accountID
func (h *UserHandler) DeleteUserAccount(c echo.Context) error {
	userID, err := getID(c, keyUserID)
	if err != nil {
		return err
	}

	accountID, err := getID(c, keyUserAccountID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	if err := h.s.DeleteAccount(ctx, userID, accountID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// GetUserProjects GET /users/:userID/projects
func (h *UserHandler) GetUserProjects(c echo.Context) error {
	userID, err := getID(c, keyUserID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	projects, err := h.s.GetUserProjects(ctx, userID)
	if err != nil {
		return err
	}
	res := make([]schema.UserProject, len(projects))
	for i, v := range projects {
		res[i] = newUserProject(
			v.ID,
			v.Name,
			schema.ConvertDuration(v.Duration),
			schema.ConvertDuration(v.UserDuration),
		)
	}

	return c.JSON(http.StatusOK, res)
}

// GetUserContests GET /users/:userID/contests
func (h *UserHandler) GetUserContests(c echo.Context) error {
	userID, err := getID(c, keyUserID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	contests, err := h.s.GetUserContests(ctx, userID)
	if err != nil {
		return err
	}

	res := make([]schema.UserContest, len(contests))
	for i, c := range contests {
		teams := make([]schema.ContestTeamWithoutMembers, len(c.Teams))
		for j, ct := range c.Teams {
			teams[j] = newContestTeamWithoutMembers(ct.ID, ct.Name, ct.Result)
		}
		res[i] = newUserContest(
			newContest(c.ID, c.Name, c.TimeStart, c.TimeEnd),
			teams,
		)
	}

	return c.JSON(http.StatusOK, res)
}

// GetUserGroups GET /users/:userID/groups
func (h *UserHandler) GetUserGroups(c echo.Context) error {
	userID, err := getID(c, keyUserID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	groups, err := h.s.GetGroupsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	res := make([]schema.UserGroup, len(groups))
	for i, group := range groups {
		res[i] = newUserGroup(
			newGroup(group.ID, group.Name),
			schema.ConvertDuration(group.Duration),
		)
	}
	return c.JSON(http.StatusOK, res)
}

// GetUserEvents GET /users/:userID/events
func (h *UserHandler) GetUserEvents(c echo.Context) error {
	userID, err := getID(c, keyUserID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	events, err := h.s.GetUserEvents(ctx, userID)
	if err != nil {
		return err
	}

	res := make([]schema.Event, len(events))
	for i, v := range events {
		res[i] = newEvent(v.ID, v.Name, v.Level, v.TimeStart, v.TimeEnd)
	}

	return c.JSON(http.StatusOK, res)
}

func newUser(id uuid.UUID, name string, realName string) schema.User {
	return schema.User{
		Id:       id,
		Name:     name,
		RealName: realName,
	}
}

func newUserDetail(user schema.User, accounts []schema.Account, bio string, state domain.TraQState) schema.UserDetail {
	return schema.UserDetail{
		Accounts: accounts,
		Bio:      bio,
		Id:       user.Id,
		Name:     user.Name,
		RealName: user.RealName,
		State:    schema.UserAccountState(state),
	}
}

func newAccount(id uuid.UUID, displayName string, atype schema.AccountType, url string, prPermitted bool) schema.Account {
	return schema.Account{
		Id:          id,
		DisplayName: displayName,
		Type:        atype,
		Url:         url,
		PrPermitted: schema.PrPermitted(prPermitted),
	}
}

func newUserProject(id uuid.UUID, name string, duration schema.YearWithSemesterDuration, userDuration schema.YearWithSemesterDuration) schema.UserProject {
	return schema.UserProject{
		Duration:     duration,
		Id:           id,
		Name:         name,
		UserDuration: userDuration,
	}
}

func newUserContest(contest schema.Contest, teams []schema.ContestTeamWithoutMembers) schema.UserContest {
	return schema.UserContest{
		Id:       contest.Id,
		Name:     contest.Name,
		Duration: contest.Duration,
		Teams:    teams,
	}
}

func newGroup(id uuid.UUID, name string) schema.Group {
	return schema.Group{
		Id:   id,
		Name: name,
	}
}

func newUserGroup(group schema.Group, Duration schema.YearWithSemesterDuration) schema.UserGroup {
	return schema.UserGroup{
		Duration: Duration,
		Id:       group.Id,
		Name:     group.Name,
	}
}
