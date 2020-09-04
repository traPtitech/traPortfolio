package controllers

import (
	"github.com/traPtitech/traPortfolio/interfaces/repository"
	"github.com/traPtitech/traPortfolio/usecase/interactor"
)

type PingController struct {
	Interactor interactor.PingInteractor
}

func NewPingController() *PingController {
	return &PingController{
		Interactor: interactor.PingInteractor{
			PingRepository: &repository.PingRepository{},
		},
	}
}

func (controller *PingController) Ping(c Context) (err error) {
	controller.Interactor.Ping()
	c.String(200, "pong")
	return
}
