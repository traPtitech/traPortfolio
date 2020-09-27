//+build wireinject

package infrastructure

import (
	"github.com/google/wire"
	impl "github.com/traPtitech/traPortfolio/interfaces/repository"
	"github.com/traPtitech/traPortfolio/usecases/handler"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	service "github.com/traPtitech/traPortfolio/usecases/service/user_service"
	"github.com/traPtitech/traPortfolio/usecases/usecase"
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
	wire.Bind(new(usecase.PingUsecase), new(*handler.PingHandler)),
)

var userSet = wire.NewSet(
	impl.NewUserRepository,
	service.NewUserService,
	handler.NewUserHandler,
	wire.Bind(new(repository.UserRepository), new(*impl.UserRepository)),
	wire.Bind(new(usecase.UserUsecase), new(*handler.UserHandler)),
)

var sqlSet = wire.NewSet(NewSQLHandler)

var apiSet = wire.NewSet(handler.NewAPI)

func InjectAPIServer(traQToken impl.TraQToken, portalToken impl.PortalToken) (handler.API, error) {
	wire.Build(
		pingSet,
		userSet,
		sqlSet,
		apiSet,
		portalSet,
		traQSet,
	)
	return handler.API{}, nil
}
