//+build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/traPtitech/traPortfolio/interfaces/controllers"
	impl "github.com/traPtitech/traPortfolio/interfaces/repository"
	"github.com/traPtitech/traPortfolio/usecase/interactor"
	"github.com/traPtitech/traPortfolio/usecase/repository"
)

var repositorySet = wire.NewSet(
	impl.NewPingRepository,
)

var interactorSet = wire.NewSet(
	interactor.NewPingInteractor,
)

var controllerSet = wire.NewSet(
	controllers.NewPingController,
)

var apiSet = wire.NewSet(controllers.NewAPI)

func InjectAPIServer() controllers.API {
	wire.Build(
		repositorySet,
		interactorSet,
		controllerSet,
		apiSet,
		wire.Bind(new(repository.PingRepository), new(*impl.PingRepository)),
	)
	return controllers.API{}
}
