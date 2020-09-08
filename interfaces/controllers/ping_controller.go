package controllers

import "github.com/traPtitech/traPortfolio/usecases/usecase"

type PingController struct {
	Interactor usecase.PingUsecase
}

func NewPingController(it usecase.PingUsecase) *PingController {
	return &PingController{
		Interactor: it,
	}
}

func (controller *PingController) Ping(c Context) (err error) {
	err = controller.Interactor.Ping()
	if err != nil {
		return err
	}
	return c.String(200, "pong")
}
