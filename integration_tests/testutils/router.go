package testutils

import (
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traPortfolio/infrastructure"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/handler"
)

type initDBFunc func(database.SQLHandler) error

func SetupRoutes(t *testing.T, e *echo.Echo, dbName string, f initDBFunc) (*handler.API, error) {
	t.Helper()

	db := SetupDB(t, dbName)
	if f != nil {
		if err := f(db); err != nil {
			return nil, err
		}
	}

	s := infrastructure.NewSQLConfig("root", "password", "localhost", testDBName(dbName), 3307)
	tr := infrastructure.NewTraQConfig("", "", true)
	p := infrastructure.NewPortalConfig("", "", true)
	k := infrastructure.NewKnoqConfig("", "", true)

	api, err := infrastructure.InjectAPIServer(&s, &tr, &p, &k)
	if err != nil {
		return nil, err
	}

	if err := handler.Setup(e, api); err != nil {
		return nil, err
	}

	return &api, nil
}
