package repository

import (
	"context"

	"github.com/traPtitech/traPortfolio/domain"
)

type TraQRepository struct {
	token string
}

func NewTraQRepository(token string) *TraQRepository {
	return &TraQRepository{token: token}
}

func (repo *TraQRepository) GetUser(ctx context.Context, name string) (user *domain.TraQUser, err error) {
	// TODO
	return
}
