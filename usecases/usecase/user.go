package usecase

import (
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/input"
)

type UserUsecase interface {
	UserByID(ipt input.GetUser) (user domain.User, err error)
	Users() (users []domain.User, err error)
	Add(ipt input.AddUser) (user domain.User, err error)
	Update(ipt input.UpdateUser) (user domain.User, err error)
	DeleteByID(ipt input.DeleteUser) (err error)
}
