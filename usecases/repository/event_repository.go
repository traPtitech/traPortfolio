//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/util/optional"
)

type CreateEventLevelArgs struct {
	EventID uuid.UUID
	Level   domain.EventLevel
}

type UpdateEventLevelArgs struct {
	Level optional.Uint8
}

type EventRepository interface {
	GetEvents(ctx context.Context) ([]*domain.Event, error)
	GetEvent(ctx context.Context, eventID uuid.UUID) (*domain.EventDetail, error)
	CreateEventLevel(ctx context.Context, args *CreateEventLevelArgs) error
	UpdateEventLevel(ctx context.Context, eventID uuid.UUID, args *UpdateEventLevelArgs) error
	GetUserEvents(ctx context.Context, userID uuid.UUID) ([]*domain.Event, error)
}
