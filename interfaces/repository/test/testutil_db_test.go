//go:build integration && db

package repository_test

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/infrastructure"
	"github.com/traPtitech/traPortfolio/infrastructure/migration"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/optional"
	"github.com/traPtitech/traPortfolio/util/random"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	dbPrefix = "portfolio_test_repo_"
)

func TestMain(m *testing.M) {
	dbUser := getEnvOrDefault("MYSQL_USERNAME", "root")
	dbPass := getEnvOrDefault("MYSQL_PASSWORD", "password")
	dbHost := getEnvOrDefault("MYSQL_HOSTNAME", "127.0.0.1")
	dbPort := getEnvOrDefault("MYSQL_PORT", "3306")

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
			panic(fmt.Errorf("faild to connect to DB"))
		}
		err = conn.Ping()
		log.Println(err)
		if err == nil {
			break
		}
		time.Sleep(time.Second * 10)
	}

	m.Run()
}

func setup(t *testing.T, dbName string) database.SQLHandler {
	t.Helper()

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
	dbUser := getEnvOrDefault("MYSQL_USERNAME", "root")
	dbPass := getEnvOrDefault("MYSQL_PASSWORD", "password")
	dbHost := getEnvOrDefault("MYSQL_HOSTNAME", "127.0.0.1")
	dbPort := getEnvOrDefault("MYSQL_PORT", "3306")

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

	allTables := []interface{}{"migrations"}
	allTables = append(allTables, AllTables()...)

	err := db.Migrator().DropTable(allTables...)
	require.NoError(t, err)
}

func AllTables() []interface{} {
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

func mustMakeContest(t *testing.T, repo repository.ContestRepository, args *repository.CreateContestArgs) *domain.ContestDetail {
	t.Helper()

	if args == nil {
		var link optional.String
		if rand.Intn(2) == 1 {
			link = optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true)
		}

		var t optional.Time
		if rand.Intn(2) == 1 {
			t = optional.NewTime(random.Time(), true)
		}

		args = &repository.CreateContestArgs{
			Name:        random.AlphaNumeric(rand.Intn(30) + 1),
			Description: random.AlphaNumeric(rand.Intn(30) + 1),
			Link:        link,
			Since:       random.Time(),
			Until:       t,
		}
	}

	contest, err := repo.CreateContest(args)
	require.NoError(t, err)
	return contest
}

func mustMakeContestTeam(t *testing.T, repo repository.ContestRepository, contestID uuid.UUID, args *repository.CreateContestTeamArgs) *domain.ContestTeamDetail {
	t.Helper()

	if args == nil {
		var result optional.String
		if rand.Intn(2) == 0 {
			result = optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true)
		}

		var link optional.String
		if rand.Intn(2) == 0 {
			link = optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true)
		}

		args = &repository.CreateContestTeamArgs{
			Name:        random.AlphaNumeric(rand.Intn(30) + 1),
			Result:      result,
			Link:        link,
			Description: random.AlphaNumeric(rand.Intn(30) + 1),
		}
	}

	var err error
	var contestTeamDetail *domain.ContestTeamDetail
	if contestID == uuid.Nil {
		contest := mustMakeContest(t, nil, nil)
		contestTeamDetail, err = repo.CreateContestTeam(contest.ID, args)
	} else {
		contestTeamDetail, err = repo.CreateContestTeam(contestID, args)
	}

	require.NoError(t, err)

	return contestTeamDetail
}

func mustMakeUser(t *testing.T, repo repository.UserRepository, args *repository.CreateUserArgs) *domain.UserDetail {
	t.Helper()

	if args == nil {
		args = &repository.CreateUserArgs{
			Description: random.AlphaNumeric(rand.Intn(30) + 1),
			Check:       random.Bool(),
			Name:        random.AlphaNumeric(rand.Intn(30) + 1),
		}
	}

	user, err := repo.CreateUser(*args)
	require.NoError(t, err)

	return user
}

func getEnvOrDefault(env string, def string) string {
	s := os.Getenv(env)
	if len(s) == 0 {
		return def
	}
	return s
}
