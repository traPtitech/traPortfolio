package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type ContestRepository struct {
	h database.SQLHandler
}

func NewContestRepository(sql database.SQLHandler) *ContestRepository {
	return &ContestRepository{h: sql}
}

func (repo *ContestRepository) Create(contest *model.Contest) (*model.Contest, error) {
	err := repo.h.Create(contest).Error()
	if err != nil {
		return nil, err
	}
	return contest, nil
}

func (repo *ContestRepository) Update(id uuid.UUID, changes map[string]interface{}) error {
	if id == uuid.Nil {
		return repository.ErrNilID
	}

	var (
		old model.Contest
		new model.Contest
	)

	tx := repo.h.Begin()
	if err := tx.First(&old, &model.Contest{ID: id}).Error(); err != nil {
		return err
	}
	if err := tx.Model(&old).Updates(changes).Error(); err != nil {
		return err
	}
	if err := tx.Where(&model.Contest{ID: id}).First(&new).Error(); err != nil {
		return err
	}
	tx.Commit()
	return nil
}
