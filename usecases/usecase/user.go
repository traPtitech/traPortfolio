package usecase

import (
	"github.com/labstack/echo/v4"
)

type UserUsecase interface {
	Update(c echo.Context) (err error)
	DeleteByID(c echo.Context) (err error)
}
