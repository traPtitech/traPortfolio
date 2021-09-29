package mock_external //nolint:revive

import (
	"fmt"

	"github.com/traPtitech/traPortfolio/interfaces/external"
)

var (
	portalUserMap = map[string]*external.PortalUserResponse{
		"user1": {
			TraQID:         "user1",
			RealName:       "ユーザー1 ユーザー1",
			AlphabeticName: "user1 user1",
		},
		"user2": {
			TraQID:         "user2",
			RealName:       "ユーザー2 ユーザー2",
			AlphabeticName: "user2 user2",
		},
		"lolico": {
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
	res := make([]*external.PortalUserResponse, 0)
	for _, v := range portalUserMap {
		res = append(res, v)
	}

	return res, nil
}

func (m *MockPortalAPI) GetByID(traQID string) (*external.PortalUserResponse, error) {
	if res, ok := portalUserMap[traQID]; ok {
		return res, nil
	}

	return nil, fmt.Errorf("GET /user/%v failed: 404", traQID)
}
