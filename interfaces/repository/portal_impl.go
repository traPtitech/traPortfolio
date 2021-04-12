package repository

import (
	"context"

	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type PortalRepository struct {
}

func NewPortalRepository() *PortalRepository {
	return &PortalRepository{}
}

func (repo *PortalRepository) GetUser(ctx context.Context, name string) (user *domain.PortalUser, err error) {
	// TODO
	return
}

func (repo *PortalRepository) GetUsers(ctx context.Context) (users []*domain.PortalUser, err error) {
	return
}

// Interface guards
var (
	_ repository.PortalRepository = (*PortalRepository)(nil)
)
