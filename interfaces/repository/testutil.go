package repository

import (
	"database/sql/driver"
	"errors"
	"math/rand"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/util/random"
)

var (
	sampleTime = time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)

	errUnexpected = errors.New("unexpected error")
)

type anyTime struct{}

func (a anyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type anyUUID struct{}

func (a anyUUID) Match(v driver.Value) bool {
	vstr, ok := v.(string)
	return ok && validUUIDStr(vstr)
}

func validUUIDStr(s string) bool {
	_, err := uuid.FromString(s)
	return err == nil
}

func makePortalUsers(users []*domain.User) []*external.PortalUserResponse {
	res := make([]*external.PortalUserResponse, len(users))
	for i, u := range users {
		res[i] = &external.PortalUserResponse{
			TraQID:         u.Name,
			RealName:       u.RealName,
			AlphabeticName: random.AlphaNumeric(rand.Intn(30) + 1),
		}
	}

	return res
}

// Interface guards
var (
	_ sqlmock.Argument = anyTime{}
	_ sqlmock.Argument = anyUUID{}
)
