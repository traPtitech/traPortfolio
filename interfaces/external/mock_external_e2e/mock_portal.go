package mock_external_e2e //nolint:revive

import (
	"fmt"

	"github.com/traPtitech/traPortfolio/interfaces/external"
)

var (
	mockPortalUsers = []*external.PortalUserResponse{
		{
			TraQID:         "user1",
			RealName:       "ユーザー1 ユーザー1",
			AlphabeticName: "user1 user1",
		},
		{
			TraQID:         "user2",
			RealName:       "ユーザー2 ユーザー2",
			AlphabeticName: "user2 user2",
		},
		{
			TraQID:         "lolico",
			RealName:       "東 工子",
			AlphabeticName: "Noriko Azuma",
		},
	}
)

type MockPortalAPI struct{}

func NewMockPortalAPI() *MockPortalAPI {
	return &MockPortalAPI{}
}

func (m *MockPortalAPI) GetAll() ([]*external.PortalUserResponse, error) {
	return mockPortalUsers, nil
}

func (m *MockPortalAPI) GetByID(traQID string) (*external.PortalUserResponse, error) {
	for _, v := range mockPortalUsers {
		if v.TraQID == traQID {
			return v, nil
		}
	}

	return nil, fmt.Errorf("GET /user/%v failed: 404", traQID)
}
