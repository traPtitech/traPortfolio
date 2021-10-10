//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
)

type GroupRepository interface {
	GetGroupsByID(userID uuid.UUID) ([]*domain.Groups, error)
	GetAllGroups() ([]*domain.Groups, error)
	GetGroup(groupID uuid.UUID) (*domain.Groups, error)
}
