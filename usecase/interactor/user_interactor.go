package interactor

import (
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecase/input"
	"github.com/traPtitech/traPortfolio/usecase/repository"
)

type UserInteractor struct {
	UserRepository repository.UserRepository
}

func NewUserInteractor(repo repository.UserRepository) UserInteractor {
	return UserInteractor{UserRepository: repo}
}

func (interactor *UserInteractor) UserByID(ipt input.GetUser) (user domain.User, err error) {
	user, err = interactor.UserRepository.FindByID(ipt.ID)
	return
}

func (interactor *UserInteractor) Users() (users []domain.User, err error) {
	users, err = interactor.UserRepository.FindAll()
	return
}

func (interactor *UserInteractor) Add(ipt input.AddUser) (user domain.User, err error) {
	user, err = interactor.UserRepository.Store(ipt.User)
	return
}

func (interactor *UserInteractor) Update(ipt input.UpdateUser) (user domain.User, err error) {
	user, err = interactor.UserRepository.Update(ipt.User)
	return
}

func (interactor *UserInteractor) DeleteByID(ipt input.DeleteUser) (err error) {
	err = interactor.UserRepository.DeleteByID(ipt.ID)
	return
}
