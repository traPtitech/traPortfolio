package interactor

import "github.com/traPtitech/traPortfolio/usecase/repository"

type PingInteractor struct {
	PingRepository repository.PingRepository
}

func (interactor *PingInteractor) Ping() (err error) {
	err = interactor.PingRepository.Ping()
	return
}
