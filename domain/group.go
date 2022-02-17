package domain

import (
	"github.com/gofrs/uuid"
)

// GroupUser indicates Group which User belongs
type GroupUser struct {
	ID       uuid.UUID // Group ID
	Name     string    // Group name
	Duration YearWithSemesterDuration
}

type GroupDuration struct {
	Since YearWithSemester
	Until YearWithSemester
}

type Group struct {
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
