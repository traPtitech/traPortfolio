//go:generate go run go.uber.org/mock/mockgen@latest -typed -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/internal/domain"
	"github.com/traPtitech/traPortfolio/internal/pkgs/optional"
)

type GetGroupsArgs struct {
	Limit optional.Of[int]
}

type GroupRepository interface {
	GetGroups(ctx context.Context, args *GetGroupsArgs) ([]*domain.Group, error)
	GetGroup(ctx context.Context, groupID uuid.UUID) (*domain.GroupDetail, error)
}
