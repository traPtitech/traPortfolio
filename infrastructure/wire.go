//+build wireinject

package infrastructure

import (
	"github.com/google/wire"
	impl "github.com/traPtitech/traPortfolio/interfaces/repository"
	"github.com/traPtitech/traPortfolio/usecases/handler"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/usecases/usecase"
)

var pingSet = wire.NewSet(
	handler.NewPingHandler,
	wire.Bind(new(usecase.PingUsecase), new(*handler.PingHandler)),
)

var userSet = wire.NewSet(
	impl.NewUserRepository,
	handler.NewUserHandler,
	wire.Bind(new(repository.UserRepository), new(*impl.UserRepository)),
	wire.Bind(new(usecase.UserUsecase), new(*handler.UserHandler)),
)

var sqlSet = wire.NewSet(NewSQLHandler)

var apiSet = wire.NewSet(handler.NewAPI)

func InjectAPIServer() (handler.API, error) {
	wire.Build(
		pingSet,
		userSet,
		sqlSet,
		apiSet,
	)
	return handler.API{}, nil
}
