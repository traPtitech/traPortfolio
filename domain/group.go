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
	Link        string
	Admin       []*User
	Members     []*GroupMember
	Description string
}

// GroupMember indicates User who belongs to Group
type GroupMember struct {
	ID       uuid.UUID // User ID
	Name     string    // User Name
	RealName string
	Duration YearWithSemesterDuration
}
