//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/util/optional"
)

type CreateProjectArgs struct {
	Name          string
	Description   string
	Link          optional.Of[string]
	SinceYear     int
	SinceSemester int
	UntilYear     int
	UntilSemester int
}

type UpdateProjectArgs struct {
	Name          optional.Of[string]
	Description   optional.Of[string]
	Link          optional.Of[string]
	SinceYear     optional.Of[int64]
	SinceSemester optional.Of[int64]
	UntilYear     optional.Of[int64]
	UntilSemester optional.Of[int64]
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
	GetProject(projectID uuid.UUID) (*domain.ProjectDetail, error)
	CreateProject(args *CreateProjectArgs) (*domain.ProjectDetail, error)
	UpdateProject(projectID uuid.UUID, args *UpdateProjectArgs) error
	GetProjectMembers(projectID uuid.UUID) ([]*domain.UserWithDuration, error)
	AddProjectMembers(projectID uuid.UUID, args []*CreateProjectMemberArgs) error
	DeleteProjectMembers(projectID uuid.UUID, memberIDs []uuid.UUID) error
}
