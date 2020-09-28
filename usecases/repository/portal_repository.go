package repository

import (
	"context"

	"github.com/traPtitech/traPortfolio/domain"
)

type PortalRepository interface {
	GetUser(context.Context, string) (*domain.PortalUser, error)
}
