package handler

import (
	"github.com/traPtitech/traPortfolio/usecases/service/mock_service"
)

// wire が_testファイル用に使えない(https://github.com/google/wire/issues/48)のでここに定義している。

type TestService struct {
	*mock_service.MockContestService
	*mock_service.MockEventService
	*mock_service.MockProjectService
	*mock_service.MockUserService
}

type TestHandlers struct {
	API     API
	Service TestService
}
