package model

import (
	"time"

	"github.com/traPtitech/traPortfolio/internal/domain"

	"github.com/gofrs/uuid"
)

type EventLevelRelation struct {
	ID        uuid.UUID         `gorm:"type:char(36);not null;primaryKey"`
	Level     domain.EventLevel `gorm:"type:tinyint unsigned;not null;default:0"`
	CreatedAt time.Time         `gorm:"precision:6"`
	UpdatedAt time.Time         `gorm:"precision:6"`
}

func (*EventLevelRelation) TableName() string {
	return "event_level_relations"
}
