// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package handler

import (
	"github.com/golang/mock/gomock"
	"github.com/google/wire"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/usecases/repository/mock_repository"
	"github.com/traPtitech/traPortfolio/usecases/service"
)

// Injectors from test_wire.go:

func SetupTestApi(ctrl *gomock.Controller) TestHandlers {
	pingHandler := NewPingHandler()
	mockUserRepository := mock_repository.NewMockUserRepository(ctrl)
	mockEventRepository := mock_repository.NewMockEventRepository(ctrl)
	userService := service.NewUserService(mockUserRepository, mockEventRepository)
	userHandler := NewUserHandler(userService)
	mockProjectRepository := mock_repository.NewMockProjectRepository(ctrl)
	mockPortalRepository := mock_repository.NewMockPortalRepository(ctrl)
	projectService := service.NewProjectService(mockProjectRepository, mockPortalRepository)
	projectHandler := NewProjectHandler(projectService)
	eventService := service.NewEventService(mockEventRepository)
	eventHandler := NewEventHandler(eventService)
	mockContestRepository := mock_repository.NewMockContestRepository(ctrl)
	contestService := service.NewContestService(mockContestRepository)
	contestHandler := NewContestHandler(contestService)
	api := NewAPI(pingHandler, userHandler, projectHandler, eventHandler, contestHandler)
	mockKnoqRepository := mock_repository.NewMockKnoqRepository(ctrl)
	mockTraQRepository := mock_repository.NewMockTraQRepository(ctrl)
	testRepository := TestRepository{
		MockContestRepository: mockContestRepository,
		MockEventRepository:   mockEventRepository,
		MockKnoqRepository:    mockKnoqRepository,
		MockPortalRepository:  mockPortalRepository,
		MockProjectRepository: mockProjectRepository,
		MockTraQRepository:    mockTraQRepository,
		MockUserRepository:    mockUserRepository,
	}
	testHandlers := TestHandlers{
		Api:        api,
		Repository: testRepository,
	}
	return testHandlers
}

// test_wire.go:

var mockPortalSet = wire.NewSet(mock_repository.NewMockPortalRepository, wire.Bind(new(repository.PortalRepository), new(*mock_repository.MockPortalRepository)))

var mockTraQSet = wire.NewSet(mock_repository.NewMockTraQRepository, wire.Bind(new(repository.TraQRepository), new(*mock_repository.MockTraQRepository)))

var mockPingSet = wire.NewSet(
	NewPingHandler,
)

var mockUserSet = wire.NewSet(mock_repository.NewMockUserRepository, service.NewUserService, NewUserHandler, wire.Bind(new(repository.UserRepository), new(*mock_repository.MockUserRepository)))

var mockProjectSet = wire.NewSet(mock_repository.NewMockProjectRepository, service.NewProjectService, NewProjectHandler, wire.Bind(new(repository.ProjectRepository), new(*mock_repository.MockProjectRepository)))

var mockKnoQSet = wire.NewSet(mock_repository.NewMockKnoqRepository, wire.Bind(new(repository.KnoqRepository), new(*mock_repository.MockKnoqRepository)))

var mockEventSet = wire.NewSet(
	mockKnoQSet, mock_repository.NewMockEventRepository, service.NewEventService, NewEventHandler, wire.Bind(new(repository.EventRepository), new(*mock_repository.MockEventRepository)),
)

var mockContestSet = wire.NewSet(mock_repository.NewMockContestRepository, service.NewContestService, NewContestHandler, wire.Bind(new(repository.ContestRepository), new(*mock_repository.MockContestRepository)))

var mockApiSet = wire.NewSet(NewAPI)
