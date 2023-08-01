package repository

import (
	"context"
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

	h := testutils.SetupSQLHandler(t, sqlConf)
	repo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())
	contest1 := mustMakeContest(t, repo, nil)
	contest2 := mustMakeContest(t, repo, nil)
	contests, err := repo.GetContests(context.Background())
	assert.NoError(t, err)

	expected := []*domain.Contest{&contest1.Contest, &contest2.Contest}

	assert.ElementsMatch(t, expected, contests)
}

func TestContestRepository_GetContest(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("contest_repository_get_contest")
	sqlConf := conf.SQLConf()

	h := testutils.SetupSQLHandler(t, sqlConf)
	repo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())
	contest1 := mustMakeContest(t, repo, nil)
	contest2 := mustMakeContest(t, repo, nil)

	gotContest1, err := repo.GetContest(context.Background(), contest1.ID)
	assert.NoError(t, err)
	assert.Equal(t, contest1, gotContest1)

	gotContest2, err := repo.GetContest(context.Background(), contest2.ID)
	assert.NoError(t, err)
	assert.Equal(t, contest2, gotContest2)
}

func TestContestRepository_UpdateContest(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("contest_repository_update_contest")
	sqlConf := conf.SQLConf()

	h := testutils.SetupSQLHandler(t, sqlConf)
	repo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())
	tests := []struct {
		name string
		ctx  context.Context
		args *urepository.UpdateContestArgs
	}{
		{
			name: "all fields",
			ctx:  context.Background(),
			args: random.UpdateContestArgs(),
		},
		{
			name: "partial fields",
			ctx:  context.Background(),
			args: random.UpdateContestArgs(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			contest1 := mustMakeContest(t, repo, nil)
			contest2 := mustMakeContest(t, repo, nil)

			args := tt.args
			contest1.Name = args.Name.ValueOr(contest1.Name)
			contest1.Description = args.Description.ValueOr(contest1.Description)
			contest1.Link = args.Link.ValueOr(contest1.Link)
			contest1.TimeStart = args.Since.ValueOr(contest1.TimeStart)
			contest1.TimeEnd = args.Until.ValueOr(contest1.TimeEnd)

			err := repo.UpdateContest(context.Background(), contest1.ID, args)
			assert.NoError(t, err)

			gotContest1, err := repo.GetContest(context.Background(), contest1.ID)
			assert.NoError(t, err)
			gotContest2, err := repo.GetContest(context.Background(), contest2.ID)
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
		})
	}
}

func TestContestRepository_DeleteContest(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("contest_repository_delete_contest")
	sqlConf := conf.SQLConf()

	h := testutils.SetupSQLHandler(t, sqlConf)
	repo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())
	contest1 := mustMakeContest(t, repo, nil)
	contest2 := mustMakeContest(t, repo, nil)

	gotContest1, err := repo.GetContest(context.Background(), contest1.ID)
	assert.NoError(t, err)
	assert.Equal(t, contest1, gotContest1)

	gotContest2, err := repo.GetContest(context.Background(), contest2.ID)
	assert.NoError(t, err)
	assert.Equal(t, contest2, gotContest2)

	err = repo.DeleteContest(context.Background(), contest1.ID)
	assert.NoError(t, err)

	deletedContest1, err := repo.GetContest(context.Background(), contest1.ID)
	assert.Nil(t, deletedContest1)
	assert.Equal(t, err, urepository.ErrNotFound)

	gotContest2, err = repo.GetContest(context.Background(), contest2.ID)
	assert.NoError(t, err)
	assert.Equal(t, contest2, gotContest2)
}

func TestContestRepository_GetContestTeams(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("contest_repository_get_contest_teams")
	sqlConf := conf.SQLConf()

	h := testutils.SetupSQLHandler(t, sqlConf)
	repo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())

	contest1 := mustMakeContest(t, repo, nil)
	contest2 := mustMakeContest(t, repo, nil)
	team1 := mustMakeContestTeam(t, repo, contest1.ID, &urepository.CreateContestTeamArgs{
		Name:        random.AlphaNumeric(),
		Result:      random.Optional(random.AlphaNumeric()),
		Link:        random.Optional(random.RandURLString()),
		Description: random.AlphaNumeric(),
	})
	team2 := mustMakeContestTeam(t, repo, contest1.ID, &urepository.CreateContestTeamArgs{
		Name:        random.AlphaNumeric(),
		Result:      random.Optional(random.AlphaNumeric()),
		Link:        random.Optional(random.RandURLString()),
		Description: random.AlphaNumeric(),
	})

	expected1 := []*domain.ContestTeam{&team1.ContestTeam, &team2.ContestTeam}
	gotTeams1, err := repo.GetContestTeams(context.Background(), contest1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected1, gotTeams1)

	expected2 := []*domain.ContestTeam{}
	gotTeams2, err := repo.GetContestTeams(context.Background(), contest2.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected2, gotTeams2)
}

func TestContestRepository_GetContestTeam(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("contest_repository_get_contest_team")
	sqlConf := conf.SQLConf()
	h := testutils.SetupSQLHandler(t, sqlConf)
	repo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())

	contest1 := mustMakeContest(t, repo, nil)
	contest2 := mustMakeContest(t, repo, nil)
	team1 := mustMakeContestTeam(t, repo, contest1.ID, &urepository.CreateContestTeamArgs{
		Name:        random.AlphaNumeric(),
		Result:      random.Optional(random.AlphaNumeric()),
		Link:        random.Optional(random.RandURLString()),
		Description: random.AlphaNumeric(),
	})
	team2 := mustMakeContestTeam(t, repo, contest1.ID, &urepository.CreateContestTeamArgs{
		Name:        random.AlphaNumeric(),
		Result:      random.Optional(random.AlphaNumeric()),
		Link:        random.Optional(random.RandURLString()),
		Description: random.AlphaNumeric(),
	})

	gotTeams1, err := repo.GetContestTeam(context.Background(), contest1.ID, team1.ID)
	assert.NoError(t, err)
	assert.Equal(t, team1, gotTeams1)

	gotTeams2, err := repo.GetContestTeam(context.Background(), contest2.ID, team2.ID)
	assert.Error(t, err)
	assert.Nil(t, gotTeams2)
}

// func TestCreateContestTeam(t *testing.T) {
// }

func TestContestRepository_UpdateContestTeam(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("contest_repository_update_contest_teams")
	sqlConf := conf.SQLConf()
	h := testutils.SetupSQLHandler(t, sqlConf)
	repo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())

	tests := []struct {
		name string
		ctx  context.Context
		args *urepository.UpdateContestTeamArgs
	}{
		{
			name: "all fields",
			ctx:  context.Background(),
			args: random.UpdateContestTeamArgs(),
		},
		{
			name: "partial fields",
			ctx:  context.Background(),
			args: random.OptUpdateContestTeamArgs(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			contest := mustMakeContest(t, repo, nil)
			team := mustMakeContestTeam(t, repo, contest.ID, &urepository.CreateContestTeamArgs{
				Name:        random.AlphaNumeric(),
				Result:      random.Optional(random.AlphaNumeric()),
				Link:        random.Optional(random.RandURLString()),
				Description: random.AlphaNumeric(),
			})

			args := tt.args
			team.Name = args.Name.ValueOr(team.Name)
			team.Result = args.Result.ValueOr(team.Result)
			team.Link = args.Link.ValueOr(team.Link)
			team.Description = args.Description.ValueOr(team.Description)

			err := repo.UpdateContestTeam(tt.ctx, team.ID, args)
			assert.NoError(t, err)

			got, err := repo.GetContestTeam(tt.ctx, contest.ID, team.ID)
			assert.NoError(t, err)
			assert.Equal(t, team, got)
		})
	}
}

func TestContestRepository_DeleteContestTeam(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("contest_repository_delete_contest_teams")
	sqlConf := conf.SQLConf()

	h := testutils.SetupSQLHandler(t, sqlConf)
	repo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())

	contest1 := mustMakeContest(t, repo, nil)
	contest2 := mustMakeContest(t, repo, nil)
	team1 := mustMakeContestTeam(t, repo, contest1.ID, &urepository.CreateContestTeamArgs{
		Name:        random.AlphaNumeric(),
		Result:      random.Optional(random.AlphaNumeric()),
		Link:        random.Optional(random.RandURLString()),
		Description: random.AlphaNumeric(),
	})
	team2 := mustMakeContestTeam(t, repo, contest1.ID, &urepository.CreateContestTeamArgs{
		Name:        random.AlphaNumeric(),
		Result:      random.Optional(random.AlphaNumeric()),
		Link:        random.Optional(random.RandURLString()),
		Description: random.AlphaNumeric(),
	})

	expected1 := []*domain.ContestTeam{&team1.ContestTeam, &team2.ContestTeam}
	gotTeams1, err := repo.GetContestTeams(context.Background(), contest1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected1, gotTeams1)

	expected2 := []*domain.ContestTeam{}
	gotTeams2, err := repo.GetContestTeams(context.Background(), contest2.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected2, gotTeams2)

	err = repo.DeleteContestTeam(context.Background(), contest1.ID, team1.ID)
	assert.NoError(t, err)
	expected3 := []*domain.ContestTeam{&team2.ContestTeam}
	gotTeams3, err := repo.GetContestTeams(context.Background(), contest1.ID)
	assert.NoError(t, err)
	assert.Equal(t, expected3, gotTeams3)

	err = repo.DeleteContestTeam(context.Background(), contest1.ID, team2.ID)
	assert.NoError(t, err)
	expected4 := []*domain.ContestTeam{}
	gotTeams4, err := repo.GetContestTeams(context.Background(), contest1.ID)
	assert.NoError(t, err)
	assert.Equal(t, expected4, gotTeams4)
}

func TestContestRepository_GetContestTeamMembers(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("contest_repository_get_contest_team_members")
	sqlConf := conf.SQLConf()
	h := testutils.SetupSQLHandler(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
	repo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())

	contest1 := mustMakeContest(t, repo, nil)
	contest2 := mustMakeContest(t, repo, nil)
	team1 := mustMakeContestTeam(t, repo, contest1.ID, &urepository.CreateContestTeamArgs{
		Name:        random.AlphaNumeric(),
		Result:      random.Optional(random.AlphaNumeric()),
		Link:        random.Optional(random.RandURLString()),
		Description: random.AlphaNumeric(),
	})
	team2 := mustMakeContestTeam(t, repo, contest1.ID, &urepository.CreateContestTeamArgs{
		Name:        random.AlphaNumeric(),
		Result:      random.Optional(random.AlphaNumeric()),
		Link:        random.Optional(random.RandURLString()),
		Description: random.AlphaNumeric(),
	})
	team3 := mustMakeContestTeam(t, repo, contest1.ID, &urepository.CreateContestTeamArgs{
		Name:        random.AlphaNumeric(),
		Result:      random.Optional(random.AlphaNumeric()),
		Link:        random.Optional(random.RandURLString()),
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
		domain.NewUser(user1.ID, user1.Name, portalUser1.RealName, user1.Check),
		domain.NewUser(user2.ID, user2.Name, portalUser2.RealName, user2.Check),
	}
	users1, err := repo.GetContestTeamMembers(context.Background(), contest1.ID, team1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected1, users1)

	expected2 := []*domain.User{
		domain.NewUser(user2.ID, user2.Name, portalUser2.RealName, user2.Check),
	}
	users2, err := repo.GetContestTeamMembers(context.Background(), contest1.ID, team2.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected2, users2)

	expected3 := []*domain.User{}
	users3, err := repo.GetContestTeamMembers(context.Background(), contest2.ID, team3.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected3, users3)
}

// func TestContestRepository_AddContestTeamMembers(t *testing.T) {
// }

func TestContestRepository_EditContestTeamMembers(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("contest_repository_edit_contest_team_members")
	sqlConf := conf.SQLConf()
	h := testutils.SetupSQLHandler(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
	repo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())

	contest1 := mustMakeContest(t, repo, nil)
	contest2 := mustMakeContest(t, repo, nil)
	team1 := mustMakeContestTeam(t, repo, contest1.ID, &urepository.CreateContestTeamArgs{
		Name:        random.AlphaNumeric(),
		Result:      random.Optional(random.AlphaNumeric()),
		Link:        random.Optional(random.RandURLString()),
		Description: random.AlphaNumeric(),
	})
	team2 := mustMakeContestTeam(t, repo, contest1.ID, &urepository.CreateContestTeamArgs{
		Name:        random.AlphaNumeric(),
		Result:      random.Optional(random.AlphaNumeric()),
		Link:        random.Optional(random.RandURLString()),
		Description: random.AlphaNumeric(),
	})
	team3 := mustMakeContestTeam(t, repo, contest1.ID, &urepository.CreateContestTeamArgs{
		Name:        random.AlphaNumeric(),
		Result:      random.Optional(random.AlphaNumeric()),
		Link:        random.Optional(random.RandURLString()),
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
		domain.NewUser(user1.ID, user1.Name, portalUser1.RealName, user1.Check),
		domain.NewUser(user2.ID, user2.Name, portalUser2.RealName, user2.Check),
	}
	users1, err := repo.GetContestTeamMembers(context.Background(), contest1.ID, team1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected1, users1)

	expected2 := []*domain.User{
		domain.NewUser(user2.ID, user2.Name, portalUser2.RealName, user2.Check),
	}
	users2, err := repo.GetContestTeamMembers(context.Background(), contest1.ID, team2.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected2, users2)

	expected3 := []*domain.User{}
	users3, err := repo.GetContestTeamMembers(context.Background(), contest2.ID, team3.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected3, users3)

	expected4 := []*domain.User{
		domain.NewUser(user2.ID, user2.Name, portalUser2.RealName, user2.Check),
	}
	err = repo.EditContestTeamMembers(context.Background(), team1.ID, []uuid.UUID{user2.ID})
	assert.NoError(t, err)
	users4, err := repo.GetContestTeamMembers(context.Background(), contest1.ID, team1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected4, users4)

	expected5 := []*domain.User{}
	err = repo.EditContestTeamMembers(context.Background(), team1.ID, []uuid.UUID{})
	assert.NoError(t, err)
	users5, err := repo.GetContestTeamMembers(context.Background(), contest1.ID, team1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected5, users5)
}
