package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
)

type PortalRepository interface {
	GetUser(context.Context, uuid.UUID) (*domain.PortalUser, error)
}
