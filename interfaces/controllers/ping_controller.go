package controllers

import (
	"github.com/traPtitech/traPortfolio/interfaces/repository"
	"github.com/traPtitech/traPortfolio/usecase"
)

type PingController struct {
	Interactor usecase.PingInteractor
}

func NewPingController() *PingController {
	return &PingController{
		Interactor: usecase.PingInteractor{
			PingRepository: &repository.PingRepository{},
		},
	}
}

func (controller *PingController) Ping(c Context) (err error) {
	controller.Interactor.Ping()
	c.String(200, "pong")
	return
}
