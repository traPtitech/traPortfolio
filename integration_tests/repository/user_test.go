package repository

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external_e2e"
	irepository "github.com/traPtitech/traPortfolio/interfaces/repository"

	urepository "github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/mockdata"
	"github.com/traPtitech/traPortfolio/util/optional"
	"github.com/traPtitech/traPortfolio/util/random"
)

func TestUserRepository_GetUsers(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("user_repository_get_users")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
	repo := irepository.NewUserRepository(h, mock_external_e2e.NewMockPortalAPI(), mock_external_e2e.NewMockTraQAPI())

	type args struct {
		args *urepository.GetUsersArgs
	}
	cases := []struct {
		name      string
		args      args
		expected  []*domain.User
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "All NotIncludeSuspended",
			args: args{args: &urepository.GetUsersArgs{}},
			expected: []*domain.User{
				{
					ID:       mockdata.MockTraQUsers[0].User.ID,
					Name:     mockdata.MockTraQUsers[0].Name,
					RealName: mockdata.MockPortalUsers[0].RealName,
				},
				{
					ID:       mockdata.MockTraQUsers[2].User.ID,
					Name:     mockdata.MockTraQUsers[2].Name,
					RealName: mockdata.MockPortalUsers[2].RealName,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "All IncludeSuspended",
			args: args{args: &urepository.GetUsersArgs{
				IncludeSuspended: optional.NewBool(true, true),
			}},
			expected: []*domain.User{
				{
					ID:       mockdata.MockTraQUsers[0].User.ID,
					Name:     mockdata.MockTraQUsers[0].Name,
					RealName: mockdata.MockPortalUsers[0].RealName,
				},
				{
					ID:       mockdata.MockTraQUsers[1].User.ID,
					Name:     mockdata.MockTraQUsers[1].Name,
					RealName: mockdata.MockPortalUsers[1].RealName,
				},
				{
					ID:       mockdata.MockTraQUsers[2].User.ID,
					Name:     mockdata.MockTraQUsers[2].Name,
					RealName: mockdata.MockPortalUsers[2].RealName,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Name",
			args: args{args: &urepository.GetUsersArgs{
				Name: optional.NewString(mockdata.MockTraQUsers[0].Name, true),
			}},
			expected: []*domain.User{
				{
					ID:       mockdata.MockTraQUsers[0].User.ID,
					Name:     mockdata.MockTraQUsers[0].Name,
					RealName: mockdata.MockPortalUsers[0].RealName,
				},
			},
			assertion: assert.NoError,
		},
		{
			name: "Invalid arg",
			args: args{args: &urepository.GetUsersArgs{
				Name:             optional.NewString(mockdata.MockTraQUsers[0].Name, true),
				IncludeSuspended: optional.NewBool(true, true),
			}},
			expected:  nil,
			assertion: assert.Error,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			users, err := repo.GetUsers(tc.args.args)
			tc.assertion(t, err)
			assert.ElementsMatch(t, tc.expected, users)
		})
	}
}

func TestUserRepository_GetUser(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("user_repository_get_user")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
	repo := irepository.NewUserRepository(h, mock_external_e2e.NewMockPortalAPI(), mock_external_e2e.NewMockTraQAPI())

	type args struct {
		userID uuid.UUID
	}

	cases := []struct {
		name      string
		args      args
		expected  *domain.UserDetail
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "No account",
			args: args{
				mockdata.MockTraQUsers[2].User.ID,
			},
			expected: &domain.UserDetail{
				User: domain.User{
					ID:       mockdata.MockTraQUsers[2].User.ID,
					Name:     mockdata.MockTraQUsers[2].Name,
					RealName: mockdata.MockPortalUsers[2].RealName,
				},
				State:    mockdata.MockTraQUsers[2].User.State,
				Bio:      mockdata.MockUsers[2].Description,
				Accounts: []*domain.Account{},
			},
			assertion: assert.NoError,
		},
		{
			name: "With account",
			args: args{
				mockdata.MockTraQUsers[0].User.ID,
			},
			expected: &domain.UserDetail{
				User: domain.User{
					ID:       mockdata.MockTraQUsers[0].User.ID,
					Name:     mockdata.MockTraQUsers[0].Name,
					RealName: mockdata.MockPortalUsers[0].RealName,
				},
				State: mockdata.MockTraQUsers[0].User.State,
				Bio:   mockdata.MockUsers[0].Description,
				Accounts: []*domain.Account{
					{
						ID:          mockdata.MockAccounts[0].ID,
						DisplayName: mockdata.MockAccounts[0].Name,
						Type:        mockdata.MockAccounts[0].Type,
						PrPermitted: mockdata.MockAccounts[0].Check,
						URL:         mockdata.MockAccounts[0].URL,
					},
				},
			},
			assertion: assert.NoError,
		},
		{
			name:     "Not exist",
			expected: nil,
			args: args{
				random.UUID(),
			},
			assertion: assert.Error,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			users, err := repo.GetUser(tc.args.userID)
			tc.assertion(t, err)
			assert.Equal(t, tc.expected, users)
		})
	}
}

func TestUserRepository_CreateUser(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("user_repository_create_user")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
	repo := irepository.NewUserRepository(h, mock_external_e2e.NewMockPortalAPI(), mock_external_e2e.NewMockTraQAPI())

	type args struct {
		args *urepository.CreateUserArgs
	}

	description := random.AlphaNumeric()
	cases := []struct {
		name      string
		args      args
		expected  *domain.UserDetail
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				args: &urepository.CreateUserArgs{
					Description: description,
					Check:       random.Bool(),
					Name:        mockdata.MockUsers[1].Name,
				},
			},
			expected: &domain.UserDetail{
				User: domain.User{
					// ID is replaced by generated one.
					Name:     mockdata.MockUsers[1].Name,
					RealName: mockdata.MockPortalUsers[1].RealName,
				},
				State:    mockdata.MockTraQUsers[1].User.State,
				Bio:      description,
				Accounts: []*domain.Account{},
			},
			assertion: assert.NoError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := repo.CreateUser(tc.args.args)
			tc.expected.ID = user.ID
			tc.assertion(t, err)
			assert.Equal(t, tc.expected, user)
		})
	}
}

func TestUserRepository_UpdateUser(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("user_repository_update_user")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
	repo := irepository.NewUserRepository(h, mock_external_e2e.NewMockPortalAPI(), mock_external_e2e.NewMockTraQAPI())

	idx := 1
	user := mockdata.MockUsers[idx]
	portalUser := mockdata.MockPortalUsers[idx]
	traqUser := mockdata.MockTraQUsers[idx]
	args := &urepository.UpdateUserArgs{
		Description: random.OptAlphaNumeric(),
		Check:       random.OptBool(),
	}

	err = repo.UpdateUser(user.ID, args)
	assert.NoError(t, err)

	var bio string
	if args.Description.Valid {
		bio = args.Description.String
	} else {
		bio = user.Description
	}

	expected := &domain.UserDetail{
		User: domain.User{
			ID:       user.ID,
			Name:     user.Name,
			RealName: portalUser.RealName,
		},
		State:    traqUser.User.State,
		Bio:      bio,
		Accounts: []*domain.Account{},
	}
	got, err := repo.GetUser(user.ID)
	assert.NoError(t, err)

	assert.Equal(t, expected, got)
}

func TestUserRepository_GetAccounts(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("user_repository_get_accounts")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
	repo := irepository.NewUserRepository(h, mock_external_e2e.NewMockPortalAPI(), mock_external_e2e.NewMockTraQAPI())

	idx := 1
	user := mockdata.MockUsers[idx]
	account1 := mustMakeAccount(t, repo, user.ID, nil)
	account2 := mustMakeAccount(t, repo, user.ID, nil)
	expected := []*domain.Account{account1, account2}

	got, err := repo.GetAccounts(user.ID)
	assert.NoError(t, err)

	assert.ElementsMatch(t, expected, got)
}

func TestUserRepository_GetAccount(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("user_repository_get_account")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
	repo := irepository.NewUserRepository(h, mock_external_e2e.NewMockPortalAPI(), mock_external_e2e.NewMockTraQAPI())

	idx := 1
	user := mockdata.MockUsers[idx]
	account1 := mustMakeAccount(t, repo, user.ID, nil)
	mustMakeAccount(t, repo, user.ID, nil)

	got, err := repo.GetAccount(user.ID, account1.ID)
	assert.NoError(t, err)
	assert.Equal(t, account1, got)
}

// func TestUserRepository_CreateAccount(t *testing.T) {
// }

func TestUserRepository_UpdateAccount(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("user_repository_update_account")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
	repo := irepository.NewUserRepository(h, mock_external_e2e.NewMockPortalAPI(), mock_external_e2e.NewMockTraQAPI())

	idx := 1
	user := mockdata.MockUsers[idx]
	account1 := mustMakeAccount(t, repo, user.ID, nil)
	mustMakeAccount(t, repo, user.ID, nil)

	args := &urepository.UpdateAccountArgs{
		DisplayName: random.OptAlphaNumeric(),
		Type:        random.OptInt64n(int64(domain.AccountLimit)),
		URL:         random.OptAlphaNumeric(),
		PrPermitted: random.OptBool(),
	}
	if args.DisplayName.Valid {
		account1.DisplayName = args.DisplayName.String
	}
	if args.Type.Valid {
		account1.Type = uint(args.Type.Int64)
	}
	if args.URL.Valid {
		account1.URL = args.URL.String
	}
	if args.PrPermitted.Valid {
		account1.PrPermitted = args.PrPermitted.Bool
	}
	err = repo.UpdateAccount(user.ID, account1.ID, args)
	assert.NoError(t, err)

	got, err := repo.GetAccount(user.ID, account1.ID)
	assert.NoError(t, err)
	assert.Equal(t, account1, got)
}

func TestUserRepository_DeleteAccount(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("user_repository_delete_account")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
	repo := irepository.NewUserRepository(h, mock_external_e2e.NewMockPortalAPI(), mock_external_e2e.NewMockTraQAPI())

	idx := 1
	user := mockdata.MockUsers[idx]
	account1 := mustMakeAccount(t, repo, user.ID, nil)
	account2 := mustMakeAccount(t, repo, user.ID, nil)

	err = repo.DeleteAccount(user.ID, account1.ID)
	assert.NoError(t, err)

	expected := []*domain.Account{account2}

	got, err := repo.GetAccounts(user.ID)
	assert.NoError(t, err)

	assert.ElementsMatch(t, expected, got)
}

func TestUserRepository_GetUserProjects(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("user_repository_get_user_projects")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
	userRepo := irepository.NewUserRepository(h, mock_external_e2e.NewMockPortalAPI(), mock_external_e2e.NewMockTraQAPI())
	projectRepo := irepository.NewProjectRepository(h, mock_external_e2e.NewMockPortalAPI())

	project1 := mustMakeProject(t, projectRepo, nil)
	project2 := mustMakeProject(t, projectRepo, nil)
	user1 := mockdata.MockUsers[2]

	expected1 := []*domain.UserWithDuration{}
	expected2 := []*domain.UserWithDuration{}
	users1, err := projectRepo.GetProjectMembers(project1.ID)
	assert.NoError(t, err)
	users2, err := projectRepo.GetProjectMembers(project2.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected1, users1)
	assert.ElementsMatch(t, expected2, users2)

	args1 := mustAddProjectMember(t, projectRepo, project1.ID, user1.ID, nil)
	args2 := mustAddProjectMember(t, projectRepo, project2.ID, user1.ID, nil)

	expected3 := []*domain.UserProject{newUserProject(args1, project1), newUserProject(args2, project2)}
	projects1, err := userRepo.GetProjects(user1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected3, projects1)
}

func TestUserRepository_GetContests(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("user_repository_get_contests")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
	userRepo := irepository.NewUserRepository(h, mock_external_e2e.NewMockPortalAPI(), mock_external_e2e.NewMockTraQAPI())
	contestRepo := irepository.NewContestRepository(h, mock_external_e2e.NewMockPortalAPI())

	contestNum := 3
	var contests []*domain.ContestDetail
	for i := 0; i < contestNum; i++ {
		contests = append(contests, mustMakeContest(t, contestRepo, nil))
	}

	contestTeamNum := 10
	var contestTeams []*domain.ContestTeamDetail
	for i := 0; i < contestTeamNum; i++ {
		contestTeams = append(contestTeams, mustMakeContestTeam(t, contestRepo, contests[i%contestNum].ID, nil))
	}

	team1 := contestTeams[0]
	team2 := contestTeams[contestNum-1]
	contest1 := contests[0]
	contest2 := contests[(contestNum-1)%contestNum]
	user1 := mockdata.MockUsers[1]
	user2 := mockdata.MockUsers[2]

	mustAddContestTeamMembers(t, contestRepo, team1.ID, []uuid.UUID{user1.ID, user2.ID})
	mustAddContestTeamMembers(t, contestRepo, team2.ID, []uuid.UUID{user1.ID})

	expected1 := []*domain.UserContest{newUserContest(&contest1.Contest, &team1.ContestTeam), newUserContest(&contest2.Contest, &team2.ContestTeam)}
	projects1, err := userRepo.GetContests(user1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected1, projects1)

	expected2 := []*domain.UserContest{newUserContest(&contest1.Contest, &team1.ContestTeam)}
	projects2, err := userRepo.GetContests(user2.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected2, projects2)
}

// func TestUserRepository_GetGroupsByUserID(t *testing.T) {
// }

func newUserProject(args *urepository.CreateProjectMemberArgs, project *domain.Project) *domain.UserProject {
	return &domain.UserProject{
		ID:       project.ID,
		Name:     project.Name,
		Duration: project.Duration,
		UserDuration: domain.YearWithSemesterDuration{
			Since: domain.YearWithSemester{
				Year:     args.SinceYear,
				Semester: args.SinceSemester,
			},
			Until: domain.YearWithSemester{
				Year:     args.UntilYear,
				Semester: args.UntilSemester,
			},
		},
	}
}

func newUserContest(contest *domain.Contest, team *domain.ContestTeam) *domain.UserContest {
	return &domain.UserContest{
		ID:          team.ID,
		Name:        team.Name,
		Result:      team.Result,
		ContestName: contest.Name,
	}
}
