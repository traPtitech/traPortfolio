package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
)

type UserRepository interface {
	GetUsers() ([]*domain.User, error)
	GetUser(uuid.UUID) (*domain.UserDetail, error)
	GetAccounts(uuid.UUID) ([]*domain.Account, error)
	Update(*domain.EditUser) error
}
