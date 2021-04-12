package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
)

type EventRepository interface {
	GetEvents() ([]*domain.Event, error)
	GetEvent(id uuid.UUID) (*domain.EventDetail, error)
}
