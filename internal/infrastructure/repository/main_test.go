package repository

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"strings"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/internal/domain"
	"github.com/traPtitech/traPortfolio/internal/pkgs/config"
	"github.com/traPtitech/traPortfolio/internal/pkgs/optional"
	"github.com/traPtitech/traPortfolio/internal/pkgs/random"
	"github.com/traPtitech/traPortfolio/internal/pkgs/testutils"
	"github.com/traPtitech/traPortfolio/internal/usecases/repository"
	"gorm.io/gorm"
)

var (
	testDB        *sql.DB
	testSQLConfig config.SQLConfig
)

func TestMain(m *testing.M) {
	// disable mysql driver logging
	_ = mysql.SetLogger(mysql.Logger(log.New(io.Discard, "", 0)))

	user, pass, host := "root", "password", "localhost"
	db, port, closeFunc, err := testutils.RunMySQLContainerOnDocker(user, pass, host)
	if err != nil {
		panic(err)
	}

	testDB = db
	testSQLConfig = config.SQLConfig{
		User: user,
		Pass: pass,
		Host: host,
		Name: "",
		Port: port,
	}

	m.Run()

	if err := closeFunc(); err != nil {
		panic(err)
	}
}

func SetupTestGormDB(t *testing.T) *gorm.DB {
	t.Helper()

	dbName := fmt.Sprintf("portfolio_test_%s", strings.ToLower(t.Name()))
	_, err := testDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName))
	assert.NoError(t, err)

	sqlConfig := testSQLConfig
	sqlConfig.Name = dbName

	gormDB, err := NewGormDB(sqlConfig)
	assert.NoError(t, err)

	return gormDB
}

func mustMakeContest(t *testing.T, repo repository.ContestRepository, args *repository.CreateContestArgs) *domain.ContestDetail {
	t.Helper()

	if args == nil {
		since, until := random.SinceAndUntil()
		args = &repository.CreateContestArgs{
			Name:        random.AlphaNumeric(),
			Description: random.AlphaNumeric(),
			Links:       random.Array(random.RandURLString, 1, 3),
			Since:       since,
			Until:       optional.New(until, random.Bool()),
		}
	}

	contest, err := repo.CreateContest(context.Background(), args)
	assert.NoError(t, err)
	return contest
}

func mustMakeContestTeam(t *testing.T, repo repository.ContestRepository, contestID uuid.UUID, args *repository.CreateContestTeamArgs) *domain.ContestTeamDetail {
	t.Helper()

	if args == nil {
		args = &repository.CreateContestTeamArgs{
			Name:        random.AlphaNumeric(),
			Result:      random.Optional(random.AlphaNumeric()),
			Links:       random.Array(random.RandURLString, 1, 3),
			Description: random.AlphaNumeric(),
		}
	}

	var err error
	var contestTeamDetail *domain.ContestTeamDetail
	if contestID == uuid.Nil {
		contest := mustMakeContest(t, nil, nil)
		contestTeamDetail, err = repo.CreateContestTeam(context.Background(), contest.ID, args)
	} else {
		contestTeamDetail, err = repo.CreateContestTeam(context.Background(), contestID, args)
	}

	assert.NoError(t, err)

	return contestTeamDetail
}

func mustMakeEventLevel(t *testing.T, repo repository.EventRepository, args *repository.CreateEventLevelArgs) *repository.CreateEventLevelArgs {
	t.Helper()

	if args == nil {
		args = &repository.CreateEventLevelArgs{
			EventID: random.UUID(),
			Level:   rand.N(domain.EventLevelLimit),
		}
	}

	err := repo.CreateEventLevel(context.Background(), args)
	assert.NoError(t, err)

	return args
}

func mustMakeAccount(t *testing.T, repo repository.UserRepository, userID uuid.UUID, args *repository.CreateAccountArgs) *domain.Account {
	t.Helper()

	if args == nil {
		accountType := rand.N(domain.AccountLimit)
		args = &repository.CreateAccountArgs{
			DisplayName: random.AlphaNumeric(),
			Type:        accountType,
			URL:         random.AccountURLString(accountType),
			PrPermitted: random.Bool(),
		}
	}

	assert.NotEmpty(t, userID)
	account, err := repo.CreateAccount(context.Background(), userID, args)
	assert.NoError(t, err)

	return account
}

func mustMakeProject(t *testing.T, repo repository.ProjectRepository, args *repository.CreateProjectArgs) *domain.Project {
	t.Helper()

	project := mustMakeProjectDetail(t, repo, args)

	return &project.Project
}

func mustMakeProjectDetail(t *testing.T, repo repository.ProjectRepository, args *repository.CreateProjectArgs) *domain.ProjectDetail {
	t.Helper()

	if args == nil {
		duration := random.Duration()
		args = &repository.CreateProjectArgs{
			Name:          random.AlphaNumeric(),
			Description:   random.AlphaNumeric(),
			Links:         random.Array(random.RandURLString, 1, 3),
			SinceYear:     duration.Since.Year,
			SinceSemester: duration.Since.Semester,
			UntilYear:     duration.Until.ValueOrZero().Year,
			UntilSemester: duration.Until.ValueOrZero().Semester,
		}
	}

	project, err := repo.CreateProject(context.Background(), args)
	assert.NoError(t, err)

	return project
}

func mustExistProjectMember(t *testing.T, repo repository.ProjectRepository, projectID uuid.UUID, projectDuration domain.YearWithSemesterDuration, users []*repository.EditProjectMemberArgs) {
	t.Helper()

	assert.NotEmpty(t, projectID)
	assert.NotEmpty(t, projectDuration)
	assert.NotEmpty(t, users)
	assert.True(t, projectDuration.IsValid())
	for _, user := range users {
		assert.NotEmpty(t, user.UserID)
		userDuration := domain.NewYearWithSemesterDuration(user.SinceYear, user.SinceSemester, user.UntilYear, user.UntilSemester)
		assert.True(t, userDuration.IsValid())
	}

	var duration = random.DurationBetween(projectDuration.Since, projectDuration.Until.ValueOrZero())

	assert.True(t, duration.IsValid())

	err := repo.EditProjectMembers(context.Background(), projectID, users)
	assert.NoError(t, err)
}

func mustExistContestTeamMembers(t *testing.T, repo repository.ContestRepository, teamID uuid.UUID, userIDs []uuid.UUID) {
	t.Helper()

	for _, id := range userIDs {
		assert.NotEmpty(t, id)
	}

	assert.NotEmpty(t, teamID)

	err := repo.EditContestTeamMembers(context.Background(), teamID, userIDs)
	assert.NoError(t, err)
}
