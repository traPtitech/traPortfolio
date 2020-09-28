package repository

import (
	"context"

	"github.com/traPtitech/traPortfolio/domain"
)

type PortalRepository struct {
	token string
}

type PortalToken string

func NewPortalRepository(portalToken PortalToken) *PortalRepository {
	return &PortalRepository{token: string(portalToken)}
}

func (repo *PortalRepository) GetUser(ctx context.Context, name string) (user *domain.PortalUser, err error) {
	// TODO
	return
}
