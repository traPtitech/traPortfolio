package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
)

type TraQRepository struct {
	token string
}

type TraQToken string

func NewTraQRepository(traQToken TraQToken) *TraQRepository {
	return &TraQRepository{token: string(traQToken)}
}

func (repo *TraQRepository) GetUser(ctx context.Context, id uuid.UUID) (user *domain.TraQUser, err error) {
	// TODO
	return
}
