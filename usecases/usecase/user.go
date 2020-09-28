package usecase

import (
	"github.com/labstack/echo/v4"
)

type UserUsecase interface {
	Get(c echo.Context) (err error)
	Update(c echo.Context) (err error)
}
