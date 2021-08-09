package mock_external //nolint:revive

import "github.com/traPtitech/traPortfolio/interfaces/external"

var (
	samplePortalUserResponse = external.PortalUserResponse{
		TraQID:         "lolico",
		RealName:       "東 工子",
		AlphabeticName: "Noriko Azuma",
	}
)

type MockPortalAPI struct{}

func (m *MockPortalAPI) GetAll() ([]*external.PortalUserResponse, error) {
	return []*external.PortalUserResponse{&samplePortalUserResponse}, nil
}

func (m *MockPortalAPI) GetByID(traQID string) (*external.PortalUserResponse, error) {
	return &samplePortalUserResponse, nil
}
