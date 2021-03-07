package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/external"
)

type TraQRepository struct {
	token string
	api   external.TraQAPI
}

type TraQToken string

func NewTraQRepository(api external.TraQAPI, traQToken TraQToken) *TraQRepository {
	return &TraQRepository{
		token: string(traQToken),
		api:   api,
	}
}

func (repo *TraQRepository) GetUser(ctx context.Context, id uuid.UUID) (*domain.TraQUser, error) {
	ures, err := repo.api.GetByID(id, repo.token)
	if err != nil {
		return nil, err
	}
	return &domain.TraQUser{
		State:       ures.State,
		Bot:         ures.Bot,
		DisplayName: ures.DisplayName,
		Name:        ures.Name,
	}, nil
}
