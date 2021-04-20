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

type EditUserRequest struct {
	Bio   optional.String `json:"bio"`
	Check optional.Bool   `json:"check"`
}

type UserHandler struct {
	srv service.UserService
}

// userResponse Portfolioのレスポンスで使うイベント情報
type userResponse struct {
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

type Account struct {
	ID          string `json:"id"`
	Type        uint   `json:"type"`
	URL         string `json:"url"`
	PrPermitted bool   `json:"prPermitted"`
}

type EditAccountRequest struct {
	ID          optional.String `json:"id"` // traqID
	Type        optional.Int64  `json:"type"`
	URL         optional.String `json:"url"`
	PrPermitted optional.Bool   `json:"prPermitted"`
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{srv: s}
}

// GetAll GET /users
func (handler *UserHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	users, err := handler.srv.GetUsers(ctx)
	if err != nil {
		return err
	}

	res := make([]*userResponse, 0, len(users))
	for _, user := range users {
		res = append(res, &userResponse{
			ID:       user.ID,
			Name:     user.Name,
			RealName: user.RealName,
		})
	}
	return c.JSON(http.StatusOK, res)
}

// GetByID GET /users/:userID
func (handler *UserHandler) GetByID(c echo.Context) error {
	_id := c.Param("userID")
	if _id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user id must not be blank")
	}

	id := uuid.FromStringOrNil(_id)
	if id == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid uuid")
	}
	ctx := c.Request().Context()
	user, err := handler.srv.GetUser(ctx, id)
	if err == repository.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if err != nil {
		return err
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
func (handler *UserHandler) Update(c echo.Context) error {
	_id := c.Param("userID")
	if _id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user id must not be blank")
	}

	id := uuid.FromStringOrNil(_id)
	if id == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid uuid")
	}
	req := EditUserRequest{}
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	ctx := c.Request().Context()
	u := repository.UpdateUserArgs{
		Description: req.Bio,
		Check:       req.Check,
	}
	err = handler.srv.Update(ctx, id, &u)
	if err == repository.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

// GetAccounts GET /users/:userID/accounts
func (handler *UserHandler) GetAccounts(_c echo.Context) error {
	c := Context{_c}
	_ = c.Request().Context()
	_id := c.Param("userID")
	id := uuid.FromStringOrNil(_id)

	accounts, err := handler.srv.GetAccounts(id)
	if err != nil {
		return convertError(err)
	}

	return c.JSON(http.StatusOK, accounts)
}

// GetAccount GET /users/:userID/accounts/:accountID
func (handler *UserHandler) GetAccount(_c echo.Context) error {
	c := Context{_c}
	_ = c.Request().Context()
	_id := c.Param("userID")
	userID := uuid.FromStringOrNil(_id)
	_id = c.Param("accountID")
	accountID := uuid.FromStringOrNil(_id)

	account, err := handler.srv.GetAccount(userID, accountID)
	if err != nil {
		return convertError(err)
	}

	return c.JSON(http.StatusOK, account)
}

// AddAccount POST /users/:userID/accounts
func (handler *UserHandler) AddAccount(c echo.Context) error {
	_id := c.Param("userID")
	if _id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user id must not be blank")
	}

	id := uuid.FromStringOrNil(_id)
	if id == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid uuid")
	}

	req := Account{}
	err := c.Bind(&req)
	if err != nil {
		return err
	}

	args := repository.CreateAccountArgs{
		ID:          req.ID,
		Type:        req.Type,
		PrPermitted: req.PrPermitted,
	}

	account, err := handler.srv.CreateAccount(c.Request().Context(), id, &args)
	if err == repository.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, account)
}

// PatchAccount PATCH /users/:userID/accounts/:accountID
func (handler *UserHandler) PatchAccount(_c echo.Context) error {
	c := Context{_c}
	ctx := c.Request().Context()
	_id := c.Param("userID")
	userID := uuid.FromStringOrNil(_id)
	_id = c.Param("accountID")
	accountID := uuid.FromStringOrNil(_id)
	req := EditAccountRequest{}
	err := c.BindAndValidate(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	args := repository.UpdateAccountArgs{
		ID:          req.ID,
		Type:        req.Type,
		URL:         req.URL,
		PrPermitted: req.PrPermitted,
	}
	err = handler.srv.EditAccount(ctx, accountID, userID, &args)
	if err != nil {
		return convertError(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteAccount DELETE /users/:userID/accounts/:accountID
func (handler *UserHandler) DeleteAccount(c echo.Context) error {
	_accountid := c.Param("accountID")
	if _accountid == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user id must not be blank")
	}

	accountid := uuid.FromStringOrNil(_accountid)
	if accountid == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid uuid")
	}

	_userid := c.Param("userID")
	if _userid == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user id must not be blank")
	}

	userid := uuid.FromStringOrNil(_userid)
	if userid == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid uuid")
	}

	err := handler.srv.DeleteAccount(c.Request().Context(), accountid, userid)
	if err == repository.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
