package service

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
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

func (s *ProjectService) GetProjects(ctx context.Context) ([]*model.Project, error) {
	res, err := s.repo.GetProjects()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *ProjectService) GetProject(ctx context.Context, id uuid.UUID) (*model.ProjectDetail, error) {
	project, err := s.repo.GetProject(id)
	if err != nil {
		return nil, err
	}
	projectMembers, err := s.repo.GetProjectMembers(id)
	if err != nil {
		return nil, err
	}
	portalUsers, err := s.portal.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	members := make([]*model.ProjectMemberDetail, 0, len(projectMembers))
	for _, v := range projectMembers {
		for _, pu := range portalUsers {
			if v.Name == pu.ID {
				v.RealName = pu.Name
				members = append(members, v)
			}
		}
	}
	res := &model.ProjectDetail{
		ID:          project.ID,
		Name:        project.Name,
		Link:        project.Link,
		Description: project.Description,
		Members:     members,
		Since:       project.Since,
		Until:       project.Until,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}
	return res, nil
}

func (s *ProjectService) CreateProject(ctx context.Context, args *repository.CreateProjectArgs) (*model.Project, error) {
	uid := uuid.Must(uuid.NewV4())
	project := &model.Project{
		ID:          uid,
		Name:        args.Name,
		Description: args.Description,
		Link:        args.Link,
		Since:       args.Since,
		Until:       args.Until,
	}
	project, err := s.repo.CreateProject(project)
	if err != nil {
		return nil, err
	}
	return project, nil
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
