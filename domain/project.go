package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type Project struct {
	ID       uuid.UUID
	Name     string
	Duration ProjectDuration
}

type ProjectDuration struct {
	Since YearWithSemester
	Until YearWithSemester
}

type YearWithSemester struct {
	Year     uint
	Semester uint
}

type ProjectDetail struct {
	ID          uuid.UUID
	Name        string
	Duration    ProjectDuration
	Link        string
	Description string
	Members     []*ProjectMember
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProjectMember struct {
	ID       uuid.UUID
	Name     string
	RealName string
	Duration ProjectDuration
}
