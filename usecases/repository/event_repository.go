package repository

import (
	"context"

	"github.com/gofrs/uuid"

	"github.com/traPtitech/traPortfolio/domain"
)

type EventRepository interface {
	GetEventByID(context.Context, uuid.UUID) (*domain.Event, error)
	GetEvents(context.Context) ([]*domain.Event, error)
}
