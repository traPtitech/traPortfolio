package handler

import (
	"io"
	"log"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
	"github.com/traPtitech/traPortfolio/util/config"
)

func TestMain(m *testing.M) {
	_c, err := config.Load(config.LoadOpts{SkipReadFromFiles: true})
	if err != nil {
		panic(err)
	}

	testutils.Config = _c

	// disable mysql driver logging
	_ = mysql.SetLogger(mysql.Logger(log.New(io.Discard, "", 0)))
	_db, closeFunc, err := testutils.RunMySQLContainerOnDocker(testutils.Config.DB)
	if err != nil {
		panic(err)
	}

	testutils.DB = _db

	defer func() {
		if err := closeFunc(); err != nil {
			panic(err)
		}
	}()

	m.Run()
}
