package domain

import (
	"time"

	"github.com/gofrs/uuid"
)

type Project struct {
	ID          uuid.UUID
	Name        string
	Since       time.Time
	Until       time.Time
	Description string
	Link        string
	Members     []*ProjectMember
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProjectMember struct {
	UserID   uuid.UUID
	Name     string
	RealName string
	Since    time.Time
	Until    time.Time
}

type ProjectDuration struct {
	Since YearWithSemester
	Until YearWithSemester
}

type YearWithSemester struct {
	Year     uint
	Semester uint
}
