package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
)

type TraQRepository interface {
	GetUser(context.Context, uuid.UUID) (*domain.TraQUser, error)
}
