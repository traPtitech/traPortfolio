package handler

import (
	"database/sql"
	"io"
	"log"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/traPtitech/traPortfolio/util/config"
	"github.com/traPtitech/traPortfolio/util/testutils"
)

var (
	testConfig *config.Config
	testDB     *sql.DB
)

func TestMain(m *testing.M) {
	// TODO: loadをやめてテスト固有のconfigを使う
	c, err := config.Load(config.LoadOpts{SkipReadFromFiles: true})
	if err != nil {
		panic(err)
	}

	// disable mysql driver logging
	_ = mysql.SetLogger(mysql.Logger(log.New(io.Discard, "", 0)))
	db, port, closeFunc, err := testutils.RunMySQLContainerOnDocker(c.DB.User, c.DB.Pass, c.DB.Host)
	if err != nil {
		panic(err)
	}

	c.DB.Port = port

	defer func() {
		if err := closeFunc(); err != nil {
			panic(err)
		}
	}()

	testConfig = c
	testDB = db

	m.Run()
}
