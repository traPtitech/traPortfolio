//go:generate go run github.com/google/wire/cmd/wire@v0.5.0
//go:build wireinject

package infrastructure

import (
	"github.com/google/wire"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	impl "github.com/traPtitech/traPortfolio/interfaces/repository"
	"github.com/traPtitech/traPortfolio/usecases/service"
)

var portalSet = wire.NewSet(
	NewPortalAPI,
	impl.NewPortalRepository,
)

var traQSet = wire.NewSet(
	NewTraQAPI,
	impl.NewTraQRepository,
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

var knoQSet = wire.NewSet(
	NewKnoqAPI,
	impl.NewKnoqRepository,
)

var eventSet = wire.NewSet(
	knoQSet,
	impl.NewEventRepository,
	service.NewEventService,
	handler.NewEventHandler,
)

var contestSet = wire.NewSet(
	impl.NewContestRepository,
	service.NewContestService,
	handler.NewContestHandler,
)

var sqlSet = wire.NewSet(
	NewSQLHandler,
)

var apiSet = wire.NewSet(handler.NewAPI)

func InjectAPIServer(s *SQLConfig, t *TraQConfig, p *PortalConfig, k *KnoQConfig) (handler.API, error) {
	wire.Build(
		pingSet,
		userSet,
		projectSet,
		eventSet,
		sqlSet,
		apiSet,
		portalSet,
		traQSet,
		contestSet,
	)
	return handler.API{}, nil
}
