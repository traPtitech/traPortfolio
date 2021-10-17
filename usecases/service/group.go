package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type GroupService struct {
	repo repository.GroupRepository
}

func NewGroupService(repo repository.GroupRepository) GroupService {
	return GroupService{
		repo,
	}
}

func (s *GroupService) GetGroupsByID(ctx context.Context, userID uuid.UUID) ([]*domain.GroupUser, error) {
	return s.repo.GetGroupsByID(userID)
}

func (s *GroupService) GetAllGroups(ctx context.Context) ([]*domain.Groups, error) {
	return s.repo.GetAllGroups()
}

func (s *GroupService) GetGroup(ctx context.Context, groupID uuid.UUID) (*domain.GroupDetail, error) {
	return s.repo.GetGroup(groupID)
}
