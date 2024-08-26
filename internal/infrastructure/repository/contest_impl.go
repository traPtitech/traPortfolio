package repository

import (
	"context"
	"errors"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/internal/domain"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/external"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/repository/model"
	"github.com/traPtitech/traPortfolio/internal/usecases/repository"
	"gorm.io/gorm"
)

type ContestRepository struct {
	h      *gorm.DB
	portal external.PortalAPI
}

func NewContestRepository(sql *gorm.DB, portal external.PortalAPI) repository.ContestRepository {
	return &ContestRepository{h: sql, portal: portal}
}

func (r *ContestRepository) GetContests(ctx context.Context) ([]*domain.Contest, error) {
	contests := make([]*model.Contest, 10)
	err := r.h.WithContext(ctx).Find(&contests).Error
	if err != nil {
		return nil, err
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

func (r *ContestRepository) GetContest(ctx context.Context, contestID uuid.UUID) (*domain.ContestDetail, error) {
	return r.getContest(ctx, contestID)
}

// Teamsは別途GetContestTeamsで取得するためここではnilのまま返す
func (r *ContestRepository) getContest(ctx context.Context, contestID uuid.UUID) (*domain.ContestDetail, error) {
	contest := new(model.Contest)
	if err := r.h.
		WithContext(ctx).
		Where(&model.Contest{ID: contestID}).
		First(contest).
		Error; err != nil {
		return nil, err
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

func (r *ContestRepository) CreateContest(ctx context.Context, args *repository.CreateContestArgs) (*domain.ContestDetail, error) {
	contest := &model.Contest{
		ID:          uuid.Must(uuid.NewV4()),
		Name:        args.Name,
		Description: args.Description,
		Link:        args.Link.ValueOrZero(),
		Since:       args.Since,
		Until:       args.Until.ValueOrZero(),
	}

	// 既に同名のコンテストが存在するか
	err := r.h.
		WithContext(ctx).
		Where(&model.Contest{Name: contest.Name}).
		First(&model.Contest{}).
		Error
	if err == nil {
		return nil, repository.ErrAlreadyExists
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	err = r.h.WithContext(ctx).Create(contest).Error
	if err != nil {
		return nil, err
	}

	result, err := r.getContest(ctx, contest.ID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ContestRepository) UpdateContest(ctx context.Context, contestID uuid.UUID, args *repository.UpdateContestArgs) error {
	changes := map[string]interface{}{}
	if v, ok := args.Name.V(); ok {
		changes["name"] = v
	}
	if v, ok := args.Description.V(); ok {
		changes["description"] = v
	}
	if v, ok := args.Link.V(); ok {
		changes["link"] = v
	}
	if v, ok := args.Since.V(); ok {
		changes["since"] = v
	}
	if v, ok := args.Until.V(); ok {
		changes["until"] = v
	}

	if len(changes) == 0 {
		return nil
	}

	var c model.Contest
	err := r.h.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			WithContext(ctx).
			Where(&model.Contest{ID: contestID}).
			First(&model.Contest{}).
			Error; err != nil {
			return err
		}

		if err := tx.
			WithContext(ctx).
			Model(&model.Contest{ID: contestID}).
			Updates(changes).
			Error; err != nil {
			return err
		}

		if err := tx.
			WithContext(ctx).
			Where(&model.Contest{ID: contestID}).
			First(&c).
			Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *ContestRepository) DeleteContest(ctx context.Context, contestID uuid.UUID) error {
	err := r.h.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			WithContext(ctx).
			Where(&model.Contest{ID: contestID}).
			First(&model.Contest{}).
			Error; err != nil {
			return err
		}

		if err := tx.
			WithContext(ctx).
			Where(&model.Contest{ID: contestID}).
			Delete(&model.Contest{}).
			Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *ContestRepository) GetContestTeams(ctx context.Context, contestID uuid.UUID) ([]*domain.ContestTeam, error) {
	//IDがcontestIDであるようなcontestが存在するかチェック
	if err := r.h.
		WithContext(ctx).
		Where(&model.Contest{ID: contestID}).
		First(&model.Contest{}).
		Error; err != nil {
		return nil, err
	}

	//ContestIDがcontestIDであるようなcontestTeamを10件まで列挙する
	teams := make([]*model.ContestTeam, 10)
	err := r.h.
		WithContext(ctx).
		Where(&model.ContestTeam{ContestID: contestID}).
		Find(&teams).
		Error
	if err != nil {
		return nil, err
	}

	//teamsの要素vについてTeamIDがv.IDである(TeamIDがteamsIDListに入っているID)ようなContestTeamUserBelongingを列挙する
	var teamsIDList = make([]uuid.UUID, len(teams))
	for i, v := range teams {
		teamsIDList[i] = v.ID
	}

	var belongings = []*model.ContestTeamUserBelonging{}
	err = r.h.
		WithContext(ctx).
		Preload("User").
		Where("team_id IN (?)", teamsIDList).
		Find(&belongings).
		Error
	if err != nil {
		return nil, err
	}

	belongingMap := make(map[uuid.UUID]([]*model.ContestTeamUserBelonging))
	for _, v := range belongings {
		belongingMap[v.TeamID] = append(belongingMap[v.TeamID], v)
	}

	realNameMap, err := external.GetRealNameMap(r.portal)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.ContestTeam, 0, len(teams))
	for _, v := range teams {
		members := make([]*domain.User, len(belongingMap[v.ID]))
		for i, w := range belongingMap[v.ID] {
			u := w.User
			members[i] = domain.NewUser(u.ID, u.Name, realNameMap[u.Name], u.Check)
		}

		result = append(result, &domain.ContestTeam{
			ContestTeamWithoutMembers: domain.ContestTeamWithoutMembers{
				ID:        v.ID,
				ContestID: v.ContestID,
				Name:      v.Name,
				Result:    v.Result,
			},
			Members: members,
		})
	}

	// if len(result) == 0 {
	// 	result = []*domain.ContestTeam{}
	// }

	return result, nil
}

// Membersは別途GetContestTeamMembersで取得するためここではnilのまま返す
func (r *ContestRepository) GetContestTeam(ctx context.Context, contestID uuid.UUID, teamID uuid.UUID) (*domain.ContestTeamDetail, error) {
	var team model.ContestTeam
	if err := r.h.
		WithContext(ctx).
		Where(&model.ContestTeam{ID: teamID, ContestID: contestID}).
		First(&team).
		Error; err != nil {
		return nil, err
	}

	var belongings []*model.ContestTeamUserBelonging
	err := r.h.
		WithContext(ctx).
		Preload("User").
		Where(&model.ContestTeamUserBelonging{TeamID: teamID}).
		Find(&belongings).
		Error
	if err != nil {
		return nil, err
	}

	realNameMap, err := external.GetRealNameMap(r.portal)
	if err != nil {
		return nil, err
	}

	members := make([]*domain.User, len(belongings))
	for i, w := range belongings {
		u := w.User
		members[i] = domain.NewUser(u.ID, u.Name, realNameMap[u.Name], u.Check)
	}

	res := &domain.ContestTeamDetail{
		ContestTeam: domain.ContestTeam{
			ContestTeamWithoutMembers: domain.ContestTeamWithoutMembers{
				ID:        team.ID,
				ContestID: team.ContestID,
				Name:      team.Name,
				Result:    team.Result,
			},
			Members: members,
		},
		Link:        team.Link,
		Description: team.Description,
	}
	return res, nil
}

func (r *ContestRepository) CreateContestTeam(ctx context.Context, contestID uuid.UUID, _contestTeam *repository.CreateContestTeamArgs) (*domain.ContestTeamDetail, error) {
	if err := r.h.
		WithContext(ctx).
		Where(&model.Contest{ID: contestID}).
		First(&model.Contest{}).
		Error; err != nil {
		return nil, err
	}

	contestTeam := &model.ContestTeam{
		ID:          uuid.Must(uuid.NewV4()),
		ContestID:   contestID,
		Name:        _contestTeam.Name,
		Description: _contestTeam.Description,
		Result:      _contestTeam.Result.ValueOrZero(),
		Link:        _contestTeam.Link.ValueOrZero(),
	}

	err := r.h.WithContext(ctx).Create(contestTeam).Error
	if err != nil {
		return nil, err
	}

	result := &domain.ContestTeamDetail{
		ContestTeam: domain.ContestTeam{
			ContestTeamWithoutMembers: domain.ContestTeamWithoutMembers{
				ID:        contestTeam.ID,
				ContestID: contestTeam.ContestID,
				Name:      contestTeam.Name,
				Result:    contestTeam.Result,
			},
			Members: make([]*domain.User, 0),
		},
		Link:        contestTeam.Link,
		Description: contestTeam.Description,
	}
	return result, nil
}

func (r *ContestRepository) UpdateContestTeam(ctx context.Context, teamID uuid.UUID, args *repository.UpdateContestTeamArgs) error {
	changes := map[string]interface{}{}
	if v, ok := args.Name.V(); ok {
		changes["name"] = v
	}
	if v, ok := args.Description.V(); ok {
		changes["description"] = v
	}
	if v, ok := args.Link.V(); ok {
		changes["link"] = v
	}
	if v, ok := args.Result.V(); ok {
		changes["result"] = v
	}

	if len(changes) == 0 {
		return nil
	}

	var ct model.ContestTeam
	err := r.h.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			WithContext(ctx).
			Where(&model.ContestTeam{ID: teamID}).
			First(&model.ContestTeam{}).
			Error; err != nil {
			return err
		}

		if err := tx.
			WithContext(ctx).
			Model(&model.ContestTeam{ID: teamID}).
			Updates(changes).
			Error; err != nil {
			return err
		}

		if err := tx.
			WithContext(ctx).
			Where(&model.ContestTeam{ID: teamID}).
			First(&ct).
			Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *ContestRepository) DeleteContestTeam(ctx context.Context, contestID uuid.UUID, teamID uuid.UUID) error {
	if err := r.h.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 存在確認
		if err := tx.
			WithContext(ctx).
			Where(&model.ContestTeam{ID: teamID, ContestID: contestID}).
			First(&model.ContestTeam{}).
			Error; err != nil {
			return err
		}

		if err := tx.
			WithContext(ctx).
			Where(&model.ContestTeam{ID: teamID}).
			Delete(&model.ContestTeam{}).
			Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (r *ContestRepository) GetContestTeamMembers(ctx context.Context, contestID uuid.UUID, teamID uuid.UUID) ([]*domain.User, error) {
	// 存在チェック
	err := r.h.
		WithContext(ctx).
		Where(&model.Contest{ID: contestID}).
		First(&model.Contest{}).
		Error
	if err != nil {
		return nil, err
	}
	err = r.h.
		WithContext(ctx).
		Where(&model.ContestTeam{ID: teamID}).
		First(&model.ContestTeam{}).
		Error
	if err != nil {
		return nil, err
	}

	var belongings []*model.ContestTeamUserBelonging
	err = r.h.
		WithContext(ctx).
		Preload("User").
		Where(&model.ContestTeamUserBelonging{TeamID: teamID}).
		Find(&belongings).
		Error
	if err != nil {
		return nil, err
	}

	realNameMap, err := external.GetRealNameMap(r.portal)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.User, len(belongings))
	for i, v := range belongings {
		u := v.User
		result[i] = domain.NewUser(u.ID, u.Name, realNameMap[u.Name], u.Check)
	}
	return result, nil
}

func (r *ContestRepository) EditContestTeamMembers(ctx context.Context, teamID uuid.UUID, members []uuid.UUID) error {
	// 存在チェック
	err := r.h.
		WithContext(ctx).
		Where(&model.ContestTeam{ID: teamID}).
		First(&model.ContestTeam{}).
		Error
	if err != nil {
		return err
	}

	belongings := make(map[uuid.UUID]struct{}, len(members))
	_belongings := make([]*model.ContestTeamUserBelonging, 0, len(members))
	err = r.h.
		WithContext(ctx).
		Where(&model.ContestTeamUserBelonging{TeamID: teamID}).
		Find(&_belongings).
		Error
	if err != nil {
		return err
	}
	for _, v := range _belongings {
		belongings[v.UserID] = struct{}{}
	}

	membersMap := make(map[uuid.UUID]struct{}, len(members))
	for _, v := range members {
		membersMap[v] = struct{}{}
	}

	err = r.h.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		//チームに所属していなくて渡された配列に入っているメンバーをチームに追加
		membersToBeAdded := make([]*model.ContestTeamUserBelonging, 0, len(members))
		for _, memberID := range members {
			if _, ok := belongings[memberID]; !ok {
				membersToBeAdded = append(membersToBeAdded, &model.ContestTeamUserBelonging{TeamID: teamID, UserID: memberID})
			}
		}
		if len(membersToBeAdded) > 0 {
			err = tx.WithContext(ctx).Create(&membersToBeAdded).Error
			if err != nil {
				return err
			}
		}
		//チームに所属していて渡された配列に入っていないメンバーをチームから削除
		membersToBeRemoved := make([]uuid.UUID, 0, len(members))
		for _, belonging := range _belongings {
			if _, ok := membersMap[belonging.UserID]; !ok {
				membersToBeRemoved = append(membersToBeRemoved, belonging.UserID)
			}
		}
		if len(membersToBeRemoved) > 0 {
			err = tx.
				WithContext(ctx).
				Where(&model.ContestTeamUserBelonging{TeamID: teamID}).
				Where("`contest_team_user_belongings`.`user_id` IN (?)", membersToBeRemoved).
				Delete(&model.ContestTeamUserBelonging{}).
				Error
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

// Interface guards
var (
	_ repository.ContestRepository = (*ContestRepository)(nil)
)
