//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/util/optional"
)

type GetUsersArgs struct {
	IncludeSuspended optional.Of[bool]
	Name             optional.Of[string]
	Limit            optional.Of[int]
}

type UpdateUserArgs struct {
	Description optional.Of[string]
	Check       optional.Of[bool]
}

type CreateAccountArgs struct {
	DisplayName string // 外部アカウントの表示名
	Type        domain.AccountType
	URL         string
	PrPermitted bool
}

type UpdateAccountArgs struct {
	DisplayName optional.Of[string] // 外部アカウントの表示名
	Type        optional.Of[domain.AccountType]
	URL         optional.Of[string]
	PrPermitted optional.Of[bool]
}

type UserRepository interface {
	GetUsers(ctx context.Context, args *GetUsersArgs) ([]*domain.User, error)
	SyncUsers(ctx context.Context) error
	GetUser(ctx context.Context, userID uuid.UUID) (*domain.UserDetail, error)
	UpdateUser(ctx context.Context, userID uuid.UUID, args *UpdateUserArgs) error
	GetAccounts(ctx context.Context, userID uuid.UUID) ([]*domain.Account, error)
	GetAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) (*domain.Account, error)
	CreateAccount(ctx context.Context, userID uuid.UUID, args *CreateAccountArgs) (*domain.Account, error)
	UpdateAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID, args *UpdateAccountArgs) error
	DeleteAccount(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) error
	GetProjects(ctx context.Context, userID uuid.UUID) ([]*domain.UserProject, error)
	GetContests(ctx context.Context, userID uuid.UUID) ([]*domain.UserContest, error)
	GetGroupsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.UserGroup, error)
}
