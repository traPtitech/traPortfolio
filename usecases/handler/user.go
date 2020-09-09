package handler

import (
	"net/http"
	"strconv"

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

func (handler *UserHandler) UserByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	user, err := handler.UserRepository.FindByID(id)
	if err != nil {
		return c.JSON(http.StatusOK, user)
	}
	return c.JSON(http.StatusOK, user)
}

func (handler *UserHandler) Users(c echo.Context) error {
	users, err := handler.UserRepository.FindAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, users)
}

func (handler *UserHandler) Add(c echo.Context) error {
	u := domain.User{}
	err := c.Bind(&u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	_, err = handler.UserRepository.Store(&u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}

func (handler *UserHandler) Update(c echo.Context) error {
	u := domain.User{}
	err := c.Bind(&u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	_, err = handler.UserRepository.Update(&u)
	return c.NoContent(http.StatusOK)
}

func (handler *UserHandler) DeleteByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err = handler.UserRepository.DeleteByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusOK)
}
