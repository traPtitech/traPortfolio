package repository

import (
	"context"

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

func (repo *EventRepository) GetEventLevels() ([]*domain.EventLevelRelation, error) {
	elv := make([]*domain.EventLevelRelation, 0)
	err := repo.Find(&elv).Error()
	if err != nil {
		return nil, err
	}
	return elv, nil
}

func (repo *EventRepository) GetEvents(ctx context.Context) (events []*domain.Event, err error) {
	er, err := repo.KnoqAPI.GetAll()
	if err != nil {
		return nil, err
	}

	elvs, err := repo.GetEventLevels()
	if err == repository.ErrNotFound {
		elvs = make([]*domain.EventLevelRelation, 0)
	}
	if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	elvsmp := make(map[uuid.UUID]*domain.EventLevelRelation)
	for _, v := range elvs {
		v := v
		elvsmp[v.ID] = v
	}
	result := make([]*domain.Event, 0, len(er))
	for _, v := range er {
		_level, ok := elvsmp[v.ID]
		var level domain.EventLevel = 1
		if ok {
			level = _level.Level
		}
		result = append(result, &domain.Event{
			ID:          v.ID,
			Description: v.Description,
			GroupID:     v.GroupID,
			Level:       level,
			Name:        v.Name,
			RoomID:      v.RoomID,
			SharedRoom:  v.SharedRoom,
			TimeEnd:     v.TimeEnd,
			TimeStart:   v.TimeStart,
		})
	}

	return result, nil
}

func (repo *EventRepository) GetEventByID(ctx context.Context, ID uuid.UUID) (*domain.Event, error) {
	er, err := repo.KnoqAPI.GetByID(ID)
	if err != nil {
		return nil, err
	}

	elv := domain.EventLevelRelation{ID: ID}
	err = repo.First(&elv).Error()
	if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	var level domain.EventLevel = 1
	if err == nil {
		level = elv.Level
	}
	result := &domain.Event{
		ID:          er.ID,
		Description: er.Description,
		GroupID:     er.GroupID,
		Level:       level,
		Name:        er.Name,
		RoomID:      er.RoomID,
		SharedRoom:  er.SharedRoom,
		TimeEnd:     er.TimeEnd,
		TimeStart:   er.TimeStart,
	}

	return result, nil
}
