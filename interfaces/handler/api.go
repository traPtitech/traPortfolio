package handler

import (
	"github.com/labstack/echo/v4"
)

type API struct {
	Ping    *PingHandler
	User    *UserHandler
	Event   *EventHandler
	Contest *ContestHandler
}

func NewAPI(ping *PingHandler, user *UserHandler, event *EventHandler, contest *ContestHandler) API {
	return API{
		Ping:    ping,
		User:    user,
		Event:   event,
		Contest: contest,
	}
}

type Context struct {
	echo.Context
}

func (c *Context) BindAndValidate(i interface{}) error {
	if err := c.Bind(i); err != nil {
		return err
	}
	if err := c.Validate(i); err != nil {
		return err
	}
	return nil
}
