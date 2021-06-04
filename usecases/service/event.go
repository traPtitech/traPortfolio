package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type EventService struct {
	repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) EventService {
	return EventService{
		repo,
	}
}

func (s *EventService) GetEvents(ctx context.Context) ([]*domain.Event, error) {
	return s.repo.GetEvents()
}

func (s *EventService) GetEventByID(ctx context.Context, id uuid.UUID) (*domain.EventDetail, error) {
	return s.repo.GetEvent(id)
}

func (s *EventService) UpdateEvent(ctx context.Context, id uuid.UUID, arg *repository.UpdateEventArg) error {
	elv := model.EventLevelRelation{
		ID: id,
		Level: arg.Level,
	}
	return s.repo.UpdateEvent(&elv)
}
