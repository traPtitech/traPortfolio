//go:build integration && db && repository

package repository_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external_e2e"
	"github.com/traPtitech/traPortfolio/testutils"

	irepository "github.com/traPtitech/traPortfolio/interfaces/repository"
)

func TestContestRepositoryDB_GetContests(t *testing.T) {
	t.Parallel()

	h := testutils.Setup(t, "get_contests")
	repo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())
	contest1 := mustMakeContest(t, repo, nil)
	contest2 := mustMakeContest(t, repo, nil)
	contests, err := repo.GetContests()
	require.NoError(t, err)

	expected := []*domain.Contest{&contest1.Contest, &contest2.Contest}
	sort.SliceStable(expected, func(i, j int) bool {
		return expected[i].ID.String() > expected[j].ID.String()
	})
	sort.SliceStable(contests, func(i, j int) bool {
		return contests[i].ID.String() > contests[j].ID.String()
	})

	assert.Equal(t, expected, contests)
}
