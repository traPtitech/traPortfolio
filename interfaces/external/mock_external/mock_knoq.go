package mock_external //nolint:revive

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/external"
)

var (
	sampleUUID          = uuid.FromStringOrNil("3fa85f64-5717-4562-b3fc-2c963f66afa6")
	sampleTime          = time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)
	sampleEventResponse = external.EventResponse{
		ID:          sampleUUID,
		Name:        "第n回進捗回",
		Description: "第n回の進捗会です。",
		GroupID:     sampleUUID,
		RoomID:      sampleUUID,
		TimeStart:   sampleTime,
		TimeEnd:     sampleTime,
		SharedRoom:  true,
	}
)

type MockKnoqAPI struct{}

func (m *MockKnoqAPI) GetAll() ([]*external.EventResponse, error) {
	return []*external.EventResponse{&sampleEventResponse}, nil
}

func (m *MockKnoqAPI) GetByID(id uuid.UUID) (*external.EventResponse, error) {
	return &sampleEventResponse, nil
}

func (m *MockKnoqAPI) GetByUserID(userID uuid.UUID) ([]*external.EventResponse, error) {
	return []*external.EventResponse{&sampleEventResponse}, nil
}
