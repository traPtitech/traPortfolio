package mock_external_e2e //nolint:revive

import (
	"github.com/traPtitech/traPortfolio/internal/domain"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/external"
	"github.com/traPtitech/traPortfolio/internal/pkgs/mockdata"
)

type MockTraQAPI struct{}

func NewMockTraQAPI() *MockTraQAPI {
	return &MockTraQAPI{}
}

func (m *MockTraQAPI) GetUsers(args *external.TraQGetAllArgs) ([]*external.TraQUserResponse, error) {
	users := make([]*external.TraQUserResponse, 0, len(mockdata.MockTraQUsers))
	for _, u := range mockdata.MockTraQUsers {
		if args.Name == u.Name {
			users = append(users, u)

			return users, nil
		}

		if args.IncludeSuspended || u.State == domain.TraqStateActive {
			users = append(users, u)
		}
	}

	return users, nil
}
