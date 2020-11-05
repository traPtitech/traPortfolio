package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

// Event knoQ上のイベント情報
type Event struct {
	ID          uuid.UUID
	Name        string
	Description string
	GroupID     uuid.UUID
	RoomID      uuid.UUID
	TimeStart   time.Time
	TimeEnd     time.Time
	SharedRoom  bool
	Level       EventLevel
}

// EventLevel 0 匿名で公開, 1 全て公開, 2 部内にのみ公開
type EventLevel uint

const (
	EventLevelAnonymous = iota
	EventLevelPublic
	EventLevelPrivate
)

type EventLevelRelation struct {
	ID    uuid.UUID  `gorm:"type:char(36);not null;primary_key"`
	Level EventLevel `gorm:"type:tinyint unsigned;not null;default:0"`
}
