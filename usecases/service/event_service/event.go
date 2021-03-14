package service

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

// Event Portfolioのレスポンスで使うイベント情報
type Event struct {
	ID          uuid.UUID        `json:"eventId"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	GroupID     uuid.UUID        `json:"groupId"`
	RoomID      uuid.UUID        `json:"roomId"`
	TimeStart   time.Time        `json:"timeStart"`
	TimeEnd     time.Time        `json:"timeEnd"`
	SharedRoom  bool             `json:"sharedRoom"`
	Level       model.EventLevel `json:"eventLevel"`
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

func (s *EventService) GetEvents(ctx context.Context) ([]*model.Event, error) {
	er, err := s.knoQ.GetAll()
	if err != nil {
		return nil, err
	}

	elvsmp, err := s.repo.GetEventLevels()
	if err != nil {
		return nil, err
	}

	result := make([]*model.Event, 0, len(er))
	for _, v := range er {
		e := &model.Event{
			ID:          v.ID,
			Description: v.Description,
			GroupID:     v.GroupID,
			Name:        v.Name,
			RoomID:      v.RoomID,
			SharedRoom:  v.SharedRoom,
			TimeEnd:     v.TimeEnd,
		}
		if level, ok := elvsmp[v.ID]; ok {
			e.Level = level.Level
		}
		result = append(result, e)
	}

	return result, nil
}

func (s *EventService) GetEventByID(ctx context.Context, id uuid.UUID) (*model.Event, error) {
	er, err := s.knoQ.GetByID(id)
	if err != nil {
		return nil, err
	}

	elv, err := s.repo.GetEventLevelByID(id)
	if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	result := &model.Event{
		ID:          er.ID,
		Description: er.Description,
		GroupID:     er.GroupID,
		Name:        er.Name,
		RoomID:      er.RoomID,
		SharedRoom:  er.SharedRoom,
		TimeEnd:     er.TimeEnd,
		TimeStart:   er.TimeStart,
	}

	if err == nil {
		result.Level = elv.Level
	}

	return result, nil
}
