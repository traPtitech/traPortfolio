package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
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
