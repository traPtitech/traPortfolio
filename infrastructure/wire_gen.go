// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package infrastructure

import (
	"github.com/google/wire"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/interfaces/repository"
	repository2 "github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/usecases/service"
)

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Injectors from wire.go:

func InjectAPIServer(traQToken repository.TraQToken, portalToken repository.PortalToken) (handler.API, error) {
	pingHandler := handler.NewPingHandler()
	sqlHandler, err := NewSQLHandler()
	if err != nil {
		return handler.API{}, err
	}
	userRepository := repository.NewUserRepository(sqlHandler)
	traQAPI, err := NewTraQAPI()
	if err != nil {
		return handler.API{}, err
	}
	traQRepository := repository.NewTraQRepository(traQAPI, traQToken)
	portalRepository := repository.NewPortalRepository(portalToken)
	userService := service.NewUserService(userRepository, traQRepository, portalRepository)
	userHandler := handler.NewUserHandler(userService)
	projectRepository := repository.NewProjectRepository(sqlHandler)
	projectService := service.NewProjectService(projectRepository, traQRepository)
	projectHandler := handler.NewProjectHandler(projectService)
	knoqAPI, err := NewKnoqAPI()
	if err != nil {
		return handler.API{}, err
	}
	eventRepository := repository.NewEventRepository(sqlHandler, knoqAPI)
	knoqRepository := repository.NewKnoqRepository(knoqAPI)
	eventService := service.NewEventService(eventRepository, knoqRepository)
	eventHandler := handler.NewEventHandler(eventService)
	contestRepository := repository.NewContestRepository(sqlHandler)
	contestService := service.NewContestService(contestRepository)
	contestHandler := handler.NewContestHandler(contestService)
	api := handler.NewAPI(pingHandler, userHandler, projectHandler, eventHandler, contestHandler)
	return api, nil
}

// wire.go:

var portalSet = wire.NewSet(repository.NewPortalRepository, wire.Bind(new(repository2.PortalRepository), new(*repository.PortalRepository)))

var traQSet = wire.NewSet(
	NewTraQAPI, repository.NewTraQRepository, wire.Bind(new(external.TraQAPI), new(*TraQAPI)), wire.Bind(new(repository2.TraQRepository), new(*repository.TraQRepository)),
)

var pingSet = wire.NewSet(handler.NewPingHandler)

var userSet = wire.NewSet(repository.NewUserRepository, service.NewUserService, handler.NewUserHandler, wire.Bind(new(repository2.UserRepository), new(*repository.UserRepository)))

var projectSet = wire.NewSet(repository.NewProjectRepository, service.NewProjectService, handler.NewProjectHandler, wire.Bind(new(repository2.ProjectRepository), new(*repository.ProjectRepository)))

var knoQSet = wire.NewSet(
	NewKnoqAPI, repository.NewKnoqRepository, wire.Bind(new(external.KnoqAPI), new(*KnoqAPI)), wire.Bind(new(repository2.KnoqRepository), new(*repository.KnoqRepository)),
)

var eventSet = wire.NewSet(
	knoQSet, repository.NewEventRepository, service.NewEventService, handler.NewEventHandler, wire.Bind(new(repository2.EventRepository), new(*repository.EventRepository)),
)

var contestSet = wire.NewSet(repository.NewContestRepository, service.NewContestService, handler.NewContestHandler, wire.Bind(new(repository2.ContestRepository), new(*repository.ContestRepository)))

var sqlSet = wire.NewSet(
	NewSQLHandler, wire.Bind(new(database.SQLHandler), new(*SQLHandler)),
)

var apiSet = wire.NewSet(handler.NewAPI)
