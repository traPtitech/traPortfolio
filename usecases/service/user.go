//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type UserService interface {
	GetUsers(ctx context.Context, args *repository.GetUsersArgs) ([]*domain.User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*domain.UserDetail, error)
	Update(ctx context.Context, userID uuid.UUID, args *repository.UpdateUserArgs) error
	GetAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) (*domain.Account, error)
	GetAccounts(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error)
	CreateAccount(ctx context.Context, userID uuid.UUID, account *repository.CreateAccountArgs) (*domain.Account, error)
	EditAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, args *repository.UpdateAccountArgs) error
	DeleteAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error
	GetUserProjects(ctx context.Context, userID uuid.UUID) ([]*domain.UserProject, error)
	GetUserContests(ctx context.Context, userID uuid.UUID) ([]*domain.UserContest, error)
	GetGroupsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.UserGroup, error)
	GetUserEvents(ctx context.Context, userID uuid.UUID) ([]*domain.Event, error)
}

type userService struct {
	user  repository.UserRepository
	event repository.EventRepository
}

func NewUserService(userRepository repository.UserRepository, eventRepository repository.EventRepository) UserService {
	return &userService{user: userRepository, event: eventRepository}
}

func (s *userService) GetUsers(ctx context.Context, args *repository.GetUsersArgs) ([]*domain.User, error) {
	users, err := s.user.GetUsers(ctx, args)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) GetUser(ctx context.Context, userID uuid.UUID) (*domain.UserDetail, error) {
	user, err := s.user.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Update(ctx context.Context, userID uuid.UUID, args *repository.UpdateUserArgs) error {
	return s.user.UpdateUser(ctx, userID, args)
}

func (s *userService) GetAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) (*domain.Account, error) {
	return s.user.GetAccount(ctx, userID, accountID)
}

func (s *userService) GetAccounts(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error) {
	return s.user.GetAccounts(ctx, userID)
}

func (s *userService) CreateAccount(ctx context.Context, userID uuid.UUID, account *repository.CreateAccountArgs) (*domain.Account, error) {
	return s.user.CreateAccount(ctx, userID, account)
}

func (s *userService) EditAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, args *repository.UpdateAccountArgs) error {
	return s.user.UpdateAccount(ctx, userID, accountID, args)
}

func (s *userService) DeleteAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error {
	//TODO
	/*userのaccount.type番目のアカウントを削除する処理をしたい*/

	return s.user.DeleteAccount(ctx, userID, accountID)
}

func (s *userService) GetUserProjects(ctx context.Context, userID uuid.UUID) ([]*domain.UserProject, error) {
	projects, err := s.user.GetProjects(ctx, userID)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (s *userService) GetUserContests(ctx context.Context, userID uuid.UUID) ([]*domain.UserContest, error) {
	contests, err := s.user.GetContests(ctx, userID)
	if err != nil {
		return nil, err
	}
	return contests, nil
}

func (s *userService) GetGroupsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.UserGroup, error) {
	groups, err := s.user.GetGroupsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (s *userService) GetUserEvents(ctx context.Context, userID uuid.UUID) ([]*domain.Event, error) {
	events, err := s.event.GetUserEvents(ctx, userID)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// Interface guards
var (
	_ UserService = (*userService)(nil)
)
