package repository

import (
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
)

type UserRepository struct {
	database.SqlHandler
}

func NewUserRepository(sql database.SqlHandler) UserRepository {
	return UserRepository{SqlHandler: sql}
}

func (repo *UserRepository) FindById(id int) (user domain.User, err error) {
	if err = repo.Find(&user, id).Error(); err != nil {
		return
	}
	return
}

func (repo *UserRepository) FindAll() (users domain.User, err error) {
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
	if err = repo.Save(&u).Error(); err != nil {
		return
	}
	user = u
	return
}

func (repo *UserRepository) DeleteById(id int) (err error) {
	user := domain.User{}
	if err = repo.Find(&user, id).Error(); err != nil {
		return
	}
	if err = repo.Delete(&user).Error(); err != nil {
		return
	}
	return
}
