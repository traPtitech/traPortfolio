//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
)

type KnoqRepository interface {
	GetAll() ([]*domain.KnoQEvent, error)
	GetByID(id uuid.UUID) (*domain.KnoQEvent, error)
}
