package repository

import (
	"context"

	"github.com/traPtitech/traPortfolio/domain"
)

type PortalRepository struct {
	token string
}

func NewPortalRepository(token string) *PortalRepository {
	return &PortalRepository{token: token}
}

func (repo *PortalRepository) GetUser(ctx context.Context, name string) (user *domain.PortalUser, err error) {
	// TODO
	return
}
