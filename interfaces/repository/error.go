package repository

import (
	"errors"

	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

func convertError(err error) error {
	if errors.Is(err, database.ErrNoRows) {
		return repository.ErrNotFound
	}
	return err
}
