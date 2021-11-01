package mock_external //nolint:revive

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/external"
)

var (
	traQUserMap = map[uuid.UUID]*external.TraQUserResponse{
		uuid.FromStringOrNil("11111111-1111-1111-1111-111111111111"): {
			State:       domain.TraqStateActive,
			Bot:         false,
			DisplayName: "user1",
			Name:        "user1",
		},
		uuid.FromStringOrNil("22222222-2222-2222-2222-222222222222"): {
			State:       domain.TraqStateDeactivated,
			Bot:         false,
			DisplayName: "凍結ユーザー",
			Name:        "user2",
		},
		uuid.FromStringOrNil("33333333-3333-3333-3333-333333333333"): {
			State:       domain.TraqStateActive,
			Bot:         true,
			DisplayName: "Noriko Azuma",
			Name:        "lolico",
		},
	}
)

type MockTraQAPI struct{}

func NewMockTraQAPI() *MockTraQAPI {
	return &MockTraQAPI{}
}

func (m *MockTraQAPI) GetByID(id uuid.UUID) (*external.TraQUserResponse, error) {
	if res, ok := traQUserMap[id]; ok {
		return res, nil
	}

	return nil, fmt.Errorf("GET /users/%v failed: 404", id)
}
