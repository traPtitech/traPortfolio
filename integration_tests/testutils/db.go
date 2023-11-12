package testutils

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/infrastructure/migration"
	"github.com/traPtitech/traPortfolio/infrastructure/repository"
	"github.com/traPtitech/traPortfolio/util/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TODO: integration_tests/handler以下に置く
var (
	DB   *sql.DB
	Port int
)

func SetupGormDB(t *testing.T, sqlConf *config.SQLConfig) *gorm.DB {
	t.Helper()

	db := establishTestDBConnection(t, sqlConf)
	dropAll(t, db)
	init, err := migration.Migrate(db, migration.AllTables())
	assert.True(t, init)
	assert.NoError(t, err)

	return db
}

func establishTestDBConnection(t *testing.T, sqlConf *config.SQLConfig) *gorm.DB {
	t.Helper()

	_, err := DB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", sqlConf.Name))
	assert.NoError(t, err)

	sqlConf.Port = Port
	db, err := repository.NewGormDB(sqlConf)
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
	db, err := gorm.Open(mysql.New(mysql.Config{DSNConfig: conf.DsnConfig()}), conf.GormConfig())
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	// DBとの接続が確立できるまで待つ
	for i := 0; ; i++ {
		log.Println(i)
		if i > 10 {
			panic(fmt.Errorf("failed to connect to DB"))
		}
		err = sqlDB.Ping()
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

func RunMySQLContainerOnDocker(sqlConf *config.SQLConfig) (*sql.DB, func() error, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("create docker pool: %w", err)
	}

	pool.MaxWait = 1 * time.Minute

	if err := pool.Client.Ping(); err != nil {
		return nil, nil, fmt.Errorf("ping docker: %w", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mariadb",
		Tag:        "10.6.4",
		Env:        []string{"MYSQL_ROOT_PASSWORD=" + sqlConf.Pass},
	}, func(hc *docker.HostConfig) {
		hc.AutoRemove = true
		hc.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, nil, fmt.Errorf("run docker container: %w", err)
	}

	if err := resource.Expire(60); err != nil {
		return nil, nil, fmt.Errorf("expire docker container: %w", err)
	}

	sqlConf.Name = ""
	sqlConf.Port, _ = strconv.Atoi(resource.GetPort("3306/tcp"))

	var db *sql.DB

	Port, _ = strconv.Atoi(resource.GetPort("3306/tcp"))

	if err := pool.Retry(func() error {
		_db, err := sql.Open("mysql", sqlConf.DsnConfig().FormatDSN())
		if err != nil {
			return fmt.Errorf("open mysql: %w", err)
		}

		if err := _db.Ping(); err != nil {
			return fmt.Errorf("ping mysql: %w", err)
		}

		db = _db

		return nil
	}); err != nil {
		return nil, nil, fmt.Errorf("retry: %w", err)
	}

	closeFunc := func() error {
		if err := db.Close(); err != nil {
			return fmt.Errorf("close mysql: %w", err)
		}

		if err := pool.Purge(resource); err != nil {
			return fmt.Errorf("purge docker container: %w", err)
		}

		return nil
	}

	return db, closeFunc, nil
}
