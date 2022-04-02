package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type Context struct {
	echo.Context
}

func (c *Context) BindAndValidate(i interface{}) error {
	if err := c.Bind(i); err != nil {
		c.Logger().Error(err)
		return repository.ErrBind
	}
	if err := c.Validate(i); err != nil {
		c.Logger().Error(err)
		return repository.ErrValidate
	}

	return nil
}
