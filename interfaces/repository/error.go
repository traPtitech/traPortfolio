package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

func convertError(err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return repository.ErrNotFound
	default:
		return err
	}
}
