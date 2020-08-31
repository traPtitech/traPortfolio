package usecase

type PingInteractor struct {
	PingRepository PingRepository
}

func (interactor *PingInteractor) Ping() (err error) {
	err = interactor.PingRepository.Ping()
	return
}
