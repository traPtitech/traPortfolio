package service

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/traPtitech/traPortfolio/interfaces/repository/model"

	"github.com/gofrs/uuid"

	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type ProjectService struct {
	repo repository.ProjectRepository
	traQ repository.TraQRepository
}

func NewProjectService(projectRepository repository.ProjectRepository, traQRepository repository.TraQRepository) ProjectService {
	return ProjectService{
		repo: projectRepository,
		traQ: traQRepository,
	}
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
	project, err := s.repo.Create(project)
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
		err := s.repo.Update(id, changes)
		if err != nil && err == repository.ErrNotFound {
			return echo.NewHTTPError(http.StatusNotFound)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
