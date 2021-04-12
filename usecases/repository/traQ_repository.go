package repository

import (
	"context"

	"github.com/traPtitech/traPortfolio/domain"

	"github.com/gofrs/uuid"
)

type TraQRepository interface {
	GetUser(context.Context, uuid.UUID) (*domain.TraQUser, error)
}
