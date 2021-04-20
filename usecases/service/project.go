package service

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type ProjectService struct {
	repo   repository.ProjectRepository
	portal repository.PortalRepository
}

func NewProjectService(projectRepository repository.ProjectRepository, portalRepository repository.PortalRepository) ProjectService {
	return ProjectService{
		repo:   projectRepository,
		portal: portalRepository,
	}
}

func (s *ProjectService) GetProjects(ctx context.Context) ([]*domain.Project, error) {
	res, err := s.repo.GetProjects()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ProjectService) GetProject(ctx context.Context, id uuid.UUID) (*domain.Project, error) {
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

func (s *ProjectService) CreateProject(ctx context.Context, args *repository.CreateProjectArgs) (*domain.Project, error) {
	uid := uuid.Must(uuid.NewV4())
	project := &model.Project{
		ID:          uid,
		Name:        args.Name,
		Description: args.Description,
		Link:        args.Link,
		Since:       args.Since,
		Until:       args.Until,
	}
	res, err := s.repo.CreateProject(project)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, id uuid.UUID, args *repository.UpdateProjectArgs) error {
	if id == uuid.Nil {
		return repository.ErrInvalidID
	}
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
	if args.Since.Valid {
		changes["since"] = args.Since.Time
	}
	if args.Until.Valid {
		changes["until"] = args.Until.Time
	}
	if len(changes) > 0 {
		err := s.repo.UpdateProject(id, changes)
		if err != nil && err == repository.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ProjectService) GetProjectMembers(ctx context.Context, id uuid.UUID) ([]*domain.User, error) {
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
