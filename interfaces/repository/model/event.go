package model

import (
	"time"

	"github.com/traPtitech/traPortfolio/domain"

	"github.com/gofrs/uuid"
)

// Event knoQ上のイベント情報
type Event struct {
	ID          uuid.UUID
	Name        string
	TimeStart   time.Time
	TimeEnd     time.Time
	Description string
	Place       string
	HostName    []*User
	GroupID     uuid.UUID
	RoomID      uuid.UUID
}

type EventLevelRelation struct {
	ID    uuid.UUID         `gorm:"type:char(36);not null;primary_key"`
	Level domain.EventLevel `gorm:"type:tinyint unsigned;not null;default:0"`
}
