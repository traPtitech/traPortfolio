package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

func convertError(err error) error {
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
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("bad request: %w", err))

	case errors.Is(err, repository.ErrAlreadyExists):
		return echo.NewHTTPError(http.StatusConflict, fmt.Errorf("conflicts: %w", err))

	case errors.Is(err, repository.ErrForbidden):
		return echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("forbideen: %w", err))

	case errors.Is(err, repository.ErrNotFound):
		return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("not found: %w", err))

	default:
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("%w", err))
	}
}
