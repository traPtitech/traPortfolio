package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

func convertError(err error) error {
	var (
		code int
		msg  string
	)

	switch {
	case errors.Is(err, repository.ErrNilID):
		fallthrough
	case errors.Is(err, repository.ErrInvalidID):
		fallthrough
	case errors.Is(err, repository.ErrInvalidArg):
		fallthrough
	case errors.Is(err, repository.ErrBind):
		fallthrough
	case errors.Is(err, repository.ErrValidate):
		code = http.StatusBadRequest

	case errors.Is(err, repository.ErrAlreadyExists):
		code = http.StatusConflict

	case errors.Is(err, repository.ErrForbidden):
		code = http.StatusForbidden

	case errors.Is(err, repository.ErrNotFound):
		code = http.StatusNotFound

	default:
		code = http.StatusInternalServerError
		msg = http.StatusText(code)
	}

	if len(msg) == 0 {
		msg = fmt.Sprintf("%s: %s", http.StatusText(code), err.Error())
	}

	return echo.NewHTTPError(code, msg).SetInternal(err)
}
