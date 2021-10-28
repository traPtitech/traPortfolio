//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
)

type GroupRepository interface {
	GetGroupsByUserID(userID uuid.UUID) ([]*domain.GroupUser, error)
	GetAllGroups() ([]*domain.Group, error)
	GetGroup(groupID uuid.UUID) (*domain.GroupDetail, error)
}
