package repository

import (
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
)

type ContestRepository struct {
	h database.SQLHandler
}

func NewContestRepository(sql database.SQLHandler) *ContestRepository {
	return &ContestRepository{h: sql}
}

func (repo *ContestRepository) Create(contest *domain.Contest) (*domain.Contest, error) {
	err := repo.h.Create(contest).Error()
	if err != nil {
		return nil, err
	}
	return contest, nil
}

func (repo *ContestRepository) Update(changes map[string]interface{}) error {
	err := repo.h.Updates(changes).Error()
	if err != nil {
		return err
	}
	return nil
}
