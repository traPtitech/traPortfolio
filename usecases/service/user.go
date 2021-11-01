package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type UserService struct {
	repo  repository.UserRepository
	event repository.EventRepository
}

func NewUserService(userRepository repository.UserRepository, eventRepository repository.EventRepository) UserService {
	return UserService{repo: userRepository, event: eventRepository}
}

func (s *UserService) GetUsers(ctx context.Context) ([]*domain.User, error) {
	users, err := s.repo.GetUsers()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*domain.UserDetail, error) {
	user, err := s.repo.GetUser(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Update(ctx context.Context, id uuid.UUID, args *repository.UpdateUserArgs) error {
	changes := map[string]interface{}{}
	if args.Description.Valid {
		changes["description"] = args.Description.String
	}
	if args.Check.Valid {
		changes["check"] = args.Check.Bool
	}
	if len(changes) > 0 {
		err := s.repo.Update(id, changes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *UserService) GetAccount(userID uuid.UUID, accountID uuid.UUID) (*domain.Account, error) {
	return s.repo.GetAccount(userID, accountID)
}

func (s *UserService) GetAccounts(userID uuid.UUID) ([]*domain.Account, error) {
	return s.repo.GetAccounts(userID)
}

func (s *UserService) CreateAccount(ctx context.Context, id uuid.UUID, account *repository.CreateAccountArgs) (*domain.Account, error) {

	/*userのaccount.type番目のアカウントを追加する処理をしたい*/

	if len(account.ID) == 0 {
		return nil, repository.ErrInvalidArg
	}

	if account.Type >= domain.AccountLimit {
		return nil, repository.ErrInvalidArg
	}

	//implに実装は書く
	//accountの構造体たりないので補う
	//ここらへんのコメントアウトはリファクタのときにでも消す

	return s.repo.CreateAccount(id, account)

}

func (s *UserService) EditAccount(ctx context.Context, accountID uuid.UUID, userID uuid.UUID, args *repository.UpdateAccountArgs) error {
	changes := map[string]interface{}{}
	if args.Name.Valid {
		changes["name"] = args.Name.String
	}
	if args.URL.Valid {
		changes["url"] = args.URL.String
	}
	if args.PrPermitted.Valid {
		changes["check"] = args.PrPermitted.Bool
	}
	if args.Type.Valid {
		changes["type"] = args.Type.Int64
	}
	if len(changes) > 0 {
		err := s.repo.UpdateAccount(accountID, userID, changes)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *UserService) DeleteAccount(ctx context.Context, accountid uuid.UUID, userid uuid.UUID) error {

	//TODO
	/*userのaccount.type番目のアカウントを削除する処理をしたい*/

	err := s.repo.DeleteAccount(accountid, userid)

	return err

}

func (s *UserService) GetUserProjects(ctx context.Context, userID uuid.UUID) ([]*domain.UserProject, error) {
	projects, err := s.repo.GetProjects(userID)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (s *UserService) GetUserContests(ctx context.Context, userID uuid.UUID) ([]*domain.UserContest, error) {
	contests, err := s.repo.GetContests(userID)
	if err != nil {
		return nil, err
	}
	return contests, nil
}

func (s *UserService) GetUserEvents(ctx context.Context, userID uuid.UUID) ([]*domain.Event, error) {
	events, err := s.event.GetUserEvents(userID)
	if err != nil {
		return nil, err
	}
	return events, nil
}
