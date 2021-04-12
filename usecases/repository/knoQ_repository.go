package repository

import (
	"github.com/gofrs/uuid"
	"github.com/traPtitech/traPortfolio/domain"
)

type KnoqRepository interface {
	GetAll() ([]*domain.KnoQEvent, error)
	GetByID(id uuid.UUID) (*domain.KnoQEvent, error)
}
