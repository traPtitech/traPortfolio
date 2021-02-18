package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
)

type UserRepository interface {
	GetUsers() ([]*domain.User, error)
	GetUser(uuid.UUID) (*domain.User, error)
	GetAccounts(uuid.UUID) ([]*domain.Account, error)
	Update(uuid.UUID, *domain.EditUserRequest) (*domain.User, error)
}
