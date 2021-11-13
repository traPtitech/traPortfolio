package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type KnoqRepository struct {
	api external.KnoqAPI
}

func NewKnoqRepository(api external.KnoqAPI) repository.KnoqRepository {
	return &KnoqRepository{api}
}

func (repo *KnoqRepository) GetAll() ([]*domain.KnoQEvent, error) {
	eres, err := repo.api.GetAll()
	if err != nil {
		return nil, convertError(err)
	}
	result := make([]*domain.KnoQEvent, 0, len(eres))
	for _, v := range eres {
		result = append(result, &domain.KnoQEvent{
			ID:          v.ID,
			Description: v.Description,
			GroupID:     v.GroupID,
			Name:        v.Name,
			RoomID:      v.RoomID,
			SharedRoom:  v.SharedRoom,
			TimeEnd:     v.TimeEnd,
			TimeStart:   v.TimeStart,
		})
	}
	return result, nil
}

func (repo *KnoqRepository) GetByID(id uuid.UUID) (*domain.KnoQEvent, error) {
	eres, err := repo.api.GetByID(id)
	if err != nil {
		return nil, convertError(err)
	}

	return &domain.KnoQEvent{
		ID:          eres.ID,
		Description: eres.Description,
		GroupID:     eres.GroupID,
		Name:        eres.Name,
		RoomID:      eres.RoomID,
		SharedRoom:  eres.SharedRoom,
		TimeEnd:     eres.TimeEnd,
		TimeStart:   eres.TimeStart,
	}, nil
}

// Interface guards
var (
	_ repository.KnoqRepository = (*KnoqRepository)(nil)
)
