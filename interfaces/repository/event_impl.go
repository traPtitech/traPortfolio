package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/external"
)

type EventRepository struct {
	database.SQLHandler
	external.KnoqAPI
}

func NewEventRepository(sql database.SQLHandler, knoq external.KnoqAPI) *EventRepository {
	return &EventRepository{SQLHandler: sql, KnoqAPI: knoq}
}

func (repo *EventRepository) GetEventLevels() ([]*domain.EventLevelRelation, error) {
	elvs := make([]*domain.EventLevelRelation, 0)
	err := repo.Find(&elvs).Error()
	if err != nil {
		return nil, err
	}
	return elvs, nil
}

func (repo *EventRepository) GetEventLevelByID(id uuid.UUID) (*domain.EventLevelRelation, error) {
	elv := domain.EventLevelRelation{ID: id}
	err := repo.First(&elv).Error()
	if err != nil {
		return nil, err
	}
	return &elv, nil
}
