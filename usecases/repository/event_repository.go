//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
)

type UpdateEventLevelArgs struct {
	Level domain.EventLevel
}

type EventRepository interface {
	GetEvents() ([]*domain.Event, error)
	GetEvent(eventID uuid.UUID) (*domain.EventDetail, error)
	UpdateEventLevel(eventID uuid.UUID, args *UpdateEventLevelArgs) error
	GetUserEvents(userID uuid.UUID) ([]*domain.Event, error)
}
