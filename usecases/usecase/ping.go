package usecase

import "github.com/labstack/echo/v4"

type PingUsecase interface {
	Ping(e echo.Context) error
}
