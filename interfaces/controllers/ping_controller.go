package controllers

import (
	"github.com/traPtitech/traPortfolio/usecase/interactor"
)

type PingController struct {
	Interactor interactor.PingInteractor
}

func NewPingController(it interactor.PingInteractor) *PingController {
	return &PingController{
		Interactor: it,
	}
}

func (controller *PingController) Ping(c Context) (err error) {
	controller.Interactor.Ping()
	c.String(200, "pong")
	return
}
