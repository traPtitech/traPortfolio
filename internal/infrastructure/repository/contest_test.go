package repository

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/internal/domain"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/external"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/external/mock_external"
	"go.uber.org/mock/gomock"

	"github.com/traPtitech/traPortfolio/internal/pkgs/optional"
	"github.com/traPtitech/traPortfolio/internal/pkgs/random"
	"github.com/traPtitech/traPortfolio/internal/usecases/repository"
)

func Test_GetContests(t *testing.T) {
	t.Parallel()

	db := SetupTestGormDB(t)
	portalAPI := mock_external.NewMockPortalAPI(gomock.NewController(t))
	repo := NewContestRepository(db, portalAPI)

	contest1, err := repo.CreateContest(context.Background(), random.CreateContestArgs())
	assert.NoError(t, err)
	contest2, err := repo.CreateContest(context.Background(), random.CreateContestArgs())
	assert.NoError(t, err)

	t.Run("get all contests", func(t *testing.T) {
		gotContests, err := repo.GetContests(context.Background())
		assert.NoError(t, err)

		expectedContests := []*domain.Contest{&contest1.Contest, &contest2.Contest}
		assert.ElementsMatch(t, expectedContests, gotContests)
	})
}

func Test_GetContest(t *testing.T) {
	t.Parallel()

	db := SetupTestGormDB(t)
	portalAPI := mock_external.NewMockPortalAPI(gomock.NewController(t))
	repo := NewContestRepository(db, portalAPI)

	contest1, err := repo.CreateContest(context.Background(), random.CreateContestArgs())
	assert.NoError(t, err)
	contest2, err := repo.CreateContest(context.Background(), random.CreateContestArgs())
	assert.NoError(t, err)

	t.Run("get contest1", func(t *testing.T) {
		gotContest, err := repo.GetContest(context.Background(), contest1.ID)
		assert.NoError(t, err)
		assert.Equal(t, contest1, gotContest)
	})

	t.Run("get contest2", func(t *testing.T) {
		gotContest, err := repo.GetContest(context.Background(), contest2.ID)
		assert.NoError(t, err)
		assert.Equal(t, contest2, gotContest)
	})
}

func Test_CreateContest(t *testing.T) {}

func Test_UpdateContest(t *testing.T) {
	t.Parallel()

	db := SetupTestGormDB(t)
	portalAPI := mock_external.NewMockPortalAPI(gomock.NewController(t))
	repo := NewContestRepository(db, portalAPI)

	contest, err := repo.CreateContest(context.Background(), random.CreateContestArgs())
	assert.NoError(t, err)

	t.Run("update all fields", func(t *testing.T) {
		args := random.UpdateContestArgs()
		err := repo.UpdateContest(context.Background(), contest.ID, args)
		assert.NoError(t, err)

		gotContest, err := repo.GetContest(context.Background(), contest.ID)
		assert.NoError(t, err)

		contest.Name = args.Name.ValueOr(contest.Name)
		contest.Description = args.Description.ValueOr(contest.Description)
		contest.Link = args.Link.ValueOr(contest.Link)
		contest.TimeStart = args.Since.ValueOr(contest.TimeStart)
		contest.TimeEnd = args.Until.ValueOr(contest.TimeEnd)
		assert.Equal(t, contest, gotContest)
	})

	t.Run("update no fields", func(t *testing.T) {
		args := &repository.UpdateContestArgs{}
		args.Until = optional.New(contest.TimeEnd, contest.TimeEnd == time.Time{})
		err := repo.UpdateContest(context.Background(), contest.ID, args)
		assert.NoError(t, err)

		gotContest, err := repo.GetContest(context.Background(), contest.ID)
		assert.NoError(t, err)

		assert.Equal(t, contest, gotContest)
	})

	t.Run("update until to nil", func(t *testing.T) {
		argWithUntil := random.CreateContestArgs()
		argWithUntil.Until = random.UpdateContestArgs().Since
		contest, err := repo.CreateContest(context.Background(), argWithUntil)
		assert.NoError(t, err)

		argWithoutUntil := random.UpdateContestArgs()
		argWithoutUntil.Until = optional.Of[time.Time]{}
		err = repo.UpdateContest(context.Background(), contest.ID, argWithoutUntil)
		assert.NoError(t, err)
	})
}

func Test_DeleteContest(t *testing.T) {
	t.Parallel()

	db := SetupTestGormDB(t)
	portalAPI := mock_external.NewMockPortalAPI(gomock.NewController(t))
	repo := NewContestRepository(db, portalAPI)

	contest1, err := repo.CreateContest(context.Background(), random.CreateContestArgs())
	assert.NoError(t, err)
	contest2, err := repo.CreateContest(context.Background(), random.CreateContestArgs())
	assert.NoError(t, err)

	t.Run("delete contest1", func(t *testing.T) {
		gotContest1, err := repo.GetContest(context.Background(), contest1.ID)
		assert.NoError(t, err)
		assert.Equal(t, contest1, gotContest1)

		err = repo.DeleteContest(context.Background(), contest1.ID)
		assert.NoError(t, err)

		deletedContest1, err := repo.GetContest(context.Background(), contest1.ID)
		assert.Nil(t, deletedContest1)
		assert.Equal(t, err, repository.ErrNotFound)
	})

	t.Run("delete of contest1 doesn't affect contest2", func(t *testing.T) {
		gotContest2, err := repo.GetContest(context.Background(), contest2.ID)
		assert.NoError(t, err)
		assert.Equal(t, contest2, gotContest2)
	})
}

func Test_GetContestTeams(t *testing.T) {
	t.Parallel()

	db := SetupTestGormDB(t)
	portalAPI := mock_external.NewMockPortalAPI(gomock.NewController(t))
	repo := NewContestRepository(db, portalAPI)

	// contest1 has two teams (team1, team2)
	contest1, err := repo.CreateContest(context.Background(), random.CreateContestArgs())
	assert.NoError(t, err)
	team1, err := repo.CreateContestTeam(context.Background(), contest1.ID, random.CreateContestTeamArgs())
	assert.NoError(t, err)
	team2, err := repo.CreateContestTeam(context.Background(), contest1.ID, random.CreateContestTeamArgs())
	assert.NoError(t, err)
	// contest2 has no teams
	contest2, err := repo.CreateContest(context.Background(), random.CreateContestArgs())
	assert.NoError(t, err)

	t.Run("get teams of contest1 (two teams belongs)", func(t *testing.T) {
		expectedTeams := []*domain.ContestTeam{&team1.ContestTeam, &team2.ContestTeam}
		portalAPI.EXPECT().GetUsers().Return([]*external.PortalUserResponse{}, nil)
		gotTeams, err := repo.GetContestTeams(context.Background(), contest1.ID)
		assert.NoError(t, err)
		assert.ElementsMatch(t, expectedTeams, gotTeams)
	})

	t.Run("get teams of contest2 (no teams belongs)", func(t *testing.T) {
		expectedTeams := []*domain.ContestTeam{}
		portalAPI.EXPECT().GetUsers().Return([]*external.PortalUserResponse{}, nil)
		gotTeams, err := repo.GetContestTeams(context.Background(), contest2.ID)
		assert.NoError(t, err)
		assert.ElementsMatch(t, expectedTeams, gotTeams)
	})
}

func Test_GetContestTeam(t *testing.T) {
	t.Parallel()

	db := SetupTestGormDB(t)
	portalAPI := mock_external.NewMockPortalAPI(gomock.NewController(t))
	repo := NewContestRepository(db, portalAPI)

	// contest1 has a team (team1)
	contest1, err := repo.CreateContest(context.Background(), random.CreateContestArgs())
	assert.NoError(t, err)
	team1, err := repo.CreateContestTeam(context.Background(), contest1.ID, random.CreateContestTeamArgs())
	assert.NoError(t, err)
	// contest2 has no teams
	contest2, err := repo.CreateContest(context.Background(), random.CreateContestArgs())
	assert.NoError(t, err)

	t.Run("get team1 (belongs to contest1)", func(t *testing.T) {
		portalAPI.EXPECT().GetUsers().Return([]*external.PortalUserResponse{}, nil)
		gotTeam1, err := repo.GetContestTeam(context.Background(), contest1.ID, team1.ID)
		assert.NoError(t, err)
		assert.Equal(t, team1, gotTeam1)
	})

	t.Run("cannot get team1 (doesn't belong to contest2)", func(t *testing.T) {
		_, err := repo.GetContestTeam(context.Background(), contest2.ID, team1.ID)
		assert.Error(t, err)
		assert.Equal(t, err, repository.ErrNotFound)
	})
}

func Test_CreateContestTeam(t *testing.T) {}

func Test_UpdateContestTeam(t *testing.T) {
	t.Parallel()

	db := SetupTestGormDB(t)
	portalAPI := mock_external.NewMockPortalAPI(gomock.NewController(t))
	repo := NewContestRepository(db, portalAPI)

	contest, err := repo.CreateContest(context.Background(), random.CreateContestArgs())
	assert.NoError(t, err)
	team, err := repo.CreateContestTeam(context.Background(), contest.ID, random.CreateContestTeamArgs())
	assert.NoError(t, err)

	t.Run("update all fields", func(t *testing.T) {
		args := random.UpdateContestTeamArgs()
		err := repo.UpdateContestTeam(context.Background(), team.ID, args)
		assert.NoError(t, err)

		portalAPI.EXPECT().GetUsers().Return([]*external.PortalUserResponse{}, nil)
		gotTeam, err := repo.GetContestTeam(context.Background(), contest.ID, team.ID)
		assert.NoError(t, err)

		team.Name = args.Name.ValueOr(team.Name)
		team.Result = args.Result.ValueOr(team.Result)
		team.Link = args.Link.ValueOr(team.Link)
		team.Description = args.Description.ValueOr(team.Description)
		assert.Equal(t, team, gotTeam)
	})

	t.Run("update no fields", func(t *testing.T) {
		args := &repository.UpdateContestTeamArgs{}
		err := repo.UpdateContestTeam(context.Background(), team.ID, args)
		assert.NoError(t, err)

		portalAPI.EXPECT().GetUsers().Return([]*external.PortalUserResponse{}, nil)
		gotTeam, err := repo.GetContestTeam(context.Background(), contest.ID, team.ID)
		assert.NoError(t, err)

		assert.Equal(t, team, gotTeam)
	})
}

func Test_DeleteContestTeam(t *testing.T) {
	t.Parallel()

	db := SetupTestGormDB(t)
	portalAPI := mock_external.NewMockPortalAPI(gomock.NewController(t))
	repo := NewContestRepository(db, portalAPI)

	contest1, err := repo.CreateContest(context.Background(), random.CreateContestArgs())
	assert.NoError(t, err)
	team1, err := repo.CreateContestTeam(context.Background(), contest1.ID, random.CreateContestTeamArgs())
	assert.NoError(t, err)
	team2, err := repo.CreateContestTeam(context.Background(), contest1.ID, random.CreateContestTeamArgs())
	assert.NoError(t, err)
	contest2, err := repo.CreateContest(context.Background(), random.CreateContestArgs())
	assert.NoError(t, err)

	t.Run("delete team1 (belongs to contest1)", func(t *testing.T) {
		portalAPI.EXPECT().GetUsers().Return([]*external.PortalUserResponse{}, nil)
		gotTeam1, err := repo.GetContestTeam(context.Background(), contest1.ID, team1.ID)
		assert.NoError(t, err)
		assert.Equal(t, team1, gotTeam1)

		err = repo.DeleteContestTeam(context.Background(), contest1.ID, team1.ID)
		assert.NoError(t, err)

		deletedTeam1, err := repo.GetContestTeam(context.Background(), contest1.ID, team1.ID)
		assert.Nil(t, deletedTeam1)
		assert.Equal(t, err, repository.ErrNotFound)
	})

	t.Run("delete of team1 doesn't affect team2", func(t *testing.T) {
		portalAPI.EXPECT().GetUsers().Return([]*external.PortalUserResponse{}, nil)
		gotTeam2, err := repo.GetContestTeam(context.Background(), contest1.ID, team2.ID)
		assert.NoError(t, err)
		assert.Equal(t, team2, gotTeam2)
	})

	t.Run("cannot delete team1 (doesn't belong to contest2)", func(t *testing.T) {
		err := repo.DeleteContestTeam(context.Background(), contest2.ID, team1.ID)
		assert.Error(t, err)
		assert.Equal(t, err, repository.ErrNotFound)
	})
}

func Test_GetContestTeamMembers(t *testing.T) {
	t.Parallel()

	db := SetupTestGormDB(t)
	portalAPI := mock_external.NewMockPortalAPI(gomock.NewController(t))
	repo := NewContestRepository(db, portalAPI)
	traqAPI := mock_external.NewMockTraQAPI(gomock.NewController(t))
	userRepo := NewUserRepository(db, portalAPI, traqAPI)

	// contest1 has a team (team1)
	// team1 has two members
	contest1, err := repo.CreateContest(context.Background(), random.CreateContestArgs())
	assert.NoError(t, err)

	team1, err := repo.CreateContestTeam(context.Background(), contest1.ID, random.CreateContestTeamArgs())
	assert.NoError(t, err)

	traqUsers := []*external.TraQUserResponse{
		{
			ID:    random.UUID(),
			Name:  "user1",
			State: domain.TraqStateActive,
			Bot:   false,
		},
		{
			ID:    random.UUID(),
			Name:  "user2",
			State: domain.TraqStateActive,
			Bot:   false,
		},
	}
	traqAPI.EXPECT().GetUsers(&external.TraQGetAllArgs{IncludeSuspended: true}).Return(traqUsers, nil)
	err = userRepo.SyncUsers(context.Background())
	assert.NoError(t, err)

	memberIDs := []uuid.UUID{traqUsers[0].ID, traqUsers[1].ID}
	err = repo.EditContestTeamMembers(context.Background(), team1.ID, memberIDs)
	assert.NoError(t, err)

	t.Run("get team1 members", func(t *testing.T) {
		portalUsers := lo.Map(traqUsers, func(u *external.TraQUserResponse, _ int) *external.PortalUserResponse {
			return &external.PortalUserResponse{
				TraQID:         u.Name,
				RealName:       u.Name + " real",
				AlphabeticName: u.Name + " alphabetic",
			}
		})
		portalAPI.EXPECT().GetUsers().Return(portalUsers, nil)
		expectedMembers := []*domain.User{
			domain.NewUser(traqUsers[0].ID, traqUsers[0].Name, portalUsers[0].RealName, false),
			domain.NewUser(traqUsers[1].ID, traqUsers[1].Name, portalUsers[1].RealName, false),
		}
		gotMembers, err := repo.GetContestTeamMembers(context.Background(), contest1.ID, team1.ID)
		assert.NoError(t, err)
		assert.ElementsMatch(t, expectedMembers, gotMembers)
	})
}

func Test_EditContestTeamMembers(t *testing.T) {
	t.Parallel()

	db := SetupTestGormDB(t)
	portalAPI := mock_external.NewMockPortalAPI(gomock.NewController(t))
	repo := NewContestRepository(db, portalAPI)
	traqAPI := mock_external.NewMockTraQAPI(gomock.NewController(t))
	userRepo := NewUserRepository(db, portalAPI, traqAPI)

	contest1, err := repo.CreateContest(context.Background(), random.CreateContestArgs())
	assert.NoError(t, err)

	team1, err := repo.CreateContestTeam(context.Background(), contest1.ID, random.CreateContestTeamArgs())
	assert.NoError(t, err)

	traqUsers := []*external.TraQUserResponse{
		{
			ID:    random.UUID(),
			Name:  "user1",
			State: domain.TraqStateActive,
			Bot:   false,
		},
		{
			ID:    random.UUID(),
			Name:  "user2",
			State: domain.TraqStateActive,
			Bot:   false,
		},
	}
	traqAPI.EXPECT().GetUsers(&external.TraQGetAllArgs{IncludeSuspended: true}).Return(traqUsers, nil)
	err = userRepo.SyncUsers(context.Background())
	assert.NoError(t, err)

	portalUsers := lo.Map(traqUsers, func(u *external.TraQUserResponse, _ int) *external.PortalUserResponse {
		return &external.PortalUserResponse{
			TraQID:         u.Name,
			RealName:       u.Name + " real",
			AlphabeticName: u.Name + " alphabetic",
		}
	})
	portalAPI.EXPECT().GetUsers().Return(portalUsers, nil)
	users, err := userRepo.GetUsers(context.Background(), &repository.GetUsersArgs{})
	assert.NoError(t, err)

	t.Run("add a user to team1", func(t *testing.T) {
		memberIDs := []uuid.UUID{users[0].ID}
		err = repo.EditContestTeamMembers(context.Background(), team1.ID, memberIDs)
		assert.NoError(t, err)

		expectedMembers := []*domain.User{users[0]}
		portalAPI.EXPECT().GetUsers().Return(portalUsers, nil)
		gotMembers, err := repo.GetContestTeamMembers(context.Background(), contest1.ID, team1.ID)
		assert.NoError(t, err)
		assert.ElementsMatch(t, expectedMembers, gotMembers)
	})

	t.Run("add another user to team1", func(t *testing.T) {
		memberIDs := []uuid.UUID{users[0].ID, users[1].ID}
		err = repo.EditContestTeamMembers(context.Background(), team1.ID, memberIDs)
		assert.NoError(t, err)

		expectedMembers := []*domain.User{users[0], users[1]}
		portalAPI.EXPECT().GetUsers().Return(portalUsers, nil)
		gotMembers, err := repo.GetContestTeamMembers(context.Background(), contest1.ID, team1.ID)
		assert.NoError(t, err)
		assert.ElementsMatch(t, expectedMembers, gotMembers)
	})

	t.Run("delete all users from team1", func(t *testing.T) {
		memberIDs := []uuid.UUID{}
		err = repo.EditContestTeamMembers(context.Background(), team1.ID, memberIDs)
		assert.NoError(t, err)

		expectedMembers := []*domain.User{}
		portalAPI.EXPECT().GetUsers().Return(portalUsers, nil)
		gotMembers, err := repo.GetContestTeamMembers(context.Background(), contest1.ID, team1.ID)
		assert.NoError(t, err)
		assert.ElementsMatch(t, expectedMembers, gotMembers)
	})
}
