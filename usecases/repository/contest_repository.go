package repository

import (
	"time"

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

type ContestRepository interface {
	//GetAll() ([]*domain.Contest, error)
	//GetByID(ID uuid.UUID) (*domain.Contest, error)
	Create(contest *model.Contest) (*model.Contest, error)
	Update(id uuid.UUID, changes map[string]interface{}) error
}
