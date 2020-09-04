package controllers

import (
	"net/http"
	"strconv"

	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/repository"
	"github.com/traPtitech/traPortfolio/usecase/input"
	"github.com/traPtitech/traPortfolio/usecase/interactor"
)

type UserController struct {
	Interactor interactor.UserInteractor
}

func NewUserController(sqlHandler database.SqlHandler) *UserController {
	return &UserController{
		Interactor: interactor.UserInteractor{
			UserRepository: &repository.UserRepository{
				SqlHandler: sqlHandler,
			},
		},
	}
}

func (controller *UserController) Show(c Context) (err error) {
	id, _ := strconv.Atoi(c.Param("id"))
	ipt := input.GetUser{
		Id: id,
	}
	user, err := controller.Interactor.UserById(ipt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}
	c.JSON(http.StatusOK, user)
	return
}

func (controller *UserController) Index(c Context) (err error) {
	users, err := controller.Interactor.Users()
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}
	c.JSON(http.StatusOK, users)
	return
}

func (controller *UserController) Create(c Context) (err error) {
	u := domain.User{}
	c.Bind(&u)
	ipt := input.AddUser{User: u}
	user, err := controller.Interactor.Add(ipt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}
	c.JSON(http.StatusCreated, user)
	return
}

func (controller *UserController) Save(c Context) (err error) {
	u := domain.User{}
	c.Bind(&u)
	ipt := input.UpdateUser{User: u}
	user, err := controller.Interactor.Update(ipt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}
	c.JSON(http.StatusCreated, user)
	return
}

func (controller *UserController) Delete(c Context) (err error) {
	id, _ := strconv.Atoi(c.Param("id"))
	ipt := input.DeleteUser{Id: id}
	err = controller.Interactor.DeleteById(ipt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError(err))
		return
	}
	c.NoContent(http.StatusNoContent)
	return
}
