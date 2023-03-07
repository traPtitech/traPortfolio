//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"context"
	"time"

	"github.com/traPtitech/traPortfolio/domain"

	"github.com/gofrs/uuid"

	"github.com/traPtitech/traPortfolio/util/optional"
)

type CreateContestArgs struct {
	Name        string
	Description string
	Link        optional.String
	Since       time.Time
	Until       optional.Time
}

type UpdateContestArgs struct {
	Name        optional.String
	Description optional.String
	Link        optional.String
	Since       optional.Time
	Until       optional.Time
}

type CreateContestTeamArgs struct {
	Name        string
	Result      optional.String
	Link        optional.String
	Description string
}

type UpdateContestTeamArgs struct {
	Name        optional.String
	Result      optional.String
	Link        optional.String
	Description optional.String
}

type ContestRepository interface {
	GetContests(ctx context.Context) ([]*domain.Contest, error)
	GetContest(ctx context.Context, contestID uuid.UUID) (*domain.ContestDetail, error)
	CreateContest(ctx context.Context, args *CreateContestArgs) (*domain.ContestDetail, error)
	UpdateContest(ctx context.Context, contestID uuid.UUID, args *UpdateContestArgs) error
	DeleteContest(ctx context.Context, contestID uuid.UUID) error
	GetContestTeams(ctx context.Context, contestID uuid.UUID) ([]*domain.ContestTeam, error)
	GetContestTeam(ctx context.Context, contestID uuid.UUID, teamID uuid.UUID) (*domain.ContestTeamDetail, error)
	CreateContestTeam(ctx context.Context, contestID uuid.UUID, args *CreateContestTeamArgs) (*domain.ContestTeamDetail, error)
	UpdateContestTeam(ctx context.Context, teamID uuid.UUID, args *UpdateContestTeamArgs) error
	DeleteContestTeam(ctx context.Context, contestID uuid.UUID, teamID uuid.UUID) error
	GetContestTeamMembers(ctx context.Context, contestID uuid.UUID, teamID uuid.UUID) ([]*domain.User, error)
	AddContestTeamMembers(ctx context.Context, teamID uuid.UUID, memberIDs []uuid.UUID) error
	EditContestTeamMembers(ctx context.Context, teamID uuid.UUID, memberIDs []uuid.UUID) error
}
