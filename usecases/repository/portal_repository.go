//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"context"

	"github.com/traPtitech/traPortfolio/domain"
)

type PortalRepository interface {
	GetUser(context.Context, string) (*domain.PortalUser, error)
	GetUsers(context.Context) ([]*domain.PortalUser, error)
}
