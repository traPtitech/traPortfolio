package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type UserHandler struct {
	UserRepository repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{UserRepository: repo}
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
