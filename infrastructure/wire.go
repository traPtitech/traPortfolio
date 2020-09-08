//+build wireinject

package infrastructure

import (
	"github.com/google/wire"
	"github.com/traPtitech/traPortfolio/interfaces/controllers"
	impl "github.com/traPtitech/traPortfolio/interfaces/repository"
	"github.com/traPtitech/traPortfolio/usecases/interactor"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/usecases/usecase"
)

var pingSet = wire.NewSet(
	interactor.NewPingInteractor,
	controllers.NewPingController,
	wire.Bind(new(usecase.PingUsecase), new(*interactor.PingInteractor)),
)

var userSet = wire.NewSet(
	impl.NewUserRepository,
	interactor.NewUserInteractor,
	controllers.NewUserController,
	wire.Bind(new(repository.UserRepository), new(*impl.UserRepository)),
	wire.Bind(new(usecase.UserUsecase), new(*interactor.UserInteractor)),
)

var sqlSet = wire.NewSet(NewSQLHandler)

var apiSet = wire.NewSet(controllers.NewAPI)

func InjectAPIServer() (controllers.API, error) {
	wire.Build(
		pingSet,
		userSet,
		sqlSet,
		apiSet,
	)
	return controllers.API{}, nil
}
