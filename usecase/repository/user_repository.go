package repository

import "github.com/traPtitech/traPortfolio/domain"

type UserRepository interface {
	FindById(id int) (domain.User, error)
	FindAll() (domain.User, error)
	Store(domain.User) (domain.User, error)
	Update(domain.User) (domain.User, error)
	DeleteById(id int) error
}
