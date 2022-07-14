package repository

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external_e2e"

	irepository "github.com/traPtitech/traPortfolio/interfaces/repository"
	urepository "github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/mockdata"
	"github.com/traPtitech/traPortfolio/util/random"
)

func TestContestRepository_GetContests(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("contest_repository_get_contests")
	sqlConf := conf.SQLConf()

	h := testutils.SetupDB(t, sqlConf)
	repo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())
	contest1 := mustMakeContest(t, repo, nil)
	contest2 := mustMakeContest(t, repo, nil)
	contests, err := repo.GetContests()
	assert.NoError(t, err)

	expected := []*domain.Contest{&contest1.Contest, &contest2.Contest}

	assert.ElementsMatch(t, expected, contests)
}

func TestContestRepository_GetContest(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("contest_repository_get_contest")
	sqlConf := conf.SQLConf()

	h := testutils.SetupDB(t, sqlConf)
	repo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())
	contest1 := mustMakeContest(t, repo, nil)
	contest2 := mustMakeContest(t, repo, nil)

	gotContest1, err := repo.GetContest(contest1.ID)
	assert.NoError(t, err)
	assert.Equal(t, contest1, gotContest1)

	gotContest2, err := repo.GetContest(contest2.ID)
	assert.NoError(t, err)
	assert.Equal(t, contest2, gotContest2)
}

func TestContestRepository_UpdateContest(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("contest_repository_update_contest")
	sqlConf := conf.SQLConf()

	h := testutils.SetupDB(t, sqlConf)
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
	assert.NoError(t, err)
	gotContest2, err := repo.GetContest(contest2.ID)
	assert.NoError(t, err)

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

	conf := testutils.GetConfigWithDBName("contest_repository_delete_contest")
	sqlConf := conf.SQLConf()

	h := testutils.SetupDB(t, sqlConf)
	repo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())
	contest1 := mustMakeContest(t, repo, nil)
	contest2 := mustMakeContest(t, repo, nil)

	gotContest1, err := repo.GetContest(contest1.ID)
	assert.NoError(t, err)
	assert.Equal(t, contest1, gotContest1)

	gotContest2, err := repo.GetContest(contest2.ID)
	assert.NoError(t, err)
	assert.Equal(t, contest2, gotContest2)

	err = repo.DeleteContest(contest1.ID)
	assert.NoError(t, err)

	deletedContest1, err := repo.GetContest(contest1.ID)
	assert.Nil(t, deletedContest1)
	assert.Equal(t, err, urepository.ErrNotFound)

	gotContest2, err = repo.GetContest(contest2.ID)
	assert.NoError(t, err)
	assert.Equal(t, contest2, gotContest2)
}

func TestContestRepository_GetContestTeams(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("contest_repository_get_contest_teams")
	sqlConf := conf.SQLConf()

	h := testutils.SetupDB(t, sqlConf)
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

	expected1 := []*domain.ContestTeam{&team1.ContestTeam, &team2.ContestTeam}
	gotTeams1, err := repo.GetContestTeams(contest1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected1, gotTeams1)

	expected2 := []*domain.ContestTeam{}
	gotTeams2, err := repo.GetContestTeams(contest2.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected2, gotTeams2)
}

func TestContestRepository_GetContestTeam(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("contest_repository_get_contest_team")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
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

	gotTeams1, err := repo.GetContestTeam(contest1.ID, team1.ID)
	assert.NoError(t, err)
	assert.Equal(t, team1, gotTeams1)

	gotTeams2, err := repo.GetContestTeam(contest2.ID, team2.ID)
	assert.Error(t, err)
	assert.Nil(t, gotTeams2)
}

// func TestCreateContestTeam(t *testing.T) {
// }

func TestContestRepository_UpdateContestTeam(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("contest_repository_update_contest_teams")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
	repo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())

	contest1 := mustMakeContest(t, repo, nil)
	mustMakeContest(t, repo, nil)
	team1 := mustMakeContestTeam(t, repo, contest1.ID, &urepository.CreateContestTeamArgs{
		Name:        random.AlphaNumeric(),
		Result:      random.OptAlphaNumeric(),
		Link:        random.OptURLString(),
		Description: random.AlphaNumeric(),
	})
	mustMakeContestTeam(t, repo, contest1.ID, &urepository.CreateContestTeamArgs{
		Name:        random.AlphaNumeric(),
		Result:      random.OptAlphaNumeric(),
		Link:        random.OptURLString(),
		Description: random.AlphaNumeric(),
	})

	args1 := &urepository.UpdateContestTeamArgs{
		Name:        random.OptAlphaNumeric(),
		Result:      random.OptAlphaNumeric(),
		Link:        random.OptURLString(),
		Description: random.OptAlphaNumeric(),
	}
	if args1.Name.Valid {
		team1.Name = args1.Name.String
	}
	if args1.Result.Valid {
		team1.Result = args1.Result.String
	}
	if args1.Link.Valid {
		team1.Link = args1.Link.String
	}
	if args1.Description.Valid {
		team1.Description = args1.Description.String
	}

	err := repo.UpdateContestTeam(team1.ID, args1)
	assert.NoError(t, err)

	got, err := repo.GetContestTeam(contest1.ID, team1.ID)
	assert.NoError(t, err)
	assert.Equal(t, team1, got)
}

func TestContestRepository_DeleteContestTeam(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("contest_repository_delete_contest_teams")
	sqlConf := conf.SQLConf()

	h := testutils.SetupDB(t, sqlConf)
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

	expected1 := []*domain.ContestTeam{&team1.ContestTeam, &team2.ContestTeam}
	gotTeams1, err := repo.GetContestTeams(contest1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected1, gotTeams1)

	expected2 := []*domain.ContestTeam{}
	gotTeams2, err := repo.GetContestTeams(contest2.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected2, gotTeams2)

	err = repo.DeleteContestTeam(contest1.ID, team1.ID)
	assert.NoError(t, err)
	expected3 := []*domain.ContestTeam{&team2.ContestTeam}
	gotTeams3, err := repo.GetContestTeams(contest1.ID)
	assert.NoError(t, err)
	assert.Equal(t, expected3, gotTeams3)

	err = repo.DeleteContestTeam(contest1.ID, team2.ID)
	assert.NoError(t, err)
	expected4 := []*domain.ContestTeam{}
	gotTeams4, err := repo.GetContestTeams(contest1.ID)
	assert.NoError(t, err)
	assert.Equal(t, expected4, gotTeams4)
}

func TestContestRepository_GetContestTeamMembers(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("contest_repository_get_contest_team_members")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
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
	team3 := mustMakeContestTeam(t, repo, contest1.ID, &urepository.CreateContestTeamArgs{
		Name:        random.AlphaNumeric(),
		Result:      random.OptAlphaNumeric(),
		Link:        random.OptURLString(),
		Description: random.AlphaNumeric(),
	})
	user1 := mockdata.MockUsers[1]
	user2 := mockdata.MockUsers[2]
	portalUser1 := mockdata.MockPortalUsers[1]
	portalUser2 := mockdata.MockPortalUsers[2]

	mustAddContestTeamMembers(t, repo, team1.ID, []uuid.UUID{user1.ID})
	mustAddContestTeamMembers(t, repo, team1.ID, []uuid.UUID{user2.ID})
	mustAddContestTeamMembers(t, repo, team2.ID, []uuid.UUID{user2.ID})

	expected1 := []*domain.User{
		{
			ID:       user1.ID,
			Name:     user1.Name,
			RealName: portalUser1.RealName,
		},
		{
			ID:       user2.ID,
			Name:     user2.Name,
			RealName: portalUser2.RealName,
		},
	}
	users1, err := repo.GetContestTeamMembers(contest1.ID, team1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected1, users1)

	expected2 := []*domain.User{
		{
			ID:       user2.ID,
			Name:     user2.Name,
			RealName: portalUser2.RealName,
		},
	}
	users2, err := repo.GetContestTeamMembers(contest1.ID, team2.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected2, users2)

	expected3 := []*domain.User{}
	users3, err := repo.GetContestTeamMembers(contest2.ID, team3.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected3, users3)
}

// func TestContestRepository_AddContestTeamMembers(t *testing.T) {
// }

func TestContestRepository_EditContestTeamMembers(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("contest_repository_edit_contest_team_members")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
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
	team3 := mustMakeContestTeam(t, repo, contest1.ID, &urepository.CreateContestTeamArgs{
		Name:        random.AlphaNumeric(),
		Result:      random.OptAlphaNumeric(),
		Link:        random.OptURLString(),
		Description: random.AlphaNumeric(),
	})
	user1 := mockdata.MockUsers[1]
	user2 := mockdata.MockUsers[2]
	portalUser1 := mockdata.MockPortalUsers[1]
	portalUser2 := mockdata.MockPortalUsers[2]

	mustAddContestTeamMembers(t, repo, team1.ID, []uuid.UUID{user1.ID})
	mustAddContestTeamMembers(t, repo, team1.ID, []uuid.UUID{user2.ID})
	mustAddContestTeamMembers(t, repo, team2.ID, []uuid.UUID{user2.ID})

	expected1 := []*domain.User{
		{
			ID:       user1.ID,
			Name:     user1.Name,
			RealName: portalUser1.RealName,
		},
		{
			ID:       user2.ID,
			Name:     user2.Name,
			RealName: portalUser2.RealName,
		},
	}
	users1, err := repo.GetContestTeamMembers(contest1.ID, team1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected1, users1)

	expected2 := []*domain.User{
		{
			ID:       user2.ID,
			Name:     user2.Name,
			RealName: portalUser2.RealName,
		},
	}
	users2, err := repo.GetContestTeamMembers(contest1.ID, team2.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected2, users2)

	expected3 := []*domain.User{}
	users3, err := repo.GetContestTeamMembers(contest2.ID, team3.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected3, users3)

	expected4 := []*domain.User{
		{
			ID:       user2.ID,
			Name:     user2.Name,
			RealName: portalUser2.RealName,
		},
	}
	err = repo.EditContestTeamMembers(team1.ID, []uuid.UUID{user2.ID})
	assert.NoError(t, err)
	users4, err := repo.GetContestTeamMembers(contest1.ID, team1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected4, users4)

	expected5 := []*domain.User{}
	err = repo.EditContestTeamMembers(team1.ID, []uuid.UUID{})
	assert.NoError(t, err)
	users5, err := repo.GetContestTeamMembers(contest1.ID, team1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected5, users5)
}
