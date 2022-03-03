//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

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
	ID          string // traqID
	Type        uint
	URL         string
	PrPermitted bool
}

type UpdateAccountArgs struct {
	Name        optional.String // Twitter等のアカウントID
	Type        optional.Int64
	URL         optional.String
	PrPermitted optional.Bool
}

type UserRepository interface {
	GetUsers(args *GetUsersArgs) ([]*domain.User, error)
	GetUser(id uuid.UUID) (*domain.UserDetail, error)
	CreateUser(args CreateUserArgs) (*domain.UserDetail, error)
	UpdateUser(id uuid.UUID, args *UpdateUserArgs) error
	GetAccounts(id uuid.UUID) ([]*domain.Account, error)
	GetAccount(id uuid.UUID, accountID uuid.UUID) (*domain.Account, error)
	CreateAccount(id uuid.UUID, args *CreateAccountArgs) (*domain.Account, error)
	UpdateAccount(id uuid.UUID, accountID uuid.UUID, args *UpdateAccountArgs) error
	DeleteAccount(id uuid.UUID, accountID uuid.UUID) error
	GetProjects(id uuid.UUID) ([]*domain.UserProject, error)
	GetContests(id uuid.UUID) ([]*domain.UserContest, error)
	GetGroupsByUserID(id uuid.UUID) ([]*domain.GroupUser, error)
}
