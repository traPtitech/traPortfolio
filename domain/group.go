package domain

import (
	"github.com/gofrs/uuid"
)

type GroupUser struct {
	ID       uuid.UUID
	name     string
	duration ProjectDuration
}

type Groups struct {
	ID   uuid.UUID
	Name string
}

type GroupDetail struct {
	ID      uuid.UUID
	name    string
	link    string
	leader  *User
	Members []*UserGroup
}
