// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package infrastructure

import (
	"github.com/google/wire"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/interfaces/repository"
	"github.com/traPtitech/traPortfolio/usecases/service"
)

// Injectors from wire.go:

func InjectAPIServer(s *SQLConfig, t *TraQConfig, p *PortalConfig, k *KnoQConfig) (handler.API, error) {
	pingHandler := handler.NewPingHandler()
	sqlHandler, err := NewSQLHandler(s)
	if err != nil {
		return handler.API{}, err
	}
	portalAPI, err := NewPortalAPI(p)
	if err != nil {
		return handler.API{}, err
	}
	traQAPI, err := NewTraQAPI(t)
	if err != nil {
		return handler.API{}, err
	}
	userRepository := repository.NewUserRepository(sqlHandler, portalAPI, traQAPI)
	knoqAPI, err := NewKnoqAPI(k)
	if err != nil {
		return handler.API{}, err
	}
	eventRepository := repository.NewEventRepository(sqlHandler, knoqAPI)
	userService := service.NewUserService(userRepository, eventRepository)
	userHandler := handler.NewUserHandler(userService)
	projectRepository := repository.NewProjectRepository(sqlHandler)
	portalRepository := repository.NewPortalRepository(portalAPI)
	projectService := service.NewProjectService(projectRepository, portalRepository)
	projectHandler := handler.NewProjectHandler(projectService)
	eventService := service.NewEventService(eventRepository)
	eventHandler := handler.NewEventHandler(eventService)
	contestRepository := repository.NewContestRepository(sqlHandler, portalRepository)
	contestService := service.NewContestService(contestRepository)
	contestHandler := handler.NewContestHandler(contestService)
	api := handler.NewAPI(pingHandler, userHandler, projectHandler, eventHandler, contestHandler)
	return api, nil
}

// wire.go:

var portalSet = wire.NewSet(
	NewPortalAPI, repository.NewPortalRepository,
)

var traQSet = wire.NewSet(
	NewTraQAPI, repository.NewTraQRepository,
)

var pingSet = wire.NewSet(handler.NewPingHandler)

var userSet = wire.NewSet(repository.NewUserRepository, service.NewUserService, handler.NewUserHandler)

var projectSet = wire.NewSet(repository.NewProjectRepository, service.NewProjectService, handler.NewProjectHandler)

var knoQSet = wire.NewSet(
	NewKnoqAPI, repository.NewKnoqRepository,
)

var eventSet = wire.NewSet(
	knoQSet, repository.NewEventRepository, service.NewEventService, handler.NewEventHandler,
)

var contestSet = wire.NewSet(repository.NewContestRepository, service.NewContestService, handler.NewContestHandler)

var sqlSet = wire.NewSet(
	NewSQLHandler,
)

var apiSet = wire.NewSet(handler.NewAPI)
