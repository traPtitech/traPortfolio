package repository

import (
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
)

type UserRepository struct {
	database.SQLHandler
}

func NewUserRepository(sql database.SQLHandler) *UserRepository {
	return &UserRepository{SQLHandler: sql}
}

func (repo *UserRepository) FindByID(id int) (user domain.User, err error) {
	if err = repo.Find(&user, id).Error(); err != nil {
		return
	}
	return
}

func (repo *UserRepository) FindAll() (users []domain.User, err error) {
	if err = repo.Find(&users).Error(); err != nil {
		return
	}
	return
}

func (repo *UserRepository) Store(u domain.User) (user domain.User, err error) {
	if err = repo.Create(&u).Error(); err != nil {
		return
	}
	user = u
	return
}

func (repo *UserRepository) Update(u domain.User) (user domain.User, err error) {
	if err = repo.Save(&user).Error(); err != nil {
		return
	}
	return
}

func (repo *UserRepository) DeleteByID(id int) (err error) {
	user := domain.User{}
	if err = repo.Find(&user, id).Error(); err != nil {
		return
	}
	if err = repo.Delete(&user).Error(); err != nil {
		return
	}
	return
}
