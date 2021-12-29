package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type ProjectRepository struct {
	h database.SQLHandler
}

func NewProjectRepository(sql database.SQLHandler) repository.ProjectRepository {
	return &ProjectRepository{h: sql}
}

func (repo *ProjectRepository) GetProjects() ([]*domain.Project, error) {
	projects := make([]*model.Project, 0)
	err := repo.h.Find(&projects).Error()
	if err != nil {
		return nil, convertError(err)
	}
	res := make([]*domain.Project, 0, len(projects))
	for _, v := range projects {
		p := &domain.Project{
			ID:          v.ID,
			Name:        v.Name,
			Since:       v.Since,
			Until:       v.Until,
			Description: v.Description,
			Link:        v.Link,
		}
		res = append(res, p)
	}
	return res, nil
}

func (repo *ProjectRepository) GetProject(id uuid.UUID) (*domain.Project, error) {
	project := new(model.Project)
	if err := repo.h.First(project, &model.Project{ID: id}).Error(); err != nil {
		return nil, convertError(err)
	}

	members := make([]*model.ProjectMember, 0)
	err := repo.h.
		Preload("User").
		Where(model.ProjectMember{ProjectID: id}).
		Find(&members).
		Error()
	if err != nil {
		return nil, convertError(err)
	}

	m := make([]*domain.ProjectMember, 0, len(members))
	for _, v := range members {
		m = append(m, &domain.ProjectMember{
			UserID: v.UserID,
			Name:   v.User.Name,
			Since:  v.Since,
			Until:  v.Until,
		})
	}
	res := &domain.Project{
		ID:          id,
		Name:        project.Name,
		Since:       project.Since,
		Until:       project.Until,
		Description: project.Description,
		Link:        project.Link,
		Members:     m,
	}
	return res, nil
}

func (repo *ProjectRepository) CreateProject(project *model.Project) (*domain.Project, error) {
	err := repo.h.Create(project).Error()
	if err != nil {
		return nil, convertError(err)
	}
	res := &domain.Project{
		ID:          project.ID,
		Name:        project.Name,
		Since:       project.Since,
		Until:       project.Until,
		Description: project.Description,
		Link:        project.Link,
	}
	return res, nil
}

func (repo *ProjectRepository) UpdateProject(id uuid.UUID, changes map[string]interface{}) error {
	var (
		old model.Project
		new model.Project
	)

	err := repo.h.Transaction(func(tx database.SQLHandler) error {
		if err := tx.First(&old, model.Project{ID: id}).Error(); err != nil {
			return convertError(err)
		}
		if err := tx.Model(&old).Updates(changes).Error(); err != nil {
			return convertError(err)
		}
		if err := tx.Where(&model.Project{ID: id}).First(&new).Error(); err != nil {
			return convertError(err)
		}

		return nil
	})
	if err != nil {
		return convertError(err)
	}

	return nil
}

func (repo *ProjectRepository) GetProjectMembers(id uuid.UUID) ([]*domain.User, error) {
	members := make([]*model.ProjectMember, 0)
	err := repo.h.
		Preload("User").
		Where(&model.ProjectMember{ProjectID: id}).
		Find(&members).
		Error()
	if err != nil {
		return nil, convertError(err)
	}

	res := make([]*domain.User, 0, len(members))
	for _, v := range members {
		res = append(res, &domain.User{
			ID:   v.UserID,
			Name: v.User.Name,
		})
	}

	return res, nil
}

func (repo *ProjectRepository) AddProjectMembers(projectID uuid.UUID, projectMembers []*repository.CreateProjectMemberArgs) error {
	if len(projectMembers) == 0 {
		return repository.ErrInvalidArg
	}

	// プロジェクトの存在チェック
	err := repo.h.First(&model.Project{}, &model.Project{ID: projectID}).Error()
	if err != nil {
		return convertError(err)
	}

	mmbsMp := make(map[uuid.UUID]struct{}, len(projectMembers))
	_mmbs := make([]*model.ProjectMember, 0, len(projectMembers))
	err = repo.h.Where(&model.ProjectMember{ProjectID: projectID}).Find(&_mmbs).Error()
	if err != nil {
		return convertError(err)
	}
	for _, v := range _mmbs {
		mmbsMp[v.UserID] = struct{}{}
	}

	members := make([]*model.ProjectMember, 0, len(projectMembers))
	for _, v := range projectMembers {
		uid := uuid.Must(uuid.NewV4())
		m := &model.ProjectMember{
			ID:        uid,
			ProjectID: projectID,
			UserID:    v.UserID,
			Since:     v.Since,
			Until:     v.Until,
		}
		members = append(members, m)
	}

	err = repo.h.Transaction(func(tx database.SQLHandler) error {
		for _, v := range members {
			if _, ok := mmbsMp[v.UserID]; ok {
				continue
			}
			err = tx.Create(v).Error()
			if err != nil {
				return convertError(err)
			}
		}
		return nil
	})
	if err != nil {
		return convertError(err)
	}

	return nil
}

func (repo *ProjectRepository) DeleteProjectMembers(projectID uuid.UUID, members []uuid.UUID) error {
	// 存在チェック
	err := repo.h.First(&model.ProjectMember{}, &model.ProjectMember{ProjectID: projectID}).Error()
	if err != nil {
		return convertError(err)
	}

	mmbsMp := make(map[uuid.UUID]struct{}, len(members))
	_mmbs := make([]*model.ProjectMember, 0, len(members))
	err = repo.h.Where(&model.ProjectMember{ProjectID: projectID}).Find(&_mmbs).Error()
	if err != nil {
		return convertError(err)
	}
	for _, v := range _mmbs {
		mmbsMp[v.UserID] = struct{}{}
	}

	err = repo.h.Transaction(func(tx database.SQLHandler) error {
		for _, memberID := range members {
			if _, ok := mmbsMp[memberID]; ok {
				err = tx.Delete(&model.ProjectMember{}, &model.ProjectMember{ProjectID: projectID, UserID: memberID}).Error()
				if err != nil {
					return convertError(err)
				}
			}
		}
		return nil
	})
	if err != nil {
		return convertError(err)
	}
	return nil
}

// Interface guards
var (
	_ repository.ProjectRepository = (*ProjectRepository)(nil)
)
