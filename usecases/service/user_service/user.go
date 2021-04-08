package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

// User Portfolioのレスポンスで使うユーザー情報
type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	RealName string    `json:"realName"`
}

type UserDetail struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	RealName string    `json:"realName"`
	State    uint8     `json:"state"`
	Bio      string    `json:"bio"`
	Accounts []Account `json:"accounts"`
}

type Account struct {
	ID          string `json:"id"`
	Type        uint   `json:"type"`
	URL         string `gorm:"type:text"`
	PrPermitted bool   `json:"prPermitted"`
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
	traqUser, err := s.traQ.GetUser(ctx, user.ID)
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
			ID:          v.Name,
			Type:        v.Type,
			PrPermitted: v.Check,
		})
	}
	return &UserDetail{
		ID:       user.ID,
		Name:     user.Name,
		RealName: portalUser.Name,
		State:    traqUser.State,
		Bio:      user.Description,
		Accounts: accounts,
	}, nil
}

func (s *UserService) CreateAccount(ctx context.Context, id uuid.UUID, account *repository.CreateAccountArgs) (*domain.Account, error) {

	/*userのaccount.type番目のアカウントを追加する処理をしたい*/

	if id == uuid.Nil {
		return nil, repository.ErrInvalidArg
	}

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

func (s *UserService) DeleteAccount(ctx context.Context, accountid uuid.UUID, userid uuid.UUID) error {

	//TODO
	/*userのaccount.type番目のアカウントを削除する処理をしたい*/

	if accountid == uuid.Nil {
		return repository.ErrInvalidArg
	}

	if userid == uuid.Nil {
		return repository.ErrInvalidArg
	}

	err := s.repo.DeleteAccount(accountid, userid)

	return err

}
