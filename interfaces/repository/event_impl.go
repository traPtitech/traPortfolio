package repository

import (
	"errors"

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

func NewEventRepository(sql database.SQLHandler, knoq external.KnoqAPI) repository.EventRepository {
	return &EventRepository{h: sql, api: knoq}
}

func (repo *EventRepository) GetEvents() ([]*domain.Event, error) {
	events, err := repo.api.GetAll()
	if err != nil {
		return nil, convertError(err)
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
		return nil, convertError(err)
	}

	// IDのリストだけ取得、Name,RealNameはPortalから取得する
	hostName := make([]*domain.User, 0, len(er.Admins))
	for _, aid := range er.Admins {
		hostName = append(hostName, &domain.User{ID: aid})
	}

	result := &domain.EventDetail{
		Event: domain.Event{
			ID:        er.ID,
			Name:      er.Name,
			TimeStart: er.TimeStart,
			TimeEnd:   er.TimeEnd,
		},
		Description: er.Description,
		Place:       er.Place,
		// Level:
		HostName: hostName,
		GroupID:  er.GroupID,
		RoomID:   er.RoomID,
	}

	elv, err := repo.getEventLevelByID(id)
	if err == nil {
		result.Level = elv.Level
	} else if errors.Is(err, repository.ErrNotFound) {
		result.Level = domain.EventLevelAnonymous
	} else {
		return nil, convertError(err)
	}

	return result, nil
}

func (repo *EventRepository) UpdateEventLevel(id uuid.UUID, arg *repository.UpdateEventLevelArg) error {
	err := repo.h.Transaction(func(tx database.SQLHandler) error {
		if elv, err := repo.getEventLevelByID(id); err != nil {
			return convertError(err)
		} else if elv.Level == arg.Level {
			return nil // updateする必要がないのでここでcommitする
		}

		if err := tx.Model(&model.EventLevelRelation{ID: id}).Update("level", arg.Level).Error(); err != nil {
			return convertError(err)
		}

		return nil
	})
	if err != nil {
		return convertError(err)
	}

	return nil
}

func (repo *EventRepository) GetUserEvents(userID uuid.UUID) ([]*domain.Event, error) {
	events, err := repo.api.GetByUserID(userID)
	if err != nil {
		return nil, convertError(err)
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

func (repo *EventRepository) getEventLevelByID(id uuid.UUID) (*model.EventLevelRelation, error) {
	elv := &model.EventLevelRelation{}
	err := repo.h.First(elv, &model.EventLevelRelation{ID: id}).Error()
	if err != nil {
		return nil, convertError(err)
	}
	return elv, nil
}

// Interface guards
var (
	_ repository.EventRepository = (*EventRepository)(nil)
)
