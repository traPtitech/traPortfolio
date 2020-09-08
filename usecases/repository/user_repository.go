package repository

import "github.com/traPtitech/traPortfolio/domain"

type UserRepository interface {
	FindByID(id int) (*domain.User, error)
	FindAll() ([]domain.User, error)
	Store(*domain.User) (*domain.User, error)
	Update(*domain.User) (*domain.User, error)
	DeleteByID(id int) error
}
