package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type EventRepository struct {
	h   database.SQLHandler
	api external.KnoqAPI
}

func NewEventRepository(sql database.SQLHandler, knoq external.KnoqAPI) *EventRepository {
	return &EventRepository{h: sql, api: knoq}
}

func (repo *EventRepository) GetEvents() ([]*domain.Event, error) {
	events, err := repo.api.GetAll()
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Event, 0, len(events))
	for _, v := range events {
		e := &domain.Event{
			ID:        v.ID,
			Name:      v.Name,
			TimeStart: v.TimeStart,
			TimeEnd:   v.TimeEnd,
		}
		result = append(result, e)
	}

	return result, nil
}

func (repo *EventRepository) GetEvent(id uuid.UUID) (*domain.EventDetail, error) {
	er, err := repo.api.GetByID(id)
	if err != nil {
		return nil, err
	}

	elv, err := repo.getEventLevelByID(id)
	if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	result := &domain.EventDetail{
		Event: domain.Event{
			ID:        er.ID,
			Name:      er.Name,
			TimeStart: er.TimeStart,
			TimeEnd:   er.TimeEnd,
		},
		Description: er.Description,
		GroupID:     er.GroupID,
		RoomID:      er.RoomID,
	}

	if err == nil {
		result.Level = elv.Level
	}

	return result, nil
}

func (repo *EventRepository) getEventLevelByID(id uuid.UUID) (*model.EventLevelRelation, error) {
	elv := &model.EventLevelRelation{}
	err := repo.h.First(elv, &model.EventLevelRelation{ID: id}).Error()
	if err != nil {
		return nil, err
	}
	return elv, nil
}

// Interface guards
var (
	_ repository.EventRepository = (*EventRepository)(nil)
)
