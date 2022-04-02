//go:build integration && db

package testutils

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/traPortfolio/infrastructure"
	"github.com/traPtitech/traPortfolio/infrastructure/migration"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
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
	require.NoError(t, err)
	h := infrastructure.FromDB(db)
	return h
}

func establishTestDBConnection(t *testing.T, dbName string) *gorm.DB {
	t.Helper()
	dbUser := GetEnvOrDefault("MYSQL_USERNAME", "root")
	dbPass := GetEnvOrDefault("MYSQL_PASSWORD", "password")
	dbHost := GetEnvOrDefault("MYSQL_HOSTNAME", "127.0.0.1")
	dbPort := GetEnvOrDefault("MYSQL_PORT", "3306")

	dbDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort)
	conn, err := sql.Open("mysql", dbDsn)
	require.NoError(t, err)
	defer conn.Close()
	_, err = conn.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName))
	assert.NoError(t, err)

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), config)
	require.NoError(t, err)

	return db
}

func dropAll(t *testing.T, db *gorm.DB) {
	t.Helper()

	tables := []interface{}{"migrations"}
	tables = append(tables, allTables()...)

	err := db.Migrator().DropTable(tables...)
	require.NoError(t, err)
}

func allTables() []interface{} {
	return []interface{}{
		model.User{},
		model.Account{},
		model.Project{},
		model.ProjectMember{},
		model.EventLevelRelation{},
		model.Contest{},
		model.ContestTeam{},
		model.ContestTeamUserBelonging{},
		model.Group{},
		model.GroupUserBelonging{},
	}
}

func GetEnvOrDefault(env string, def string) string {
	s := os.Getenv(env)
	if len(s) == 0 {
		return def
	}
	return s
}
