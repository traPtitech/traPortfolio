package handler

import (
	"io"
	"log"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
	"github.com/traPtitech/traPortfolio/internal/util/config"
)

func TestMain(m *testing.M) {
	c, err := config.Load(config.LoadOpts{SkipReadFromFiles: true})
	if err != nil {
		panic(err)
	}

	// disable mysql driver logging
	_ = mysql.SetLogger(mysql.Logger(log.New(io.Discard, "", 0)))
	db, closeFunc, err := testutils.RunMySQLContainerOnDocker(c.DB)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := closeFunc(); err != nil {
			panic(err)
		}
	}()

	testutils.Config = c
	testutils.DB = db

	m.Run()
}
