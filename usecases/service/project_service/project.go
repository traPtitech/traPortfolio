package service

import (
	"context"

	"github.com/traPtitech/traPortfolio/domain"
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

func (s *ProjectService) PostProject(ctx context.Context, p *domain.ProjectDetail) (*model.Project, error) {
	return s.repo.PostProject(p)
}
