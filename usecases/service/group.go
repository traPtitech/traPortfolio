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
	group repository.GroupRepository
	user  repository.UserRepository
}

func NewGroupService(group repository.GroupRepository, user repository.UserRepository) GroupService {
	return &groupService{
		group, user,
	}
}

func (s *groupService) GetAllGroups(ctx context.Context) ([]*domain.Group, error) {
	return s.group.GetAllGroups()
}

func (s *groupService) GetGroup(ctx context.Context, groupID uuid.UUID) (*domain.GroupDetail, error) {
	groups, err := s.group.GetGroup(groupID)
	if err != nil {
		return nil, err
	}

	// pick all users info
	users, err := s.user.GetUsers(&repository.GetUsersArgs{}) // TODO: IncludeSuspendedをtrueにするか考える
	if err != nil {
		return nil, err
	}

	umap := make(map[uuid.UUID]*domain.User)
	for _, u := range users {
		umap[u.ID] = u
	}

	// fill members info
	for i, v := range groups.Members {
		if u, ok := umap[v.ID]; ok {
			groups.Members[i].Name = u.Name
			groups.Members[i].RealName = u.RealName
		}
	}

	// fill leader info
	for i, v := range groups.Admin {
		if u, ok := umap[v.ID]; ok {
			groups.Admin[i].Name = u.Name
			groups.Admin[i].RealName = u.RealName
		}
	}

	return groups, nil
}

// Interface guards
var (
	_ GroupService = (*groupService)(nil)
)
