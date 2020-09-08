package usecase

import "github.com/labstack/echo"

type PingUsecase interface {
	Ping(e echo.Context) error
}
