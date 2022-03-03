package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/random"
)

type ProjectRepository struct {
	h      database.SQLHandler
	portal external.PortalAPI
}

func NewProjectRepository(h database.SQLHandler, portal external.PortalAPI) repository.ProjectRepository {
	return &ProjectRepository{h, portal}
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
			Duration:    domain.NewYearWithSemesterDuration(v.SinceYear, v.SinceSemester, v.UntilYear, v.UntilSemester),
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

	portalUsers, err := repo.portal.GetAll()
	if err != nil {
		return nil, err
	}

	nameMap := make(map[string]string, len(portalUsers))
	for _, v := range portalUsers {
		nameMap[v.TraQID] = v.RealName
	}

	m := make([]*domain.ProjectMember, len(members))
	for i, v := range members {
		pm := domain.ProjectMember{
			UserID:   v.UserID,
			Name:     v.User.Name,
			Duration: domain.NewYearWithSemesterDuration(v.SinceYear, v.SinceSemester, v.UntilYear, v.UntilSemester),
		}

		if rn, ok := nameMap[v.User.Name]; ok {
			pm.RealName = rn
		}

		m[i] = &pm
	}

	res := &domain.Project{
		ID:          id,
		Name:        project.Name,
		Duration:    domain.NewYearWithSemesterDuration(project.SinceYear, project.SinceSemester, project.UntilYear, project.UntilSemester),
		Description: project.Description,
		Link:        project.Link,
		Members:     m,
	}
	return res, nil
}

func (repo *ProjectRepository) CreateProject(args *repository.CreateProjectArgs) (*domain.Project, error) {
	p := model.Project{
		ID:            random.UUID(),
		Name:          args.Name,
		Description:   args.Description,
		SinceYear:     args.SinceYear,
		SinceSemester: args.SinceSemester,
		UntilYear:     args.UntilYear,
		UntilSemester: args.UntilSemester,
	}
	if args.Link.Valid {
		p.Link = args.Link.String
	}

	err := repo.h.Create(&p).Error()
	if err != nil {
		return nil, convertError(err)
	}

	res := &domain.Project{
		ID:          p.ID,
		Name:        p.Name,
		Duration:    domain.NewYearWithSemesterDuration(p.SinceYear, p.SinceSemester, p.UntilYear, p.UntilSemester),
		Description: p.Description,
		Link:        p.Link,
	}

	return res, nil
}

func (repo *ProjectRepository) UpdateProject(id uuid.UUID, changes map[string]interface{}) error {
	err := repo.h.
		Model(&model.Project{}).
		Where(model.Project{ID: id}).
		Updates(changes).
		Error()
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

	portalUsers, err := repo.portal.GetAll()
	if err != nil {
		return nil, err
	}

	nameMap := make(map[string]string, len(portalUsers))
	for _, v := range portalUsers {
		nameMap[v.TraQID] = v.RealName
	}

	res := make([]*domain.User, len(members))
	for i, v := range members {
		u := domain.User{
			ID:   v.UserID,
			Name: v.User.Name,
		}

		if rn, ok := nameMap[v.User.Name]; ok {
			u.RealName = rn
		}

		res[i] = &u
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
			ID:            uid,
			ProjectID:     projectID,
			UserID:        v.UserID,
			SinceYear:     v.SinceYear,
			SinceSemester: v.SinceSemester,
			UntilYear:     v.UntilYear,
			UntilSemester: v.UntilSemester,
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
	if len(members) == 0 {
		return repository.ErrInvalidArg
	}

	// プロジェクトの存在チェック
	err := repo.h.First(&model.Project{}, &model.Project{ID: projectID}).Error()
	if err != nil {
		return convertError(err)
	}

	err = repo.h.
		Where(&model.ProjectMember{ProjectID: projectID}).
		Where("user_id IN ?", members).
		Delete(&model.ProjectMember{}).
		Error()
	if err != nil {
		return convertError(err)
	}

	return nil
}

// Interface guards
var (
	_ repository.ProjectRepository = (*ProjectRepository)(nil)
)
