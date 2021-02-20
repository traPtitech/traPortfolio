package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
)

type EditUserRequest struct {
	Bio          string `json:"bio"`
	HideRealName bool   `json:"hideRealName"`
}

type UserRepository interface {
	GetUsers() ([]*domain.User, error)
	GetUser(uuid.UUID) (*domain.User, error)
	GetAccounts(uuid.UUID) ([]*domain.Account, error)
	Update(uuid.UUID, *EditUserRequest) (*domain.User, error)
}
