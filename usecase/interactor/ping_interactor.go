package interactor

import "github.com/traPtitech/traPortfolio/usecase/repository"

type PingInteractor struct {
	PingRepository repository.PingRepository
}

func NewPingInteractor(repo repository.PingRepository) PingInteractor {
	return PingInteractor{PingRepository: repo}
}

func (interactor *PingInteractor) Ping() (err error) {
	err = interactor.PingRepository.Ping()
	return
}
