package repository

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/samber/lo"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/infrastructure/external"
	"github.com/traPtitech/traPortfolio/infrastructure/repository/model"
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

func (r *EventRepository) GetEvents(ctx context.Context) ([]*domain.Event, error) {
	events, err := r.knoq.GetEvents()
	if err != nil {
		return nil, err
	}

	result, err := r.completeEventLevels(ctx, events)
	if err != nil {
		return nil, err
	}
	result = filterAccessibleEvents(result)
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

	ed := domain.EventDetail{
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
		ed.EventLevel = elv.Level
	} else if errors.Is(err, repository.ErrNotFound) {
		ed.EventLevel = domain.EventLevelAnonymous
	} else {
		return nil, err
	}

	res := domain.ApplyEventLevel(ed)
	if v, ok := res.V(); ok {
		return &v, nil
	}
	return nil, repository.ErrNotFound
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

func (r *EventRepository) GetUserEvents(ctx context.Context, userID uuid.UUID) ([]*domain.Event, error) {
	events, err := r.knoq.GetEventsByUserID(userID)
	if err != nil {
		return nil, err
	}

	result, err := r.completeEventLevels(ctx, events)
	if err != nil {
		return nil, err
	}
	// TODO: 自分のイベントで非公開のものが見られない
	result = filterAccessibleEvents(result)
	return result, nil
}

func (r *EventRepository) completeEventLevels(ctx context.Context, events []*external.EventResponse) ([]*domain.Event, error) {
	ids := lo.Map(events, func(e *external.EventResponse, _ int) uuid.UUID {
		return e.ID
	})
	rels := make([]*model.EventLevelRelation, 0, len(events))
	err := r.h.
		WithContext(ctx).
		Where("id IN ?", ids).
		Find(&rels).
		Error
	if err != nil {
		return nil, err
	}
	relById := lo.Associate(rels, func(r *model.EventLevelRelation) (uuid.UUID, domain.EventLevel) {
		return r.ID, r.Level
	})
	result := lo.Map(events, func(e *external.EventResponse, _ int) *domain.Event {
		level, ok := relById[e.ID]
		if !ok {
			level = domain.EventLevelAnonymous
		}
		return &domain.Event{
			ID:         e.ID,
			Name:       e.Name,
			TimeStart:  e.TimeStart,
			TimeEnd:    e.TimeEnd,
			EventLevel: level,
		}
	})
	return result, nil
}

func filterAccessibleEvents(events []*domain.Event) []*domain.Event {
	// privateのものだけ除外する
	return lo.Filter(events, func(e *domain.Event, _ int) bool {
		return e.EventLevel != domain.EventLevelPrivate
	})
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
