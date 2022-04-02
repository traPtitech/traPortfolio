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
	"github.com/traPtitech/traPortfolio/util/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Setup(t *testing.T, dbName string) database.SQLHandler {
	t.Helper()
	const dbPrefix = "portfolio_test_repo_"

	db := establishTestDBConnection(t, dbPrefix+dbName)
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

func establishTestDBConnection(t *testing.T, dbName string) *gorm.DB {
	t.Helper()
	dbUser := config.GetEnvOrDefault("MYSQL_USERNAME", "root")
	dbPass := config.GetEnvOrDefault("MYSQL_PASSWORD", "password")
	dbHost := config.GetEnvOrDefault("MYSQL_HOSTNAME", "127.0.0.1")
	dbPort := config.GetEnvOrDefault("MYSQL_PORT", "3307")

	dbDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort)
	conn, err := sql.Open("mysql", dbDsn)
	assert.NoError(t, err)
	defer conn.Close()
	_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName))
	assert.NoError(t, err)

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), config)
	assert.NoError(t, err)

	return db
}

func WaitTestDBConnection() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		waitTestDBConnection()
		close(ch)
	}()

	return ch
}

func waitTestDBConnection() {
	dbUser := config.GetEnvOrDefault("MYSQL_USERNAME", "root")
	dbPass := config.GetEnvOrDefault("MYSQL_PASSWORD", "password")
	dbHost := config.GetEnvOrDefault("MYSQL_HOSTNAME", "127.0.0.1")
	dbPort := config.GetEnvOrDefault("MYSQL_PORT", "3307")

	dbDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort)
	log.Println(dbDsn)
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
