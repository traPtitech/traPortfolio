package repository

import (
	"context"

	"github.com/traPtitech/traPortfolio/domain"

	"github.com/gofrs/uuid"
)

type TraQRepository interface {
	GetUser(context.Context, uuid.UUID) (*domain.TraQUser, error)
}

/*
https://q.trap.jp/api/v3/users
名前で探してる
https://md.trap.jp/yfLG73ctSgG-wudmKnGAAw?view#

*/
