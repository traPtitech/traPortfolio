package service

import (
	"context"

	"github.com/traPtitech/traPortfolio/usecases/repository"
)

// User Portfolioのレスポンスで使うユーザー情報
type User struct{}

type UserService struct {
	repo   repository.UserRepository
	traQ   repository.TraQRepository
	portal repository.PortalRepository
}

func NewUserService(userRepository repository.UserRepository, traQRepository repository.TraQRepository, portalRepository repository.PortalRepository) *UserService {
	return &UserService{
		repo:   userRepository,
		traQ:   traQRepository,
		portal: portalRepository,
	}
}

func (s *UserService) GetUser(ctx context.Context, name string) User {
	_, _ = s.traQ.GetUser(ctx, name)
	_, _ = s.portal.GetUser(ctx, name)
	_, _ = s.repo.Get(name)
	return User{}
}
