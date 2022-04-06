//go:build integration && db

package testutils

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/infrastructure"
	"github.com/traPtitech/traPortfolio/infrastructure/migration"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDB(t *testing.T, sqlConf *infrastructure.SQLConfig) database.SQLHandler {
	t.Helper()

	db := establishTestDBConnection(t, sqlConf)
	dropAll(t, db)
	init, err := migration.Migrate(db, migration.AllTables())
	assert.True(t, init)

	if err != nil {
		panic(err)
	}
	assert.NoError(t, err)
	h := infrastructure.FromDB(db)
	return h
}

func establishTestDBConnection(t *testing.T, sqlConf *infrastructure.SQLConfig) *gorm.DB {
	t.Helper()

	dbDsn := sqlConf.DsnWithoutName()
	conn, err := sql.Open("mysql", dbDsn)
	assert.NoError(t, err)
	defer conn.Close()
	_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", sqlConf.Name()))
	assert.NoError(t, err)

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	dsn := sqlConf.Dsn()
	db, err := gorm.Open(mysql.Open(dsn), config)
	assert.NoError(t, err)

	return db
}

func WaitTestDBConnection(conf *infrastructure.SQLConfig) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		waitTestDBConnection(conf)
		close(ch)
	}()

	return ch
}

func waitTestDBConnection(conf *infrastructure.SQLConfig) {
	dbDsn := conf.Dsn()
	conn, err := sql.Open("mysql", dbDsn)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// DBとの接続が確立できるまで待つ
	for i := 0; ; i++ {
		log.Println(i)
		if i > 10 {
			panic(fmt.Errorf("failed to connect to DB"))
		}
		err = conn.Ping()
		log.Println(err)
		if err == nil {
			break
		}
		time.Sleep(time.Second * 10)
	}
}

func dropAll(t *testing.T, db *gorm.DB) {
	t.Helper()

	tables := []interface{}{"migrations"}
	tables = append(tables, migration.AllTables()...)

	err := db.Migrator().DropTable(tables...)
	assert.NoError(t, err)
}

func testDBName(dbName string) string {
	const dbPrefix = "portfolio_test_repo_"

	return dbPrefix + dbName
}
