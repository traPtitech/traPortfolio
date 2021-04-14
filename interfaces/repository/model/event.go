package model

import (
	"github.com/traPtitech/traPortfolio/domain"

	"github.com/gofrs/uuid"
)

type EventLevelRelation struct {
	ID    uuid.UUID         `gorm:"type:char(36);not null;primary_key"`
	Level domain.EventLevel `gorm:"type:tinyint unsigned;not null;default:0"`
}
