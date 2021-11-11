//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE

package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
)

type KnoqRepository interface {
	GetAll() ([]*domain.KnoQEvent, error)
	GetByID(id uuid.UUID) (*domain.KnoQEvent, error)
}
