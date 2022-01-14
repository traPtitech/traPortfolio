//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/util/optional"
)

type CreateProjectArgs struct {
	Name          string
	Description   string
	Link          optional.String
	SinceYear     int
	SinceSemester int
	UntilYear     int
	UntilSemester int
}

type UpdateProjectArgs struct {
	Name          optional.String
	Description   optional.String
	Link          optional.String
	SinceYear     optional.Int64
	SinceSemester optional.Int64
	UntilYear     optional.Int64
	UntilSemester optional.Int64
}

type CreateProjectMemberArgs struct {
	UserID        uuid.UUID
	SinceYear     int
	SinceSemester int
	UntilYear     int
	UntilSemester int
}

type ProjectRepository interface {
	GetProjects() ([]*domain.Project, error)
	GetProject(id uuid.UUID) (*domain.Project, error)
	CreateProject(project *model.Project) (*domain.Project, error)
	UpdateProject(id uuid.UUID, changes map[string]interface{}) error
	GetProjectMembers(id uuid.UUID) ([]*domain.User, error)
	AddProjectMembers(id uuid.UUID, args []*CreateProjectMemberArgs) error
	DeleteProjectMembers(id uuid.UUID, memberIDs []uuid.UUID) error
}
