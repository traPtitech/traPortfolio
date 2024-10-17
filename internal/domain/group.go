package domain

import (
	"github.com/gofrs/uuid"
)

type Group struct {
	ID   uuid.UUID
	Name string
}

type GroupDetail struct {
	ID          uuid.UUID
	Name        string
	Links       []string
	Admin       []*User
	Members     []*UserWithDuration
	Description string
}
