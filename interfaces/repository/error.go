package repository

import (
	"errors"

	"github.com/traPtitech/traPortfolio/usecases/repository"
	"gorm.io/gorm"
)

func convertError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.ErrNotFound
	}
	return err
}
