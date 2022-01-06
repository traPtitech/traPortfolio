//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package service

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type GroupService interface {
	GetAllGroups(ctx context.Context) ([]*domain.Group, error)
	GetGroup(ctx context.Context, groupID uuid.UUID) (*domain.GroupDetail, error)
}

type groupService struct {
	repo repository.GroupRepository
}

func NewGroupService(repo repository.GroupRepository) GroupService {
	return &groupService{
		repo,
	}
}

func (s *groupService) GetAllGroups(ctx context.Context) ([]*domain.Group, error) {
	return s.repo.GetAllGroups()
}

func (s *groupService) GetGroup(ctx context.Context, groupID uuid.UUID) (*domain.GroupDetail, error) {
	return s.repo.GetGroup(groupID)
}

// Interface guards
var (
	_ GroupService = (*groupService)(nil)
)
