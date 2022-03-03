package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type ContestRepository struct {
	h      database.SQLHandler
	portal external.PortalAPI
}

func NewContestRepository(sql database.SQLHandler, portal external.PortalAPI) repository.ContestRepository {
	return &ContestRepository{h: sql, portal: portal}
}

func (repo *ContestRepository) GetContests() ([]*domain.Contest, error) {
	contests := make([]*model.Contest, 10)
	err := repo.h.Find(&contests).Error()
	if err != nil {
		return nil, convertError(err)
	}

	result := make([]*domain.Contest, 0, len(contests))

	for _, v := range contests {
		result = append(result, &domain.Contest{
			ID:        v.ID,
			Name:      v.Name,
			TimeStart: v.Since,
			TimeEnd:   v.Until,
		})
	}
	return result, nil
}

// Teamsは別途GetContestTeamsで取得するためここではnilのまま返す
func (repo *ContestRepository) GetContest(id uuid.UUID) (*domain.ContestDetail, error) {
	contest := &model.Contest{ID: id}
	err := repo.h.First(contest).Error()
	if err != nil {
		return nil, convertError(err)
	}

	res := &domain.ContestDetail{
		Contest: domain.Contest{
			ID:        contest.ID,
			Name:      contest.Name,
			TimeStart: contest.Since,
			TimeEnd:   contest.Until,
		},
		Link:        contest.Link,
		Description: contest.Description,
		// Teams:
	}

	return res, nil
}

func (repo *ContestRepository) CreateContest(args *repository.CreateContestArgs) (*domain.Contest, error) {
	contest := &model.Contest{
		ID:          uuid.Must(uuid.NewV4()),
		Name:        args.Name,
		Description: args.Description,
		Link:        args.Link.ValueOrZero(),
		Since:       args.Since,
		Until:       args.Until.ValueOrZero(),
	}

	err := repo.h.Create(contest).Error()
	if err != nil {
		return nil, convertError(err)
	}

	result := &domain.Contest{
		ID:        contest.ID,
		Name:      contest.Name,
		TimeStart: contest.Since,
		TimeEnd:   contest.Until,
	}

	return result, nil
}

func (repo *ContestRepository) UpdateContest(id uuid.UUID, args *repository.UpdateContestArgs) error {
	changes := map[string]interface{}{}
	if args.Name.Valid {
		changes["name"] = args.Name.String
	}
	if args.Description.Valid {
		changes["description"] = args.Description.String
	}
	if args.Link.Valid {
		changes["link"] = args.Link.String
	}
	if args.Since.Valid {
		changes["since"] = args.Since.Time
	}
	if args.Until.Valid {
		changes["until"] = args.Until.Time
	}

	if len(changes) == 0 {
		return nil
	}

	var (
		old model.Contest
		new model.Contest
	)

	err := repo.h.Transaction(func(tx database.SQLHandler) error {
		if err := tx.
			Where(&model.Contest{ID: id}).
			First(&old).
			Error(); err != nil {
			return convertError(err)
		}
		if err := tx.Model(&old).Updates(changes).Error(); err != nil {
			return convertError(err)
		}
		err := tx.
			Where(&model.Contest{ID: id}).
			First(&new).
			Error()

		return convertError(err)
	})
	if err != nil {
		return convertError(err)
	}
	return nil
}

func (repo *ContestRepository) DeleteContest(id uuid.UUID) error {
	err := repo.h.Transaction(func(tx database.SQLHandler) error {
		if err := repo.h.
			Where(&model.Contest{ID: id}).
			First(&model.Contest{}).
			Error(); err != nil {
			return convertError(err)
		}

		if err := tx.
			Where(&model.Contest{ID: id}).
			Delete(&model.Contest{}).
			Error(); err != nil {
			return convertError(err)
		}

		return nil
	})
	if err != nil {
		return convertError(err)
	}

	return nil
}

func (repo *ContestRepository) GetContestTeams(contestID uuid.UUID) ([]*domain.ContestTeam, error) {
	if err := repo.h.
		Where(&model.Contest{ID: contestID}).
		First(&model.Contest{}).
		Error(); err != nil {
		return nil, convertError(err)
	}

	teams := make([]*model.ContestTeam, 10)
	err := repo.h.
		Where(&model.ContestTeam{ContestID: contestID}).
		Find(&teams).
		Error()
	if err != nil {
		return nil, convertError(err)
	}
	result := make([]*domain.ContestTeam, 0, len(teams))
	for _, v := range teams {
		result = append(result, &domain.ContestTeam{
			ID:        v.ID,
			ContestID: v.ContestID,
			Name:      v.Name,
			Result:    v.Result,
		})
	}
	return result, nil
}

// Membersは別途GetContestTeamMembersで取得するためここではnilのまま返す
func (repo *ContestRepository) GetContestTeam(contestID uuid.UUID, teamID uuid.UUID) (*domain.ContestTeamDetail, error) {
	var team model.ContestTeam
	if err := repo.h.
		Where(&model.ContestTeam{ID: teamID, ContestID: contestID}).
		First(&team).
		Error(); err != nil {
		return nil, convertError(err)
	}

	res := &domain.ContestTeamDetail{
		ContestTeam: domain.ContestTeam{
			ID:        team.ID,
			ContestID: team.ContestID,
			Name:      team.Name,
			Result:    team.Result,
		},
		Link:        team.Link,
		Description: team.Description,
		// Members:
	}
	return res, nil
}

func (repo *ContestRepository) CreateContestTeam(contestID uuid.UUID, _contestTeam *repository.CreateContestTeamArgs) (*domain.ContestTeamDetail, error) {
	contestTeam := &model.ContestTeam{
		ID:          uuid.Must(uuid.NewV4()),
		ContestID:   contestID,
		Name:        _contestTeam.Name,
		Description: _contestTeam.Description,
		Result:      _contestTeam.Result.ValueOrZero(),
		Link:        _contestTeam.Link.ValueOrZero(),
	}

	err := repo.h.Create(contestTeam).Error()
	if err != nil {
		return nil, convertError(err)
	}
	result := &domain.ContestTeamDetail{
		ContestTeam: domain.ContestTeam{
			ID:        contestTeam.ID,
			ContestID: contestTeam.ContestID,
			Name:      contestTeam.Name,
			Result:    contestTeam.Result,
		},
		Link:        contestTeam.Link,
		Description: contestTeam.Description,
		Members:     nil,
	}
	return result, nil
}

func (repo *ContestRepository) UpdateContestTeam(teamID uuid.UUID, args *repository.UpdateContestTeamArgs) error {
	changes := map[string]interface{}{}
	if args.Name.Valid {
		changes["name"] = args.Name.String
	}
	if args.Description.Valid {
		changes["description"] = args.Description.String
	}
	if args.Link.Valid {
		changes["link"] = args.Link.String
	}
	if args.Result.Valid {
		changes["result"] = args.Result.String
	}

	if len(changes) == 0 {
		return nil
	}

	var (
		old model.ContestTeam
		new model.ContestTeam
	)

	err := repo.h.Transaction(func(tx database.SQLHandler) error {
		if err := tx.
			Where(&model.ContestTeam{ID: teamID}).
			First(&old).
			Error(); err != nil {
			return convertError(err)
		}
		if err := tx.Model(&old).Updates(changes).Error(); err != nil {
			return convertError(err)
		}
		err := tx.
			Where(&model.ContestTeam{ID: teamID}).
			First(&new).
			Error()

		return convertError(err)
	})
	if err != nil {
		return convertError(err)
	}
	return nil
}

func (repo *ContestRepository) DeleteContestTeam(contestID uuid.UUID, teamID uuid.UUID) error {
	err := repo.h.
		Where(&model.ContestTeam{ID: teamID}).
		First(&model.ContestTeam{}).
		Error()
	if err != nil {
		return convertError(err)
	}

	err = repo.h.Transaction(func(tx database.SQLHandler) error {
		err = tx.
			Where(&model.ContestTeam{ID: teamID}).
			Delete(&model.ContestTeam{}).
			Error()
		if err != nil {
			return convertError(err)
		}
		return nil
	})

	if err != nil {
		return convertError(err)
	}

	return nil
}

func (repo *ContestRepository) GetContestTeamMembers(contestID uuid.UUID, teamID uuid.UUID) ([]*domain.User, error) {
	var belongings []*model.ContestTeamUserBelonging
	err := repo.h.
		Preload("User").
		Where(&model.ContestTeamUserBelonging{TeamID: teamID}).
		Find(&belongings).
		Error()
	if err != nil {
		return nil, convertError(err)
	}

	nameMap, err := repo.makeUserNameMap()
	if err != nil {
		return nil, convertError(err)
	}

	result := make([]*domain.User, len(belongings))
	for i, v := range belongings {
		u := v.User
		newUser := domain.User{
			ID:   u.ID,
			Name: u.Name,
		}

		if rn, ok := nameMap[u.Name]; ok {
			newUser.RealName = rn
		}

		result[i] = &newUser
	}
	return result, nil
}

func (repo *ContestRepository) AddContestTeamMembers(teamID uuid.UUID, members []uuid.UUID) error {
	if len(members) == 0 {
		return repository.ErrInvalidArg
	}

	// 存在チェック
	err := repo.h.
		Where(&model.ContestTeam{ID: teamID}).
		First(&model.ContestTeam{}).
		Error()
	if err != nil {
		return convertError(err)
	}

	// 既に所属しているメンバーを検索
	belongingsMap := make(map[uuid.UUID]struct{}, len(members))
	_belongings := make([]*model.ContestTeamUserBelonging, 0, len(members))
	err = repo.h.
		Where(&model.ContestTeamUserBelonging{TeamID: teamID}).
		Find(&_belongings).
		Error()
	if err != nil {
		return convertError(err)
	}
	for _, v := range _belongings {
		belongingsMap[v.UserID] = struct{}{}
	}

	err = repo.h.Transaction(func(tx database.SQLHandler) error {
		for _, memberID := range members {
			if _, ok := belongingsMap[memberID]; ok {
				continue
			}
			err = tx.Create(&model.ContestTeamUserBelonging{TeamID: teamID, UserID: memberID}).Error()
			if err != nil {
				return convertError(err)
			}
		}
		return nil
	})
	if err != nil {
		return convertError(err)
	}
	return nil

}

func (repo *ContestRepository) DeleteContestTeamMembers(teamID uuid.UUID, members []uuid.UUID) error {
	// 存在チェック
	err := repo.h.
		Where(&model.ContestTeam{ID: teamID}).
		First(&model.ContestTeam{}).
		Error()
	if err != nil {
		return convertError(err)
	}

	belongings := make(map[uuid.UUID]struct{}, len(members))
	_belongings := make([]*model.ContestTeamUserBelonging, 0, len(members))
	err = repo.h.
		Where(&model.ContestTeamUserBelonging{TeamID: teamID}).
		Find(&_belongings).
		Error()
	if err != nil {
		return convertError(err)
	}
	for _, v := range _belongings {
		belongings[v.UserID] = struct{}{}
	}

	err = repo.h.Transaction(func(tx database.SQLHandler) error {
		for _, memberID := range members {
			if _, ok := belongings[memberID]; ok {
				err = tx.
					Where(&model.ContestTeamUserBelonging{TeamID: teamID, UserID: memberID}).
					Delete(&model.ContestTeamUserBelonging{}).
					Error()
				if err != nil {
					return convertError(err)
				}
			}
		}
		return nil
	})
	if err != nil {
		return convertError(err)
	}
	return nil

}

func (repo *ContestRepository) makeUserNameMap() (map[string]string, error) {
	users, err := repo.portal.GetAll()
	if err != nil {
		return nil, convertError(err)
	}

	mp := make(map[string]string, len(users))

	for _, v := range users {
		mp[v.TraQID] = v.RealName
	}
	return mp, nil
}

// Interface guards
var (
	_ repository.ContestRepository = (*ContestRepository)(nil)
)
