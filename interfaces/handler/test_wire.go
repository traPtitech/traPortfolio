//go:generate go run github.com/google/wire/cmd/wire@v0.5.0
//go:build wireinject

package handler

import (
	"github.com/golang/mock/gomock"
	"github.com/google/wire"
	"github.com/traPtitech/traPortfolio/usecases/service"
	"github.com/traPtitech/traPortfolio/usecases/service/mock_service"
)

var mockPingSet = wire.NewSet(
	NewPingHandler,
)

var mockUserSet = wire.NewSet(
	mock_service.NewMockUserService,
	NewUserHandler,
	wire.Bind(new(service.UserService), new(*mock_service.MockUserService)),
)

var mockProjectSet = wire.NewSet(
	mock_service.NewMockProjectService,
	NewProjectHandler,
	wire.Bind(new(service.ProjectService), new(*mock_service.MockProjectService)),
)

var mockEventSet = wire.NewSet(
	mock_service.NewMockEventService,
	NewEventHandler,
	wire.Bind(new(service.EventService), new(*mock_service.MockEventService)),
)

var mockGroupSet = wire.NewSet(
	mock_service.NewMockGroupService,
	NewGroupHandler,
	wire.Bind(new(service.GroupService), new(*mock_service.MockGroupService)),
)

var mockContestSet = wire.NewSet(
	mock_service.NewMockContestService,
	NewContestHandler,
	wire.Bind(new(service.ContestService), new(*mock_service.MockContestService)),
)

var mockApiSet = wire.NewSet(NewAPI)

func SetupTestApi(ctrl *gomock.Controller) TestHandlers {
	wire.Build(
		wire.Struct(new(TestService), "*"),
		wire.Struct(new(TestHandlers), "*"),
		mockPingSet,
		mockUserSet,
		mockProjectSet,
		mockEventSet,
		mockApiSet,
		mockContestSet,
		mockGroupSet,
	)
	return TestHandlers{}
}
