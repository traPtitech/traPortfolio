//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
)

type EventRepository interface {
	GetEvents() ([]*domain.Event, error)
	GetEvent(id uuid.UUID) (*domain.EventDetail, error)
	GetUserEvents(userID uuid.UUID) ([]*domain.Event, error)
}
