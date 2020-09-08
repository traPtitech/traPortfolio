package interactor

import "log"

type PingInteractor struct {
}

func NewPingInteractor() *PingInteractor {
	return &PingInteractor{}
}

func (interactor *PingInteractor) Ping() (err error) {
	log.Println("ping received")
	return
}
