package repository

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/util/optional"
)

type CreateProjectArgs struct {
	Name        string
	Description string
	Link        string
	Since       time.Time
	Until       time.Time
}

type UpdateProjectArgs struct {
	Name        optional.String
	Description optional.String
	Link        optional.String
	Since       optional.Time
	Until       optional.Time
}

type ProjectRepository interface {
	GetProjects() ([]*model.Project, error)
	GetProject(id uuid.UUID) (*model.Project, error)
	GetProjectMembers(id uuid.UUID) ([]*model.ProjectMemberDetail, error)
	CreateProject(project *model.Project) (*model.Project, error)
	UpdateProject(id uuid.UUID, changes map[string]interface{}) error
}
