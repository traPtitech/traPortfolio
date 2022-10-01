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
	"github.com/traPtitech/traPortfolio/util/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupGormDB(t *testing.T, sqlConf *config.SQLConfig) *gorm.DB {
	t.Helper()

	db := establishTestDBConnection(t, sqlConf)
	dropAll(t, db)
	init, err := migration.Migrate(db, migration.AllTables())
	assert.True(t, init)

	if err != nil {
		panic(err)
	}
	assert.NoError(t, err)

	return db
}

func SetupDB(t *testing.T, sqlConf *config.SQLConfig) database.SQLHandler {
	t.Helper()

	db := SetupGormDB(t, sqlConf)

	h := infrastructure.NewSQLHandler(db)
	return h
}

func establishTestDBConnection(t *testing.T, sqlConf *config.SQLConfig) *gorm.DB {
	t.Helper()

	dbDsn := sqlConf.DsnWithoutName()
	conn, err := sql.Open("mysql", dbDsn)
	assert.NoError(t, err)
	defer conn.Close()
	_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", sqlConf.DBName()))
	assert.NoError(t, err)

	dsn := sqlConf.Dsn()
	db, err := gorm.Open(mysql.Open(dsn), sqlConf.GormConfig())
	assert.NoError(t, err)

	return db
}

func WaitTestDBConnection(conf *config.SQLConfig) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		waitTestDBConnection(conf)
		close(ch)
	}()

	return ch
}

func waitTestDBConnection(conf *config.SQLConfig) {
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
	const dbPrefix = "portfolio_test_"

	return dbPrefix + dbName
}
