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
	AddAccount(uuid.UUID, *CreateAccountArgs) error
	CreateAccount(uuid.UUID, *CreateAccountArgs) error
	DeleteAccount(uuid.UUID) error
}
