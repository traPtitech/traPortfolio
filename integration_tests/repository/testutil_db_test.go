//go:build integration && db

package repository_test

import (
	"math/rand"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/optional"
	"github.com/traPtitech/traPortfolio/util/random"
)

func TestMain(m *testing.M) {
	testutils.ParseConfig("../testdata")

	appConf := testutils.GetConfig()
	sqlConf := appConf.SQLConf()
	<-testutils.WaitTestDBConnection(sqlConf)

	m.Run()
}

func mustMakeContest(t *testing.T, repo repository.ContestRepository, args *repository.CreateContestArgs) *domain.ContestDetail {
	t.Helper()

	if args == nil {
		var link optional.String
		if rand.Intn(2) == 1 {
			link = optional.NewString(random.AlphaNumeric(), true)
		}

		var t optional.Time
		if rand.Intn(2) == 1 {
			t = optional.NewTime(random.Time(), true)
		}

		args = &repository.CreateContestArgs{
			Name:        random.AlphaNumeric(),
			Description: random.AlphaNumeric(),
			Link:        link,
			Since:       random.Time(),
			Until:       t,
		}
	}

	contest, err := repo.CreateContest(args)
	assert.NoError(t, err)
	return contest
}

func mustMakeContestTeam(t *testing.T, repo repository.ContestRepository, contestID uuid.UUID, args *repository.CreateContestTeamArgs) *domain.ContestTeamDetail {
	t.Helper()

	if args == nil {
		args = &repository.CreateContestTeamArgs{
			Name:        random.AlphaNumeric(),
			Result:      random.OptAlphaNumeric(),
			Link:        random.OptURLString(),
			Description: random.AlphaNumeric(),
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

	assert.NoError(t, err)

	return contestTeamDetail
}

// func mustMakeUser(t *testing.T, repo repository.UserRepository, args *repository.CreateUserArgs) *domain.UserDetail {
// 	t.Helper()

// 	if args == nil {
// 		args = &repository.CreateUserArgs{
// 			Description: random.AlphaNumeric(),
// 			Check:       random.Bool(),
// 			Name:        random.AlphaNumeric(),
// 		}
// 	}

// 	user, err := repo.CreateUser(*args)
// 	assert.NoError(t, err)

// 	return user
// }
