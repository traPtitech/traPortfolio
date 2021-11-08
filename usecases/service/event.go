//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type EventService interface {
	GetEvents(ctx context.Context) ([]*domain.Event, error)
	GetEventByID(ctx context.Context, id uuid.UUID) (*domain.EventDetail, error)
	UpdateEventLevel(ctx context.Context, id uuid.UUID, arg *repository.UpdateEventLevelArg) error
}

type eventService struct {
	event repository.EventRepository
	user  repository.UserRepository
}

func NewEventService(event repository.EventRepository, user repository.UserRepository) EventService {
	return &eventService{event, user}
}

func (s *eventService) GetEvents(ctx context.Context) ([]*domain.Event, error) {
	return s.event.GetEvents()
}

func (s *eventService) GetEventByID(ctx context.Context, id uuid.UUID) (*domain.EventDetail, error) {
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

func (s *eventService) UpdateEventLevel(ctx context.Context, id uuid.UUID, arg *repository.UpdateEventLevelArg) error {
	return s.event.UpdateEventLevel(id, arg)
}

// Interface guards
var (
	_ EventService = (*eventService)(nil)
)
