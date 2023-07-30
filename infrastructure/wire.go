//go:generate go run github.com/google/wire/cmd/wire@latest
//go:build wireinject

package infrastructure

import (
	"github.com/google/wire"
	impl "github.com/traPtitech/traPortfolio/infrastructure/repository"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	"github.com/traPtitech/traPortfolio/usecases/service"
	"github.com/traPtitech/traPortfolio/util/config"
	"gorm.io/gorm"
)

var pingSet = wire.NewSet(
	handler.NewPingHandler,
)

var userSet = wire.NewSet(
	impl.NewUserRepository,
	service.NewUserService,
	handler.NewUserHandler,
)

var projectSet = wire.NewSet(
	impl.NewProjectRepository,
	service.NewProjectService,
	handler.NewProjectHandler,
)

var eventSet = wire.NewSet(
	impl.NewEventRepository,
	service.NewEventService,
	handler.NewEventHandler,
)

var groupSet = wire.NewSet(
	impl.NewGroupRepository,
	service.NewGroupService,
	handler.NewGroupHandler,
)

var contestSet = wire.NewSet(
	impl.NewContestRepository,
	service.NewContestService,
	handler.NewContestHandler,
)

var sqlSet = wire.NewSet(
	NewSQLHandler,
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

func provideIsDevelopMent(c *config.Config) bool              { return !c.IsProduction }
func provideTraqConf(c *config.Config) *config.TraqConfig     { return c.TraqConf() }
func provideKnoqConf(c *config.Config) *config.KnoqConfig     { return c.KnoqConf() }
func providePortalConf(c *config.Config) *config.PortalConfig { return c.PortalConf() }

func InjectAPIServer(c *config.Config, db *gorm.DB) (handler.API, error) {
	wire.Build(
		pingSet,
		userSet,
		projectSet,
		eventSet,
		groupSet,
		contestSet,
		sqlSet,
		externalSet,
		apiSet,
		confSet,
	)
	return handler.API{}, nil
}
