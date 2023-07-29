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

type UserHandler struct {
	srv service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{srv: s}
}

// GetUsers GET /users
func (h *UserHandler) GetUsers(_c echo.Context) error {
	c := _c.(*Context)
	req := GetUsersParams{}
	if err := c.BindAndValidate(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	args := repository.GetUsersArgs{
		IncludeSuspended: optional.FromPtr((*bool)(req.IncludeSuspended)),
		Name:             optional.FromPtr((*string)(req.Name)),
		Limit:            optional.FromPtr((*int)(req.Limit)),
	}

	users, err := h.srv.GetUsers(ctx, &args)
	if err != nil {
		return err
	}

	res := make([]User, len(users))
	for i, v := range users {
		res[i] = newUser(v.ID, v.Name, v.RealName())
	}

	return c.JSON(http.StatusOK, res)
}

// GetUser GET /users/:userID
func (h *UserHandler) GetUser(_c echo.Context) error {
	c := _c.(*Context)

	userID, err := c.getID(keyUserID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	user, err := h.srv.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	accounts := make([]Account, len(user.Accounts))
	for i, v := range user.Accounts {
		accounts[i] = newAccount(v.ID, v.DisplayName, AccountType(v.Type), v.URL, v.PrPermitted)
	}

	return c.JSON(http.StatusOK, newUserDetail(
		newUser(user.ID, user.Name, user.RealName()),
		accounts,
		user.Bio,
		user.State,
	))
}

// UpdateUser PATCH /users/:userID
func (h *UserHandler) UpdateUser(_c echo.Context) error {
	c := _c.(*Context)

	userID, err := c.getID(keyUserID)
	if err != nil {
		return err
	}

	req := EditUserJSONRequestBody{}
	if err := c.BindAndValidate(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	u := repository.UpdateUserArgs{
		Description: optional.FromPtr(req.Bio),
		Check:       optional.FromPtr(req.Check),
	}

	if err := h.srv.Update(ctx, userID, &u); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// GetUserAccounts GET /users/:userID/accounts
func (h *UserHandler) GetUserAccounts(_c echo.Context) error {
	c := _c.(*Context)

	userID, err := c.getID(keyUserID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	accounts, err := h.srv.GetAccounts(ctx, userID)
	if err != nil {
		return err
	}

	res := make([]Account, len(accounts))
	for i, v := range accounts {
		res[i] = newAccount(v.ID, v.DisplayName, AccountType(v.Type), v.URL, v.PrPermitted)
	}

	return c.JSON(http.StatusOK, res)
}

// GetUserAccount GET /users/:userID/accounts/:accountID
func (h *UserHandler) GetUserAccount(_c echo.Context) error {
	c := _c.(*Context)

	userID, err := c.getID(keyUserID)
	if err != nil {
		return err
	}

	accountID, err := c.getID(keyUserAccountID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	account, err := h.srv.GetAccount(ctx, userID, accountID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, newAccount(account.ID, account.DisplayName, AccountType(account.Type), account.URL, account.PrPermitted))
}

// AddUserAccount POST /users/:userID/accounts
func (h *UserHandler) AddUserAccount(_c echo.Context) error {
	c := _c.(*Context)

	userID, err := c.getID(keyUserID)
	if err != nil {
		return err
	}

	req := AddUserAccountJSONRequestBody{}
	if err := c.BindAndValidate(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	args := repository.CreateAccountArgs{
		DisplayName: req.DisplayName,
		Type:        domain.AccountType(req.Type),
		PrPermitted: bool(req.PrPermitted),
		URL:         req.Url,
	}
	account, err := h.srv.CreateAccount(ctx, userID, &args)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, newAccount(account.ID, account.DisplayName, AccountType(account.Type), account.URL, account.PrPermitted))
}

// EditUserAccount PATCH /users/:userID/accounts/:accountID
func (h *UserHandler) EditUserAccount(_c echo.Context) error {
	c := _c.(*Context)

	userID, err := c.getID(keyUserID)
	if err != nil {
		return err
	}

	accountID, err := c.getID(keyUserAccountID)
	if err != nil {
		return err
	}

	req := EditUserAccountJSONRequestBody{}
	if err := c.BindAndValidate(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()

	args := repository.UpdateAccountArgs{
		DisplayName: optional.FromPtr(req.DisplayName),
		Type:        optional.FromPtr((*domain.AccountType)(req.Type)),
		URL:         optional.FromPtr(req.Url),
		PrPermitted: optional.FromPtr(req.PrPermitted),
	}

	err = h.srv.EditAccount(ctx, userID, accountID, &args)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteUserAccount DELETE /users/:userID/accounts/:accountID
func (h *UserHandler) DeleteUserAccount(_c echo.Context) error {
	c := _c.(*Context)

	userID, err := c.getID(keyUserID)
	if err != nil {
		return err
	}

	accountID, err := c.getID(keyUserAccountID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	if err := h.srv.DeleteAccount(ctx, userID, accountID); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// GetUserProjects GET /users/:userID/projects
func (h *UserHandler) GetUserProjects(_c echo.Context) error {
	c := _c.(*Context)

	userID, err := c.getID(keyUserID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	projects, err := h.srv.GetUserProjects(ctx, userID)
	if err != nil {
		return err
	}
	res := make([]UserProject, len(projects))
	for i, v := range projects {
		res[i] = newUserProject(
			v.ID,
			v.Name,
			ConvertDuration(v.Duration),
			ConvertDuration(v.UserDuration),
		)
	}

	return c.JSON(http.StatusOK, res)
}

// GetUserContests GET /users/:userID/contests
func (h *UserHandler) GetUserContests(_c echo.Context) error {
	c := _c.(*Context)

	userID, err := c.getID(keyUserID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	contests, err := h.srv.GetUserContests(ctx, userID)
	if err != nil {
		return err
	}

	res := make([]UserContest, len(contests))
	for i, c := range contests {
		teams := make([]ContestTeam, len(c.Teams))
		for j, ct := range c.Teams {
			teams[j] = newContestTeam(ct.ID, ct.Name, ct.Result)
		}
		res[i] = newUserContest(
			newContest(c.ID, c.Name, c.TimeStart, c.TimeEnd),
			teams,
		)
	}

	return c.JSON(http.StatusOK, res)
}

// GetUserGroups GET /users/:userID/groups
func (h *UserHandler) GetUserGroups(_c echo.Context) error {
	c := _c.(*Context)

	userID, err := c.getID(keyUserID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	groups, err := h.srv.GetGroupsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	res := make([]UserGroup, len(groups))
	for i, group := range groups {
		res[i] = newUserGroup(
			newGroup(group.ID, group.Name),
			ConvertDuration(group.Duration),
		)
	}
	return c.JSON(http.StatusOK, res)
}

// GetUserEvents GET /users/:userID/events
func (h *UserHandler) GetUserEvents(_c echo.Context) error {
	c := _c.(*Context)

	userID, err := c.getID(keyUserID)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	events, err := h.srv.GetUserEvents(ctx, userID)
	if err != nil {
		return err
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
		Accounts: accounts,
		Bio:      bio,
		Id:       user.Id,
		Name:     user.Name,
		RealName: user.RealName,
		State:    UserAccountState(state),
	}
}

func newAccount(id uuid.UUID, displayName string, atype AccountType, url string, prPermitted bool) Account {
	return Account{
		Id:          id,
		DisplayName: displayName,
		Type:        atype,
		Url:         url,
		PrPermitted: PrPermitted(prPermitted),
	}
}

func newUserProject(id uuid.UUID, name string, duration YearWithSemesterDuration, userDuration YearWithSemesterDuration) UserProject {
	return UserProject{
		Duration:     duration,
		Id:           id,
		Name:         name,
		UserDuration: userDuration,
	}
}

func newUserContest(contest Contest, teams []ContestTeam) UserContest {
	return UserContest{
		Id:       contest.Id,
		Name:     contest.Name,
		Duration: contest.Duration,
		Teams:    teams,
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
		Duration: Duration,
		Id:       group.Id,
		Name:     group.Name,
	}
}
