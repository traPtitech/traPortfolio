//+build wireinject

package di

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

var apiSet = wire.NewSet(controllers.NewAPI)

func InjectAPIServer() controllers.API {
	wire.Build(
		pingSet,
		apiSet,
	)
	return controllers.API{}
}
