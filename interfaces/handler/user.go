package handler

import (
	"net/http"

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

type UserHandler struct {
	srv service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{srv: s}
}

// GetAll GET /users
func (handler *UserHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	users, err := handler.srv.GetUsers(ctx)
	if err != nil {
		return convertError(err)
	}

	res := make([]*User, 0, len(users))
	for _, user := range users {
		res = append(res, &User{
			Id:       user.ID,
			Name:     user.Name,
			RealName: &user.RealName,
		})
	}
	return c.JSON(http.StatusOK, res)
}

// GetByID GET /users/:userID
func (handler *UserHandler) GetByID(_c echo.Context) error {
	c := Context{_c}
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
		accounts[i] = Account{
			Id:          v.ID,
			Name:        v.Name,
			Type:        AccountType(v.Type),
			Url:         v.URL,
			PrPermitted: PrPermitted(v.PrPermitted),
		}
	}

	return c.JSON(http.StatusOK, &UserDetail{
		User: User{
			Id:       user.ID,
			Name:     user.Name,
			RealName: &user.RealName,
		},
		State:    UserAccountState(user.State),
		Bio:      user.Bio,
		Accounts: accounts,
	})
}

// Update PATCH /users/:userID
func (handler *UserHandler) Update(_c echo.Context) error {
	c := Context{_c}
	req := struct {
		UserIDInPath
		EditUserJSONRequestBody
	}{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	u := repository.UpdateUserArgs{
		Description: optional.StringFrom(*req.Bio), //TODO: valid: falseを追加する
		Check:       optional.BoolFrom(*req.Check),
	}
	err := handler.srv.Update(ctx, req.UserID, &u)
	if err != nil {
		return convertError(err)
	}
	return c.NoContent(http.StatusOK)
}

// GetAccounts GET /users/:userID/accounts
func (handler *UserHandler) GetAccounts(_c echo.Context) error {
	c := Context{_c}
	req := UserIDInPath{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	accounts, err := handler.srv.GetAccounts(req.UserID)
	if err != nil {
		return convertError(err)
	}

	return c.JSON(http.StatusOK, accounts)
}

// GetAccount GET /users/:userID/accounts/:accountID
func (handler *UserHandler) GetAccount(_c echo.Context) error {
	c := Context{_c}
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

	return c.JSON(http.StatusOK, account)
}

// AddAccount POST /users/:userID/accounts
func (handler *UserHandler) AddAccount(_c echo.Context) error {
	c := Context{_c}
	req := struct {
		UserIDInPath
		AddAccountJSONRequestBody
	}{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	args := repository.CreateAccountArgs{
		ID:          *req.Id,
		Type:        uint(*req.Type),
		PrPermitted: bool(*req.PrPermitted),
		URL:         *req.Url,
	}
	account, err := handler.srv.CreateAccount(ctx, req.UserID, &args)
	if err != nil {
		return convertError(err)
	}
	return c.JSON(http.StatusOK, account)
}

// PatchAccount PATCH /users/:userID/accounts/:accountID
func (handler *UserHandler) PatchAccount(_c echo.Context) error {
	c := Context{_c}
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
		Name:        optional.StringFrom(*req.Id), // TODO
		Type:        optional.Int64From(int64(*req.Type)),
		URL:         optional.StringFrom(*req.Url),
		PrPermitted: optional.BoolFrom(bool(*req.PrPermitted)),
	}
	err = handler.srv.EditAccount(ctx, req.AccountID, req.UserID, &args)
	if err != nil {
		return convertError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteAccount DELETE /users/:userID/accounts/:accountID
func (handler *UserHandler) DeleteAccount(_c echo.Context) error {
	c := Context{_c}
	req := struct {
		UserIDInPath
		AccountIDInPath
	}{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	err := handler.srv.DeleteAccount(ctx, req.AccountID, req.UserID)
	if err != nil {
		return convertError(err)
	}
	return c.NoContent(http.StatusOK)
}

// GetProjects GET /users/:userID/projects
func (handler *UserHandler) GetProjects(_c echo.Context) error {
	c := Context{_c}
	req := UserIDInPath{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	projects, err := handler.srv.GetUserProjects(ctx, req.UserID)
	if err != nil {
		return convertError(err)
	}
	res := make([]*UserProject, 0, len(projects))
	for _, v := range projects {
		up := &UserProject{
			Project: Project{
				Id:       v.ID,
				Name:     v.Name,
				Duration: convertToProjectDuration(v.Since, v.Until),
			},
			UserDuration: []ProjectDuration{convertToProjectDuration(v.UserSince, v.UserUntil)}, //TODO: objectでいいはず
		}
		res = append(res, up)
	}
	return c.JSON(http.StatusOK, res)
}

// GetContests GET /users/:userID/contests
func (handler *UserHandler) GetContests(_c echo.Context) error {
	c := Context{_c}
	req := UserIDInPath{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	contests, err := handler.srv.GetUserContests(ctx, req.UserID)
	if err != nil {
		return convertError(err)
	}
	res := make([]*ContestTeamWithContestName, 0, len(contests))
	for _, v := range contests {
		uc := &ContestTeamWithContestName{
			ContestTeam: ContestTeam{
				Id:     v.ID,
				Name:   v.Name,
				Result: &v.Result,
			},
			ContestName: v.ContestName,
		}
		res = append(res, uc)
	}
	return c.JSON(http.StatusOK, res)
}

// GetEvents GET /users/:userID/events
func (handler *UserHandler) GetEvents(_c echo.Context) error {
	c := Context{_c}
	req := UserIDInPath{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	events, err := handler.srv.GetUserEvents(ctx, req.UserID)
	if err != nil {
		return convertError(err)
	}
	res := make([]*Event, 0, len(events))
	for _, v := range events {
		e := &Event{
			Id:   v.ID,
			Name: v.Name,
			Duration: Duration{
				Since: v.TimeStart,
				Until: &v.TimeEnd,
			},
		}
		res = append(res, e)
	}
	return c.JSON(http.StatusOK, res)
}
