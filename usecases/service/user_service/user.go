package service

import (
	"context"

	"github.com/traPtitech/traPortfolio/usecases/repository"
)

// User Portfolioのレスポンスで使うユーザー情報
type User struct {
	ID       uint      `json:"id"`
	Name     string    `json:"name"`
	RealName string    `json:"realName"`
	State    uint8     `json:"state"`
	Bio      string    `json:"bio"`
	Accounts []Account `json:"accounts"`
}

type Account struct {
	ID          uint `json:"id"`
	Type        uint `json:"type"`
	PrPermitted bool `json:"prPermitted"`
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

func (s *UserService) GetUser(ctx context.Context, name string) User {
	traqUser, _ := s.traQ.GetUser(ctx, name)
	portalUser, _ := s.portal.GetUser(ctx, name)
	user, userAccounts, _ := s.repo.Get(name)
	accounts := make([]Account, len(userAccounts))
	for _, v := range userAccounts {
		accounts = append(accounts, Account{
			ID:          v.ID,
			Type:        v.Type,
			PrPermitted: v.Check,
		})
	}
	return User{
		ID:       user.ID,
		Name:     user.Name,
		RealName: portalUser.Name,
		State:    traqUser.State,
		Bio:      user.Description,
		Accounts: accounts,
	}
}
