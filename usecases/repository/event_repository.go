package repository

import (
	"github.com/gofrs/uuid"

	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
)

type EventRepository interface {
	GetEventLevels() (map[uuid.UUID]*model.EventLevelRelation, error)
	GetEventLevelByID(ID uuid.UUID) (*model.EventLevelRelation, error)
}
