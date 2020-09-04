//+build wireinject

package infrastructure

import (
	"github.com/google/wire"
	"github.com/traPtitech/traPortfolio/interfaces/controllers"
	impl "github.com/traPtitech/traPortfolio/interfaces/repository"
	"github.com/traPtitech/traPortfolio/usecase/interactor"
	"github.com/traPtitech/traPortfolio/usecase/repository"
)

var pingSet = wire.NewSet(
	impl.NewPingRepository,
	interactor.NewPingInteractor,
	controllers.NewPingController,
	wire.Bind(new(repository.PingRepository), new(*impl.PingRepository)),
)

var userSet = wire.NewSet(
	impl.NewUserRepository,
	interactor.NewUserInteractor,
	controllers.NewUserController,
	wire.Bind(new(repository.UserRepository), new(*impl.UserRepository)),
)

var sqlSet = wire.NewSet(NewSQLHandler)

var apiSet = wire.NewSet(controllers.NewAPI)

func InjectAPIServer() controllers.API {
	wire.Build(
		pingSet,
		userSet,
		sqlSet,
		apiSet,
	)
	return controllers.API{}
}
