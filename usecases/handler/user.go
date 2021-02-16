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
	p := c.Param("id")
	id, err := uuid.FromString(p)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	ctx := context.Background()
	result, err := handler.UserService.GetUser(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, result)
}

func (handler *UserHandler) Update(c echo.Context) error {
	u := domain.User{}
	err := c.Bind(&u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	_, err = handler.UserRepository.Update(&u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}
