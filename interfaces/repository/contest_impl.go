package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
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

func (repo *ContestRepository) CreateContest(args *repository.CreateContestArgs) (*domain.Contest, error) {
	contest := model.Contest{
		ID:          uuid.Must(uuid.NewV4()),
		Name:        args.Name,
		Description: args.Description,
		Link:        args.Link,
		Since:       args.Since,
		Until:       args.Until,
	}
	err := repo.h.Create(contest).Error()
	if err != nil {
		return nil, err
	}

	result := &domain.Contest{
		ID:        contest.ID,
		Name:      contest.Name,
		TimeStart: contest.Since,
		TimeEnd:   contest.Until,
		CreatedAt: contest.CreatedAt,
		UpdatedAt: contest.UpdatedAt,
	}

	return result, nil
}

func (repo *ContestRepository) UpdateContest(id uuid.UUID, changes map[string]interface{}) error {
	if id == uuid.Nil {
		return repository.ErrNilID
	}

	var (
		old model.Contest
		new model.Contest
	)

	err := repo.h.Transaction(func(tx database.SQLHandler) error {
		if err := tx.First(&old, &model.Contest{ID: id}).Error(); err != nil {
			return err
		}
		if err := tx.Model(&old).Updates(changes).Error(); err != nil {
			return err
		}
		if err := tx.Where(&model.Contest{ID: id}).First(&new).Error(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (repo *ContestRepository) CreateContestTeam(contestID uuid.UUID, _contestTeam *repository.CreateContestTeamArgs) (*domain.ContestTeamDetail, error) {
	contestTeam := model.ContestTeam{
		ID:          uuid.Must(uuid.NewV4()),
		ContestID:   contestID,
		Name:        _contestTeam.Name,
		Description: _contestTeam.Description,
		Result:      _contestTeam.Result,
		Link:        _contestTeam.Link,
	}
	err := repo.h.Create(contestTeam).Error()
	if err != nil {
		return nil, err
	}
	result := &domain.ContestTeamDetail{
		ContestTeam: domain.ContestTeam{
			ID:        contestTeam.ID,
			ContestID: contestTeam.ContestID,
			Name:      contestTeam.Name,
			Result:    contestTeam.Result,
			CreatedAt: contestTeam.CreatedAt,
			UpdatedAt: contestTeam.UpdatedAt,
		},
		Link:        contestTeam.Link,
		Description: contestTeam.Description,
		Members:     nil,
	}
	return result, nil
}

//func (repo *ContestRepository) getTeamMember(teamID uuid.UUID) ([]*model.User, error) {
//	members := make([]*model.User, 0)
//	userTableName := (&model.User{}).TableName()
//	relationTableName := (&model.ContestTeamUserBelonging{}).TableName()
//
//	err := repo.h.Model(&model.User{}).
//		Joins(fmt.Sprintf("INNER JOIN %s ON %s.user_id = %s.id", relationTableName, relationTableName, userTableName)).
//		Where(fmt.Sprintf("%s.team_id = ?", relationTableName), teamID).
//		Find(&members).
//		Error()
//
//	if err != nil {
//		return nil, err
//	}
//	return members, nil
//}

func (repo *ContestRepository) UpdateContestTeam(teamID uuid.UUID, changes map[string]interface{}) error {
	if teamID == uuid.Nil {
		return repository.ErrNilID
	}

	var (
		old model.Contest
		new model.Contest
	)

	err := repo.h.Transaction(func(tx database.SQLHandler) error {
		if err := tx.First(&old, &model.ContestTeam{ID: teamID}).Error(); err != nil {
			return err
		}
		if err := tx.Model(&old).Updates(changes).Error(); err != nil {
			return err
		}
		if err := tx.Where(&model.ContestTeam{ID: teamID}).First(&new).Error(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (repo *ContestRepository) AddContestTeamMember(teamID uuid.UUID, members []uuid.UUID) error {
	if teamID == uuid.Nil {
		return repository.ErrNilID
	}

	// 存在チェック
	err := repo.h.First(&model.ContestTeam{}, &model.ContestTeam{ID: teamID}).Error()
	if err != nil {
		return err
	}

	curMp := make(map[uuid.UUID]struct{}, len(members))
	_cur := make([]*model.ContestTeamUserBelonging, 0, len(members))
	err = repo.h.Where(&model.ContestTeamUserBelonging{TeamID: teamID}).Find(_cur).Error()
	if err != nil {
		return err
	}
	for _, v := range _cur {
		curMp[v.UserID] = struct{}{}
	}

	err = repo.h.Transaction(func(tx database.SQLHandler) error {
		for _, memberID := range members {
			if _, ok := curMp[memberID]; ok {
				continue
			}
			err = tx.Create(model.ContestTeamUserBelonging{TeamID: teamID, UserID: memberID}).Error()
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil

}

func (repo *ContestRepository) DeleteContestTeamMember(teamID uuid.UUID, members []uuid.UUID) error {
	if teamID == uuid.Nil {
		return repository.ErrNilID
	}

	// 存在チェック
	err := repo.h.First(&model.ContestTeam{}, &model.ContestTeam{ID: teamID}).Error()
	if err != nil {
		return err
	}

	curMp := make(map[uuid.UUID]struct{}, len(members))
	_cur := make([]*model.ContestTeamUserBelonging, 0, len(members))
	err = repo.h.Where(&model.ContestTeamUserBelonging{TeamID: teamID}).Find(_cur).Error()
	if err != nil {
		return err
	}
	for _, v := range _cur {
		curMp[v.UserID] = struct{}{}
	}

	err = repo.h.Transaction(func(tx database.SQLHandler) error {
		for _, memberID := range members {
			if _, ok := curMp[memberID]; ok {
				err = tx.Delete(&model.ContestTeamUserBelonging{}, &model.ContestTeamUserBelonging{TeamID: teamID, UserID: memberID}).Error()
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil

}

// Interface guards
var (
	_ repository.ContestRepository = (*ContestRepository)(nil)
)
