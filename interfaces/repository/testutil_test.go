package repository_test

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

func makeKnoqEvents(events []*domain.Event) []*external.EventResponse {
	res := make([]*external.EventResponse, len(events))
	for i, e := range events {
		res[i] = &external.EventResponse{
			ID:        e.ID,
			Name:      e.Name,
			TimeStart: e.TimeStart,
			TimeEnd:   e.TimeEnd,
		}
	}

	return res
}

func makeKnoqEvent(event *domain.EventDetail) *external.EventResponse {
	admins := make([]uuid.UUID, len(event.HostName))
	for i, h := range event.HostName {
		admins[i] = h.ID
	}

	return &external.EventResponse{
		ID:          event.ID,
		Name:        event.Name,
		Description: event.Description,
		Place:       event.Place,
		GroupID:     event.GroupID,
		RoomID:      event.RoomID,
		TimeStart:   event.TimeStart,
		TimeEnd:     event.TimeEnd,
		Admins:      admins,
	}
}

func makePortalUsers(users []*domain.User) []*external.PortalUserResponse {
	res := make([]*external.PortalUserResponse, len(users))
	for i, u := range users {
		res[i] = makePortalUser(u)
	}

	return res
}

func makePortalUser(user *domain.User) *external.PortalUserResponse {
	return &external.PortalUserResponse{
		TraQID:         user.Name,
		RealName:       user.RealName,
		AlphabeticName: random.AlphaNumeric(rand.Intn(30) + 1),
	}
}

func makeTraqUser(user *domain.UserDetail) *external.TraQUserResponse {
	return &external.TraQUserResponse{
		State:       user.State,
		Bot:         false,
		DisplayName: random.AlphaNumeric(rand.Intn(30) + 1),
		Name:        user.Name,
	}
}

// Interface guards
var (
	_ sqlmock.Argument = anyTime{}
	_ sqlmock.Argument = anyUUID{}
)
