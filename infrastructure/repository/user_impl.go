package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/samber/lo"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/infrastructure/external"
	"github.com/traPtitech/traPortfolio/infrastructure/repository/model"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	h      *gorm.DB
	portal external.PortalAPI
	traQ   external.TraQAPI
}

func NewUserRepository(h *gorm.DB, portalAPI external.PortalAPI, traQAPI external.TraQAPI) repository.UserRepository {
	return &UserRepository{
		h:      h,
		portal: portalAPI,
		traQ:   traQAPI,
	}
}

func (r *UserRepository) GetUsers(ctx context.Context, args *repository.GetUsersArgs) ([]*domain.User, error) {
	limit := args.Limit.ValueOr(-1)
	tx := r.h.WithContext(ctx).Limit(limit)
	includeSuspended, iok := args.IncludeSuspended.V()
	name, nok := args.Name.V()
	if iok && nok {
		return nil, fmt.Errorf("%w: you must not specify both includeSuspended and name", repository.ErrInvalidArg)
	} else if nok {
		tx = tx.Where(&model.User{Name: name})
	} else if !(iok && includeSuspended) {
		tx = tx.Where(&model.User{State: domain.TraqStateActive})
	}

	users := make([]*model.User, 0)
	if err := tx.Find(&users).Error; err != nil {
		return nil, err
	}

	if l := len(users); l == 0 {
		return []*domain.User{}, nil
	} else if l == 1 {
		portalUser, err := r.portal.GetUserByTraqID(users[0].Name)
		if err != nil {
			return nil, err
		}

		return []*domain.User{
			domain.NewUser(
				users[0].ID,
				users[0].Name,
				portalUser.RealName,
				users[0].Check,
			),
		}, nil
	} else {
		realNameMap, err := external.GetRealNameMap(r.portal)
		if err != nil {
			return nil, err
		}

		result := make([]*domain.User, 0, l)
		for _, v := range users {
			result = append(result, domain.NewUser(
				v.ID,
				v.Name,
				realNameMap[v.Name],
				v.Check,
			))
		}

		return result, nil
	}
}

func (r *UserRepository) SyncUsers(ctx context.Context) error {
	traqUsers, err := r.traQ.GetUsers(&external.TraQGetAllArgs{IncludeSuspended: true})
	if err != nil {
		return err
	}

	users := lo.FilterMap(traqUsers, func(u *external.TraQUserResponse, _ int) (*model.User, bool) {
		if u.Bot {
			return nil, false
		}

		return &model.User{
			ID:    u.ID,
			Name:  u.Name,
			State: u.State,
		}, true
	})

	err = r.h.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "state", "updated_at"}),
		}).
		Create(&users).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUser(ctx context.Context, userID uuid.UUID) (*domain.UserDetail, error) {
	user := new(model.User)
	err := r.h.
		WithContext(ctx).
		Preload("Accounts").
		Where(&model.User{ID: userID}).
		First(user).
		Error
	if err != nil {
		return nil, err
	}

	accounts := make([]*domain.Account, 0, len(user.Accounts))
	for _, v := range user.Accounts {
		accounts = append(accounts, &domain.Account{
			ID:          v.ID,
			DisplayName: v.Name,
			Type:        domain.AccountType(v.Type),
			PrPermitted: v.Check,
			URL:         v.URL,
		})
	}

	portalUser, err := r.portal.GetUserByTraqID(user.Name)
	if err != nil {
		return nil, err
	}

	result := domain.UserDetail{
		User: *domain.NewUser(
			user.ID,
			user.Name,
			portalUser.RealName,
			user.Check,
		),
		State:    user.State,
		Bio:      user.Description,
		Accounts: accounts,
	}

	return &result, nil
}

// TODO: テスト用にしか使われていないので消すかテスト用であることを明示する
func (r *UserRepository) CreateUser(ctx context.Context, args *repository.CreateUserArgs) (*domain.UserDetail, error) {
	portalUser, err := r.portal.GetUserByTraqID(args.Name)
	if err != nil {
		return nil, err
	}

	user := model.User{
		// TODO: traQのUUIDを使うべきかも
		ID:          uuid.Must(uuid.NewV4()),
		Description: args.Description,
		Check:       args.Check,
		Name:        args.Name,
		State:       domain.TraqStateActive,
	}

	err = r.h.WithContext(ctx).Create(&user).Error
	if err != nil {
		return nil, err
	}

	result := &domain.UserDetail{
		User: *domain.NewUser(
			user.ID,
			user.Name,
			portalUser.RealName,
			user.Check,
		),
		State:    0,
		Bio:      user.Description,
		Accounts: []*domain.Account{},
	}
	return result, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, userID uuid.UUID, args *repository.UpdateUserArgs) error {
	changes := map[string]interface{}{}
	if v, ok := args.Description.V(); ok {
		changes["description"] = v
	}
	if v, ok := args.Check.V(); ok {
		changes["check"] = v
	}

	if len(changes) == 0 {
		return nil
	}

	err := r.h.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		user := new(model.User)
		err := tx.
			WithContext(ctx).Where(&model.User{ID: userID}).
			First(user).
			Error
		if err != nil {
			return err
		}

		err = tx.WithContext(ctx).Model(user).Updates(changes).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetAccounts(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error) {
	err := r.h.
		WithContext(ctx).
		Where(&model.User{ID: userID}).
		First(&model.User{}).
		Error
	if err != nil {
		return nil, err
	}

	accounts := make([]*model.Account, 0)
	err = r.h.
		WithContext(ctx).
		Where(&model.Account{UserID: userID}).
		Find(&accounts).
		Error
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Account, 0, len(accounts))
	for _, v := range accounts {
		result = append(result, &domain.Account{
			ID:          v.ID,
			Type:        domain.AccountType(v.Type),
			PrPermitted: v.Check,
			DisplayName: v.Name,
			URL:         v.URL,
		})
	}
	return result, nil
}

func (r *UserRepository) GetAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) (*domain.Account, error) {
	account := &model.Account{}
	err := r.h.
		WithContext(ctx).
		Where(&model.Account{ID: accountID, UserID: userID}).
		First(account).
		Error
	if err != nil {
		return nil, err
	}

	result := &domain.Account{
		ID:          account.ID,
		Type:        domain.AccountType(account.Type),
		PrPermitted: account.Check,
		DisplayName: account.Name,
		URL:         account.URL,
	}

	return result, nil
}

func (r *UserRepository) CreateAccount(ctx context.Context, userID uuid.UUID, args *repository.CreateAccountArgs) (*domain.Account, error) {
	if !domain.IsValidAccountURL(args.Type, args.URL) {
		return nil, repository.ErrInvalidArg
	}

	if err := r.h.
		Where("`accounts`.`user_id` = ? AND `accounts`.`type` = ?", userID, uint8(args.Type)).
		First(&model.Account{}).
		Error; err == nil {
		return nil, repository.ErrAlreadyExists
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	account := model.Account{
		ID:     uuid.Must(uuid.NewV4()),
		Type:   uint8(args.Type),
		Name:   args.DisplayName,
		URL:    args.URL,
		UserID: userID,
		Check:  args.PrPermitted,
	}
	err := r.h.WithContext(ctx).Create(&account).Error
	if err != nil {
		return nil, err
	}

	ver := new(model.Account)
	if err := r.h.
		WithContext(ctx).
		Where(&model.Account{ID: account.ID}).
		First(ver).
		Error; err != nil {
		return nil, err
	}

	return &domain.Account{
		ID:          ver.ID,
		DisplayName: ver.Name,
		Type:        domain.AccountType(ver.Type),
		PrPermitted: ver.Check,
		URL:         ver.URL,
	}, nil
}

func (r *UserRepository) UpdateAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, args *repository.UpdateAccountArgs) error {
	changes := map[string]interface{}{}
	if v, ok := args.DisplayName.V(); ok {
		changes["name"] = v
	}
	if v, ok := args.URL.V(); ok {
		changes["url"] = v
	}
	if v, ok := args.PrPermitted.V(); ok {
		changes["check"] = v
	}
	if v, ok := args.Type.V(); ok {
		changes["type"] = v
	}

	if len(changes) == 0 {
		return nil
	}

	err := r.h.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		account := new(model.Account)
		err := tx.
			WithContext(ctx).Where(&model.Account{ID: accountID, UserID: userID}).
			First(account).
			Error
		if err != nil {
			return err
		}

		// 同タイプ生成回避
		av, aok := args.Type.V()
		if aok && av != domain.AccountType(account.Type) {
			if err := tx.
				WithContext(ctx).
				Where("`accounts`.`user_id` = ? AND `accounts`.`type` = ?", userID, uint8(av)).
				First(&model.Account{}).
				Error; err == nil {
				return repository.ErrAlreadyExists
			} else if !errors.Is(err, repository.ErrNotFound) {
				return err
			}
		}

		// URLのvalidation
		tv, tok := args.Type.V()
		uv, uok := args.URL.V()
		if tok && uok {
			if !domain.IsValidAccountURL(domain.AccountType(tv), uv) {
				return repository.ErrInvalidArg
			}
		} else if !tok && uok {
			if !domain.IsValidAccountURL(domain.AccountType(account.Type), uv) {
				return repository.ErrInvalidArg
			}
		} else if tok && !uok {
			if !domain.IsValidAccountURL(tv, account.URL) {
				return repository.ErrInvalidArg
			}
		}

		err = tx.WithContext(ctx).Model(account).Updates(changes).Error
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (r *UserRepository) DeleteAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error {
	if err := r.h.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			WithContext(ctx).Where(&model.Account{ID: accountID, UserID: userID}).
			First(&model.Account{}).
			Error; err != nil {
			return err
		}

		if err := tx.
			WithContext(ctx).Where(&model.Account{ID: accountID, UserID: userID}).
			Delete(&model.Account{}).
			Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetProjects(ctx context.Context, userID uuid.UUID) ([]*domain.UserProject, error) {
	err := r.h.
		WithContext(ctx).
		Where(&model.User{ID: userID}).
		First(&model.User{}).
		Error
	if err != nil {
		return nil, err
	}

	projects := make([]*model.ProjectMember, 0)
	err = r.h.
		WithContext(ctx).
		Preload("Project").
		Where(&model.ProjectMember{UserID: userID}).
		Find(&projects).
		Error
	if err != nil {
		return nil, err
	}

	res := make([]*domain.UserProject, 0, len(projects))
	for _, v := range projects {
		p := v.Project
		res = append(res, &domain.UserProject{
			ID:           v.Project.ID,
			Name:         v.Project.Name,
			Duration:     domain.NewYearWithSemesterDuration(p.SinceYear, p.SinceSemester, p.UntilYear, p.UntilSemester),
			UserDuration: domain.NewYearWithSemesterDuration(v.SinceYear, v.SinceSemester, v.UntilYear, v.UntilSemester),
		})
	}
	return res, nil
}

func (r *UserRepository) GetGroupsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.UserGroup, error) {
	err := r.h.
		WithContext(ctx).
		Where(&model.User{ID: userID}).
		First(&model.User{}).
		Error
	if err != nil {
		return nil, err
	}

	groups := make([]*model.GroupUserBelonging, 0)
	err = r.h.
		WithContext(ctx).
		Preload("Group").
		Where(&model.GroupUserBelonging{UserID: userID}).
		Find(&groups).
		Error
	if err != nil {
		return nil, err
	}

	result := make([]*domain.UserGroup, 0, len(groups))
	for _, v := range groups {
		gr := v.Group
		result = append(result, &domain.UserGroup{
			ID:   gr.GroupID,
			Name: gr.Name,
			Duration: domain.NewYearWithSemesterDuration(
				v.SinceYear,
				v.SinceSemester,
				v.UntilYear,
				v.UntilSemester,
			),
		})
	}
	return result, nil
}

func (r *UserRepository) GetContests(ctx context.Context, userID uuid.UUID) ([]*domain.UserContest, error) {
	err := r.h.
		WithContext(ctx).
		Where(&model.User{ID: userID}).
		First(&model.User{}).
		Error
	if err != nil {
		return nil, err
	}

	contestTeamUserBelongings := make([]*model.ContestTeamUserBelonging, 0)
	err = r.h.
		WithContext(ctx).
		Preload("ContestTeam.Contest").
		Where(&model.ContestTeamUserBelonging{UserID: userID}).
		Find(&contestTeamUserBelongings).
		Error
	if err != nil {
		return nil, err
	}

	contestsMap := make(map[uuid.UUID]*domain.UserContest)
	for _, v := range contestTeamUserBelongings {
		ct := v.ContestTeam
		if _, ok := contestsMap[ct.ContestID]; !ok {
			contestsMap[ct.ContestID] = &domain.UserContest{
				ID:        ct.Contest.ID,
				Name:      ct.Contest.Name,
				TimeStart: ct.Contest.Since,
				TimeEnd:   ct.Contest.Until,
				Teams:     []*domain.ContestTeamWithoutMembers{},
			}
		}
	}

	for _, v := range contestTeamUserBelongings {
		if userID == v.UserID {
			ct := v.ContestTeam
			contestsMap[ct.ContestID].Teams = append(contestsMap[ct.ContestID].Teams, &domain.ContestTeamWithoutMembers{
				ID:        ct.ID,
				ContestID: ct.ContestID,
				Name:      ct.Name,
				Result:    ct.Result,
			})
		}
	}

	res := make([]*domain.UserContest, 0, len(contestTeamUserBelongings))
	for _, v := range contestsMap {
		res = append(res, v)
	}

	return res, nil
}

// Interface guards
var (
	_ repository.UserRepository = (*UserRepository)(nil)
)
