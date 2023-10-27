package infrastructure

import (
	"github.com/traPtitech/traPortfolio/infrastructure/external"
	"github.com/traPtitech/traPortfolio/infrastructure/external/mock_external_e2e"
	"github.com/traPtitech/traPortfolio/infrastructure/repository"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/usecases/service"
	"github.com/traPtitech/traPortfolio/util/config"
	"gorm.io/gorm"
)

func InjectAPIServer(c *config.Config, db *gorm.DB) (handler.API, error) {
	// external API
	var (
		portalAPI external.PortalAPI
		traQAPI   external.TraQAPI
		knoqAPI   external.KnoqAPI
	)

	if c.IsProduction {
		var err error

		portalAPI, err = external.NewPortalAPI(c.PortalConf())
		if err != nil {
			return handler.API{}, err
		}

		traQAPI, err = external.NewTraQAPI(c.TraqConf())
		if err != nil {
			return handler.API{}, err
		}

		knoqAPI, err = external.NewKnoqAPI(c.KnoqConf())
		if err != nil {
			return handler.API{}, err
		}
	} else {
		portalAPI = mock_external_e2e.NewMockPortalAPI()
		traQAPI = mock_external_e2e.NewMockTraQAPI()
		knoqAPI = mock_external_e2e.NewMockKnoqAPI()
	}

	// repository
	userRepo := repository.NewUserRepository(db, portalAPI, traQAPI)
	projectRepo := repository.NewProjectRepository(db, portalAPI)
	eventRepo := repository.NewEventRepository(db, knoqAPI)
	contestRepo := repository.NewContestRepository(db, portalAPI)
	groupRepo := repository.NewGroupRepository(db)

	// service, handler, API
	api := handler.NewAPI(
		handler.NewPingHandler(),
		handler.NewUserHandler(service.NewUserService(userRepo, eventRepo)),
		handler.NewProjectHandler(service.NewProjectService(projectRepo)),
		handler.NewEventHandler(service.NewEventService(eventRepo, userRepo)),
		handler.NewContestHandler(service.NewContestService(contestRepo)),
		handler.NewGroupHandler(service.NewGroupService(groupRepo, userRepo)),
	)

	return api, nil
}
