package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

// User Portfolioのレスポンスで使うユーザー情報
type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	RealName string    `json:"realName"`
}

type UserDetail struct {
	ID       uuid.UUID       `json:"id"`
	Name     string          `json:"name"`
	RealName string          `json:"realName"`
	State    model.TraQState `json:"state"`
	Bio      string          `json:"bio"`
	Accounts []Account       `json:"accounts"`
}

type Account struct {
	ID          uuid.UUID `json:"id"`
	Type        uint      `json:"type"`
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

func (s *UserService) GetUsers(ctx context.Context) ([]*User, error) {
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
	result := make([]*User, 0, len(users))
	for _, v := range portalUsers {
		if id, ok := idMap[v.ID]; ok {
			result = append(result, &User{
				ID:       id,
				Name:     v.ID,
				RealName: v.Name,
			})
		}
	}
	return result, nil
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*UserDetail, error) {
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
	accounts := make([]Account, 0, len(userAccounts))
	for _, v := range userAccounts {
		accounts = append(accounts, Account{
			ID:          v.ID,
			Type:        v.Type,
			PrPermitted: v.Check,
		})
	}
	return &UserDetail{
		ID:       user.ID,
		Name:     user.Name,
		RealName: portalUser.Name,
		State:    traQUser.State,
		Bio:      user.Description,
		Accounts: accounts,
	}, nil
}
