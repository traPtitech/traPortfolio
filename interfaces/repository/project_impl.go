package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type ProjectRepository struct {
	h database.SQLHandler
}

func NewProjectRepository(sql database.SQLHandler) *ProjectRepository {
	return &ProjectRepository{h: sql}
}

//TODO GetAllの方がいい？
func (repo *ProjectRepository) GetProjects() ([]*model.Project, error) {
	projects := []*model.Project{}
	err := repo.h.Find(&projects).Error()
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (repo *ProjectRepository) GetProject(id uuid.UUID) (*model.Project, error) {
	project := &model.Project{ID: id}
	err := repo.h.First(&project).Error()
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (repo *ProjectRepository) Create(project *model.Project) (*model.Project, error) {
	err := repo.h.Create(project).Error()
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (repo *ProjectRepository) Update(id uuid.UUID, changes map[string]interface{}) error {
	if id == uuid.Nil {
		return repository.ErrNilID
	}

	var (
		old model.Project
		new model.Project
	)

	tx := repo.h.Begin()
	if err := tx.First(&old, model.Project{ID: id}).Error(); err != nil {
		return err
	}
	if err := tx.Model(&old).Updates(changes).Error(); err != nil {
		return err
	}
	if err := tx.Where(&model.Project{ID: id}).First(&new).Error(); err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// Interface guards
var (
	_ repository.ProjectRepository = (*ProjectRepository)(nil)
)
