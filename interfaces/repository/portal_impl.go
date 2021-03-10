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

// type PortalToken string

func NewPortalRepository(api external.PortalAPI) *PortalRepository {
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
			ID:             v.ID,
			Name:           v.Name,
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
		if v.ID == name {
			return &domain.PortalUser{
				ID:             v.ID,
				Name:           v.Name,
				AlphabeticName: v.AlphabeticName,
			}, nil
		}
	}
	return nil, repository.ErrNotFound
}
