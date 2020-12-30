package service

import (
	"context"
	"time"

	"github.com/traPtitech/traPortfolio/usecases/repository"
)

// User Portfolioのレスポンスで使うユーザー情報
type User struct{
	ID          uint      `json:"id"`
	Description string    `json:"description"`
	Check       bool      `json:"check"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
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
	_, _ = s.traQ.GetUser(ctx, name)
	_, _ = s.portal.GetUser(ctx, name)
	user, _ := s.repo.Get(name)
	return User{
		ID: 			user.ID,
		Description: 	user.Description,
		Check: 			user.Check,
		Name: 			user.Name,
		CreatedAt: 		user.CreatedAt,
		UpdatedAt: 		user.UpdatedAt,
	}
}
