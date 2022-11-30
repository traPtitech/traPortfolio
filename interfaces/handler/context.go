package handler

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type Context struct {
	echo.Context
}

func (c *Context) BindAndValidate(i interface{}) error {
	if err := c.Bind(i); err != nil {
		return fmt.Errorf("%w: %s", repository.ErrBind, err.Error())
	}
	if err := c.Validate(i); err != nil {
		return fmt.Errorf("%w: %s", repository.ErrValidate, err.Error())
	}

	return nil
}

type idKey string

const (
	keyUserID        idKey = "userID"
	keyUserAccountID idKey = "accountID"
	keyProject       idKey = "projectID"
	keyEventID       idKey = "eventID"
	keyContestID     idKey = "contestID"
	keyContestTeamID idKey = "teamID"
	keyGroupID       idKey = "groupID"
)

func (c *Context) getID(key idKey) (uuid.UUID, error) {
	id, err := uuid.FromString(c.Param(string(key)))
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %s", repository.ErrInvalidID, err.Error())
	} else if id.IsNil() {
		return uuid.Nil, repository.ErrNilID
	}

	return id, nil
}
