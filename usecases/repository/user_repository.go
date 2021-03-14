package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
)

type UserRepository interface {
	GetUsers() ([]*model.User, error)
	GetUser(uuid.UUID) (*model.User, error)
	GetAccounts(uuid.UUID) ([]*model.Account, error)
	Update(*model.User) error
}
