package repository

import (
	"context"

	"github.com/traPtitech/traPortfolio/usecases/repository"

	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
)

type PortalRepository struct {
	token string
}

type PortalToken string

func NewPortalRepository(portalToken PortalToken) *PortalRepository {
	return &PortalRepository{token: string(portalToken)}
}

func (repo *PortalRepository) GetUser(ctx context.Context, name string) (user *model.PortalUser, err error) {
	// TODO
	return
}

func (repo *PortalRepository) GetUsers(ctx context.Context) (users []*model.PortalUser, err error) {
	return
}

// Interface guards
var (
	_ repository.PortalRepository = (*PortalRepository)(nil)
)
