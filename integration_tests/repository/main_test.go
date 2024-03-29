package repository

import (
	"context"
	"io"
	"log"
	"math/rand/v2"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/config"
	"github.com/traPtitech/traPortfolio/util/optional"
	"github.com/traPtitech/traPortfolio/util/random"
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

func mustMakeContest(t *testing.T, repo repository.ContestRepository, args *repository.CreateContestArgs) *domain.ContestDetail {
	t.Helper()

	if args == nil {
		since, until := random.SinceAndUntil()
		args = &repository.CreateContestArgs{
			Name:        random.AlphaNumeric(),
			Description: random.AlphaNumeric(),
			Link:        random.Optional(random.RandURLString()),
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
			Link:        random.Optional(random.RandURLString()),
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
		args = &repository.CreateProjectArgs{
			Name:          random.AlphaNumeric(),
			Description:   random.AlphaNumeric(),
			Link:          random.Optional(random.RandURLString()),
			SinceYear:     rand.IntN(2100),
			SinceSemester: rand.IntN(2),
			UntilYear:     rand.IntN(2100),
			UntilSemester: rand.IntN(2),
		}
	}

	project, err := repo.CreateProject(context.Background(), args)
	assert.NoError(t, err)

	return project
}

func mustAddProjectMember(t *testing.T, repo repository.ProjectRepository, projectID uuid.UUID, userID uuid.UUID, args *repository.CreateProjectMemberArgs) *repository.CreateProjectMemberArgs {
	t.Helper()

	assert.NotEmpty(t, projectID)
	assert.NotEmpty(t, userID)

	if args == nil {
		args = &repository.CreateProjectMemberArgs{
			UserID:        userID,
			SinceYear:     rand.IntN(2100),
			SinceSemester: rand.IntN(2),
			UntilYear:     rand.IntN(2100),
			UntilSemester: rand.IntN(2),
		}
	}

	err := repo.AddProjectMembers(context.Background(), projectID, []*repository.CreateProjectMemberArgs{args})
	assert.NoError(t, err)

	return args
}

func mustAddContestTeamMembers(t *testing.T, repo repository.ContestRepository, teamID uuid.UUID, userIDs []uuid.UUID) {
	t.Helper()

	for _, id := range userIDs {
		assert.NotEmpty(t, id)
	}

	assert.NotEmpty(t, teamID)

	err := repo.AddContestTeamMembers(context.Background(), teamID, userIDs)
	assert.NoError(t, err)
}
