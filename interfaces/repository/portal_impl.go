package repository

import (
	"context"

	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type PortalRepository struct {
	api external.PortalAPI
}

func NewPortalRepository(api external.PortalAPI) repository.PortalRepository {
	return &PortalRepository{api}
}

func (repo *PortalRepository) GetUsers(ctx context.Context) ([]*domain.PortalUser, error) {
	users, err := repo.api.GetAll()
	if err != nil {
		return nil, err
	}
	result := make([]*domain.PortalUser, 0, len(users))
	for _, v := range users {
		result = append(result, &domain.PortalUser{
			ID:             v.TraQID,
			Name:           v.RealName,
			AlphabeticName: v.AlphabeticName,
		})
	}
	return result, nil
}

func (repo *PortalRepository) GetUser(ctx context.Context, name string) (*domain.PortalUser, error) {
	users, err := repo.api.GetAll()
	if err != nil {
		return nil, err
	}

	for _, v := range users {
		if v.TraQID == name {
			return &domain.PortalUser{
				ID:             v.TraQID,
				Name:           v.RealName,
				AlphabeticName: v.AlphabeticName,
			}, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (repo *PortalRepository) MakeUserMp() (map[string]*domain.PortalUser, error) {
	users, err := repo.api.GetAll()
	if err != nil {
		return nil, err
	}

	mp := make(map[string]*domain.PortalUser, len(users))

	for _, v := range users {
		mp[v.TraQID] = &domain.PortalUser{
			ID:             v.TraQID,
			Name:           v.RealName,
			AlphabeticName: v.AlphabeticName,
		}
	}
	return mp, nil
}

// Interface guards
var (
	_ repository.PortalRepository = (*PortalRepository)(nil)
)
