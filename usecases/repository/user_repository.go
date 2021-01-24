package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
)

type UserRepository interface {
	Get(uuid.UUID) (*domain.User, []*domain.Account, error)
	Update(*domain.User) (*domain.User, error)
}
