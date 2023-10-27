package repository

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/infrastructure/external"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/random"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	sampleTime = time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)

	errUnexpected = errors.New("unexpected error")
)

type MockSQLHandler struct {
	Conn *gorm.DB
	Mock sqlmock.Sqlmock
}

func NewMockSQLHandler() *MockSQLHandler {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(err)
	}

	conf := mysql.Config{SkipInitializeWithVersion: true}
	conf.Conn = db

	d := dialector{mysql.New(conf)}
	engine, err := gorm.Open(d, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}

	return &MockSQLHandler{engine, mock}
}

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

func makeSQLQueryRegexp(query string) string {
	return fmt.Sprintf("^%s$", regexp.QuoteMeta(query))
}

func makeKnoqEvents(t *testing.T, events []*domain.Event) []*external.EventResponse {
	t.Helper()
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

func makeKnoqEvent(t *testing.T, event *domain.EventDetail) *external.EventResponse {
	t.Helper()
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

func mustMakeTraqGetAllArgs(t *testing.T, rargs *repository.GetUsersArgs) *external.TraQGetAllArgs {
	t.Helper()

	eargs, err := makeTraqGetAllArgs(rargs)
	assert.NoError(t, err)

	return eargs
}

func makeTraqUsers(t *testing.T, users []*domain.User) []*external.TraQUserResponse {
	t.Helper()

	res := make([]*external.TraQUserResponse, len(users))
	for i, u := range users {
		res[i] = &external.TraQUserResponse{
			ID: u.ID,
		}
	}

	return res
}

func makePortalUsers(t *testing.T, users []*domain.User) []*external.PortalUserResponse {
	t.Helper()
	res := make([]*external.PortalUserResponse, len(users))
	for i, u := range users {
		res[i] = makePortalUser(t, u)
	}

	return res
}

func makePortalUser(t *testing.T, user *domain.User) *external.PortalUserResponse {
	t.Helper()
	return &external.PortalUserResponse{
		TraQID:         user.Name,
		RealName:       user.RealNameForTest(t),
		AlphabeticName: random.AlphaNumeric(),
	}
}

func makeTraqUser(t *testing.T, user *domain.UserDetail) *external.TraQUserResponse {
	t.Helper()
	return &external.TraQUserResponse{
		ID:    user.ID,
		State: user.State,
	}
}

// Interface guards
var (
	_ sqlmock.Argument = anyTime{}
	_ sqlmock.Argument = anyUUID{}
)
