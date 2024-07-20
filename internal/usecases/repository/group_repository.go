//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/internal/domain"
)

type GroupRepository interface {
	GetGroups(ctx context.Context) ([]*domain.Group, error)
	GetGroup(ctx context.Context, groupID uuid.UUID) (*domain.GroupDetail, error)
}
