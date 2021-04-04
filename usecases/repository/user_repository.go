package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/util/optional"
)

type UpdateUserArgs struct {
	Description optional.String
	Check       optional.Bool
}

type UserRepository interface {
	GetUsers() ([]*domain.User, error)
	GetUser(uuid.UUID) (*domain.UserDetail, error)
	GetAccounts(uuid.UUID) ([]*domain.Account, error)
	Update(id uuid.UUID, changes map[string]interface{}) error
}
