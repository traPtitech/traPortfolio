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
	portal external.PortalQAPI
	traQ   external.TraQAPI
}

func NewUserRepository(sql database.SQLHandler) *UserRepository {
	return &UserRepository{SQLHandler: sql}
}

func (repo *UserRepository) GetUsers() ([]*domain.User, error) {
	users, err := repo.GetUsers()
	if err != nil {
		return nil, err
	}
	portalUsers, err := repo.portal.GetAll()
	if err != nil {
		return nil, err
	}
	idMap := make(map[string]uuid.UUID, len(users))
	for _, v := range users {
		idMap[v.Name] = v.ID
	}
	result := make([]*domain.User, 0, len(users))
	for _, v := range portalUsers {
		if id, ok := idMap[v.ID]; ok {
			result = append(result, &domain.User{
				ID:       id,
				Name:     v.ID,
				RealName: v.Name,
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

	portalUser, err := repo.portal.GetByID(id)
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
		RealName: portalUser.Name,
		State:    traQUser.State,
		Bio:      user.Description,
		Accounts: accounts,
	}

	return &result, nil
}

func (repo *UserRepository) GetAccounts(id uuid.UUID) ([]*domain.Account, error) {
	accounts := make([]*model.Account, 0)
	err := repo.Find(&accounts, "user_id = ?", id).Error()
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

func (repo *UserRepository) Update(u *domain.EditUser) error {
	user := &model.User{ID: u.ID}
	err := repo.First(user).Error()
	if err == gorm.ErrRecordNotFound {
		return repository.ErrNotFound
	} else if err != nil {
		return err
	}

	user.Description = u.Description
	user.Check = u.Check
	err = repo.Save(user).Error()
	return err
}

// Interface guards
var (
	_ repository.UserRepository = (*UserRepository)(nil)
)
