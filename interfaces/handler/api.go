package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type API struct {
	Ping    *PingHandler
	User    *UserHandler
	Project *ProjectHandler
	Event   *EventHandler
	Contest *ContestHandler
}

func NewAPI(ping *PingHandler, user *UserHandler, project *ProjectHandler, event *EventHandler, contest *ContestHandler) API {
	return API{
		Ping:    ping,
		User:    user,
		Project: project,
		Event:   event,
		Contest: contest,
	}
}

type Context struct {
	echo.Context
}

func (c *Context) BindAndValidate(i interface{}) error {
	if err := c.Bind(i); err != nil {
		return repository.ErrBind
	}
	if err := c.Validate(i); err != nil {
		return repository.ErrValidate
	}

	return nil
}

func IsValidUUID(fl validator.FieldLevel) bool {
	id, ok := fl.Field().Interface().(uuid.UUID)
	return ok && id != uuid.Nil
}
