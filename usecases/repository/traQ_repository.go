package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
)

type TraQRepository interface {
	GetUser(context.Context, uuid.UUID) (*model.TraQUser, error)
}
