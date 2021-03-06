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
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.CreatedAt,
		}
		res = append(res, p)
	}
	return res, nil
}

func (repo *ProjectRepository) GetProject(id uuid.UUID) (*domain.Project, error) {
	project := &model.Project{ID: id}
	if err := repo.h.First(&project).Error(); err != nil {
		return nil, convertError(err)
	}

	members := make([]*domain.ProjectMember, 0)
	selectQuery := "users.id as user_id, users.name, project_members.since, project_members.until"
	whereQuery := "project_members.project_id = ?"
	joinQuery := "left join users on users.id = project_members.user_id"
	err := repo.h.Model(&model.ProjectMember{}).Select(selectQuery).Where(whereQuery, id).Joins(joinQuery).Scan(&members).Error()
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
		Members:     members,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.CreatedAt,
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
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.CreatedAt,
	}
	return res, nil
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

func (repo *ProjectRepository) GetProjectMembers(id uuid.UUID) ([]*domain.User, error) {
	members := make([]*domain.User, 0)
	selectQuery := "users.id as id, users.name"
	whereQuery := "project_members.project_id = ?"
	joinQuery := "left join users on users.id = project_members.user_id"
	err := repo.h.Model(&model.ProjectMember{}).Select(selectQuery).Where(whereQuery, id).Joins(joinQuery).Scan(&members).Error()
	if err != nil {
		return nil, convertError(err)
	}
	return members, nil
}

func (repo *ProjectRepository) AddProjectMembers(projectID uuid.UUID, projectMembers []*repository.CreateProjectMemberArgs) error {
	if projectID == uuid.Nil {
		return repository.ErrInvalidID
	}

	//存在チェック
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
		if v.UserID == uuid.Nil {
			return repository.ErrInvalidID
		}
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
	return nil
}

func (repo *ProjectRepository) DeleteProjectMembers(projectID uuid.UUID, members []uuid.UUID) error {
	if projectID == uuid.Nil {
		return repository.ErrNilID
	}

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
