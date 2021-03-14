package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
)

type UserRepository struct {
	database.SQLHandler
}

func NewUserRepository(sql database.SQLHandler) *UserRepository {
	return &UserRepository{SQLHandler: sql}
}

func (repo *UserRepository) GetUsers() (users []*model.User, err error) {
	err = repo.Find(&users).Error()
	return
}

func (repo *UserRepository) GetUser(id uuid.UUID) (*model.User, error) {
	user := model.User{ID: id}
	err := repo.First(&user).Error()
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepository) GetAccounts(id uuid.UUID) (accounts []*model.Account, err error) {
	err = repo.Find(&accounts, "user_id = ?", id).Error()
	return
}

func (repo *UserRepository) Update(u *model.User) error {
	err := repo.Model(&model.User{}).Updates(&u).Error()
	return err
}
