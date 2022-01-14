//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type ProjectService interface {
	GetProjects(ctx context.Context) ([]*domain.Project, error)
	GetProject(ctx context.Context, id uuid.UUID) (*domain.Project, error)
	CreateProject(ctx context.Context, args *repository.CreateProjectArgs) (*domain.Project, error)
	UpdateProject(ctx context.Context, id uuid.UUID, args *repository.UpdateProjectArgs) error
	GetProjectMembers(ctx context.Context, id uuid.UUID) ([]*domain.User, error)
	AddProjectMembers(ctx context.Context, projectID uuid.UUID, args []*repository.CreateProjectMemberArgs) error
	DeleteProjectMembers(ctx context.Context, projectID uuid.UUID, memberIDs []uuid.UUID) error
}

type projectService struct {
	repo   repository.ProjectRepository
	portal repository.PortalRepository
}

func NewProjectService(projectRepository repository.ProjectRepository, portalRepository repository.PortalRepository) ProjectService {
	return &projectService{
		repo:   projectRepository,
		portal: portalRepository,
	}
}

func (s *projectService) GetProjects(ctx context.Context) ([]*domain.Project, error) {
	res, err := s.repo.GetProjects()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *projectService) GetProject(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
	project, err := s.repo.GetProject(id)
	if err != nil {
		return nil, err
	}
	portalUsers, err := s.portal.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	NameMap := make(map[string]string, len(portalUsers))
	for _, v := range portalUsers {
		NameMap[v.ID] = v.Name
	}
	for i, v := range project.Members {
		project.Members[i].RealName = NameMap[v.Name]
	}
	return project, nil
}

func (s *projectService) CreateProject(ctx context.Context, args *repository.CreateProjectArgs) (*domain.Project, error) {
	uid := uuid.Must(uuid.NewV4())
	project := &model.Project{
		ID:            uid,
		Name:          args.Name,
		Description:   args.Description,
		SinceYear:     args.SinceYear,
		SinceSemester: args.SinceSemester,
		UntilYear:     args.UntilYear,
		UntilSemester: args.UntilSemester,
	}
	if args.Link.Valid {
		project.Link = args.Link.String
	}
	res, err := s.repo.CreateProject(project)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *projectService) UpdateProject(ctx context.Context, id uuid.UUID, args *repository.UpdateProjectArgs) error {
	changes := map[string]interface{}{}
	if args.Name.Valid {
		changes["name"] = args.Name.String
	}
	if args.Description.Valid {
		changes["description"] = args.Description.String
	}
	if args.Link.Valid {
		changes["link"] = args.Link.String
	}
	if args.SinceYear.Valid && args.SinceSemester.Valid {
		changes["since_year"] = args.SinceYear.Int64
		changes["since_semester"] = args.SinceSemester.Int64
	}
	if args.UntilYear.Valid && args.UntilSemester.Valid {
		changes["until_year"] = args.UntilYear.Int64
		changes["until_semester"] = args.UntilSemester.Int64
	}
	if len(changes) > 0 {
		err := s.repo.UpdateProject(id, changes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *projectService) GetProjectMembers(ctx context.Context, id uuid.UUID) ([]*domain.User, error) {
	members, err := s.repo.GetProjectMembers(id)
	if err != nil {
		return nil, err
	}
	portalUsers, err := s.portal.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	NameMap := make(map[string]string, len(portalUsers))
	for _, v := range portalUsers {
		NameMap[v.ID] = v.Name
	}
	for i, v := range members {
		members[i].RealName = NameMap[v.Name]
	}
	return members, nil
}

func (s *projectService) AddProjectMembers(ctx context.Context, projectID uuid.UUID, args []*repository.CreateProjectMemberArgs) error {
	err := s.repo.AddProjectMembers(projectID, args)
	if err != nil {
		return err
	}
	return nil
}

func (s *projectService) DeleteProjectMembers(ctx context.Context, projectID uuid.UUID, memberIDs []uuid.UUID) error {
	err := s.repo.DeleteProjectMembers(projectID, memberIDs)
	if err != nil {
		return err
	}
	return err
}

// Interface guards
var (
	_ ProjectService = (*projectService)(nil)
)
