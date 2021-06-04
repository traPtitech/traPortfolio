package external

import (
	"time"

	"github.com/gofrs/uuid"
)

type EventResponse struct {
	ID          uuid.UUID `json:"eventId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	GroupID     uuid.UUID `json:"groupId"`
	RoomID      uuid.UUID `json:"roomId"`
	TimeStart   time.Time `json:"timeStart"`
	TimeEnd     time.Time `json:"timeEnd"`
	SharedRoom  bool      `json:"sharedRoom"`
}

type KnoqAPI interface {
	GetAll() ([]*EventResponse, error)
	GetByID(id uuid.UUID) (*EventResponse, error)
	GetByUserID(userID uuid.UUID) ([]*EventResponse, error)
}
