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
}

type ContestDetail struct {
	Contest
	Link         string
	Description  string
	ContestTeams []*ContestTeam
}

type ContestTeamWithoutMembers struct {
	ID        uuid.UUID
	ContestID uuid.UUID
	Name      string
	Result    string
}

type ContestTeam struct {
	ContestTeamWithoutMembers
	Members []*User
}

type ContestTeamDetail struct {
	ContestTeam
	Link        string
	Description string
}
