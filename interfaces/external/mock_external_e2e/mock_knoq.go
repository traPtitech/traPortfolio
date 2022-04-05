package mock_external_e2e //nolint:revive

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/util/mockdata"
)

type MockKnoqAPI struct{}

func NewMockKnoqAPI() *MockKnoqAPI {
	return &MockKnoqAPI{}
}

func (m *MockKnoqAPI) GetAll() ([]*external.EventResponse, error) {
	return mockdata.MockKnoqEvents, nil
}

func (m *MockKnoqAPI) GetByID(id uuid.UUID) (*external.EventResponse, error) {
	for _, v := range mockdata.MockKnoqEvents {
		if v.ID == id {
			return v, nil
		}
	}

	return nil, fmt.Errorf("GET /events/%v failed: 404", id)
}

func (m *MockKnoqAPI) GetByUserID(userID uuid.UUID) ([]*external.EventResponse, error) {
	return mockdata.MockKnoqEvents, nil
}
