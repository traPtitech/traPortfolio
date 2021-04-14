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

type CreateAccountArgs struct {
	ID          string
	Type        uint
	URL         string
	PrPermitted bool
}

type UserRepository interface {
	GetUsers() ([]*domain.User, error)
	GetUser(uuid.UUID) (*domain.UserDetail, error)
	GetAccounts(uuid.UUID) ([]*domain.Account, error)
	Update(id uuid.UUID, changes map[string]interface{}) error
	CreateAccount(uuid.UUID, *CreateAccountArgs) (*domain.Account, error)
	DeleteAccount(uuid.UUID, uuid.UUID) error
}
