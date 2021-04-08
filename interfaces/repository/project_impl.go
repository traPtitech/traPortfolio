package repository

import (
	"fmt"

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
	if err := repo.h.First(&project).Error(); err != nil {
		return nil, err
	}
	return project, nil
}

func (repo *ProjectRepository) GetProjectMembers(id uuid.UUID) ([]*model.ProjectMemberDetail, error) {
	members := make([]*model.ProjectMemberDetail, 0)
	selectQuery := "project_members.project_id, users.id as user_id, users.name, project_members.since, project_members.until"
	joinQuery := "left join users on users.id = project_members.user_id"
	err := repo.h.Model(&model.ProjectMember{}).Select(selectQuery).Where("project_id = ?", id).Joins(joinQuery).Scan(&members).Error()
	fmt.Printf("%#v\n\n", *members[0])
	if err != nil {
		return nil, err
	}
	return members, nil
}

func (repo *ProjectRepository) CreateProject(project *model.Project) (*model.Project, error) {
	err := repo.h.Create(project).Error()
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (repo *ProjectRepository) UpdateProject(id uuid.UUID, changes map[string]interface{}) error {
	if id == uuid.Nil {
		return repository.ErrNilID
	}

	var (
		old model.Project
		new model.Project
	)

	tx := repo.h.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error(); err != nil {
		return err
	}
	if err := tx.First(&old, model.Project{ID: id}).Error(); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Model(&old).Updates(changes).Error(); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Where(&model.Project{ID: id}).First(&new).Error(); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// Interface guards
var (
	_ repository.ProjectRepository = (*ProjectRepository)(nil)
)
