package handler

import (
	"net/http"

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

func (handler *UserHandler) Get(c echo.Context) error {
	u := domain.User{}
	err := c.Bind(&u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	result := handler.UserService.GetUser(c.Request().Context(), u.ID)
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
