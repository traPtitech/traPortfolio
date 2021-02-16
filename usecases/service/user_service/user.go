package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

// User Portfolioのレスポンスで使うユーザー情報
type UserDetail struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	RealName string    `json:"realName"`
	State    uint8     `json:"state"`
	Bio      string    `json:"bio"`
	Accounts []Account `json:"accounts"`
}

type Account struct {
	ID          uuid.UUID `json:"id"`
	Type        uuid.UUID `json:"type"`
	PrPermitted bool      `json:"prPermitted"`
}

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

func (s *UserService) GetUsers(ctx context.Context) ([]UserDetail, error) {
	users, err := s.repo.GetUsers()
	if err != nil {
		return nil, err
	}
	result := make([]UserDetail, 0, len(users))
	for _, v := range users {
		portalUser, err := s.portal.GetUser(ctx, v.Name)
		if err != nil {
			return nil, err
		}
		result = append(result, UserDetail{
			ID:       v.ID,
			Name:     portalUser.Name,
			RealName: portalUser.AlphabeticName,
		})
	}
	return result, nil
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (UserDetail, error) {
	user, err := s.repo.GetUser(id)
	if err != nil {
		return UserDetail{}, err
	}
	userAccounts, err := s.repo.GetAccounts(id)
	if err != nil {
		return UserDetail{}, err
	}
	traqUser, err := s.traQ.GetUser(ctx, user.Name)
	if err != nil {
		return UserDetail{}, err
	}
	portalUser, err := s.portal.GetUser(ctx, user.Name)
	if err != nil {
		return UserDetail{}, err
	}
	accounts := make([]Account, 0, len(userAccounts))
	for _, v := range userAccounts {
		accounts = append(accounts, Account{
			ID:          v.ID,
			Type:        v.Type,
			PrPermitted: v.Check,
		})
	}
	return UserDetail{
		ID:       user.ID,
		Name:     user.Name,
		RealName: portalUser.Name,
		State:    traqUser.State,
		Bio:      user.Description,
		Accounts: accounts,
	}, nil
}
