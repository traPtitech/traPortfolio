package handler

import (
	"github.com/traPtitech/traPortfolio/usecases/repository/mock_repository"
)

// wire が_testファイル用に使えない(https://github.com/google/wire/issues/48)のでここに定義している。

type TestRepository struct {
	*mock_repository.MockContestRepository
	*mock_repository.MockEventRepository
	*mock_repository.MockKnoqRepository
	*mock_repository.MockPortalRepository
	*mock_repository.MockProjectRepository
	*mock_repository.MockTraQRepository
	*mock_repository.MockUserRepository
}

type TestHandlers struct {
	Api        API
	Repository TestRepository
}
