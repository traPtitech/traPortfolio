//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
)

type UpdateEventArg struct {
	Level domain.EventLevel
}

type EventRepository interface {
	GetEvents() ([]*domain.Event, error)
	GetEvent(id uuid.UUID) (*domain.EventDetail, error)
	UpdateEvent(elv *model.EventLevelRelation) error
	GetUserEvents(userID uuid.UUID) ([]*domain.Event, error)
}
