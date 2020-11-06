package service

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

// Event Portfolioのレスポンスで使うイベント情報
type Event struct {
	ID          uuid.UUID         `json:"eventId"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	GroupID     uuid.UUID         `json:"groupId"`
	RoomID      uuid.UUID         `json:"roomId"`
	TimeStart   time.Time         `json:"timeStart"`
	TimeEnd     time.Time         `json:"timeEnd"`
	SharedRoom  bool              `json:"sharedRoom"`
	Level       domain.EventLevel `json:"eventLevel"`
}

type EventService struct {
	repo repository.EventRepository
	knoQ repository.KnoqRepository
}

func NewEventService(repo repository.EventRepository, knoQ repository.KnoqRepository) EventService {
	return EventService{
		repo,
		knoQ,
	}
}

func (s *EventService) GetEvents(ctx context.Context) ([]*domain.Event, error) {
	er, err := s.knoQ.GetAll()
	if err != nil {
		return nil, err
	}

	elvs, err := s.repo.GetEventLevels()
	if err == repository.ErrNotFound {
		elvs = make([]*domain.EventLevelRelation, 0)
	}
	if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	elvsmp := make(map[uuid.UUID]*domain.EventLevelRelation)
	for _, v := range elvs {
		v := v
		elvsmp[v.ID] = v
	}
	result := make([]*domain.Event, 0, len(er))
	for _, v := range er {
		_level, ok := elvsmp[v.ID]
		var level domain.EventLevel = domain.EventLevelAnonymous
		if ok {
			level = _level.Level
		}
		result = append(result, &domain.Event{
			ID:          v.ID,
			Description: v.Description,
			GroupID:     v.GroupID,
			Level:       level,
			Name:        v.Name,
			RoomID:      v.RoomID,
			SharedRoom:  v.SharedRoom,
			TimeEnd:     v.TimeEnd,
			TimeStart:   v.TimeStart,
		})
	}

	return result, nil
}

func (s *EventService) GetEventByID(ctx context.Context, id uuid.UUID) (*domain.Event, error) {
	er, err := s.knoQ.GetByID(id)
	if err != nil {
		return nil, err
	}

	elv, err := s.repo.GetEventLevelByID(id)
	if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	var level domain.EventLevel = domain.EventLevelAnonymous
	if err == nil {
		level = elv.Level
	}
	result := &domain.Event{
		ID:          er.ID,
		Description: er.Description,
		GroupID:     er.GroupID,
		Level:       level,
		Name:        er.Name,
		RoomID:      er.RoomID,
		SharedRoom:  er.SharedRoom,
		TimeEnd:     er.TimeEnd,
		TimeStart:   er.TimeStart,
	}

	return result, nil
}
