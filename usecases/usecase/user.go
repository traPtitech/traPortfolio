package usecase

import (
	"github.com/labstack/echo"
)

type UserUsecase interface {
	UserByID(c echo.Context) (err error)
	Users(c echo.Context) (err error)
	Add(c echo.Context) (err error)
	Update(c echo.Context) (err error)
	DeleteByID(c echo.Context) (err error)
}
