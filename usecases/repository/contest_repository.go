//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
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
	GetContests() ([]*domain.Contest, error)
	GetContest(id uuid.UUID) (*domain.ContestDetail, error)
	CreateContest(args *CreateContestArgs) (*domain.Contest, error)
	UpdateContest(id uuid.UUID, args *UpdateContestArgs) error
	DeleteContest(id uuid.UUID) error
	GetContestTeams(contestID uuid.UUID) ([]*domain.ContestTeam, error)
	GetContestTeam(contestID uuid.UUID, teamID uuid.UUID) (*domain.ContestTeamDetail, error)
	CreateContestTeam(contestID uuid.UUID, args *CreateContestTeamArgs) (*domain.ContestTeamDetail, error)
	UpdateContestTeam(teamID uuid.UUID, args *UpdateContestTeamArgs) error
	DeleteContestTeam(contestID uuid.UUID, teamID uuid.UUID) error
	GetContestTeamMembers(contestID uuid.UUID, teamID uuid.UUID) ([]*domain.User, error)
	AddContestTeamMembers(teamID uuid.UUID, memberIDs []uuid.UUID) error
	DeleteContestTeamMembers(teamID uuid.UUID, memberIDs []uuid.UUID) error
}
