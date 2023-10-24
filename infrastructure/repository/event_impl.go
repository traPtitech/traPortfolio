package repository

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/infrastructure/repository/model"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"gorm.io/gorm"
)

type EventRepository struct {
	h    *gorm.DB
	knoq external.KnoqAPI
}

func NewEventRepository(sql *gorm.DB, knoq external.KnoqAPI) repository.EventRepository {
	return &EventRepository{h: sql, knoq: knoq}
}

func (r *EventRepository) GetEvents(_ context.Context) ([]*domain.Event, error) {
	events, err := r.knoq.GetEvents()
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

func (r *EventRepository) GetEvent(ctx context.Context, eventID uuid.UUID) (*domain.EventDetail, error) {
	er, err := r.knoq.GetEvent(eventID)
	if err != nil {
		return nil, err
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

	elv, err := r.getEventLevelByID(ctx, eventID)
	if err == nil {
		result.Level = elv.Level
	} else if errors.Is(err, repository.ErrNotFound) {
		result.Level = domain.EventLevelAnonymous
	} else {
		return nil, err
	}

	return result, nil
}

func (r *EventRepository) CreateEventLevel(ctx context.Context, arg *repository.CreateEventLevelArgs) error {
	_, err := r.knoq.GetEvent(arg.EventID)
	if err != nil {
		return err
	}

	relation := model.EventLevelRelation{
		ID:    arg.EventID,
		Level: arg.Level,
	}

	err = r.h.WithContext(ctx).Create(&relation).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *EventRepository) UpdateEventLevel(ctx context.Context, eventID uuid.UUID, arg *repository.UpdateEventLevelArgs) error {
	newLevel, ok := arg.Level.V()
	if !ok {
		return nil // updateする必要がないのでここでcommitする
	}

	err := r.h.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if elv, err := r.getEventLevelByID(ctx, eventID); err != nil {
			return err
		} else if elv.Level == newLevel {
			return nil // updateする必要がないのでここでcommitする
		}

		if err := tx.
			WithContext(ctx).
			Model(&model.EventLevelRelation{ID: eventID}).
			Update("level", newLevel).
			Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *EventRepository) GetUserEvents(_ context.Context, userID uuid.UUID) ([]*domain.Event, error) {
	events, err := r.knoq.GetEventsByUserID(userID)
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

func (r *EventRepository) getEventLevelByID(ctx context.Context, eventID uuid.UUID) (*model.EventLevelRelation, error) {
	elv := &model.EventLevelRelation{}
	err := r.h.
		WithContext(ctx).
		Where(&model.EventLevelRelation{ID: eventID}).
		First(elv).
		Error
	if err != nil {
		return nil, err
	}
	return elv, nil
}

// Interface guards
var (
	_ repository.EventRepository = (*EventRepository)(nil)
)
