package repository

import (
	"context"

	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
)

type PortalRepository interface {
	GetUser(context.Context, string) (*model.PortalUser, error)
	GetUsers(context.Context) ([]*model.PortalUser, error)
}
