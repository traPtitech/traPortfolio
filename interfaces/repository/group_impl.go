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
	groups := make([]*model.GroupUserBelonging, 0)
	if err := repo.h.Preload("Group").Find(&groups).Error(); err != nil {
		return nil, convertError(err)
	}

	result := make([]*domain.Group, 0, len(groups))
	for _, v := range groups {
		group := make([]*model.Group, 0)
		if err := repo.h.Where(model.Group{GroupID: v.GroupID}).Find(&group).Error(); err != nil {
			return nil, convertError(err)
		}

		result = append(result, &domain.Group{
			ID:   v.UserID,
			Name: group[0].Name,
		})
	}
	return result, nil
}

func (repo *GroupRepository) GetGroup(groupID uuid.UUID) (*domain.GroupDetail, error) {
	users := make([]*model.GroupUserBelonging, 0)
	if err := repo.h.Preload("Group").Where(model.GroupUserBelonging{GroupID: groupID}).Find(&users).Error(); err != nil {
		return nil, convertError(err)
	}

	erMembers := make([]*domain.UserGroup, 0, len(users))
	for _, v := range users {
		group := make([]*model.Group, 0)
		if err := repo.h.Where(model.Group{GroupID: v.GroupID}).Find(&group).Error(); err != nil {
			return nil, convertError(err)
		}

		// Name,RealNameはPortalから取得する
		erMembers = append(erMembers, &domain.UserGroup{
			ID: v.UserID,
			// Name:     v.Name,
			// RealName: v.RealName,
			Duration: domain.GroupDuration{
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
	group := make([]*model.Group, 0)
	if err := repo.h.Where(model.Group{GroupID: groupID}).Find(&group).Error(); err != nil {
		return nil, convertError(err)
	}

	// Name,RealNameはPortalから取得する
	result := &domain.GroupDetail{
		ID:   groupID,
		Name: group[0].Name,
		Link: group[0].Link,
		Leader: &domain.User{
			ID: group[0].Leader,
			// Name:     eres.Leader.Name,
			// RealName: eres.Leader.RealName,
		},
		Members: erMembers,
		// GroupのテーブルにDescription入れるの忘れたのでとりあえずnullで返す
		// Description: group[0].Description,
	}
	return result, nil
}
