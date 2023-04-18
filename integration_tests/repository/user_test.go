package repository

import (
	"context"
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
	h := testutils.SetupSQLHandler(t, sqlConf)
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
				domain.NewUser(
					mockdata.MockTraQUsers[0].User.ID,
					mockdata.MockTraQUsers[0].Name,
					mockdata.MockPortalUsers[0].RealName,
					mockdata.MockUsers[0].Check,
				),
				domain.NewUser(
					mockdata.MockTraQUsers[2].User.ID,
					mockdata.MockTraQUsers[2].Name,
					mockdata.MockPortalUsers[2].RealName,
					mockdata.MockUsers[2].Check,
				),
			},
			assertion: assert.NoError,
		},
		{
			name: "All IncludeSuspended",
			args: args{args: &urepository.GetUsersArgs{
				IncludeSuspended: optional.From(true),
			}},
			expected: []*domain.User{
				domain.NewUser(
					mockdata.MockUsers[0].ID,
					mockdata.MockUsers[0].Name,
					mockdata.MockPortalUsers[0].RealName,
					mockdata.MockUsers[0].Check,
				),
				domain.NewUser(
					mockdata.MockUsers[1].ID,
					mockdata.MockUsers[1].Name,
					mockdata.MockPortalUsers[1].RealName,
					mockdata.MockUsers[1].Check,
				),
				domain.NewUser(
					mockdata.MockUsers[2].ID,
					mockdata.MockUsers[2].Name,
					mockdata.MockPortalUsers[2].RealName,
					mockdata.MockUsers[2].Check,
				),
			},
			assertion: assert.NoError,
		},
		{
			name: "Name",
			args: args{args: &urepository.GetUsersArgs{
				Name: optional.From(mockdata.MockTraQUsers[0].Name),
			}},
			expected: []*domain.User{
				domain.NewUser(
					mockdata.MockUsers[0].ID,
					mockdata.MockUsers[0].Name,
					mockdata.MockPortalUsers[0].RealName,
					mockdata.MockUsers[0].Check,
				),
			},
			assertion: assert.NoError,
		},
		{
			name: "Invalid arg",
			args: args{args: &urepository.GetUsersArgs{
				Name:             optional.From(mockdata.MockTraQUsers[0].Name),
				IncludeSuspended: optional.From(true),
			}},
			expected:  nil,
			assertion: assert.Error,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			users, err := repo.GetUsers(context.Background(), tc.args.args)
			tc.assertion(t, err)
			assert.ElementsMatch(t, tc.expected, users)
		})
	}
}

func TestUserRepository_GetUser(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("user_repository_get_user")
	sqlConf := conf.SQLConf()
	h := testutils.SetupSQLHandler(t, sqlConf)
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
				User: *domain.NewUser(
					mockdata.MockUsers[2].ID,
					mockdata.MockUsers[2].Name,
					mockdata.MockPortalUsers[2].RealName,
					mockdata.MockUsers[2].Check,
				),
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
				User: *domain.NewUser(
					mockdata.MockUsers[0].ID,
					mockdata.MockUsers[0].Name,
					mockdata.MockPortalUsers[0].RealName,
					mockdata.MockUsers[0].Check,
				),
				State: mockdata.MockTraQUsers[0].User.State,
				Bio:   mockdata.MockUsers[0].Description,
				Accounts: []*domain.Account{
					{
						ID:          mockdata.MockAccounts[0].ID,
						DisplayName: mockdata.MockAccounts[0].Name,
						Type:        domain.AccountType(mockdata.MockAccounts[0].Type),
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
			users, err := repo.GetUser(context.Background(), tc.args.userID)
			tc.assertion(t, err)
			assert.Equal(t, tc.expected, users)
		})
	}
}

func TestUserRepository_CreateUser(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("user_repository_create_user")
	sqlConf := conf.SQLConf()
	h := testutils.SetupSQLHandler(t, sqlConf)
	repo := irepository.NewUserRepository(h, mock_external_e2e.NewMockPortalAPI(), mock_external_e2e.NewMockTraQAPI())

	type args struct {
		args *urepository.CreateUserArgs
	}

	check := random.Bool()
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
					Check:       check,
					Name:        mockdata.MockUsers[1].Name,
				},
			},
			expected: &domain.UserDetail{
				User: *domain.NewUser(
					uuid.Nil, // ID is replaced by generated one.
					mockdata.MockUsers[1].Name,
					mockdata.MockPortalUsers[1].RealName,
					check,
				),
				State:    mockdata.MockTraQUsers[1].User.State,
				Bio:      description,
				Accounts: []*domain.Account{},
			},
			assertion: assert.NoError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			user, err := repo.CreateUser(context.Background(), tc.args.args)
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
	h := testutils.SetupSQLHandler(t, sqlConf)
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

	err = repo.UpdateUser(context.Background(), user.ID, args)
	assert.NoError(t, err)

	var bio string
	if v, ok := args.Description.V(); ok {
		bio = v
	} else {
		bio = user.Description
	}

	var check bool
	if v, ok := args.Check.V(); ok {
		check = v
	} else {
		check = user.Check
	}

	expected := &domain.UserDetail{
		User: *domain.NewUser(
			user.ID,
			user.Name,
			portalUser.RealName,
			check,
		),
		State:    traqUser.User.State,
		Bio:      bio,
		Accounts: []*domain.Account{},
	}
	got, err := repo.GetUser(context.Background(), user.ID)
	assert.NoError(t, err)

	assert.Equal(t, expected, got)
}

func TestUserRepository_GetAccounts(t *testing.T) {
	t.Parallel()
	conf := testutils.GetConfigWithDBName("user_repository_get_accounts")
	sqlConf := conf.SQLConf()
	h := testutils.SetupSQLHandler(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
	repo := irepository.NewUserRepository(h, mock_external_e2e.NewMockPortalAPI(), mock_external_e2e.NewMockTraQAPI())

	var (
		idx          = 1
		user         = mockdata.MockUsers[idx]
		accountType1 = domain.AccountType(3)
		accountType2 = domain.AccountType(4)
	)

	account1 := mustMakeAccount(t, repo, user.ID, &urepository.CreateAccountArgs{
		DisplayName: random.AlphaNumeric(),
		Type:        accountType1,
		URL:         random.AccountURLString(accountType1),
		PrPermitted: random.Bool(),
	})
	account2 := mustMakeAccount(t, repo, user.ID, &urepository.CreateAccountArgs{
		DisplayName: random.AlphaNumeric(),
		Type:        accountType2,
		URL:         random.AccountURLString(accountType2),
		PrPermitted: random.Bool(),
	})
	expected := []*domain.Account{account1, account2}

	got, err := repo.GetAccounts(context.Background(), user.ID)
	assert.NoError(t, err)

	assert.ElementsMatch(t, expected, got)
}

func TestUserRepository_GetAccount(t *testing.T) {
	t.Parallel()
	conf := testutils.GetConfigWithDBName("user_repository_get_account")
	sqlConf := conf.SQLConf()
	h := testutils.SetupSQLHandler(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
	repo := irepository.NewUserRepository(h, mock_external_e2e.NewMockPortalAPI(), mock_external_e2e.NewMockTraQAPI())

	var (
		idx          = 1
		user         = mockdata.MockUsers[idx]
		accountType1 = domain.AccountType(3)
		accountType2 = domain.AccountType(4)
	)

	account1 := mustMakeAccount(t, repo, user.ID, &urepository.CreateAccountArgs{
		DisplayName: random.AlphaNumeric(),
		Type:        accountType1,
		URL:         random.AccountURLString(accountType1),
		PrPermitted: random.Bool(),
	})
	mustMakeAccount(t, repo, user.ID, &urepository.CreateAccountArgs{
		DisplayName: random.AlphaNumeric(),
		Type:        accountType2,
		URL:         random.AccountURLString(accountType2),
		PrPermitted: random.Bool(),
	})

	got, err := repo.GetAccount(context.Background(), user.ID, account1.ID)
	assert.NoError(t, err)
	assert.Equal(t, account1, got)
}

// func TestUserRepository_CreateAccount(t *testing.T) {
// }

func TestUserRepository_UpdateAccount(t *testing.T) {
	t.Parallel()
	conf := testutils.GetConfigWithDBName("user_repository_update_account")
	sqlConf := conf.SQLConf()
	h := testutils.SetupSQLHandler(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
	repo := irepository.NewUserRepository(h, mock_external_e2e.NewMockPortalAPI(), mock_external_e2e.NewMockTraQAPI())

	var (
		idx               = 1
		user              = mockdata.MockUsers[idx]
		accountType1      = domain.AccountType(3)
		accountType2      = domain.AccountType(4)
		accountType3Int64 = int64(6)
	)

	account1 := mustMakeAccount(t, repo, user.ID, &urepository.CreateAccountArgs{
		DisplayName: random.AlphaNumeric(),
		Type:        accountType1,
		URL:         random.AccountURLString(accountType1),
		PrPermitted: random.Bool(),
	})
	mustMakeAccount(t, repo, user.ID, &urepository.CreateAccountArgs{
		DisplayName: random.AlphaNumeric(),
		Type:        accountType2,
		URL:         random.AccountURLString(accountType2),
		PrPermitted: random.Bool(),
	})

	accountType := optional.From(domain.AccountType(accountType3Int64))
	args := &urepository.UpdateAccountArgs{
		DisplayName: random.OptAlphaNumeric(),
		Type:        accountType,
		URL:         random.OptAccountURLStringNotNull(domain.AccountType(accountType3Int64)),
		PrPermitted: random.OptBool(),
	}
	if v, ok := args.DisplayName.V(); ok {
		account1.DisplayName = v
	}
	if v, ok := args.Type.V(); ok {
		account1.Type = v
	}
	if v, ok := args.URL.V(); ok {
		account1.URL = v
	}
	if v, ok := args.PrPermitted.V(); ok {
		account1.PrPermitted = v
	}
	err = repo.UpdateAccount(context.Background(), user.ID, account1.ID, args)
	assert.NoError(t, err)

	got, err := repo.GetAccount(context.Background(), user.ID, account1.ID)
	assert.NoError(t, err)
	assert.Equal(t, account1, got)
}

func TestUserRepository_DeleteAccount(t *testing.T) {
	t.Parallel()
	conf := testutils.GetConfigWithDBName("user_repository_delete_account")
	sqlConf := conf.SQLConf()
	h := testutils.SetupSQLHandler(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
	repo := irepository.NewUserRepository(h, mock_external_e2e.NewMockPortalAPI(), mock_external_e2e.NewMockTraQAPI())

	var (
		idx          = 1
		user         = mockdata.MockUsers[idx]
		accountType1 = domain.AccountType(3)
		accountType2 = domain.AccountType(4)
	)

	account1 := mustMakeAccount(t, repo, user.ID, &urepository.CreateAccountArgs{
		DisplayName: random.AlphaNumeric(),
		Type:        accountType1,
		URL:         random.AccountURLString(accountType1),
		PrPermitted: random.Bool(),
	})
	account2 := mustMakeAccount(t, repo, user.ID, &urepository.CreateAccountArgs{
		DisplayName: random.AlphaNumeric(),
		Type:        accountType2,
		URL:         random.AccountURLString(accountType2),
		PrPermitted: random.Bool(),
	})

	err = repo.DeleteAccount(context.Background(), user.ID, account1.ID)
	assert.NoError(t, err)

	expected := []*domain.Account{account2}

	got, err := repo.GetAccounts(context.Background(), user.ID)
	assert.NoError(t, err)

	assert.ElementsMatch(t, expected, got)
}

func TestUserRepository_GetUserProjects(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("user_repository_get_user_projects")
	sqlConf := conf.SQLConf()
	h := testutils.SetupSQLHandler(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
	userRepo := irepository.NewUserRepository(h, mock_external_e2e.NewMockPortalAPI(), mock_external_e2e.NewMockTraQAPI())
	projectRepo := irepository.NewProjectRepository(h, mock_external_e2e.NewMockPortalAPI())

	project1 := mustMakeProject(t, projectRepo, nil)
	project2 := mustMakeProject(t, projectRepo, nil)
	user1 := mockdata.MockUsers[2]

	expected1 := []*domain.UserWithDuration{}
	expected2 := []*domain.UserWithDuration{}
	users1, err := projectRepo.GetProjectMembers(context.Background(), project1.ID)
	assert.NoError(t, err)
	users2, err := projectRepo.GetProjectMembers(context.Background(), project2.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected1, users1)
	assert.ElementsMatch(t, expected2, users2)

	args1 := mustAddProjectMember(t, projectRepo, project1.ID, user1.ID, nil)
	args2 := mustAddProjectMember(t, projectRepo, project2.ID, user1.ID, nil)

	expected3 := []*domain.UserProject{newUserProject(args1, project1), newUserProject(args2, project2)}
	projects1, err := userRepo.GetProjects(context.Background(), user1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected3, projects1)
}

func TestUserRepository_GetContests(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("user_repository_get_contests")
	sqlConf := conf.SQLConf()
	h := testutils.SetupSQLHandler(t, sqlConf)
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

	expected1 := []*domain.UserContest{newUserContest(&contest1.Contest, []*domain.ContestTeam{&team1.ContestTeam}), newUserContest(&contest2.Contest, []*domain.ContestTeam{&team2.ContestTeam})}
	contests1, err := userRepo.GetContests(context.Background(), user1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected1, contests1)

	expected2 := []*domain.UserContest{newUserContest(&contest1.Contest, []*domain.ContestTeam{&team1.ContestTeam})}
	contests2, err := userRepo.GetContests(context.Background(), user2.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected2, contests2)
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

func newUserContest(contest *domain.Contest, teams []*domain.ContestTeam) *domain.UserContest {
	return &domain.UserContest{
		ID:        contest.ID,
		Name:      contest.Name,
		TimeStart: contest.TimeStart,
		TimeEnd:   contest.TimeEnd,
		Teams:     teams,
	}
}
