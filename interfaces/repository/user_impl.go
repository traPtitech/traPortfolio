package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type UserRepository struct {
	h      database.SQLHandler
	portal external.PortalAPI
	traQ   external.TraQAPI
}

func NewUserRepository(h database.SQLHandler, portalAPI external.PortalAPI, traQAPI external.TraQAPI) repository.UserRepository {
	return &UserRepository{
		h:      h,
		portal: portalAPI,
		traQ:   traQAPI,
	}
}

func makeTraqGetAllArgs(rargs *repository.GetUsersArgs) (*external.TraQGetAllArgs, error) {
	eargs := new(external.TraQGetAllArgs)
	if iv, nv := rargs.IncludeSuspended.Valid, rargs.Name.Valid; iv && nv {
		// Ref: https://github.com/traPtitech/traQ/blob/fa8cdf17d7b4869bfb7d0864873cd3c46b7543b2/router/v3/users.go#L31-L33
		return nil, repository.ErrInvalidArg
	} else if iv {
		eargs.IncludeSuspended = rargs.IncludeSuspended.Bool
	} else if nv {
		eargs.Name = rargs.Name.String
	}

	return eargs, nil
}

func (repo *UserRepository) GetUsers(args *repository.GetUsersArgs) ([]*domain.User, error) {
	eargs, err := makeTraqGetAllArgs(args)
	if err != nil {
		return nil, err
	}

	traqUsers, err := repo.traQ.GetAll(eargs)
	if err != nil {
		return nil, convertError(err)
	}

	traqUserIDs := make([]uuid.UUID, len(traqUsers))
	for i, v := range traqUsers {
		traqUserIDs[i] = v.ID
	}

	users := make([]*model.User, 0)
	if err := repo.h.
		Where("`users`.`id` IN (?)", traqUserIDs).
		Find(&users).
		Error(); err != nil {
		return nil, convertError(err)
	}

	if l := len(users); l == 0 {
		return []*domain.User{}, nil
	} else if l == 1 {
		portalUser, err := repo.portal.GetByTraqID(users[0].Name)
		if err != nil {
			return nil, convertError(err)
		}

		return []*domain.User{
			{
				ID:       users[0].ID,
				Name:     users[0].Name,
				RealName: portalUser.RealName,
			},
		}, nil
	} else {
		idMap := make(map[string]uuid.UUID, l)
		for _, v := range users {
			idMap[v.Name] = v.ID
		}

		portalUsers, err := repo.portal.GetAll()
		if err != nil {
			return nil, convertError(err)
		}

		result := make([]*domain.User, 0, l)
		for _, v := range portalUsers {
			if id, ok := idMap[v.TraQID]; ok {
				result = append(result, &domain.User{
					ID:       id,
					Name:     v.TraQID,
					RealName: v.RealName,
				})
			}
		}

		return result, nil
	}
}

func (repo *UserRepository) GetUser(userID uuid.UUID) (*domain.UserDetail, error) {
	user := new(model.User)
	err := repo.h.
		Preload("Accounts").
		Where(&model.User{ID: userID}).
		First(user).
		Error()
	if err != nil {
		return nil, convertError(err)
	}

	accounts := make([]*domain.Account, 0, len(user.Accounts))
	for _, v := range user.Accounts {
		accounts = append(accounts, &domain.Account{
			ID:          v.ID,
			DisplayName: v.Name,
			Type:        v.Type,
			PrPermitted: v.Check,
			URL:         v.URL,
		})
	}

	portalUser, err := repo.portal.GetByTraqID(user.Name)
	if err != nil {
		return nil, convertError(err)
	}

	traQUser, err := repo.traQ.GetByUserID(userID)
	if err != nil {
		return nil, convertError(err)
	}

	result := domain.UserDetail{
		User: domain.User{
			ID:       user.ID,
			Name:     user.Name,
			RealName: portalUser.RealName,
		},
		State:    traQUser.State,
		Bio:      user.Description,
		Accounts: accounts,
	}

	return &result, nil
}

func (repo *UserRepository) CreateUser(args *repository.CreateUserArgs) (*domain.UserDetail, error) {
	portalUser, err := repo.portal.GetByTraqID(args.Name)
	if err != nil {
		return nil, err
	}

	user := model.User{
		ID:          uuid.Must(uuid.NewV4()),
		Description: args.Description,
		Check:       args.Check,
		Name:        args.Name,
	}

	err = repo.h.Create(&user).Error()
	if err != nil {
		return nil, convertError(err)
	}

	result := &domain.UserDetail{
		User: domain.User{
			ID:       user.ID,
			Name:     user.Name,
			RealName: portalUser.RealName,
		},
		State:    0,
		Bio:      user.Description,
		Accounts: []*domain.Account{},
	}
	return result, nil
}

func (repo *UserRepository) UpdateUser(userID uuid.UUID, args *repository.UpdateUserArgs) error {
	changes := map[string]interface{}{}
	if args.Description.Valid {
		changes["description"] = args.Description.String
	}
	if args.Check.Valid {
		changes["check"] = args.Check.Bool
	}

	if len(changes) == 0 {
		return nil
	}

	err := repo.h.Transaction(func(tx database.SQLHandler) error {
		user := new(model.User)
		err := repo.h.
			Where(&model.User{ID: userID}).
			First(user).
			Error()
		if err != nil {
			return convertError(err)
		}

		err = repo.h.Model(user).Updates(changes).Error()
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

func (repo *UserRepository) GetAccounts(userID uuid.UUID) ([]*domain.Account, error) {
	accounts := make([]*model.Account, 0)
	err := repo.h.
		Where(&model.Account{UserID: userID}).
		Find(&accounts).
		Error()
	if err != nil {
		return nil, convertError(err)
	}

	result := make([]*domain.Account, 0, len(accounts))
	for _, v := range accounts {
		result = append(result, &domain.Account{
			ID:          v.ID,
			Type:        v.Type,
			PrPermitted: v.Check,
			DisplayName: v.Name,
			URL:         v.URL,
		})
	}
	return result, nil
}

func (repo *UserRepository) GetAccount(userID uuid.UUID, accountID uuid.UUID) (*domain.Account, error) {
	account := &model.Account{}
	err := repo.h.
		Where(&model.Account{ID: accountID, UserID: userID}).
		First(account).
		Error()
	if err != nil {
		return nil, convertError(err)
	}

	result := &domain.Account{
		ID:          account.ID,
		Type:        account.Type,
		PrPermitted: account.Check,
		DisplayName: account.Name,
		URL:         account.URL,
	}

	return result, nil
}

func (repo *UserRepository) CreateAccount(userID uuid.UUID, args *repository.CreateAccountArgs) (*domain.Account, error) {
	account := model.Account{
		ID:     uuid.Must(uuid.NewV4()),
		Type:   args.Type,
		Name:   args.DisplayName,
		URL:    args.URL,
		UserID: userID,
		Check:  args.PrPermitted,
	}
	err := repo.h.Create(&account).Error()
	if err != nil {
		return nil, convertError(err)
	}

	ver := new(model.Account)
	if err := repo.h.
		Where(&model.Account{ID: account.ID}).
		First(ver).
		Error(); err != nil {
		return nil, convertError(err)
	}

	return &domain.Account{
		ID:          ver.ID,
		DisplayName: ver.Name,
		Type:        ver.Type,
		PrPermitted: ver.Check,
		URL:         ver.URL,
	}, nil
}

func (repo *UserRepository) UpdateAccount(userID uuid.UUID, accountID uuid.UUID, args *repository.UpdateAccountArgs) error {
	changes := map[string]interface{}{}
	if args.DisplayName.Valid {
		changes["name"] = args.DisplayName.String
	}
	if args.URL.Valid {
		changes["url"] = args.URL.String
	}
	if args.PrPermitted.Valid {
		changes["check"] = args.PrPermitted.Bool
	}
	if args.Type.Valid {
		changes["type"] = args.Type.Int64
	}

	if len(changes) == 0 {
		return nil
	}

	err := repo.h.Transaction(func(tx database.SQLHandler) error {
		account := new(model.Account)
		err := repo.h.
			Where(&model.Account{ID: accountID, UserID: userID}).
			First(account).
			Error()
		if err != nil {
			return convertError(err)
		}

		err = repo.h.Model(account).Updates(changes).Error()
		if err != nil {
			return convertError(err)
		}
		return nil
	})
	return convertError(err)
}

func (repo *UserRepository) DeleteAccount(userID uuid.UUID, accountID uuid.UUID) error {
	if err := repo.h.Transaction(func(tx database.SQLHandler) error {
		if err := tx.
			Where(&model.Account{ID: accountID, UserID: userID}).
			First(&model.Account{}).
			Error(); err != nil {
			return convertError(err)
		}

		if err := tx.
			Where(&model.Account{ID: accountID, UserID: userID}).
			Delete(&model.Account{}).
			Error(); err != nil {
			return convertError(err)
		}

		return nil
	}); err != nil {
		return convertError(err)
	}

	return nil
}

func (repo *UserRepository) GetProjects(userID uuid.UUID) ([]*domain.UserProject, error) {
	projects := make([]*model.ProjectMember, 0)
	err := repo.h.
		Preload("Project").
		Where(&model.ProjectMember{UserID: userID}).
		Find(&projects).
		Error()
	if err != nil {
		return nil, convertError(err)
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

func (repo *UserRepository) GetGroupsByUserID(userID uuid.UUID) ([]*domain.GroupUser, error) {
	groups := make([]*model.GroupUserBelonging, 0)
	err := repo.h.
		Preload("Group").
		Where(&model.GroupUserBelonging{UserID: userID}).
		Find(&groups).
		Error()
	if err != nil {
		return nil, err
	}

	result := make([]*domain.GroupUser, 0, len(groups))
	for _, v := range groups {
		gr := v.Group
		result = append(result, &domain.GroupUser{
			ID:   gr.GroupID,
			Name: gr.Name,
			Duration: domain.YearWithSemesterDuration{
				Since: domain.YearWithSemester{
					Year:     v.SinceYear,
					Semester: v.SinceSemester,
				},
				Until: domain.YearWithSemester{
					Year:     v.UntilYear,
					Semester: v.UntilSemester,
				},
			},
		})
	}
	return result, nil
}

func (repo *UserRepository) GetContests(userID uuid.UUID) ([]*domain.UserContest, error) {
	contests := make([]*model.ContestTeamUserBelonging, 0)
	err := repo.h.
		Preload("ContestTeam.Contest").
		Where(&model.ContestTeamUserBelonging{UserID: userID}).
		Find(&contests).
		Error()
	if err != nil {
		return nil, convertError(err)
	}

	res := make([]*domain.UserContest, 0, len(contests))
	for _, v := range contests {
		ct := v.ContestTeam
		res = append(res, &domain.UserContest{
			ID:          ct.ID,
			Name:        ct.Name,
			Result:      ct.Result,
			ContestName: ct.Contest.Name,
		})
	}

	return res, nil
}

// Interface guards
var (
	_ repository.UserRepository = (*UserRepository)(nil)
)
