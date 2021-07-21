package service

import (
	"github.com/golang/mock/gomock"
	"github.com/traPtitech/traPortfolio/usecases/repository/mock_repository"
)

type MockRepository struct {
	contest *mock_repository.MockContestRepository
	event   *mock_repository.MockEventRepository
	knoq    *mock_repository.MockKnoqRepository
	portal  *mock_repository.MockPortalRepository
	project *mock_repository.MockProjectRepository
	traq    *mock_repository.MockTraQRepository
	user    *mock_repository.MockUserRepository
}

func newMockRepository(ctrl *gomock.Controller) *MockRepository {
	return &MockRepository{
		contest: mock_repository.NewMockContestRepository(ctrl),
		event:   mock_repository.NewMockEventRepository(ctrl),
		knoq:    mock_repository.NewMockKnoqRepository(ctrl),
		portal:  mock_repository.NewMockPortalRepository(ctrl),
		project: mock_repository.NewMockProjectRepository(ctrl),
		traq:    mock_repository.NewMockTraQRepository(ctrl),
		user:    mock_repository.NewMockUserRepository(ctrl),
	}
}
