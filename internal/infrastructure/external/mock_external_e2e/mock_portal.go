package mock_external_e2e //nolint:revive

import (
	"fmt"

	"github.com/traPtitech/traPortfolio/internal/infrastructure/external"
	"github.com/traPtitech/traPortfolio/internal/pkgs/mockdata"
)

type MockPortalAPI struct{}

func NewMockPortalAPI() *MockPortalAPI {
	return &MockPortalAPI{}
}

func (m *MockPortalAPI) GetUsers() ([]*external.PortalUserResponse, error) {
	return mockdata.MockPortalUsers, nil
}

func (m *MockPortalAPI) GetUserByTraqID(traQID string) (*external.PortalUserResponse, error) {
	for _, v := range mockdata.MockPortalUsers {
		if v.TraQID == traQID {
			return v, nil
		}
	}

	return nil, fmt.Errorf("GET /user/%v failed: 404", traQID)
}
