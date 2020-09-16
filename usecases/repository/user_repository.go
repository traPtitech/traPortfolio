package repository

import "github.com/traPtitech/traPortfolio/domain"

type UserRepository interface {
	Update(*domain.User) (*domain.User, error)
	DeleteByID(id int) error
}
