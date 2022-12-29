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

func (r *UserRepository) GetUsers(args *repository.GetUsersArgs) ([]*domain.User, error) {
	eargs, err := makeTraqGetAllArgs(args)
	if err != nil {
		return nil, err
	}

	traqUsers, err := r.traQ.GetAll(eargs)
	if err != nil {
		return nil, convertError(err)
	}

	traqUserIDs := make([]uuid.UUID, len(traqUsers))
	for i, v := range traqUsers {
		traqUserIDs[i] = v.ID
	}

	users := make([]*model.User, 0)
	if err := r.h.
		Where("`users`.`id` IN (?)", traqUserIDs).
		Find(&users).
		Error(); err != nil {
		return nil, convertError(err)
	}

	if l := len(users); l == 0 {
		return []*domain.User{}, nil
	} else if l == 1 {
		portalUser, err := r.portal.GetByTraqID(users[0].Name)
		if err != nil {
			return nil, convertError(err)
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
		userMap := make(map[string]*model.User, l)
		for _, v := range users {
			userMap[v.Name] = v
		}

		portalUsers, err := r.portal.GetAll()
		if err != nil {
			return nil, convertError(err)
		}

		result := make([]*domain.User, 0, l)
		for _, v := range portalUsers {
			if u, ok := userMap[v.TraQID]; ok {
				result = append(result, domain.NewUser(
					u.ID,
					u.Name,
					v.RealName,
					u.Check,
				))
			}
		}

		return result, nil
	}
}

func (r *UserRepository) GetUser(userID uuid.UUID) (*domain.UserDetail, error) {
	user := new(model.User)
	err := r.h.
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

	portalUser, err := r.portal.GetByTraqID(user.Name)
	if err != nil {
		return nil, convertError(err)
	}

	traQUser, err := r.traQ.GetByUserID(userID)
	if err != nil {
		return nil, convertError(err)
	}

	result := domain.UserDetail{
		User: *domain.NewUser(
			user.ID,
			user.Name,
			portalUser.RealName,
			user.Check,
		),
		State:    traQUser.State,
		Bio:      user.Description,
		Accounts: accounts,
	}

	return &result, nil
}

func (r *UserRepository) CreateUser(args *repository.CreateUserArgs) (*domain.UserDetail, error) {
	portalUser, err := r.portal.GetByTraqID(args.Name)
	if err != nil {
		return nil, err
	}

	user := model.User{
		ID:          uuid.Must(uuid.NewV4()),
		Description: args.Description,
		Check:       args.Check,
		Name:        args.Name,
	}

	err = r.h.Create(&user).Error()
	if err != nil {
		return nil, convertError(err)
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

func (r *UserRepository) UpdateUser(userID uuid.UUID, args *repository.UpdateUserArgs) error {
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

	err := r.h.Transaction(func(tx database.SQLHandler) error {
		user := new(model.User)
		err := tx.
			Where(&model.User{ID: userID}).
			First(user).
			Error()
		if err != nil {
			return convertError(err)
		}

		err = tx.Model(user).Updates(changes).Error()
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

func (r *UserRepository) GetAccounts(userID uuid.UUID) ([]*domain.Account, error) {
	err := r.h.
		Where(&model.User{ID: userID}).
		First(&model.User{}).
		Error()
	if err != nil {
		return nil, convertError(err)
	}

	accounts := make([]*model.Account, 0)
	err = r.h.
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

func (r *UserRepository) GetAccount(userID uuid.UUID, accountID uuid.UUID) (*domain.Account, error) {
	account := &model.Account{}
	err := r.h.
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

func (r *UserRepository) CreateAccount(userID uuid.UUID, args *repository.CreateAccountArgs) (*domain.Account, error) {
	if !domain.IsValidAccountURL(domain.AccountType(args.Type), args.URL) {
		return nil, repository.ErrInvalidArg
	}

	account := model.Account{
		ID:     uuid.Must(uuid.NewV4()),
		Type:   args.Type,
		Name:   args.DisplayName,
		URL:    args.URL,
		UserID: userID,
		Check:  args.PrPermitted,
	}
	err := r.h.Create(&account).Error()
	if err != nil {
		return nil, convertError(err)
	}

	ver := new(model.Account)
	if err := r.h.
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

func (r *UserRepository) UpdateAccount(userID uuid.UUID, accountID uuid.UUID, args *repository.UpdateAccountArgs) error {
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

	err := r.h.Transaction(func(tx database.SQLHandler) error {
		account := new(model.Account)
		err := tx.
			Where(&model.Account{ID: accountID, UserID: userID}).
			First(account).
			Error()
		if err != nil {
			return convertError(err)
		}

		// URLのvalidation
		if args.Type.Valid && args.URL.Valid {
			if !domain.IsValidAccountURL(domain.AccountType(args.Type.Int64), args.URL.String) {
				return repository.ErrInvalidArg
			}
		} else if !args.Type.Valid && args.URL.Valid {
			if !domain.IsValidAccountURL(domain.AccountType(account.Type), args.URL.String) {
				return repository.ErrInvalidArg
			}
		} else if args.Type.Valid && !args.URL.Valid {
			if !domain.IsValidAccountURL(domain.AccountType(args.Type.Int64), account.URL) {
				return repository.ErrInvalidArg
			}
		}

		err = tx.Model(account).Updates(changes).Error()
		if err != nil {
			return convertError(err)
		}
		return nil
	})
	return convertError(err)
}

func (r *UserRepository) DeleteAccount(userID uuid.UUID, accountID uuid.UUID) error {
	if err := r.h.Transaction(func(tx database.SQLHandler) error {
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

func (r *UserRepository) GetProjects(userID uuid.UUID) ([]*domain.UserProject, error) {
	err := r.h.
		Where(&model.User{ID: userID}).
		First(&model.User{}).
		Error()
	if err != nil {
		return nil, convertError(err)
	}

	projects := make([]*model.ProjectMember, 0)
	err = r.h.
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

func (r *UserRepository) GetGroupsByUserID(userID uuid.UUID) ([]*domain.UserGroup, error) {
	err := r.h.
		Where(&model.User{ID: userID}).
		First(&model.User{}).
		Error()
	if err != nil {
		return nil, convertError(err)
	}

	groups := make([]*model.GroupUserBelonging, 0)
	err = r.h.
		Preload("Group").
		Where(&model.GroupUserBelonging{UserID: userID}).
		Find(&groups).
		Error()
	if err != nil {
		return nil, err
	}

	result := make([]*domain.UserGroup, 0, len(groups))
	for _, v := range groups {
		gr := v.Group
		result = append(result, &domain.UserGroup{
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

func (r *UserRepository) GetContests(userID uuid.UUID) ([]*domain.UserContest, error) {
	err := r.h.
		Where(&model.User{ID: userID}).
		First(&model.User{}).
		Error()
	if err != nil {
		return nil, convertError(err)
	}

	contestTeamUserBelongings := make([]*model.ContestTeamUserBelonging, 0)
	err = r.h.
		Preload("ContestTeam.Contest").
		Where(&model.ContestTeamUserBelonging{UserID: userID}).
		Find(&contestTeamUserBelongings).
		Error()
	if err != nil {
		return nil, convertError(err)
	}

	contestsMap := make(map[uuid.UUID]*domain.UserContest)
	for _, v := range contestTeamUserBelongings {
		ct := v.ContestTeam
		if c, ok := contestsMap[ct.ContestID]; ok {
			c.Teams = append(c.Teams, &domain.ContestTeam{
				ID:        ct.ID,
				ContestID: ct.ContestID,
				Name:      ct.Name,
				Result:    ct.Result,
			})
		} else {
			contestsMap[v.ContestTeam.ContestID] = &domain.UserContest{
				ID:        ct.Contest.ID,
				Name:      ct.Contest.Name,
				TimeStart: ct.Contest.Since,
				TimeEnd:   ct.Contest.Until,
				Teams: []*domain.ContestTeam{
					{
						ID:        ct.ID,
						ContestID: ct.ContestID,
						Name:      ct.Name,
						Result:    ct.Result,
					},
				},
			}
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
