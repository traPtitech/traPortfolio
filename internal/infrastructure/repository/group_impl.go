package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/internal/domain"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/repository/model"
	"github.com/traPtitech/traPortfolio/internal/usecases/repository"
	"gorm.io/gorm"
)

type GroupRepository struct {
	h *gorm.DB
}

func NewGroupRepository(sql *gorm.DB) repository.GroupRepository {
	return &GroupRepository{h: sql}
}

func (r *GroupRepository) GetGroups(ctx context.Context) ([]*domain.Group, error) {
	groups := make([]*model.Group, 0)
	err := r.h.WithContext(ctx).Find(&groups).Error
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

func (r *GroupRepository) GetGroup(ctx context.Context, groupID uuid.UUID) (*domain.GroupDetail, error) {
	group := &model.Group{}
	if err := r.h.
		WithContext(ctx).
		Where(&model.Group{GroupID: groupID}).
		First(group).
		Error; err != nil {
		return nil, err
	}

	users := make([]*model.GroupUserBelonging, 0)
	if err := r.h.
		WithContext(ctx).
		Where(&model.GroupUserBelonging{GroupID: groupID}).
		Find(&users).
		Error; err != nil {
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
			Duration: domain.NewYearWithSemesterDuration(
				v.SinceYear,
				v.SinceSemester,
				v.UntilYear,
				v.UntilSemester,
			),
		})
	}

	admins := make([]*model.GroupUserAdmin, 0)
	if err := r.h.
		WithContext(ctx).
		Where(&model.GroupUserAdmin{GroupID: groupID}).
		Find(&admins).
		Error; err != nil {
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
