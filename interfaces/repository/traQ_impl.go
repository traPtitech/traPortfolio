package repository

import (
	"context"

	"github.com/traPtitech/traPortfolio/usecases/repository"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
)

type TraQRepository struct {
	api external.TraQAPI
}

func NewTraQRepository(api external.TraQAPI) *TraQRepository {
	return &TraQRepository{
		api: api,
	}
}

func (repo *TraQRepository) GetUser(ctx context.Context, id uuid.UUID) (*model.TraQUser, error) {
	ures, err := repo.api.GetByID(id)
	if err != nil {
		return nil, err
	}
	return &model.TraQUser{
		State:       ures.State,
		Bot:         ures.Bot,
		DisplayName: ures.DisplayName,
		Name:        ures.Name,
	}, nil
}

// Interface guards
var (
	_ repository.TraQRepository = (*TraQRepository)(nil)
)
