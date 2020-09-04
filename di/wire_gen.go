// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package di

import (
	"github.com/google/wire"
	"github.com/traPtitech/traPortfolio/interfaces/controllers"
	"github.com/traPtitech/traPortfolio/interfaces/repository"
	"github.com/traPtitech/traPortfolio/usecase/interactor"
)

// Injectors from wire.go:

func InjectAPIServer() controllers.API {
	pingRepository := repository.NewPingRepository()
	pingInteractor := interactor.NewPingInteractor(pingRepository)
	pingController := controllers.NewPingController(pingInteractor)
	api := controllers.NewAPI(pingController)
	return api
}

// wire.go:

var repositorySet = wire.NewSet(repository.NewPingRepository)

var interactorSet = wire.NewSet(interactor.NewPingInteractor)

var controllerSet = wire.NewSet(controllers.NewPingController)

var apiSet = wire.NewSet(controllers.NewAPI)
