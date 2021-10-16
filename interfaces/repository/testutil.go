package repository

import (
	"database/sql/driver"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
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

// Interface guards
var (
	_ sqlmock.Argument = anyTime{}
	_ sqlmock.Argument = anyUUID{}
)
