//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package external

import (
	"time"

	"github.com/gofrs/uuid"
)

type EventResponse struct {
	ID          uuid.UUID   `json:"eventId"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Place       string      `json:"place"`
	GroupID     uuid.UUID   `json:"groupId"`
	RoomID      uuid.UUID   `json:"roomId"`
	TimeStart   time.Time   `json:"timeStart"`
	TimeEnd     time.Time   `json:"timeEnd"`
	SharedRoom  bool        `json:"sharedRoom"`
	Admins      []uuid.UUID `json:"admins"`
}

type KnoqAPI interface {
	GetAll() ([]*EventResponse, error)
	GetByID(id uuid.UUID) (*EventResponse, error)
	GetByUserID(userID uuid.UUID) ([]*EventResponse, error)
}
