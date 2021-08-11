package mock_external //nolint:revive

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/external"
)

var (
	sampleUUID   = uuid.FromStringOrNil("3fa85f64-5717-4562-b3fc-2c963f66afa6")
	sample2UUID  = uuid.FromStringOrNil("11111111-1111-1111-1111-111111111111")
	sampleTime   = time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)
	knoQEventMap = map[uuid.UUID]*external.EventResponse{
		sampleUUID: {
			ID:          sampleUUID,
			Name:        "第n回進捗回",
			Description: "第n回の進捗会です。",
			GroupID:     sampleUUID,
			RoomID:      sampleUUID,
			TimeStart:   sampleTime,
			TimeEnd:     sampleTime,
			SharedRoom:  true,
		},
		sample2UUID: {
			ID:          sample2UUID,
			Name:        "sample event",
			Description: "This is a sample event.",
			GroupID:     sample2UUID,
			RoomID:      sample2UUID,
			TimeStart:   sampleTime,
			TimeEnd:     sampleTime,
			SharedRoom:  false,
		},
	}
)

type MockKnoqAPI struct{}

func (m *MockKnoqAPI) GetAll() ([]*external.EventResponse, error) {
	res := make([]*external.EventResponse, 0)
	for _, v := range knoQEventMap {
		res = append(res, v)
	}

	return res, nil
}

func (m *MockKnoqAPI) GetByID(id uuid.UUID) (*external.EventResponse, error) {
	if res, ok := knoQEventMap[id]; ok {
		return res, nil
	}

	return nil, fmt.Errorf("GET /events/%v failed: 404", id)
}

func (m *MockKnoqAPI) GetByUserID(userID uuid.UUID) ([]*external.EventResponse, error) {
	return []*external.EventResponse{knoQEventMap[sampleUUID]}, nil
}
