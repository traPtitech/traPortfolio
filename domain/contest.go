package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type Contest struct {
	ID        uuid.UUID
	Name      string
	TimeStart time.Time
	TimeEnd   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ContestDetail struct {
	Contest
	Link        string
	Description string
	Teams       []*ContestTeam
}

type ContestTeam struct {
	ID        uuid.UUID
	ContestID uuid.UUID
	Name      string
	Result    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ContestTeamDetail struct {
	ContestTeam
	Link        string
	Description string
	Members     []*User
}
