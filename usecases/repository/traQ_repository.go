//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"context"

	"github.com/traPtitech/traPortfolio/domain"

	"github.com/gofrs/uuid"
)

type TraQRepository interface {
	GetUser(ctx context.Context, id uuid.UUID) (*domain.TraQUser, error)
}

/*
https://q.trap.jp/api/v3/users
名前で探してる
https://md.trap.jp/yfLG73ctSgG-wudmKnGAAw?view#

*/
