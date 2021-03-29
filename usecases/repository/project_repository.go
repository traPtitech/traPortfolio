package repository

import (
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
)

type ProjectRepository interface {
	PostProject(*domain.ProjectDetail) (*model.Project, error)
}
