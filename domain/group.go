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
