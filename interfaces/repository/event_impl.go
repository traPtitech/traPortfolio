package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type EventRepository struct {
	database.SQLHandler
	external.KnoqAPI
}

func NewEventRepository(sql database.SQLHandler, knoq external.KnoqAPI) *EventRepository {
	return &EventRepository{SQLHandler: sql, KnoqAPI: knoq}
}

func (repo *EventRepository) GetEventLevels() (map[uuid.UUID]*domain.EventLevelRelation, error) {
	elvs := make([]*domain.EventLevelRelation, 0)
	err := repo.Find(&elvs).Error()
	if err != nil && err != repository.ErrNotFound {
		return nil, err
	}
	elvsmp := make(map[uuid.UUID]*domain.EventLevelRelation, len(elvs))
	for _, v := range elvs {
		v := v
		elvsmp[v.ID] = v
	}
	return elvsmp, nil
}

func (repo *EventRepository) GetEventLevelByID(id uuid.UUID) (*domain.EventLevelRelation, error) {
	elv := domain.EventLevelRelation{ID: id}
	err := repo.First(&elv).Error()
	if err != nil {
		return nil, err
	}
	return &elv, nil
}
