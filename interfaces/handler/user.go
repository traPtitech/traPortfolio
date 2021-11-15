package handler

import (
	"net/http"

	"github.com/traPtitech/traPortfolio/util/optional"

	"github.com/traPtitech/traPortfolio/usecases/service"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type userParam struct {
	UserID uuid.UUID `param:"userID" validate:"is-uuid"`
}

type accountParams struct {
	UserID    uuid.UUID `param:"userID" validate:"is-uuid"`
	AccountID uuid.UUID `param:"accountID" validate:"is-uuid"`
}

type EditUserRequest struct {
	UserID uuid.UUID       `param:"userID" validate:"is-uuid"`
	Bio    optional.String `json:"bio"`
	Check  optional.Bool   `json:"check"`
}

// UserResponse Portfolioのレスポンスで使うイベント情報
type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	RealName string    `json:"realName"`
}

type userDetailResponse struct {
	ID       uuid.UUID         `json:"id"`
	Name     string            `json:"name"`
	RealName string            `json:"realName"`
	State    domain.TraQState  `json:"state"`
	Bio      string            `json:"bio"`
	Accounts []*domain.Account `json:"accounts"`
}

type AddAccountRequest struct {
	UserID      uuid.UUID `param:"userID" validate:"is-uuid"`
	ID          string    `json:"id"`
	Type        uint      `json:"type"`
	URL         string    `json:"url"`
	PrPermitted bool      `json:"prPermitted"`
}

type EditAccountRequest struct {
	UserID      uuid.UUID       `param:"userID" validate:"is-uuid"`
	AccountID   uuid.UUID       `param:"accountID" validate:"is-uuid"`
	ID          optional.String `json:"id"` // traqID
	Type        optional.Int64  `json:"type"`
	URL         optional.String `json:"url"`
	PrPermitted optional.Bool   `json:"prPermitted"`
}

type UserProjectResponse struct {
	ID           uuid.UUID              `json:"id"`
	Name         string                 `json:"name"`
	Duration     domain.ProjectDuration `json:"duration"`
	UserDuration domain.ProjectDuration `json:"user_duration"`
}

type ContestTeamWithContestNameResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Result      string    `json:"result"`
	ContestName string    `json:"contest_name"`
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

	res := make([]*UserResponse, 0, len(users))
	for _, user := range users {
		res = append(res, &UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			RealName: user.RealName,
		})
	}
	return c.JSON(http.StatusOK, res)
}

// GetByID GET /users/:userID
func (handler *UserHandler) GetByID(_c echo.Context) error {
	c := Context{_c}
	req := userParam{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	user, err := handler.srv.GetUser(ctx, req.UserID)
	if err != nil {
		return convertError(err)
	}

	return c.JSON(http.StatusOK, &userDetailResponse{
		ID:       user.ID,
		Name:     user.Name,
		RealName: user.RealName,
		State:    user.State,
		Bio:      user.Bio,
		Accounts: user.Accounts,
	})
}

// Update PATCH /users/:userID
func (handler *UserHandler) Update(_c echo.Context) error {
	c := Context{_c}
	req := EditUserRequest{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	u := repository.UpdateUserArgs{
		Description: req.Bio,
		Check:       req.Check,
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
	req := userParam{}
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
	req := accountParams{}
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
	req := AddAccountRequest{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	args := repository.CreateAccountArgs{
		ID:          req.ID,
		Type:        req.Type,
		PrPermitted: req.PrPermitted,
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
	req := EditAccountRequest{}
	err := c.BindAndValidate(&req)
	if err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	args := repository.UpdateAccountArgs{
		Name:        req.ID,
		Type:        req.Type,
		URL:         req.URL,
		PrPermitted: req.PrPermitted,
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
	req := accountParams{}
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
	req := userParam{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	projects, err := handler.srv.GetUserProjects(ctx, req.UserID)
	if err != nil {
		return convertError(err)
	}
	res := make([]*UserProjectResponse, 0, len(projects))
	for _, v := range projects {
		up := &UserProjectResponse{
			ID:           v.ID,
			Name:         v.Name,
			Duration:     convertToProjectDuration(v.Since, v.Until),
			UserDuration: convertToProjectDuration(v.UserSince, v.UserUntil),
		}
		res = append(res, up)
	}
	return c.JSON(http.StatusOK, res)
}

// GetContests GET /users/:userID/contests
func (handler *UserHandler) GetContests(_c echo.Context) error {
	c := Context{_c}
	req := userParam{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	contests, err := handler.srv.GetUserContests(ctx, req.UserID)
	if err != nil {
		return convertError(err)
	}
	res := make([]*ContestTeamWithContestNameResponse, 0, len(contests))
	for _, v := range contests {
		uc := &ContestTeamWithContestNameResponse{
			ID:          v.ID,
			Name:        v.Name,
			Result:      v.Result,
			ContestName: v.ContestName,
		}
		res = append(res, uc)
	}
	return c.JSON(http.StatusOK, res)
}

// GetEvents GET /users/:userID/events
func (handler *UserHandler) GetEvents(_c echo.Context) error {
	c := Context{_c}
	req := userParam{}
	if err := c.BindAndValidate(&req); err != nil {
		return convertError(err)
	}

	ctx := c.Request().Context()
	events, err := handler.srv.GetUserEvents(ctx, req.UserID)
	if err != nil {
		return convertError(err)
	}
	res := make([]*EventResponse, 0, len(events))
	for _, v := range events {
		e := &EventResponse{
			ID:   v.ID,
			Name: v.Name,
			Duration: Duration{
				Since: v.TimeStart,
				Until: v.TimeEnd,
			},
		}
		res = append(res, e)
	}
	return c.JSON(http.StatusOK, res)
}
