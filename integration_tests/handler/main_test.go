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
	if err := testutils.ParseConfig("../testdata"); err != nil {
		panic(err)
	}

	c := config.Load()

	// disable mysql driver logging
	_ = mysql.SetLogger(mysql.Logger(log.New(io.Discard, "", 0)))
	_db, closeFunc, err := testutils.RunMySQLContainerOnDocker(c.SQLConf())
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
