//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"context"
	"time"

	"github.com/traPtitech/traPortfolio/internal/domain"

	"github.com/gofrs/uuid"

	"github.com/traPtitech/traPortfolio/internal/util/optional"
)

type CreateContestArgs struct {
	Name        string
	Description string
	Link        optional.Of[string]
	Since       time.Time
	Until       optional.Of[time.Time]
}

type UpdateContestArgs struct {
	Name        optional.Of[string]
	Description optional.Of[string]
	Link        optional.Of[string]
	Since       optional.Of[time.Time]
	Until       optional.Of[time.Time]
}

type CreateContestTeamArgs struct {
	Name        string
	Result      optional.Of[string]
	Link        optional.Of[string]
	Description string
}

type UpdateContestTeamArgs struct {
	Name        optional.Of[string]
	Result      optional.Of[string]
	Link        optional.Of[string]
	Description optional.Of[string]
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
