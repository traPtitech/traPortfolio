package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

// Event knoQ上のイベント情報
type Event struct {
	ID          uuid.UUID  `json:"eventId"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	GroupID     uuid.UUID  `json:"groupId"`
	RoomID      uuid.UUID  `json:"roomId"`
	TimeStart   time.Time  `json:"timeStart"`
	TimeEnd     time.Time  `json:"timeEnd"`
	SharedRoom  bool       `json:"sharedRoom"`
	Level       EventLevel `json:"eventLevel"`
}

// EventLevel 0 全て公開, 1 匿名で公開
type EventLevel uint

type EventLevelRelation struct {
	ID    uuid.UUID  `gorm:"type:char(36);not null;primary_key"`
	Level EventLevel `gorm:"type:tinyint unsigned;not null;default:1"`
}
