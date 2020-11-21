package repository

import (
	"time"

	"github.com/gofrs/uuid"
)

type KnoQEvent struct {
	ID          uuid.UUID
	Name        string
	Description string
	GroupID     uuid.UUID
	RoomID      uuid.UUID
	TimeStart   time.Time
	TimeEnd     time.Time
	SharedRoom  bool
}

type KnoqRepository interface {
	GetAll() ([]*KnoQEvent, error)
	GetByID(id uuid.UUID) (*KnoQEvent, error)
}
