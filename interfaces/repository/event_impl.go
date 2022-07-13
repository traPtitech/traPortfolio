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
	h    database.SQLHandler
	knoq external.KnoqAPI
}

func NewEventRepository(sql database.SQLHandler, knoq external.KnoqAPI) repository.EventRepository {
	return &EventRepository{h: sql, knoq: knoq}
}

func (repo *EventRepository) GetEvents() ([]*domain.Event, error) {
	events, err := repo.knoq.GetAll()
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

func (repo *EventRepository) GetEvent(eventID uuid.UUID) (*domain.EventDetail, error) {
	er, err := repo.knoq.GetByEventID(eventID)
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

	elv, err := repo.getEventLevelByID(eventID)
	if err == nil {
		result.Level = elv.Level
	} else if errors.Is(err, repository.ErrNotFound) {
		result.Level = domain.EventLevelAnonymous
	} else {
		return nil, convertError(err)
	}

	return result, nil
}

func (repo *EventRepository) CreateEventLevel(arg *repository.CreateEventLevelArgs) error {
	_, err := repo.knoq.GetByEventID(arg.EventID)
	if err != nil {
		return convertError(err)
	}

	relation := model.EventLevelRelation{
		ID:    arg.EventID,
		Level: arg.Level,
	}

	err = repo.h.Create(&relation).Error()
	if err != nil {
		return convertError(err)
	}

	return nil
}

func (repo *EventRepository) UpdateEventLevel(eventID uuid.UUID, arg *repository.UpdateEventLevelArgs) error {
	err := repo.h.Transaction(func(tx database.SQLHandler) error {
		if elv, err := repo.getEventLevelByID(eventID); err != nil {
			return convertError(err)
		} else if uint(elv.Level) == uint(arg.Level.Byte) {
			return nil // updateする必要がないのでここでcommitする
		}
		/*} else if elv.Level == arg.Level {
			return nil // updateする必要がないのでここでcommitする
		}*/

		if err := tx.Model(&model.EventLevelRelation{ID: eventID}).Update("level", arg.Level).Error(); err != nil {
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
	events, err := repo.knoq.GetByUserID(userID)
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

func (repo *EventRepository) getEventLevelByID(eventID uuid.UUID) (*model.EventLevelRelation, error) {
	elv := &model.EventLevelRelation{}
	err := repo.h.
		Where(&model.EventLevelRelation{ID: eventID}).
		First(elv).
		Error()
	if err != nil {
		return nil, convertError(err)
	}
	return elv, nil
}

// Interface guards
var (
	_ repository.EventRepository = (*EventRepository)(nil)
)
