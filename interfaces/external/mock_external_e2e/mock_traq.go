package mock_external_e2e //nolint:revive

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/util/mockdata"
)

type MockTraQAPI struct{}

func NewMockTraQAPI() *MockTraQAPI {
	return &MockTraQAPI{}
}

func (m *MockTraQAPI) GetUsers(args *external.TraQGetAllArgs) ([]*external.TraQUserResponse, error) {
	users := make([]*external.TraQUserResponse, 0, len(mockdata.MockTraQUsers))
	for _, u := range mockdata.MockTraQUsers {
		if args.Name == u.Name {
			users = append(users, u.User)

			return users, nil
		}

		if args.IncludeSuspended || u.User.State == domain.TraqStateActive {
			users = append(users, u.User)
		}
	}

	return users, nil
}

func (m *MockTraQAPI) GetUserByUserID(userID uuid.UUID) (*external.TraQUserResponse, error) {
	for _, u := range mockdata.MockTraQUsers {
		if u.User.ID == userID {
			return u.User, nil
		}
	}

	return nil, fmt.Errorf("GET /users/%v failed: 404", userID)
}
