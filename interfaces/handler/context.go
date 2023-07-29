package handler

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

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

func getID(c echo.Context, key idKey) (uuid.UUID, error) {
	id, err := uuid.FromString(c.Param(string(key)))
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %s", repository.ErrInvalidID, err.Error())
	} else if id.IsNil() {
		return uuid.Nil, repository.ErrNilID
	}

	return id, nil
}
