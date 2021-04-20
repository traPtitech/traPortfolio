//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
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
	GetProjects() ([]*domain.Project, error)
	GetProject(id uuid.UUID) (*domain.Project, error)
	CreateProject(project *model.Project) (*domain.Project, error)
	UpdateProject(id uuid.UUID, changes map[string]interface{}) error
	GetProjectMembers(id uuid.UUID) ([]*domain.User, error)
}
