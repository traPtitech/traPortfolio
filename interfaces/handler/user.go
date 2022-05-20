package handler

import (
	"net/http"

	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/util/optional"

	"github.com/traPtitech/traPortfolio/usecases/service"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type UserIDInPath struct {
	UserID uuid.UUID `param:"userID" validate:"is-uuid"`
}

type AccountIDInPath struct {
	AccountID uuid.UUID `param:"accountID" validate:"is-uuid"`
}

type GroupIDInPath struct {
	GroupID uuid.UUID `param:"groupID" validate:"is-uuid"`
}

type UserHandler struct {
	srv service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{srv: s}
}

// GetUsers GET /users
func (handler *UserHandler) GetUsers(_c echo.Context) error {
	c := _c.(*Context)
	req := GetUsersParams{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	args := repository.GetUsersArgs{
		IncludeSuspended: optional.BoolFrom((*bool)(req.IncludeSuspended)),
		Name:             optional.StringFrom((*string)(req.Name)),
	}

	users, err := handler.srv.GetUsers(ctx, &args)
	if err != nil {
		return convertError(err)
	}

	res := make([]User, len(users))
	for i, v := range users {
		res[i] = newUser(v.ID, v.Name, v.RealName)
	}

	return c.JSON(http.StatusOK, res)
}

// GetUser GET /users/:userID
func (handler *UserHandler) GetUser(_c echo.Context) error {
	c := _c.(*Context)
	req := UserIDInPath{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	user, err := handler.srv.GetUser(ctx, req.UserID)
	if err != nil {
		return convertError(err)
	}

	accounts := make([]Account, len(user.Accounts))
	for i, v := range user.Accounts {
		accounts[i] = newAccount(v.ID, v.DisplayName, v.Type, v.URL, v.PrPermitted)
	}

	return c.JSON(http.StatusOK, newUserDetail(
		newUser(user.ID, user.Name, user.RealName),
		accounts,
		user.Bio,
		user.State,
	))
}

// UpdateUser PATCH /users/:userID
func (handler *UserHandler) UpdateUser(_c echo.Context) error {
	c := _c.(*Context)
	req := struct {
		UserIDInPath
		EditUserJSONRequestBody
	}{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	u := repository.UpdateUserArgs{
		Description: optional.StringFrom(req.Bio),
		Check:       optional.BoolFrom(req.Check),
	}

	err := handler.srv.Update(ctx, req.UserID, &u)
	if err != nil {
		return convertError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

// GetUserAccounts GET /users/:userID/accounts
func (handler *UserHandler) GetUserAccounts(_c echo.Context) error {
	c := _c.(*Context)
	req := UserIDInPath{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	accounts, err := handler.srv.GetAccounts(req.UserID)
	if err != nil {
		return convertError(err)
	}

	res := make([]Account, len(accounts))
	for i, v := range accounts {
		res[i] = newAccount(v.ID, v.DisplayName, v.Type, v.URL, v.PrPermitted)
	}

	return c.JSON(http.StatusOK, res)
}

// GetUserAccount GET /users/:userID/accounts/:accountID
func (handler *UserHandler) GetUserAccount(_c echo.Context) error {
	c := _c.(*Context)
	req := struct {
		UserIDInPath
		AccountIDInPath
	}{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	account, err := handler.srv.GetAccount(req.UserID, req.AccountID)
	if err != nil {
		return convertError(err)
	}

	return c.JSON(http.StatusOK, newAccount(account.ID, account.DisplayName, account.Type, account.URL, account.PrPermitted))
}

// AddUserAccount POST /users/:userID/accounts
func (handler *UserHandler) AddUserAccount(_c echo.Context) error {
	c := _c.(*Context)
	req := struct {
		UserIDInPath
		AddUserAccountJSONRequestBody
	}{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	args := repository.CreateAccountArgs{
		DisplayName: req.DisplayName,
		Type:        uint(req.Type),
		PrPermitted: bool(req.PrPermitted),
		URL:         req.Url,
	}
	account, err := handler.srv.CreateAccount(ctx, req.UserID, &args)
	if err != nil {
		return convertError(err)
	}

	return c.JSON(http.StatusCreated, newAccount(account.ID, account.DisplayName, account.Type, account.URL, account.PrPermitted))
}

// EditUserAccount PATCH /users/:userID/accounts/:accountID
func (handler *UserHandler) EditUserAccount(_c echo.Context) error {
	c := _c.(*Context)
	req := struct {
		UserIDInPath
		AccountIDInPath
		EditUserAccountJSONRequestBody
	}{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	args := repository.UpdateAccountArgs{
		DisplayName: optional.StringFrom(req.DisplayName),
		Type:        optional.Int64From(((*int64)(req.Type))),
		URL:         optional.StringFrom(req.Url),
		PrPermitted: optional.BoolFrom((*bool)(req.PrPermitted)),
	}

	err = handler.srv.EditAccount(ctx, req.UserID, req.AccountID, &args)
	if err != nil {
		return convertError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteUserAccount DELETE /users/:userID/accounts/:accountID
func (handler *UserHandler) DeleteUserAccount(_c echo.Context) error {
	c := _c.(*Context)
	req := struct {
		UserIDInPath
		AccountIDInPath
	}{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	err := handler.srv.DeleteAccount(ctx, req.UserID, req.AccountID)
	if err != nil {
		return convertError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetUserProjects GET /users/:userID/projects
func (handler *UserHandler) GetUserProjects(_c echo.Context) error {
	c := _c.(*Context)
	req := UserIDInPath{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	projects, err := handler.srv.GetUserProjects(ctx, req.UserID)
	if err != nil {
		return convertError(err)
	}
	res := make([]UserProject, len(projects))
	for i, v := range projects {
		res[i] = newUserProject(
			v.ID,
			v.Name,
			convertDuration(v.Duration),
			convertDuration(v.UserDuration),
		)
	}

	return c.JSON(http.StatusOK, res)
}

// GetUserContests GET /users/:userID/contests
func (handler *UserHandler) GetUserContests(_c echo.Context) error {
	c := _c.(*Context)
	req := UserIDInPath{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	contests, err := handler.srv.GetUserContests(ctx, req.UserID)
	if err != nil {
		return convertError(err)
	}

	res := make([]ContestTeamWithContestName, len(contests))
	for i, v := range contests {
		res[i] = newContestTeamWithContestName(
			newContestTeam(v.ID, v.Name, v.Result),
			v.ContestName,
		)
	}

	return c.JSON(http.StatusOK, res)
}

// GetUserGroups GET /users/:userID/groups
func (handler *UserHandler) GetUserGroups(_c echo.Context) error {
	c := _c.(*Context)
	req := UserIDInPath{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	groups, err := handler.srv.GetGroupsByUserID(ctx, req.UserID)
	if err != nil {
		return convertError(err)
	}

	res := make([]UserGroup, len(groups))
	for i, group := range groups {
		res[i] = newUserGroup(
			newGroup(group.ID, group.Name),
			convertDuration(group.Duration),
		)
	}
	return c.JSON(http.StatusOK, res)
}

// GetUserEvents GET /users/:userID/events
func (handler *UserHandler) GetUserEvents(_c echo.Context) error {
	c := _c.(*Context)
	req := UserIDInPath{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	events, err := handler.srv.GetUserEvents(ctx, req.UserID)
	if err != nil {
		return convertError(err)
	}

	res := make([]Event, len(events))
	for i, v := range events {
		res[i] = newEvent(v.ID, v.Name, v.TimeStart, v.TimeEnd)
	}

	return c.JSON(http.StatusOK, res)
}

func newUser(id uuid.UUID, name string, realName string) User {
	return User{
		Id:       id,
		Name:     name,
		RealName: realName,
	}
}

func newUserDetail(user User, accounts []Account, bio string, state domain.TraQState) UserDetail {
	return UserDetail{
		User:     user,
		Accounts: accounts,
		Bio:      bio,
		State:    UserAccountState(state),
	}
}

func newAccount(id uuid.UUID, displayName string, atype uint, url string, prPermitted bool) Account {
	return Account{
		Id:          id,
		DisplayName: displayName,
		Type:        AccountType(atype),
		Url:         url,
		PrPermitted: PrPermitted(prPermitted),
	}
}

func newUserProject(id uuid.UUID, name string, duration YearWithSemesterDuration, userDuration YearWithSemesterDuration) UserProject {
	return UserProject{
		Project: Project{
			Id:       id,
			Name:     name,
			Duration: duration,
		},
		UserDuration: userDuration,
	}
}

// TODO: UserContestのほうがいいかも
func newContestTeamWithContestName(contestTeam ContestTeam, contestName string) ContestTeamWithContestName {
	return ContestTeamWithContestName{
		ContestTeam: contestTeam,
		ContestName: contestName,
	}
}

func newGroup(id uuid.UUID, name string) Group {
	return Group{
		Id:   id,
		Name: name,
	}
}

func newUserGroup(group Group, Duration YearWithSemesterDuration) UserGroup {
	return UserGroup{
		Group:    group,
		Duration: Duration,
	}
}
