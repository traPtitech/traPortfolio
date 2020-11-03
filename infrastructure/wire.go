//+build wireinject

package infrastructure

import (
	"github.com/google/wire"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	impl "github.com/traPtitech/traPortfolio/interfaces/repository"
	"github.com/traPtitech/traPortfolio/usecases/handler"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	service "github.com/traPtitech/traPortfolio/usecases/service/user_service"
)

var portalSet = wire.NewSet(
	impl.NewPortalRepository,
	wire.Bind(new(repository.PortalRepository), new(*impl.PortalRepository)),
)

var traQSet = wire.NewSet(
	impl.NewTraQRepository,
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

var eventSet = wire.NewSet(
	NewKnoqAPI,
	impl.NewEventRepository,
	handler.NewEventHandler,
	wire.Bind(new(external.KnoqAPI), new(*KnoqAPI)),
	wire.Bind(new(repository.EventRepository), new(*impl.EventRepository)),
)

var sqlSet = wire.NewSet(
	NewSQLHandler,
	wire.Bind(new(database.SQLHandler), new(*SQLHandler)),
)

var apiSet = wire.NewSet(handler.NewAPI)

func InjectAPIServer(traQToken impl.TraQToken, portalToken impl.PortalToken) (handler.API, error) {
	wire.Build(
		pingSet,
		userSet,
		eventSet,
		sqlSet,
		apiSet,
		portalSet,
		traQSet,
	)
	return handler.API{}, nil
}
