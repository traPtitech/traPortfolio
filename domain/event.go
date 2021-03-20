package domain

import (
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
	HostName    []*User
	GroupID     uuid.UUID
	RoomID      uuid.UUID
}
