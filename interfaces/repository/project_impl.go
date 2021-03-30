package repository

import (
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
)

type ProjectRepository struct {
	h database.SQLHandler
}

func NewProjectRepository(sql database.SQLHandler) *ProjectRepository {
	return &ProjectRepository{h: sql}
}

func (repo *ProjectRepository) Create(project *model.Project) (*model.Project, error) {
	err := repo.h.Create(project).Error()
	if err != nil {
		return nil, err
	}
	return project, nil
}
