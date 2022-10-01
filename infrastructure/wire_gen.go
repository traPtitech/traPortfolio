// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package infrastructure

import (
	"github.com/google/wire"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/interfaces/repository"
	"github.com/traPtitech/traPortfolio/usecases/service"
	"github.com/traPtitech/traPortfolio/util/config"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func InjectAPIServer(c *config.Config, db *gorm.DB) (handler.API, error) {
	pingHandler := handler.NewPingHandler()
	sqlHandler := FromDB(db)
	portalConfig := providePortalConf(c)
	bool2 := provideIsDevelopMent(c)
	portalAPI, err := NewPortalAPI(portalConfig, bool2)
	if err != nil {
		return handler.API{}, err
	}
	traqConfig := provideTraqConf(c)
	traQAPI, err := NewTraQAPI(traqConfig, bool2)
	if err != nil {
		return handler.API{}, err
	}
	userRepository := repository.NewUserRepository(sqlHandler, portalAPI, traQAPI)
	knoqConfig := provideKnoqConf(c)
	knoqAPI, err := NewKnoqAPI(knoqConfig, bool2)
	if err != nil {
		return handler.API{}, err
	}
	eventRepository := repository.NewEventRepository(sqlHandler, knoqAPI)
	userService := service.NewUserService(userRepository, eventRepository)
	userHandler := handler.NewUserHandler(userService)
	projectRepository := repository.NewProjectRepository(sqlHandler, portalAPI)
	projectService := service.NewProjectService(projectRepository)
	projectHandler := handler.NewProjectHandler(projectService)
	eventService := service.NewEventService(eventRepository, userRepository)
	eventHandler := handler.NewEventHandler(eventService)
	contestRepository := repository.NewContestRepository(sqlHandler, portalAPI)
	contestService := service.NewContestService(contestRepository)
	contestHandler := handler.NewContestHandler(contestService)
	groupRepository := repository.NewGroupRepository(sqlHandler)
	groupService := service.NewGroupService(groupRepository, userRepository)
	groupHandler := handler.NewGroupHandler(groupService)
	api := handler.NewAPI(pingHandler, userHandler, projectHandler, eventHandler, contestHandler, groupHandler)
	return api, nil
}

// wire.go:

var pingSet = wire.NewSet(handler.NewPingHandler)

var userSet = wire.NewSet(repository.NewUserRepository, service.NewUserService, handler.NewUserHandler)

var projectSet = wire.NewSet(repository.NewProjectRepository, service.NewProjectService, handler.NewProjectHandler)

var eventSet = wire.NewSet(repository.NewEventRepository, service.NewEventService, handler.NewEventHandler)

var groupSet = wire.NewSet(repository.NewGroupRepository, service.NewGroupService, handler.NewGroupHandler)

var contestSet = wire.NewSet(repository.NewContestRepository, service.NewContestService, handler.NewContestHandler)

var sqlSet = wire.NewSet(
	FromDB,
)

var externalSet = wire.NewSet(
	NewKnoqAPI,
	NewPortalAPI,
	NewTraQAPI,
)

var apiSet = wire.NewSet(handler.NewAPI)

var confSet = wire.NewSet(
	provideIsDevelopMent,
	provideTraqConf,
	provideKnoqConf,
	providePortalConf,
)

func provideIsDevelopMent(c *config.Config) bool { return c.IsDevelopment() }

func provideTraqConf(c *config.Config) *config.TraqConfig { return c.TraqConf() }

func provideKnoqConf(c *config.Config) *config.KnoqConfig { return c.KnoqConf() }

func providePortalConf(c *config.Config) *config.PortalConfig { return c.PortalConf() }
