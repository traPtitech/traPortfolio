package mock_external_e2e //nolint:revive

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/external"
)

var (
	mockTraQUsers = []*mockTraQUser{
		{
			u: &external.TraQUserResponse{
				ID:    uuid.FromStringOrNil("11111111-1111-1111-1111-111111111111"),
				State: domain.TraqStateActive,
			},
			name: "user1",
		},
		{
			u: &external.TraQUserResponse{
				ID:    uuid.FromStringOrNil("22222222-2222-2222-2222-222222222222"),
				State: domain.TraqStateDeactivated,
			},
			name: "user2",
		},
		{
			u: &external.TraQUserResponse{
				ID:    uuid.FromStringOrNil("33333333-3333-3333-3333-333333333333"),
				State: domain.TraqStateActive,
			},
			name: "lolico",
		},
	}
)

type mockTraQUser struct {
	u    *external.TraQUserResponse
	name string
}

type MockTraQAPI struct{}

func NewMockTraQAPI() *MockTraQAPI {
	return &MockTraQAPI{}
}

func (m *MockTraQAPI) GetAll(args *external.TraQGetAllArgs) ([]*external.TraQUserResponse, error) {
	users := make([]*external.TraQUserResponse, 0, len(mockTraQUsers))
	for _, u := range mockTraQUsers {
		if args.Name == u.name {
			users = append(users, u.u)

			return users, nil
		}

		if args.IncludeSuspended || u.u.State == domain.TraqStateActive {
			users = append(users, u.u)
		}
	}

	return users, nil
}

func (m *MockTraQAPI) GetByUserID(userID uuid.UUID) (*external.TraQUserResponse, error) {
	for _, u := range mockTraQUsers {
		if u.u.ID == userID {
			return u.u, nil
		}
	}

	return nil, fmt.Errorf("GET /users/%v failed: 404", userID)
}
