package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
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
