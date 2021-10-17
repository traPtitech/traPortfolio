package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type GroupRepository struct {
	api external.GroupAPI
}

func NewGroupRepository(api external.GroupAPI) repository.GroupRepository {
	return &GroupRepository{api}
}

func (repo *GroupRepository) GetGroupsByID(userID uuid.UUID) ([]*domain.GroupUser, error) {
	eres, err := repo.api.GetGroupsByID(userID)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.GroupUser, 0, len(eres))
	for _, v := range eres {
		result = append(result, &domain.GroupUser{
			ID:   v.ID,
			Name: v.Name,
			Duration: domain.ProjectDuration{
				Since: domain.YearWithSemester{
					Year:     v.Duration.Since.Year,
					Semester: v.Duration.Since.Semester,
				},
				Until: domain.YearWithSemester{
					Year:     v.Duration.Since.Year,
					Semester: v.Duration.Since.Semester,
				},
			},
		})
	}
	return result, nil
}

func (repo *GroupRepository) GetAllGroups() ([]*domain.Groups, error) {
	eres, err := repo.api.GetAllGroups()
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Groups, 0, len(eres))
	for _, v := range eres {
		result = append(result, &domain.Groups{
			ID:   v.ID,
			Name: v.Name,
		})
	}
	return result, nil
}

func (repo *GroupRepository) GetGroup(groupID uuid.UUID) (*domain.GroupDetail, error) {
	eres, err := repo.api.GetGroup(groupID)
	if err != nil {
		return nil, err
	}

	erMembers := make([]*domain.UserGroup, 0, len(eres.Members))
	for _, v := range eres.Members {
		erMembers = append(erMembers, &domain.UserGroup{
			ID:       v.ID,
			Name:     v.Name,
			RealName: v.RealName,
			Duration: domain.ProjectDuration{
				Since: domain.YearWithSemester{
					Year:     v.Duration.Since.Year,
					Semester: v.Duration.Since.Semester,
				},
				Until: domain.YearWithSemester{
					Year:     v.Duration.Since.Year,
					Semester: v.Duration.Since.Semester,
				},
			},
		})
	}

	result := &domain.GroupDetail{
		ID:   groupID,
		Name: eres.Name,
		Link: eres.Link,
		Leader: &domain.User{
			ID:       eres.Leader.ID,
			Name:     eres.Leader.Name,
			RealName: eres.Leader.RealName,
		},
		Members:     erMembers,
		Description: eres.Description,
	}
	return result, nil
}
