package controllers

import (
	"net/http"
	"strconv"

	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/input"
	"github.com/traPtitech/traPortfolio/usecases/usecase"
)

type UserController struct {
	Interactor usecase.UserUsecase
}

func NewUserController(it usecase.UserUsecase) *UserController {
	return &UserController{
		Interactor: it,
	}
}

func (controller *UserController) Show(c Context) (err error) {
	id, _ := strconv.Atoi(c.Param("id"))
	ipt := input.GetUser{
		ID: id,
	}
	user, err := controller.Interactor.UserByID(ipt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, NewError(err))
	}
	return c.JSON(http.StatusOK, user)
}

func (controller *UserController) Index(c Context) (err error) {
	users, err := controller.Interactor.Users()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, NewError(err))
	}
	return c.JSON(http.StatusOK, users)
}

func (controller *UserController) Create(c Context) (err error) {
	u := domain.User{}
	err = c.Bind(&u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, NewError(err))
	}
	ipt := input.AddUser{User: u}
	user, err := controller.Interactor.Add(ipt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, NewError(err))
	}
	return c.JSON(http.StatusCreated, user)
}

func (controller *UserController) Save(c Context) (err error) {
	u := domain.User{}
	err = c.Bind(&u)
	if err != nil {
		return
	}
	ipt := input.UpdateUser{User: u}
	user, err := controller.Interactor.Update(ipt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, NewError(err))
	}
	return c.JSON(http.StatusCreated, user)
}

func (controller *UserController) Delete(c Context) (err error) {
	id, _ := strconv.Atoi(c.Param("id"))
	ipt := input.DeleteUser{ID: id}
	err = controller.Interactor.DeleteByID(ipt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, NewError(err))
	}
	return c.NoContent(http.StatusNoContent)
}
