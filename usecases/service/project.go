//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type ProjectService interface {
	GetProjects(ctx context.Context) ([]*domain.Project, error)
	GetProject(ctx context.Context, projectID uuid.UUID) (*domain.ProjectDetail, error)
	CreateProject(ctx context.Context, args *repository.CreateProjectArgs) (*domain.ProjectDetail, error)
	UpdateProject(ctx context.Context, projectID uuid.UUID, args *repository.UpdateProjectArgs) error
	GetProjectMembers(ctx context.Context, projectID uuid.UUID) ([]*domain.UserWithDuration, error)
	AddProjectMembers(ctx context.Context, projectID uuid.UUID, args []*repository.CreateProjectMemberArgs) error
	DeleteProjectMembers(ctx context.Context, projectID uuid.UUID, memberIDs []uuid.UUID) error
}

type projectService struct {
	repo repository.ProjectRepository
}

func NewProjectService(projectRepository repository.ProjectRepository) ProjectService {
	return &projectService{
		repo: projectRepository,
	}
}

func (s *projectService) GetProjects(ctx context.Context) ([]*domain.Project, error) {
	res, err := s.repo.GetProjects()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *projectService) GetProject(ctx context.Context, projectID uuid.UUID) (*domain.ProjectDetail, error) {
	project, err := s.repo.GetProject(projectID)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (s *projectService) CreateProject(ctx context.Context, args *repository.CreateProjectArgs) (*domain.ProjectDetail, error) {
	d := domain.NewYearWithSemesterDuration(args.SinceYear, args.SinceSemester, args.UntilYear, args.UntilSemester)
	if !d.IsValid() {
		return nil, repository.ErrInvalidArg
	}

	res, err := s.repo.CreateProject(args)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *projectService) UpdateProject(ctx context.Context, projectID uuid.UUID, args *repository.UpdateProjectArgs) error {
	old, err := s.repo.GetProject(projectID)
	if err != nil {
		return err
	}

	d := old.Duration
	if args.SinceYear.Valid && args.SinceSemester.Valid {
		d.Since.Year = int(args.SinceYear.V)
		d.Since.Semester = int(args.SinceSemester.V)
	}

	if args.UntilYear.Valid && args.UntilSemester.Valid {
		d.Until.Year = int(args.UntilYear.V)
		d.Until.Semester = int(args.UntilSemester.V)
	}

	if !d.IsValid() {
		return repository.ErrInvalidArg
	}

	if err := s.repo.UpdateProject(projectID, args); err != nil {
		return err
	}

	return nil
}

func (s *projectService) GetProjectMembers(ctx context.Context, projectID uuid.UUID) ([]*domain.UserWithDuration, error) {
	members, err := s.repo.GetProjectMembers(projectID)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (s *projectService) AddProjectMembers(ctx context.Context, projectID uuid.UUID, args []*repository.CreateProjectMemberArgs) error {
	for _, v := range args {
		d := domain.NewYearWithSemesterDuration(v.SinceYear, v.SinceSemester, v.UntilYear, v.UntilSemester)
		if !d.IsValid() {
			return repository.ErrInvalidArg
		}
	}

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
