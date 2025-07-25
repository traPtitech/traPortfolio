//go:generate go run go.uber.org/mock/mockgen@latest -typed -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/internal/domain"
	"github.com/traPtitech/traPortfolio/internal/pkgs/optional"
)

type GetProjectsArgs struct {
	Limit optional.Of[int]
}

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

type EditProjectMemberArgs struct {
	UserID        uuid.UUID
	SinceYear     int
	SinceSemester int
	UntilYear     int
	UntilSemester int
}

type ProjectRepository interface {
	GetProjects(ctx context.Context, args *GetProjectsArgs) ([]*domain.Project, error)
	GetProject(ctx context.Context, projectID uuid.UUID) (*domain.ProjectDetail, error)
	CreateProject(ctx context.Context, args *CreateProjectArgs) (*domain.ProjectDetail, error)
	UpdateProject(ctx context.Context, projectID uuid.UUID, args *UpdateProjectArgs) error
	DeleteProject(ctx context.Context, projectID uuid.UUID) error
	GetProjectMembers(ctx context.Context, projectID uuid.UUID) ([]*domain.UserWithDuration, error)
	EditProjectMembers(ctx context.Context, projectID uuid.UUID, args []*EditProjectMemberArgs) error
}
