package model

import (
	"github.com/traPtitech/traPortfolio/domain"

	"github.com/gofrs/uuid"
)

type EventLevelRelation struct {
	ID    uuid.UUID         `gorm:"type:char(36);not null;primaryKey"`
	Level domain.EventLevel `gorm:"type:tinyint unsigned;not null;default:0"`
}
