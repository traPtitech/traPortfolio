package model

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

// EventLevel
type EventLevel uint

const (
	EventLevelAnonymous = iota // 匿名で公開
	EventLevelPublic           // 全て公開
	EventLevelPrivate          // 外部に非公開
)

type EventLevelRelation struct {
	ID    uuid.UUID  `gorm:"type:char(36);not null;primary_key"`
	Level EventLevel `gorm:"type:tinyint unsigned;not null;default:0"`
}
