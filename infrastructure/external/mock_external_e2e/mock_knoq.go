package mock_external_e2e //nolint:revive

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/infrastructure/external"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/mockdata"
)

type MockKnoqAPI struct{}

func NewMockKnoqAPI() *MockKnoqAPI {
	return &MockKnoqAPI{}
}

func (m *MockKnoqAPI) GetEvents() ([]*external.EventResponse, error) {
	return mockdata.MockKnoqEvents, nil
}

func (m *MockKnoqAPI) GetEvent(eventID uuid.UUID) (*external.EventResponse, error) {
	for _, v := range mockdata.MockKnoqEvents {
		if v.ID == eventID {
			return v, nil
		}
	}

	return nil, repository.ErrNotFound
}

func (m *MockKnoqAPI) GetEventsByUserID(userID uuid.UUID) ([]*external.EventResponse, error) {
	events := make([]*external.EventResponse, 0, len(mockdata.MockKnoqEvents))

	// TODO: adminsではなくattendeesを取得して判定する？
	for _, v := range mockdata.MockKnoqEvents {
		for _, admin := range v.Admins {
			if admin == userID {
				events = append(events, v)
				break
			}
		}
	}

	return events, nil
}
