//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/util/optional"
)

type CreateEventLevelArgs struct {
	EventID uuid.UUID
	Level   domain.EventLevel
}

type UpdateEventLevelArgs struct {
	Level optional.Of[domain.EventLevel]
}

type EventRepository interface {
	GetEvents() ([]*domain.Event, error)
	GetEvent(eventID uuid.UUID) (*domain.EventDetail, error)
	CreateEventLevel(args *CreateEventLevelArgs) error
	UpdateEventLevel(eventID uuid.UUID, args *UpdateEventLevelArgs) error
	GetUserEvents(userID uuid.UUID) ([]*domain.Event, error)
}
