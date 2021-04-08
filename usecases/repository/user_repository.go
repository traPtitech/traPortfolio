package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
)

type CreateAccountArgs struct {
	ID          string
	Type        uint
	URL         string
	PrPermitted bool
}

type UserRepository interface {
	GetUsers() ([]*domain.User, error)
	GetUser(uuid.UUID) (*domain.User, error)
	GetAccounts(uuid.UUID) ([]*domain.Account, error)
	Update(*domain.User) error
	CreateAccount(uuid.UUID, *CreateAccountArgs) (*domain.Account, error)
	DeleteAccount(uuid.UUID, uuid.UUID) error
}
