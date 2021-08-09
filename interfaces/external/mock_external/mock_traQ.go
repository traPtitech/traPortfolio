package mock_external //nolint:revive

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/external"
)

var (
	sampleTraQUserResponse = external.TraQUserResponse{
		State:       0,
		Bot:         true,
		DisplayName: "Noriko Azuma",
		Name:        "lolico",
	}
)

type MockTraQAPI struct{}

func (m *MockTraQAPI) GetByID(id uuid.UUID) (*external.TraQUserResponse, error) {
	return &sampleTraQUserResponse, nil
}
