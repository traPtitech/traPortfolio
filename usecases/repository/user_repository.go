//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/util/optional"
)

type GetUsersArgs struct {
	IncludeSuspended optional.Bool
	Name             optional.String
	Limit            optional.Int64
}

type CreateUserArgs struct {
	Description string
	Check       bool
	Name        string
}

type UpdateUserArgs struct {
	Description optional.String
	Check       optional.Bool
}

type CreateAccountArgs struct {
	DisplayName string // 外部アカウントの表示名
	Type        domain.AccountType
	URL         string
	PrPermitted bool
}

type UpdateAccountArgs struct {
	DisplayName optional.String // 外部アカウントの表示名
	Type        optional.Int64
	URL         optional.String
	PrPermitted optional.Bool
}

type UserRepository interface {
	GetUsers(ctx context.Context, args *GetUsersArgs) ([]*domain.User, error)
	GetUser(ctx context.Context, userID uuid.UUID) (*domain.UserDetail, error)
	CreateUser(ctx context.Context, args *CreateUserArgs) (*domain.UserDetail, error)
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
