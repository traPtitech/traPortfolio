package repository

import (
	"time"

	"github.com/traPtitech/traPortfolio/domain"

	"github.com/traPtitech/traPortfolio/interfaces/repository/model"

	"github.com/gofrs/uuid"

	"github.com/traPtitech/traPortfolio/util/optional"
)

type CreateContestArgs struct {
	Name        string
	Description string
	Link        string
	Since       time.Time
	Until       time.Time
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
	Result      string
	Link        string
	Description string
}

type UpdateContestTeamArgs struct {
	Name        optional.String
	Result      optional.String
	Link        optional.String
	Description optional.String
}

type ContestRepository interface {
	//GetAll() ([]*domain.Contest, error)
	//GetByID(ID uuid.UUID) (*domain.Contest, error)
	CreateContest(contest *model.Contest) (*model.Contest, error)
	UpdateContest(id uuid.UUID, changes map[string]interface{}) error
	CreateContestTeam(contestID uuid.UUID, args *CreateContestTeamArgs) (*domain.ContestTeamDetail, error)
	UpdateContestTeam(teamID uuid.UUID, changes map[string]interface{}) error
	AddContestTeamMember(teamID uuid.UUID, memberIDs []uuid.UUID) error
	DeleteContestTeamMember(teamID uuid.UUID, memberIDs []uuid.UUID) error
}
