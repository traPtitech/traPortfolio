//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/util/optional"
)

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
	ID          optional.String // traqID
	Type        optional.Int64
	URL         optional.String
	PrPermitted optional.Bool
}

type UserRepository interface {
	GetUsers() ([]*domain.User, error)
	GetUser(uuid.UUID) (*domain.UserDetail, error)
	Update(id uuid.UUID, changes map[string]interface{}) error
	GetAccounts(userID uuid.UUID) ([]*domain.Account, error)
	GetAccount(userID uuid.UUID, accountID uuid.UUID) (*domain.Account, error)
	CreateAccount(uuid.UUID, *CreateAccountArgs) (*domain.Account, error)
	UpdateAccount(userID uuid.UUID, accountID uuid.UUID, changes map[string]interface{}) error
	DeleteAccount(uuid.UUID, uuid.UUID) error
}
