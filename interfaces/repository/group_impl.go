package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type GroupRepository struct {
	h database.SQLHandler
}

func NewGroupRepository(sql database.SQLHandler) repository.GroupRepository {
	return &GroupRepository{h: sql}
}

func (repo *GroupRepository) GetAllGroups() ([]*domain.Group, error) {
	groups := make([]*model.Group, 0)
	err := repo.h.Find(&groups).Error()
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Group, 0, len(groups))
	for _, v := range groups {
		result = append(result, &domain.Group{
			ID:   v.GroupID,
			Name: v.Name,
		})
	}
	return result, nil
}

func (repo *GroupRepository) GetGroup(groupID uuid.UUID) (*domain.GroupDetail, error) {
	group := &model.Group{}
	if err := repo.h.
		Where(&model.Group{GroupID: groupID}).
		First(group).
		Error(); err != nil {
		return nil, err
	}

	users := make([]*model.GroupUserBelonging, 0)
	if err := repo.h.
		Where(&model.GroupUserBelonging{GroupID: groupID}).
		Find(&users).
		Error(); err != nil {
		return nil, err
	}

	// Name, RealNameはusecasesでPortalから取得する
	erMembers := make([]*domain.UserWithDuration, 0, len(users))
	for _, v := range users {
		erMembers = append(erMembers, &domain.UserWithDuration{
			User: domain.User{
				ID: v.UserID,
				// Name:     v.Name,
				// RealName: v.RealName,
			},
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

	admins := make([]*model.GroupUserAdmin, 0)
	if err := repo.h.
		Where(&model.GroupUserAdmin{GroupID: groupID}).
		Find(&admins).
		Error(); err != nil {
		return nil, err
	}

	erAdmin := make([]*domain.User, 0, len(admins))
	for _, v := range admins {
		erAdmin = append(erAdmin, &domain.User{ID: v.UserID})
	}

	// Name,RealNameはPortalから取得する
	result := &domain.GroupDetail{
		ID:          groupID,
		Name:        group.Name,
		Link:        group.Link,
		Admin:       erAdmin,
		Members:     erMembers,
		Description: group.Description,
	}
	return result, nil
}
