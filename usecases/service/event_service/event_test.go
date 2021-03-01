package service

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/usecases/repository/mock_repository"
	"github.com/traPtitech/traPortfolio/util"
)

func TestGetEvents(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	knoQEventExpected := makeKnoQEvent(20)
	knoQRepository := mock_repository.NewMockKnoqRepository(mockCtrl)
	knoQRepository.EXPECT().GetAll().Return(knoQEventExpected, nil)

	eventExpected := makeEventRelation(knoQEventExpected)
	eventRepository := mock_repository.NewMockEventRepository(mockCtrl)
	eventRepository.EXPECT().GetEventLevels().Return(eventExpected, nil)

	expected := makeEvents(knoQEventExpected, eventExpected)
	// expected[0].Description = "hoge"

	ctx := context.Background()
	service := NewEventService(eventRepository, knoQRepository)
	event, err := service.GetEvents(ctx)

	assert := assert.New(t)

	if assert.NoError(err) {
		assert.Equal(expected, event)
	}
}

func makeKnoQEvent(num int) []*repository.KnoQEvent {
	result := make([]*repository.KnoQEvent, 0, num)
	for i := 0; i < num; i++ {
		result = append(result, &repository.KnoQEvent{
			Description: util.AlphaNumeric(20),
			GroupID:     uuid.Must(uuid.NewV4()),
			ID:          uuid.Must(uuid.NewV4()),
			Name:        util.AlphaNumeric(20),
			RoomID:      uuid.Must(uuid.NewV4()),
			SharedRoom:  false,
			TimeEnd:     time.Now(),
			TimeStart:   time.Now(),
		})
	}
	return result
}

func makeEventRelation(e []*repository.KnoQEvent) map[uuid.UUID]*domain.EventLevelRelation {
	result := make(map[uuid.UUID]*domain.EventLevelRelation, len(e))
	for i, v := range e {
		result[v.ID] = &domain.EventLevelRelation{
			ID:    v.ID,
			Level: domain.EventLevelAnonymous + domain.EventLevel(i%3),
		}
	}
	return result
}

func makeEvents(knoQEvent []*repository.KnoQEvent, rel map[uuid.UUID]*domain.EventLevelRelation) []*domain.Event {
	result := make([]*domain.Event, 0, len(knoQEvent))
	for _, v := range knoQEvent {
		e := &domain.Event{
			ID:          v.ID,
			Description: v.Description,
			GroupID:     v.GroupID,
			Name:        v.Name,
			RoomID:      v.RoomID,
			SharedRoom:  v.SharedRoom,
			TimeEnd:     v.TimeEnd,
		}
		if level, ok := rel[v.ID]; ok {
			e.Level = level.Level
		}
		result = append(result, e)
	}

	return result
}
