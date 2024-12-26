package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/internal/domain"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/external"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/repository/model"
	"github.com/traPtitech/traPortfolio/internal/pkgs/random"
	"github.com/traPtitech/traPortfolio/internal/usecases/repository"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	h      *gorm.DB
	portal external.PortalAPI
}

func NewProjectRepository(h *gorm.DB, portal external.PortalAPI) repository.ProjectRepository {
	return &ProjectRepository{h, portal}
}

func (r *ProjectRepository) GetProjects(ctx context.Context) ([]*domain.Project, error) {
	projects := make([]*model.Project, 0)
	err := r.h.WithContext(ctx).Find(&projects).Error
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
		Error; err != nil {
		return nil, err
	}

	members := make([]*model.ProjectMember, 0)
	err := r.h.
		WithContext(ctx).
		Preload("User").
		Where(model.ProjectMember{ProjectID: projectID}).
		Find(&members).
		Error
	if err != nil {
		return nil, err
	}

	realNameMap, err := external.GetRealNameMap(r.portal)
	if err != nil {
		return nil, err
	}

	m := make([]*domain.UserWithDuration, len(members))
	for i, v := range members {
		pm := domain.UserWithDuration{
			User: *domain.NewUser(
				v.User.ID,
				v.User.Name,
				realNameMap[v.User.Name],
				v.User.DisplayName,
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

	// 既に同名のプロジェクトが存在するか
	err := r.h.
		WithContext(ctx).
		Where(&model.Project{Name: p.Name}).
		First(&model.Project{}).
		Error
	if err == nil {
		return nil, repository.ErrAlreadyExists
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	err = r.h.WithContext(ctx).Create(&p).Error
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
		Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ProjectRepository) DeleteProject(ctx context.Context, projectID uuid.UUID) error {
	err := r.h.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.
			WithContext(ctx).
			Where(&model.Project{ID: projectID}).
			First(&model.Project{}).
			Error
		if err != nil {
			return err
		}

		err = tx.
			WithContext(ctx).
			Where(&model.Project{ID: projectID}).
			Delete(&model.Project{}).
			Error
		if err != nil {
			return err
		}

		err = tx.
			WithContext(ctx).
			Where(&model.ProjectMember{ProjectID: projectID}).
			Delete(&model.ProjectMember{}).
			Error
		if err != nil {
			return err
		}
		return nil
	})
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
		Error
	if err != nil {
		return nil, err
	}

	realNameMap, err := external.GetRealNameMap(r.portal)
	if err != nil {
		return nil, err
	}

	res := make([]*domain.UserWithDuration, len(members))
	for i, v := range members {
		u := domain.UserWithDuration{
			User: *domain.NewUser(
				v.User.ID,
				v.User.Name,
				realNameMap[v.User.Name],
				v.User.DisplayName,
				v.User.Check,
			),
			Duration: domain.NewYearWithSemesterDuration(v.SinceYear, v.SinceSemester, v.UntilYear, v.UntilSemester),
		}

		res[i] = &u
	}

	return res, nil
}

func (r *ProjectRepository) EditProjectMembers(ctx context.Context, projectID uuid.UUID, projectMembers []*repository.EditProjectMemberArgs) error {
	p := model.Project{}

	// プロジェクトの存在チェック
	err := r.h.
		WithContext(ctx).
		Where(&model.Project{ID: projectID}).
		First(&p).
		Error
	if err != nil {
		return err
	}

	projectDuration := domain.NewYearWithSemesterDuration(p.SinceYear, p.SinceSemester, p.UntilYear, p.UntilSemester)

	// プロジェクトの期間内かどうか
	for _, v := range projectMembers {
		memberDuration := domain.NewYearWithSemesterDuration(v.SinceYear, v.SinceSemester, v.UntilYear, v.UntilSemester)
		if !projectDuration.Includes(memberDuration) {
			return fmt.Errorf("%w: exceeded duration user(project: %+v, member: %+v)", repository.ErrInvalidArg, projectDuration, memberDuration)
		}
	}

	currentProjectMembers := make([]*model.ProjectMember, 0, len(projectMembers))
	err = r.h.
		WithContext(ctx).
		Where(&model.ProjectMember{ProjectID: projectID}).
		Find(&currentProjectMembers).
		Error
	if err != nil && err != repository.ErrNotFound {
		return err
	}

	currentProjectMembersMap := make(map[uuid.UUID]*model.ProjectMember, len(projectMembers))
	for _, v := range currentProjectMembers {
		currentProjectMembersMap[v.UserID] = &model.ProjectMember{
			SinceYear:     v.SinceYear,
			SinceSemester: v.SinceSemester,
			UntilYear:     v.UntilYear,
			UntilSemester: v.UntilSemester,
		}
	}

	members := make([]*model.ProjectMember, 0, len(projectMembers))
	for _, v := range projectMembers {
		m := &model.ProjectMember{
			ProjectID:     projectID,
			UserID:        v.UserID,
			SinceYear:     v.SinceYear,
			SinceSemester: v.SinceSemester,
			UntilYear:     v.UntilYear,
			UntilSemester: v.UntilSemester,
		}
		members = append(members, m)
	}

	err = r.h.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, v := range members {
			// 既に登録されていたら更新を試し、そうでなければ新規作成
			if vdb, ok := currentProjectMembersMap[v.UserID]; ok {
				changes := map[string]interface{}{}
				if v.SinceYear != vdb.SinceYear {
					changes["since_year"] = v.SinceYear
				}
				if v.SinceSemester != vdb.SinceSemester {
					changes["since_semester"] = v.SinceSemester
				}
				if v.UntilYear != vdb.UntilYear {
					changes["until_year"] = v.UntilYear
				}
				if v.UntilSemester != vdb.UntilSemester {
					changes["until_semester"] = v.UntilSemester
				}
				if len(changes) > 0 {
					err = tx.WithContext(ctx).
						Model(&model.ProjectMember{}).
						Where(&model.ProjectMember{ProjectID: projectID, UserID: v.UserID}).
						Updates(changes).
						Error
					if err != nil {
						return err
					}
				}
				delete(currentProjectMembersMap, v.UserID)
				continue
			}
			err = tx.WithContext(ctx).Create(v).Error
			if err != nil {
				return err
			}
		}
		// 残っているものは削除
		for _, member := range currentProjectMembers {
			if _, ok := currentProjectMembersMap[member.UserID]; !ok {
				continue
			}
			err = tx.WithContext(ctx).
				Where(&model.ProjectMember{ProjectID: projectID, UserID: member.UserID}).
				Delete(&model.ProjectMember{}).
				Error
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

// Interface guards
var (
	_ repository.ProjectRepository = (*ProjectRepository)(nil)
)
