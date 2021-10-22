package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type EventService struct {
	event repository.EventRepository
	user  repository.UserRepository
}

func NewEventService(event repository.EventRepository, user repository.UserRepository) EventService {
	return EventService{event, user}
}

func (s *EventService) GetEvents(ctx context.Context) ([]*domain.Event, error) {
	return s.event.GetEvents()
}

func (s *EventService) GetEventByID(ctx context.Context, id uuid.UUID) (*domain.EventDetail, error) {
	event, err := s.event.GetEvent(id)
	if err != nil {
		return nil, err
	}

	// hostnameの詳細を取得
	users, err := s.user.GetUsers()
	if err != nil {
		return nil, err
	}

	umap := make(map[uuid.UUID]*domain.User)
	for _, u := range users {
		umap[u.ID] = u
	}

	for i, v := range event.HostName {
		if u, ok := umap[v.ID]; ok {
			event.HostName[i] = u
		}
	}

	return event, nil
}

func (s *EventService) UpdateEvent(ctx context.Context, id uuid.UUID, arg *repository.UpdateEventArg) error {
	return s.event.UpdateEvent(id, arg)
}
