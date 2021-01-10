package repository

import "github.com/traPtitech/traPortfolio/domain"

type UserRepository interface {
	Get(string) (*domain.User, []*domain.Account, error)
	Update(*domain.User) (*domain.User, error)
}
