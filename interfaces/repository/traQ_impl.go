package repository

import (
	"context"

	"github.com/traPtitech/traPortfolio/domain"

	"github.com/traPtitech/traPortfolio/usecases/repository"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/external"
)

type TraQRepository struct {
	api external.TraQAPI
}

func NewTraQRepository(api external.TraQAPI) repository.TraQRepository {
	return &TraQRepository{
		api: api,
	}
}

func (repo *TraQRepository) GetUser(ctx context.Context, id uuid.UUID) (*domain.TraQUser, error) {
	ures, err := repo.api.GetByID(id)
	if err != nil {
		return nil, convertError(err)
	}
	return &domain.TraQUser{
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
