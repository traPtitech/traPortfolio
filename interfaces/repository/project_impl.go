package repository

import (
	"context"

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

func (r *ProjectRepository) GetProjects(ctx context.Context) ([]*domain.Project, error) {
	projects := make([]*model.Project, 0)
	err := r.h.WithContext(ctx).Find(&projects).Error()
	if err != nil {
		return nil, err
	}
	res := make([]*domain.Project, 0, len(projects))
	for _, v := range projects {
		p := &domain.Project{
			ID:       v.ID,
			Name:     v.Name,
			Duration: domain.NewYearWithSemesterDuration(v.SinceYear, v.SinceSemester, v.UntilYear, v.UntilSemester),
		}
		res = append(res, p)
	}
	return res, nil
}

func (r *ProjectRepository) GetProject(ctx context.Context, projectID uuid.UUID) (*domain.ProjectDetail, error) {
	project := new(model.Project)
	if err := r.h.
		WithContext(ctx).
		Where(&model.Project{ID: projectID}).
		First(project).
		Error(); err != nil {
		return nil, err
	}

	members := make([]*model.ProjectMember, 0)
	err := r.h.
		WithContext(ctx).
		Preload("User").
		Where(model.ProjectMember{ProjectID: projectID}).
		Find(&members).
		Error()
	if err != nil {
		return nil, err
	}

	portalUsers, err := r.portal.GetPortalUsers()
	if err != nil {
		return nil, err
	}

	nameMap := make(map[string]string, len(portalUsers))
	for _, v := range portalUsers {
		nameMap[v.TraQID] = v.RealName
	}

	m := make([]*domain.UserWithDuration, len(members))
	for i, v := range members {
		realName := nameMap[v.User.Name]
		pm := domain.UserWithDuration{
			User: *domain.NewUser(
				v.User.ID,
				v.User.Name,
				realName,
				v.User.Check,
			),
			Duration: domain.NewYearWithSemesterDuration(v.SinceYear, v.SinceSemester, v.UntilYear, v.UntilSemester),
		}

		m[i] = &pm
	}

	res := &domain.ProjectDetail{
		Project: domain.Project{
			ID:       projectID,
			Name:     project.Name,
			Duration: domain.NewYearWithSemesterDuration(project.SinceYear, project.SinceSemester, project.UntilYear, project.UntilSemester),
		},
		Description: project.Description,
		Link:        project.Link,
		Members:     m,
	}
	return res, nil
}

func (r *ProjectRepository) CreateProject(ctx context.Context, args *repository.CreateProjectArgs) (*domain.ProjectDetail, error) {
	p := model.Project{
		ID:            random.UUID(),
		Name:          args.Name,
		Description:   args.Description,
		SinceYear:     args.SinceYear,
		SinceSemester: args.SinceSemester,
		UntilYear:     args.UntilYear,
		UntilSemester: args.UntilSemester,
	}
	p.Link = args.Link.ValueOr(p.Link)

	err := r.h.WithContext(ctx).Create(&p).Error()
	if err != nil {
		return nil, err
	}

	res := &domain.ProjectDetail{
		Project: domain.Project{
			ID:       p.ID,
			Name:     p.Name,
			Duration: domain.NewYearWithSemesterDuration(p.SinceYear, p.SinceSemester, p.UntilYear, p.UntilSemester),
		},
		Description: p.Description,
		Link:        p.Link,
	}

	return res, nil
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, projectID uuid.UUID, args *repository.UpdateProjectArgs) error {
	changes := map[string]interface{}{}
	if v, ok := args.Name.V(); ok {
		changes["name"] = v
	}
	if v, ok := args.Description.V(); ok {
		changes["description"] = v
	}
	if v, ok := args.Link.V(); ok {
		changes["link"] = v
	}
	if sy, ok := args.SinceYear.V(); ok {
		if ss, ok := args.SinceSemester.V(); ok {
			changes["since_year"] = sy
			changes["since_semester"] = ss
		}
	}
	if uy, ok := args.UntilYear.V(); ok {
		if us, ok := args.UntilSemester.V(); ok {
			changes["until_year"] = uy
			changes["until_semester"] = us
		}
	}

	if len(changes) == 0 {
		return nil
	}

	err := r.h.
		WithContext(ctx).
		Model(&model.Project{}).
		Where(&model.Project{ID: projectID}).
		Updates(changes).
		Error()
	if err != nil {
		return err
	}

	return nil
}

func (r *ProjectRepository) GetProjectMembers(ctx context.Context, projectID uuid.UUID) ([]*domain.UserWithDuration, error) {
	members := make([]*model.ProjectMember, 0)
	err := r.h.
		WithContext(ctx).
		Preload("User").
		Where(&model.ProjectMember{ProjectID: projectID}).
		Find(&members).
		Error()
	if err != nil {
		return nil, err
	}

	portalUsers, err := r.portal.GetPortalUsers()
	if err != nil {
		return nil, err
	}

	nameMap := make(map[string]string, len(portalUsers))
	for _, v := range portalUsers {
		nameMap[v.TraQID] = v.RealName
	}

	res := make([]*domain.UserWithDuration, len(members))
	for i, v := range members {
		realName := nameMap[v.User.Name]
		u := domain.UserWithDuration{
			User: *domain.NewUser(
				v.User.ID,
				v.User.Name,
				realName,
				v.User.Check,
			),
			Duration: domain.NewYearWithSemesterDuration(v.SinceYear, v.SinceSemester, v.UntilYear, v.UntilSemester),
		}

		res[i] = &u
	}

	return res, nil
}

func (r *ProjectRepository) AddProjectMembers(ctx context.Context, projectID uuid.UUID, projectMembers []*repository.CreateProjectMemberArgs) error {
	if len(projectMembers) == 0 {
		return repository.ErrInvalidArg
	}

	// ユーザーの重複チェック
	projectMembersMap := make(map[uuid.UUID]struct{}, len(projectMembers))
	for _, v := range projectMembers {
		if _, ok := projectMembersMap[v.UserID]; ok {
			return repository.ErrInvalidArg
		}
		projectMembersMap[v.UserID] = struct{}{}
	}

	// プロジェクトの存在チェック
	err := r.h.
		WithContext(ctx).
		Where(&model.Project{ID: projectID}).
		First(&model.Project{}).
		Error()
	if err != nil {
		return err
	}

	mmbsMp := make(map[uuid.UUID]struct{}, len(projectMembers))
	_mmbs := make([]*model.ProjectMember, 0, len(projectMembers))
	err = r.h.
		WithContext(ctx).
		Where(&model.ProjectMember{ProjectID: projectID}).
		Find(&_mmbs).
		Error()
	if err != nil {
		return err
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

	err = r.h.WithContext(ctx).Transaction(func(tx database.SQLHandler) error {
		for _, v := range members {
			if _, ok := mmbsMp[v.UserID]; ok {
				continue
			}
			err = tx.WithContext(ctx).Create(v).Error()
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *ProjectRepository) DeleteProjectMembers(ctx context.Context, projectID uuid.UUID, members []uuid.UUID) error {
	if len(members) == 0 {
		return repository.ErrInvalidArg
	}

	// プロジェクトの存在チェック
	err := r.h.
		WithContext(ctx).
		Where(&model.Project{ID: projectID}).
		First(&model.Project{}).
		Error()
	if err != nil {
		return err
	}

	err = r.h.
		WithContext(ctx).
		Where(&model.ProjectMember{ProjectID: projectID}).
		Where("`project_members`.`user_id` IN (?)", members).
		Delete(&model.ProjectMember{}).
		Error()
	if err != nil {
		return err
	}

	return nil
}

// Interface guards
var (
	_ repository.ProjectRepository = (*ProjectRepository)(nil)
)
