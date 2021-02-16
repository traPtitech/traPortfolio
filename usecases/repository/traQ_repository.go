package repository

import (
	"context"

	"github.com/traPtitech/traPortfolio/domain"
)

type TraQRepository interface {
	GetUser(context.Context, string) (*domain.TraQUser, error)
}
