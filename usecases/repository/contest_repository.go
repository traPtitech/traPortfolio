//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

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
	Link        optional.Of[string]
	Since       time.Time
	Until       optional.Time
}

type UpdateContestArgs struct {
	Name        optional.Of[string]
	Description optional.Of[string]
	Link        optional.Of[string]
	Since       optional.Time
	Until       optional.Time
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
	GetContests() ([]*domain.Contest, error)
	GetContest(contestID uuid.UUID) (*domain.ContestDetail, error)
	CreateContest(args *CreateContestArgs) (*domain.ContestDetail, error)
	UpdateContest(contestID uuid.UUID, args *UpdateContestArgs) error
	DeleteContest(contestID uuid.UUID) error
	GetContestTeams(contestID uuid.UUID) ([]*domain.ContestTeam, error)
	GetContestTeam(contestID uuid.UUID, teamID uuid.UUID) (*domain.ContestTeamDetail, error)
	CreateContestTeam(contestID uuid.UUID, args *CreateContestTeamArgs) (*domain.ContestTeamDetail, error)
	UpdateContestTeam(teamID uuid.UUID, args *UpdateContestTeamArgs) error
	DeleteContestTeam(contestID uuid.UUID, teamID uuid.UUID) error
	GetContestTeamMembers(contestID uuid.UUID, teamID uuid.UUID) ([]*domain.User, error)
	AddContestTeamMembers(teamID uuid.UUID, memberIDs []uuid.UUID) error
	EditContestTeamMembers(teamID uuid.UUID, memberIDs []uuid.UUID) error
}
