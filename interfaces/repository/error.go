package repository

import (
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

func convertError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.ErrNotFound
	}
	return err
}
