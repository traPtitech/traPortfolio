package domain

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/gofrs/uuid"
)

type Event struct {
	ID        uuid.UUID
	Name      string
	TimeStart time.Time
	TimeEnd   time.Time
}

// Event knoQ上のイベント情報
type EventDetail struct {
	Event
	Description string
	Place       string
	Level       EventLevel
	HostName    []*User
	GroupID     uuid.UUID
	RoomID      uuid.UUID
}

type EventLevel uint8

var (
	_ sql.Scanner   = (*EventLevel)(nil)
	_ driver.Valuer = EventLevel(0)
)

const (
	EventLevelAnonymous EventLevel = iota // 匿名で公開
	EventLevelPublic                      // 全て公開
	EventLevelPrivate                     // 外部に非公開
	EventLevelLimit
)

func (e *EventLevel) Scan(src interface{}) error {
	s := sql.NullByte{}
	if err := s.Scan(src); err != nil {
		return err
	}

	if s.Valid {
		newEL := EventLevel(s.Byte)
		if newEL >= EventLevelLimit {
			return ErrTooLargeEnum
		}

		*e = newEL
	}

	return nil
}

func (e EventLevel) Value() (driver.Value, error) {
	return sql.NullByte{Byte: byte(e), Valid: true}.Value()
}
