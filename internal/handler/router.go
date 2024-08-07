package handler

import (
	"errors"
	"fmt"
	"net/http"

	vd "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/traPtitech/traPortfolio/internal/usecases/repository"
)

type Option func(*echo.Echo) error

func Setup(isProduction bool, e *echo.Echo, api API, opts ...Option) error {
	e.HTTPErrorHandler = newHTTPErrorHandler(e)
	e.Binder = &binderWithValidation{}

	for _, opt := range opts {
		if err := opt(e); err != nil {
			return fmt.Errorf("apply option: %w", err)
		}
	}

	e.Use(middleware.Recover())

	apiGroup := e.Group("/api")
	setupV1API(apiGroup, api, isProduction)

	return nil
}

func newHTTPErrorHandler(e *echo.Echo) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		var (
			code int
			herr *echo.HTTPError
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

		case errors.Is(err, repository.ErrUnauthorized):
			code = http.StatusUnauthorized

		case errors.Is(err, repository.ErrForbidden):
			code = http.StatusForbidden

		case errors.Is(err, repository.ErrNotFound):
			code = http.StatusNotFound

		case errors.Is(err, repository.ErrDBInternal):
			fallthrough
		default:
			e.Logger.Error(err)
			code = http.StatusInternalServerError
			herr = echo.NewHTTPError(code, http.StatusText(code)).SetInternal(err)
		}

		if herr == nil {
			herr = echo.NewHTTPError(
				code,
				fmt.Sprintf("%s: %s", http.StatusText(code), err.Error()),
			).SetInternal(err)
		}

		e.DefaultHTTPErrorHandler(herr, c)
	}
}

type binderWithValidation struct{}

var _ echo.Binder = (*binderWithValidation)(nil)

func (b *binderWithValidation) Bind(i interface{}, c echo.Context) error {
	if err := (&echo.DefaultBinder{}).Bind(i, c); err != nil {
		return fmt.Errorf("%w: %w", repository.ErrBind, err)
	}

	if vld, ok := i.(vd.Validatable); ok {
		if err := vld.Validate(); err != nil {
			if ie, ok := err.(vd.InternalError); ok {
				c.Logger().Fatalf("ozzo-validation internal error: %s", ie.Error())
			}

			return fmt.Errorf("%w: %w", repository.ErrValidate, err)
		}
	} else {
		c.Logger().Errorf("%T is not validatable", i)
	}

	return nil
}

func WithRequestLogger() Option {
	return func(e *echo.Echo) error {
		e.Use(middleware.Logger())
		return nil
	}
}
