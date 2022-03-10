package mock_external_e2e //nolint:revive

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/external"
)

var (
	mockTraQUsers = []*external.TraQUserResponse{
		{
			ID:    uuid.FromStringOrNil("11111111-1111-1111-1111-111111111111"),
			State: domain.TraqStateActive,
		},
		{
			ID:    uuid.FromStringOrNil("22222222-2222-2222-2222-222222222222"),
			State: domain.TraqStateDeactivated,
		},
		{
			ID:    uuid.FromStringOrNil("33333333-3333-3333-3333-333333333333"),
			State: domain.TraqStateActive,
		},
	}
)

type MockTraQAPI struct{}

func NewMockTraQAPI() *MockTraQAPI {
	return &MockTraQAPI{}
}

func (m *MockTraQAPI) GetAll(args *external.TraQGetAllArgs) ([]*external.TraQUserResponse, error) {
	return mockTraQUsers, nil
}

func (m *MockTraQAPI) GetByID(id uuid.UUID) (*external.TraQUserResponse, error) {
	for _, u := range mockTraQUsers {
		if u.ID == id {
			return u, nil
		}
	}

	return nil, fmt.Errorf("GET /users/%v failed: 404", id)
}
