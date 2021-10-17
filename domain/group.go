package domain

import (
	"github.com/gofrs/uuid"
)

type GroupUser struct {
	ID       uuid.UUID
	Name     string
	Duration ProjectDuration
}

type Groups struct {
	ID   uuid.UUID
	Name string
}

type GroupDetail struct {
	ID          uuid.UUID
	Name        string
	Link        string
	Leader      *User
	Members     []*UserGroup
	Description string
}
