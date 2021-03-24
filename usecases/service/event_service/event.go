package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type EventService struct {
	repo repository.EventRepository
	knoQ repository.KnoqRepository
}

func NewEventService(repo repository.EventRepository, knoQ repository.KnoqRepository) EventService {
	return EventService{
		repo,
		knoQ,
	}
}

func (s *EventService) GetEvents(ctx context.Context) ([]*domain.Event, error) {
	er, err := s.knoQ.GetAll()
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Event, 0, len(er))
	for _, v := range er {
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

func (s *EventService) GetEventByID(ctx context.Context, id uuid.UUID) (*domain.EventDetail, error) {
	er, err := s.knoQ.GetByID(id)
	if err != nil {
		return nil, err
	}

	elv, err := s.repo.GetEventLevelByID(id)
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
