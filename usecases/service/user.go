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
	GetAccount(userID uuid.UUID, accountID uuid.UUID) (*domain.Account, error)
	GetAccounts(userID uuid.UUID) ([]*domain.Account, error)
	CreateAccount(ctx context.Context, userID uuid.UUID, account *repository.CreateAccountArgs) (*domain.Account, error)
	EditAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, args *repository.UpdateAccountArgs) error
	DeleteAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error
	GetUserProjects(ctx context.Context, userID uuid.UUID) ([]*domain.UserProject, error)
	GetUserContests(ctx context.Context, userID uuid.UUID) ([]*domain.UserContest, error)
	GetGroupsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.GroupUser, error)
	GetUserEvents(ctx context.Context, userID uuid.UUID) ([]*domain.Event, error)
}

type userService struct {
	repo  repository.UserRepository
	event repository.EventRepository
}

func NewUserService(userRepository repository.UserRepository, eventRepository repository.EventRepository) UserService {
	return &userService{repo: userRepository, event: eventRepository}
}

func (s *userService) GetUsers(ctx context.Context, args *repository.GetUsersArgs) ([]*domain.User, error) {
	users, err := s.repo.GetUsers(args)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) GetUser(ctx context.Context, userID uuid.UUID) (*domain.UserDetail, error) {
	user, err := s.repo.GetUser(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Update(ctx context.Context, userID uuid.UUID, args *repository.UpdateUserArgs) error {
	if err := s.repo.UpdateUser(userID, args); err != nil {
		return err
	}

	return nil
}

func (s *userService) GetAccount(userID uuid.UUID, accountID uuid.UUID) (*domain.Account, error) {
	return s.repo.GetAccount(userID, accountID)
}

func (s *userService) GetAccounts(userID uuid.UUID) ([]*domain.Account, error) {
	return s.repo.GetAccounts(userID)
}

func (s *userService) CreateAccount(ctx context.Context, userID uuid.UUID, account *repository.CreateAccountArgs) (*domain.Account, error) {
	return s.repo.CreateAccount(userID, account)
}

func (s *userService) EditAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, args *repository.UpdateAccountArgs) error {
	if err := s.repo.UpdateAccount(userID, accountID, args); err != nil {
		return err
	}

	return nil
}

func (s *userService) DeleteAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error {

	//TODO
	/*userのaccount.type番目のアカウントを削除する処理をしたい*/

	err := s.repo.DeleteAccount(userID, accountID)

	return err

}

func (s *userService) GetUserProjects(ctx context.Context, userID uuid.UUID) ([]*domain.UserProject, error) {
	projects, err := s.repo.GetProjects(userID)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (s *userService) GetUserContests(ctx context.Context, userID uuid.UUID) ([]*domain.UserContest, error) {
	contests, err := s.repo.GetContests(userID)
	if err != nil {
		return nil, err
	}
	return contests, nil
}

func (s *userService) GetGroupsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.GroupUser, error) {
	groups, err := s.repo.GetGroupsByUserID(userID)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

func (s *userService) GetUserEvents(ctx context.Context, userID uuid.UUID) ([]*domain.Event, error) {
	events, err := s.event.GetUserEvents(userID)
	if err != nil {
		return nil, err
	}
	return events, nil
}

// Interface guards
var (
	_ UserService = (*userService)(nil)
)
