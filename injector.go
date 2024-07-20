package main

import (
	"github.com/traPtitech/traPortfolio/internal/handler"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/external"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/external/mock_external_e2e"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/repository"
	"github.com/traPtitech/traPortfolio/internal/pkgs/config"
	"gorm.io/gorm"
)

func injectIntoAPIServer(c *config.Config, db *gorm.DB) (handler.API, error) {
	// external API
	var (
		portalAPI external.PortalAPI
		traQAPI   external.TraQAPI
		knoqAPI   external.KnoqAPI
	)

	// TODO: 初期リリースではPortalとknoQとは連携しない
	if c.IsProduction {
		var err error

		// portalAPI, err = external.NewPortalAPI(c.Portal)
		// if err != nil {
		// 	return handler.API{}, err
		// }
		portalAPI = external.NewNopPortalAPI()

		traQAPI, err = external.NewTraQAPI(c.Traq)
		if err != nil {
			return handler.API{}, err
		}

		knoqAPI, err = external.NewKnoqAPI(c.Knoq)
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
		handler.NewUserHandler(userRepo, eventRepo),
		handler.NewProjectHandler(projectRepo),
		handler.NewEventHandler(eventRepo, userRepo),
		handler.NewContestHandler(contestRepo),
		handler.NewGroupHandler(groupRepo, userRepo),
	)

	return api, nil
}
