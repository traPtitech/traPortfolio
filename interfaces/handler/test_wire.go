//go:build wireinject

package handler

import (
	"github.com/golang/mock/gomock"
	"github.com/google/wire"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/usecases/repository/mock_repository"
	"github.com/traPtitech/traPortfolio/usecases/service"
)

var mockPortalSet = wire.NewSet(
	mock_repository.NewMockPortalRepository,
	wire.Bind(new(repository.PortalRepository), new(*mock_repository.MockPortalRepository)),
)

var mockTraQSet = wire.NewSet(
	mock_repository.NewMockTraQRepository,
	wire.Bind(new(repository.TraQRepository), new(*mock_repository.MockTraQRepository)),
)

var mockPingSet = wire.NewSet(
	NewPingHandler,
)

var mockUserSet = wire.NewSet(
	mock_repository.NewMockUserRepository,
	service.NewUserService,
	NewUserHandler,
	wire.Bind(new(repository.UserRepository), new(*mock_repository.MockUserRepository)),
)

var mockProjectSet = wire.NewSet(
	mock_repository.NewMockProjectRepository,
	service.NewProjectService,
	NewProjectHandler,
	wire.Bind(new(repository.ProjectRepository), new(*mock_repository.MockProjectRepository)),
)

var mockKnoQSet = wire.NewSet(
	mock_repository.NewMockKnoqRepository,
	wire.Bind(new(repository.KnoqRepository), new(*mock_repository.MockKnoqRepository)),
)

var mockEventSet = wire.NewSet(
	mockKnoQSet,
	mock_repository.NewMockEventRepository,
	service.NewEventService,
	NewEventHandler,
	wire.Bind(new(repository.EventRepository), new(*mock_repository.MockEventRepository)),
)

var mockContestSet = wire.NewSet(
	mock_repository.NewMockContestRepository,
	service.NewContestService,
	NewContestHandler,
	wire.Bind(new(repository.ContestRepository), new(*mock_repository.MockContestRepository)),
)

var mockApiSet = wire.NewSet(NewAPI)

func SetupTestApi(ctrl *gomock.Controller) TestHandlers {
	wire.Build(
		wire.Struct(new(TestRepository), "*"),
		wire.Struct(new(TestHandlers), "*"),
		mockPingSet,
		mockUserSet,
		mockProjectSet,
		mockEventSet,
		mockApiSet,
		mockPortalSet,
		mockTraQSet,
		mockContestSet,
	)
	return TestHandlers{}
}
