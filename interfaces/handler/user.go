package handler

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	service "github.com/traPtitech/traPortfolio/usecases/service/user_service"
)

type EditUserRequest struct {
	Bio          string `json:"bio"`
	HideRealName bool   `json:"hideRealName"`
}

type UserHandler struct {
	UserService service.UserService
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

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{UserService: s}
}

// GetAll GET /users
func (handler *UserHandler) GetAll(c echo.Context) error {
	ctx := c.Request().Context()
	users, err := handler.UserService.GetUsers(ctx)
	if err != nil {
		return err
	}

	res := make([]*userResponse, len(users))
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
	user, err := handler.UserService.GetUser(ctx, id)
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
	u := domain.EditUser{
		ID:          id,
		Description: req.Bio,
		Check:       !req.HideRealName,
	}
	err = handler.UserService.Update(ctx, &u)
	if err == repository.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
