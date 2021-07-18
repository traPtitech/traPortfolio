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
	GetAllGroups() ([]*GroupsResponse, error)
}

type GroupsResponse struct {
	GroupID     uuid.UUID   `json:"groupId"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Open        bool        `json:"open"`
	IsTraQGroup bool        `json:"isTraQGroup"`
	Members     []uuid.UUID `json:"members"`
	Admins      []uuid.UUID `json:"admins"`
	CreatedBy   uuid.UUID   `json:"createdBy"`
	CreatedAt   string      `json:"createdAt"`
	UpdatedAt   string      `json:"updatedAt"`
}
