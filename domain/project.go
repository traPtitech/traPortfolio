package domain

import (
	"github.com/gofrs/uuid"
)

type Project struct {
	ID       uuid.UUID
	Name     string
	Duration YearWithSemesterDuration
}

type ProjectDetail struct {
	Project
	Description string
	Link        string
	Members     []*UserWithDuration
}
