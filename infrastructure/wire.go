//+build wireinject

package infrastructure

import (
	"github.com/google/wire"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
	impl "github.com/traPtitech/traPortfolio/interfaces/repository"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/usecases/service"
)

var portalSet = wire.NewSet(
	NewPortalAPI,
	impl.NewPortalRepository,
	wire.Bind(new(external.PortalAPI), new(*PortalAPI)),
	wire.Bind(new(repository.PortalRepository), new(*impl.PortalRepository)),
)

var traQSet = wire.NewSet(
	NewTraQAPI,
	impl.NewTraQRepository,
	wire.Bind(new(external.TraQAPI), new(*TraQAPI)),
	wire.Bind(new(repository.TraQRepository), new(*impl.TraQRepository)),
)

var pingSet = wire.NewSet(
	handler.NewPingHandler,
)

var userSet = wire.NewSet(
	impl.NewUserRepository,
	service.NewUserService,
	handler.NewUserHandler,
	wire.Bind(new(repository.UserRepository), new(*impl.UserRepository)),
)

var projectSet = wire.NewSet(
	impl.NewProjectRepository,
	service.NewProjectService,
	handler.NewProjectHandler,
	wire.Bind(new(repository.ProjectRepository), new(*impl.ProjectRepository)),
)

var knoQSet = wire.NewSet(
	NewKnoqAPI,
	impl.NewKnoqRepository,
	wire.Bind(new(external.KnoqAPI), new(*KnoqAPI)),
	wire.Bind(new(repository.KnoqRepository), new(*impl.KnoqRepository)),
)

var eventSet = wire.NewSet(
	knoQSet,
	impl.NewEventRepository,
	service.NewEventService,
	handler.NewEventHandler,
	wire.Bind(new(repository.EventRepository), new(*impl.EventRepository)),
)

var contestSet = wire.NewSet(
	impl.NewContestRepository,
	service.NewContestService,
	handler.NewContestHandler,
	wire.Bind(new(repository.ContestRepository), new(*impl.ContestRepository)),
)

var sqlSet = wire.NewSet(
	NewSQLHandler,
	wire.Bind(new(database.SQLHandler), new(*SQLHandler)),
)

var apiSet = wire.NewSet(handler.NewAPI)

func InjectAPIServer() (handler.API, error) {
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
