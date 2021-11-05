//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"context"

	"github.com/traPtitech/traPortfolio/domain"
)

type PortalRepository interface {
	GetUsers(ctx context.Context) ([]*domain.PortalUser, error)
	GetUser(ctx context.Context, traQID string) (*domain.PortalUser, error)
}
