//go:generate go run github.com/golang/mock/mockgen@latest -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/util/optional"
)

type GetUsersArgs struct {
	IncludeSuspended optional.Bool
	Name             optional.String
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
	Type        uint
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
	GetUsers(args *GetUsersArgs) ([]*domain.User, error)
	GetUser(userID uuid.UUID) (*domain.UserDetail, error)
	CreateUser(args *CreateUserArgs) (*domain.UserDetail, error)
	UpdateUser(userID uuid.UUID, args *UpdateUserArgs) error
	GetAccounts(userID uuid.UUID) ([]*domain.Account, error)
	GetAccount(userID uuid.UUID, accountID uuid.UUID) (*domain.Account, error)
	CreateAccount(userID uuid.UUID, args *CreateAccountArgs) (*domain.Account, error)
	UpdateAccount(userID uuid.UUID, accountID uuid.UUID, args *UpdateAccountArgs) error
	DeleteAccount(userID uuid.UUID, accountID uuid.UUID) error
	GetProjects(userID uuid.UUID) ([]*domain.UserProject, error)
	GetContests(userID uuid.UUID) ([]*domain.UserContest, error)
	GetGroupsByUserID(userID uuid.UUID) ([]*domain.UserGroup, error)
}
