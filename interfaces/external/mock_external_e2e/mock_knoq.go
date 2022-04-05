package mock_external_e2e //nolint:revive

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/external"
)

var (
	sampleUUID     = uuid.FromStringOrNil("3fa85f64-5717-4562-b3fc-2c963f66afa6")
	sample2UUID    = uuid.FromStringOrNil("11111111-1111-1111-1111-111111111111")
	sampleTime     = time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)
	mockKnoqEvents = []*external.EventResponse{
		{
			ID:          sampleUUID,
			Name:        "第n回進捗回",
			Description: "第n回の進捗会です。",
			Place:       "S516",
			GroupID:     sampleUUID,
			RoomID:      sampleUUID,
			TimeStart:   sampleTime,
			TimeEnd:     sampleTime,
			SharedRoom:  true,
			Admins:      []uuid.UUID{sampleUUID},
		},
		{
			ID:          sample2UUID,
			Name:        "sample event",
			Description: "This is a sample event.",
			Place:       "S516",
			GroupID:     sample2UUID,
			RoomID:      sample2UUID,
			TimeStart:   sampleTime,
			TimeEnd:     sampleTime,
			SharedRoom:  false,
			Admins:      []uuid.UUID{sample2UUID},
		},
	}
)

type MockKnoqAPI struct{}

func NewMockKnoqAPI() *MockKnoqAPI {
	return &MockKnoqAPI{}
}

func (m *MockKnoqAPI) GetAll() ([]*external.EventResponse, error) {
	return mockKnoqEvents, nil
}

func (m *MockKnoqAPI) GetByEventID(eventID uuid.UUID) (*external.EventResponse, error) {
	for _, v := range mockKnoqEvents {
		if v.ID == eventID {
			return v, nil
		}
	}

	return nil, fmt.Errorf("GET /events/%v failed: 404", eventID)
}

func (m *MockKnoqAPI) GetByUserID(userID uuid.UUID) ([]*external.EventResponse, error) {
	return mockKnoqEvents, nil
}
