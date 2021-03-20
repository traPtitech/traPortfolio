package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type UserService struct {
	repo   repository.UserRepository
	traQ   repository.TraQRepository
	portal repository.PortalRepository
}

func NewUserService(userRepository repository.UserRepository, traQRepository repository.TraQRepository, portalRepository repository.PortalRepository) UserService {
	return UserService{
		repo:   userRepository,
		traQ:   traQRepository,
		portal: portalRepository,
	}
}

func (s *UserService) GetUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := s.repo.GetUsers()
	if err != nil {
		return nil, err
	}
	portalUsers, err := s.portal.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	idMap := make(map[string]uuid.UUID, len(users))
	for _, v := range users {
		idMap[v.Name] = v.ID
	}
	result := make([]*domain.User, 0, len(users))
	for _, v := range portalUsers {
		if id, ok := idMap[v.ID]; ok {
			result = append(result, &domain.User{
				ID:       id,
				Name:     v.ID,
				RealName: v.Name,
			})
		}
	}
	return result, nil
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*domain.UserDetail, error) {
	user, err := s.repo.GetUser(id)
	if err != nil {
		return nil, err
	}
	userAccounts, err := s.repo.GetAccounts(id)
	if err != nil {
		return nil, err
	}
	traQUser, err := s.traQ.GetUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	portalUser, err := s.portal.GetUser(ctx, user.Name)
	if err != nil {
		return nil, err
	}
	accounts := make([]domain.Account, 0, len(userAccounts))
	for _, v := range userAccounts {
		accounts = append(accounts, domain.Account{
			ID:          v.ID,
			Type:        v.Type,
			PrPermitted: v.Check,
		})
	}
	return &domain.UserDetail{
		ID:       user.ID,
		Name:     user.Name,
		RealName: portalUser.Name,
		State:    traQUser.State,
		Bio:      user.Description,
		Accounts: accounts,
	}, nil
}

func (s *UserService) Update(ctx context.Context, user *model.User) error {
	return s.repo.Update(user)
}
