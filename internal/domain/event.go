package domain

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/traPtitech/traPortfolio/internal/pkgs/optional"

	"github.com/gofrs/uuid"
)

type Event struct {
	ID        uuid.UUID
	Name      string
	Level     EventLevel
	TimeStart time.Time
	TimeEnd   time.Time
}

// Event knoQ上のイベント情報
type EventDetail struct {
	Event
	Description string
	Place       string
	HostName    []*User
	GroupID     uuid.UUID
	RoomID      uuid.UUID
}

// ApplyEventLevel EventDetailのLevelに応じてEventを返す
func ApplyEventLevel(e EventDetail) optional.Of[EventDetail] {
	switch e.Level {
	case EventLevelAnonymous:
		e.HostName = nil
		return optional.From(e)
	case EventLevelPublic:
		return optional.From(e)
	case EventLevelPrivate:
		return optional.Of[EventDetail]{}
	default:
		return optional.Of[EventDetail]{}
	}
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
			return fmt.Errorf("%w: EventLevel(%d) must be less than %d", ErrTooLargeEnum, newEL, EventLevelLimit)
		}

		*e = newEL
	}

	return nil
}

func (e EventLevel) Value() (driver.Value, error) {
	return sql.NullByte{Byte: byte(e), Valid: true}.Value()
}
