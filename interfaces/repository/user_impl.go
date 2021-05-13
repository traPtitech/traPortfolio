package repository

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type UserRepository struct {
	database.SQLHandler
	portal external.PortalAPI
	traQ   external.TraQAPI
}

func NewUserRepository(sql database.SQLHandler, portalAPI external.PortalAPI, traQAPI external.TraQAPI) repository.UserRepository {
	return &UserRepository{
		SQLHandler: sql,
		portal:     portalAPI,
		traQ:       traQAPI,
	}
}

func (repo *UserRepository) GetUsers() ([]*domain.User, error) {
	users := make([]*model.User, 0)
	err := repo.Find(&users).Error()
	if err != nil {
		return nil, err
	}
	idMap := make(map[string]uuid.UUID, len(users))
	for _, v := range users {
		idMap[v.Name] = v.ID
	}

	portalUsers, err := repo.portal.GetAll()
	if err != nil {
		return nil, err
	}

	result := make([]*domain.User, 0, len(users))
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

func (repo *UserRepository) GetUser(id uuid.UUID) (*domain.UserDetail, error) {
	user := model.User{ID: id}
	err := repo.First(&user).Error()
	if err != nil {
		return nil, err
	}

	portalUser, err := repo.portal.GetByID(user.Name)
	if err != nil {
		return nil, err
	}

	traQUser, err := repo.traQ.GetByID(id)
	if err != nil {
		return nil, err
	}

	accounts, err := repo.GetAccounts(id)
	if err != nil {
		return nil, err
	}

	result := domain.UserDetail{
		ID:       user.ID,
		Name:     user.Name,
		RealName: portalUser.RealName,
		State:    traQUser.State,
		Bio:      user.Description,
		Accounts: accounts,
	}

	return &result, nil
}

func (repo *UserRepository) GetAccounts(userID uuid.UUID) ([]*domain.Account, error) {
	accounts := make([]*model.Account, 0)
	err := repo.Find(&accounts, "user_id = ?", userID).Error()
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Account, 0, len(accounts))
	for _, v := range accounts {
		result = append(result, &domain.Account{
			ID:          v.ID,
			Type:        v.Type,
			PrPermitted: v.Check,
		})
	}
	return result, nil
}

func (repo *UserRepository) GetAccount(userID uuid.UUID, accountID uuid.UUID) (*domain.Account, error) {
	account := &model.Account{ID: accountID}
	err := repo.First(account).Error()
	if err != nil {
		return nil, err
	}
	if account.UserID != userID {
		return nil, repository.ErrNotFound
	}

	result := &domain.Account{
		ID:          account.ID,
		Type:        account.Type,
		PrPermitted: account.Check,
	}

	return result, nil
}

func (repo *UserRepository) Update(id uuid.UUID, changes map[string]interface{}) error {
	err := repo.Transaction(func(tx database.SQLHandler) error {
		user := &model.User{ID: id}
		err := repo.First(user).Error()
		if err == gorm.ErrRecordNotFound {
			return repository.ErrNotFound
		} else if err != nil {
			return err
		}

		err = repo.Model(user).Updates(changes).Error()
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (repo *UserRepository) CreateAccount(id uuid.UUID, args *repository.CreateAccountArgs) (*domain.Account, error) {
	account := model.Account{
		ID:     uuid.Must(uuid.NewV4()),
		Type:   args.Type,
		Name:   args.ID,
		URL:    args.URL,
		UserID: id,
		Check:  args.PrPermitted,
	}
	err := repo.Create(account).Error()
	if err != nil {
		return nil, err
	}

	var result *domain.Account

	err = repo.First(result, domain.Account{ID: account.ID}).Error()

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *UserRepository) UpdateAccount(userID uuid.UUID, accountID uuid.UUID, changes map[string]interface{}) error {
	err := repo.Transaction(func(tx database.SQLHandler) error {
		account := &model.Account{ID: accountID}
		err := repo.First(account).Error()
		if err == gorm.ErrRecordNotFound || account.UserID != userID {
			return repository.ErrNotFound
		} else if err != nil {
			return err
		}

		err = repo.Model(account).Updates(changes).Error()
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (repo *UserRepository) DeleteAccount(accountID uuid.UUID, userID uuid.UUID) error {

	err := repo.Delete(&domain.Account{}, &model.Account{ID: accountID, UserID: userID}).Error()

	return err

}

func (repo *UserRepository) GetProjects(userID uuid.UUID) ([]*domain.UserProject, error) {
	projects := make([]*domain.UserProject, 0)
	selectQuery := "projects.id, projects.name, projects.since, projects.until, project_members.since AS user_since, project_members.until AS user_until"
	whereQuery := "project_members.user_id = ?"
	joinQuery := "LEFT JOIN projects ON project_members.project_id = projects.id"
	err := repo.Model(&model.ProjectMember{}).Select(selectQuery).Where(whereQuery, userID).Joins(joinQuery).Scan(&projects).Error()
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (repo *UserRepository) GetContests(userID uuid.UUID) ([]*domain.UserContest, error) {
	contests := make([]*domain.UserContest, 0)
	selectQuery := "contest_teams.id, contest_teams.name, contest_teams.result, contests.name AS contest_name"
	whereQuery := "contest_team_user_belongings.user_id = ?"
	joinQuery := "LEFT JOIN contests ON contest_teams.contest_id = contests.id"
	joinQuery2 := "LEFT JOIN contest_team_user_belongings ON contest_teams.id = contest_team_user_belongings.team_id"
	err := repo.Model(&model.ContestTeam{}).Select(selectQuery).Where(whereQuery, userID).Joins(joinQuery).Joins(joinQuery2).Scan(&contests).Error()
	if err != nil {
		return nil, err
	}
	return contests, nil
}

// Interface guards
var (
	_ repository.UserRepository = (*UserRepository)(nil)
)
