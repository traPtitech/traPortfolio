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
	Level       EventLevel
	HostName    []*User
	GroupID     uuid.UUID
	RoomID      uuid.UUID
}

type EventLevel uint8

const (
	EventLevelAnonymous EventLevel = iota // 匿名で公開
	EventLevelPublic                      // 全て公開
	EventLevelPrivate                     // 外部に非公開
	EventLevelLimit
)
