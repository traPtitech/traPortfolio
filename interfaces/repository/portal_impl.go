package repository

import (
	"context"
	"fmt"

	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/external"
)

type PortalRepository struct {
	token string
	api   external.PortalAPI
}

type PortalToken string

func NewPortalRepository(api external.PortalAPI, portalToken PortalToken) *PortalRepository {
	return &PortalRepository{
		token: string(portalToken),
		api:   api,
	}
}

func (repo *PortalRepository) GetUsers(ctx context.Context) ([]*domain.PortalUser, error) {
	users, err := repo.api.GetAll(repo.token)
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
	users, err := repo.api.GetAll(repo.token)
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
	return nil, fmt.Errorf("not found")
}
