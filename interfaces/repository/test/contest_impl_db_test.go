//go:build integration && db && repository

package repository_test

import (
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external_e2e"
	irepository "github.com/traPtitech/traPortfolio/interfaces/repository"
	"github.com/traPtitech/traPortfolio/testutils"
	urepository "github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/random"
)

func TestContestRepository_GetContests(t *testing.T) {
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

func TestContestRepository_GetContest(t *testing.T) {
	t.Parallel()

	h := testutils.Setup(t, "get_contest")
	repo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())
	contest1 := mustMakeContest(t, repo, nil)
	contest2 := mustMakeContest(t, repo, nil)

	gotContest1, err := repo.GetContest(contest1.ID)
	require.NoError(t, err)
	gotContest2, err := repo.GetContest(contest2.ID)
	require.NoError(t, err)

	expected := []*domain.ContestDetail{contest1, contest2}
	gots := []*domain.ContestDetail{gotContest1, gotContest2}

	for i := range expected {
		assert.Equal(t, expected[i], gots[i])
	}
}

func TestContestRepository_UpdateContest(t *testing.T) {
	t.Parallel()

	h := testutils.Setup(t, "update_contest")
	repo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())
	contest1 := mustMakeContest(t, repo, nil)
	contest2 := mustMakeContest(t, repo, nil)

	args := urepository.UpdateContestArgs{
		Name:        random.OptAlphaNumeric(),
		Description: random.OptAlphaNumeric(),
		Link:        random.OptURLString(),
		Since:       random.OptTime(),
		Until:       random.OptTime(),
	}
	if args.Name.Valid {
		contest1.Name = args.Name.String
	}
	if args.Description.Valid {
		contest1.Description = args.Description.String
	}
	if args.Link.Valid {
		contest1.Link = args.Link.String
	}
	if args.Since.Valid {
		contest1.TimeStart = args.Since.Time
	}
	if args.Until.Valid {
		contest1.TimeEnd = args.Until.Time
	}

	err := repo.UpdateContest(contest1.ID, &args)
	assert.NoError(t, err)

	gotContest1, err := repo.GetContest(contest1.ID)
	require.NoError(t, err)
	gotContest2, err := repo.GetContest(contest2.ID)
	require.NoError(t, err)

	expected := []*domain.ContestDetail{contest1, contest2}
	gots := []*domain.ContestDetail{gotContest1, gotContest2}

	for i := range expected {
		assert.True(t, expected[i].TimeStart.Equal(gots[i].TimeStart))
		assert.True(t, expected[i].TimeEnd.Equal(gots[i].TimeEnd))
		expected[i].TimeStart = time.Time{}
		expected[i].TimeEnd = time.Time{}
		gots[i].TimeStart = time.Time{}
		gots[i].TimeEnd = time.Time{}
		assert.Equal(t, expected[i], gots[i])
	}
}

func TestContestRepository_DeleteContest(t *testing.T) {
	t.Parallel()

	h := testutils.Setup(t, "delete_contest")
	repo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())
	contest1 := mustMakeContest(t, repo, nil)
	contest2 := mustMakeContest(t, repo, nil)

	gotContest1, err := repo.GetContest(contest1.ID)
	require.NoError(t, err)
	gotContest2, err := repo.GetContest(contest2.ID)
	require.NoError(t, err)

	expected := []*domain.ContestDetail{contest1, contest2}
	gots := []*domain.ContestDetail{gotContest1, gotContest2}

	for i := range expected {
		assert.Equal(t, expected[i], gots[i])
	}

	err = repo.DeleteContest(contest1.ID)
	assert.NoError(t, err)

	deletedContest1, err := repo.GetContest(contest1.ID)
	assert.Nil(t, deletedContest1)
	assert.Equal(t, err, urepository.ErrNotFound)

	gotContest2, err = repo.GetContest(contest2.ID)
	require.NoError(t, err)
	assert.Equal(t, contest2, gotContest2)
}

func TestContestRepository_GetContestTeams(t *testing.T) {
	t.Parallel()

	h := testutils.Setup(t, "get_contest_teams")
	repo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())

	contest1 := mustMakeContest(t, repo, nil)
	contest2 := mustMakeContest(t, repo, nil)
	team1 := mustMakeContestTeam(t, repo, contest1.ID, &urepository.CreateContestTeamArgs{
		Name:        random.AlphaNumeric(),
		Result:      random.OptAlphaNumeric(),
		Link:        random.OptURLString(),
		Description: random.AlphaNumeric(),
	})
	team2 := mustMakeContestTeam(t, repo, contest1.ID, &urepository.CreateContestTeamArgs{
		Name:        random.AlphaNumeric(),
		Result:      random.OptAlphaNumeric(),
		Link:        random.OptURLString(),
		Description: random.AlphaNumeric(),
	})

	gotTeams1, err := repo.GetContestTeams(contest1.ID)
	require.NoError(t, err)

	gotTeams2, err := repo.GetContestTeams(contest2.ID)
	require.NoError(t, err)

	expected := [][]*domain.ContestTeam{gotTeams1, {}}
	gots := [][]*domain.ContestTeam{{&team1.ContestTeam, &team2.ContestTeam}, gotTeams2}

	for i := range expected {
		expectedTeams := expected[i]
		sort.SliceStable(expectedTeams, func(i, j int) bool {
			return expectedTeams[i].ID.String() > expectedTeams[j].ID.String()
		})

		gotsTeams := gots[i]
		sort.SliceStable(gotsTeams, func(i, j int) bool {
			return gotsTeams[i].ID.String() > gotsTeams[j].ID.String()
		})

		assert.Equal(t, expectedTeams, gotsTeams)
	}
}
