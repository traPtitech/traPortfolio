//go:build integration && db

package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external_e2e"

	irepository "github.com/traPtitech/traPortfolio/interfaces/repository"
)

func TestContestRepositoryDB_GetContests(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("get_contests")
	sqlConf := conf.SQLConf()

	h := testutils.SetupDB(t, &sqlConf)
	repo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())
	contest1 := mustMakeContest(t, repo, nil)
	contest2 := mustMakeContest(t, repo, nil)
	contests, err := repo.GetContests()
	assert.NoError(t, err)

	expected := []*domain.Contest{&contest1.Contest, &contest2.Contest}

	assert.ElementsMatch(t, expected, contests)
}
