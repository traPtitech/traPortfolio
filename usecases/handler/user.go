package handler

import (
	"context"
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
	UserRepository repository.UserRepository
	UserService    service.UserService
}

func NewUserHandler(repo repository.UserRepository, s service.UserService) *UserHandler {
	return &UserHandler{UserRepository: repo, UserService: s}
}

func (handler *UserHandler) GetAll(c echo.Context) error {
	ctx := context.Background()
	result, err := handler.UserService.GetUsers(ctx)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func (handler *UserHandler) GetByID(c echo.Context) error {
	_id := c.Param("userID")
	if _id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user id must not be blank")
	}

	id := uuid.FromStringOrNil(_id)
	if id == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid uuid")
	}
	ctx := context.Background()
	result, err := handler.UserService.GetUser(ctx, id)
	if err == repository.ErrNotFound {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

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
	u := domain.User{
		ID:          id,
		Description: req.Bio,
		Check:       !req.HideRealName,
	}
	err = handler.UserRepository.Update(&u)
	if err == repository.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (handler *UserHandler) AddAccount(c echo.Context) error {
	_id := c.Param("userID")
	if _id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user id must not be blank")
	}

	id := uuid.FromStringOrNil(_id)
	if id == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid uuid")
	}

	a := domain.Account{}
	err := c.Bind(&a)
	if err != nil {
		return err
	}

	err = handler.UserRepository.AddAccount(a)
	if err == repository.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (handler *UserHandler) DeleteAccount(c echo.Context) error {
	_id := c.Param("userID")
	if _id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user id must not be blank")
	}

	id := uuid.FromStringOrNil(_id)
	if id == uuid.Nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid uuid")
	}

	a := domain.Account{}
	err := c.Bind(&a)
	if err != nil {
		return err
	}

	err = handler.UserRepository.DeleteAccount(a)
	if err == repository.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
