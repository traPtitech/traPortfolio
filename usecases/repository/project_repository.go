package repository

import (
	"time"

	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
)

type CreateProjectArgs struct {
	Name        string
	Description string
	Link        string
	Since       time.Time
	Until       time.Time
}

type ProjectRepository interface {
	Create(*model.Project) (*model.Project, error)
}
