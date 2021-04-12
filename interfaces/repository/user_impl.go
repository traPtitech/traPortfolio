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
		RealName: portalUser.RealName,
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

// Interface guards
var (
	_ repository.UserRepository = (*UserRepository)(nil)
)
