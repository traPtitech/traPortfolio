package repository

import (
	"github.com/gofrs/uuid"

	"github.com/traPtitech/traPortfolio/domain"
)

type EventRepository interface {
	GetEventLevels() (map[uuid.UUID]*domain.EventLevelRelation, error)
	GetEventLevelByID(ID uuid.UUID) (*domain.EventLevelRelation, error)
}
