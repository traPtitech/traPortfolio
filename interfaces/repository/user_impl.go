package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type UserRepository struct {
	database.SQLHandler
}

func NewUserRepository(sql database.SQLHandler) *UserRepository {
	return &UserRepository{SQLHandler: sql}
}

func (repo *UserRepository) GetUsers() (users []*domain.User, err error) {
	err = repo.Find(&users).Error()
	return
}

func (repo *UserRepository) GetUser(id uuid.UUID) (*domain.User, error) {
	user := domain.User{ID: id}
	err := repo.First(&user).Error()
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) GetAccounts(id uuid.UUID) (accounts []*domain.Account, err error) {
	err = repo.Find(&accounts, "user_id = ?", id).Error()
	return
}

func (repo *UserRepository) Update(u *domain.User) error {
	err := repo.Model(&domain.User{}).Updates(&u).Error()
	return err
}

func (repo *UserRepository) CreateAccount(id uuid.UUID, account *repository.CreateAccountArgs) (*domain.Account, error) {
	account_ := domain.Account{
		ID:     uuid.Must(uuid.NewV4()),
		Type:   account.Type,
		Name:   account.ID,
		URL:    account.URL,
		UserID: id,
		Check:  account.PrPermitted,
	}
	err := repo.Create(account_).Error()
	if err != nil {
		return nil, err
	}

	var result *domain.Account

	err = repo.First(result, domain.Account{ID: account_.ID}).Error()

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *UserRepository) DeleteAccount(accountid uuid.UUID, userid uuid.UUID) error {

	err := repo.Delete(&domain.Account{}, &domain.Account{ID: accountid, UserID: userid}).Error()

	return err

}
