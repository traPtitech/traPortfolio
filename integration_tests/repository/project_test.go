package repository

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	irepository "github.com/traPtitech/traPortfolio/infrastructure/repository"
	"github.com/traPtitech/traPortfolio/integration_tests/testutils"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external_e2e"
	urepository "github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/mockdata"
	"github.com/traPtitech/traPortfolio/util/random"
)

func TestProjectRepository_GetProjects(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("project_repository_get_projects")
	sqlConf := conf.SQLConf()
	h := testutils.SetupSQLHandler(t, sqlConf)
	repo := irepository.NewProjectRepository(h, mock_external_e2e.NewMockPortalAPI())

	projectNum := 4
	var projects []*domain.Project
	for i := 0; i < projectNum; i++ {
		projects = append(projects, mustMakeProject(t, repo, nil))
	}

	got, err := repo.GetProjects(context.Background())
	assert.NoError(t, err)

	assert.ElementsMatch(t, projects, got)
}

func TestProjectRepository_GetProject(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("project_repository_get_project")
	sqlConf := conf.SQLConf()
	h := testutils.SetupSQLHandler(t, sqlConf)
	repo := irepository.NewProjectRepository(h, mock_external_e2e.NewMockPortalAPI())

	projectNum := 4
	var projects []*domain.ProjectDetail
	for i := 0; i < projectNum; i++ {
		projects = append(projects, mustMakeProjectDetail(t, repo, nil))
	}

	opt := cmpopts.EquateEmpty()
	for _, p := range projects {
		got, err := repo.GetProject(context.Background(), p.ID)
		assert.NoError(t, err)

		if diff := cmp.Diff(p, got, opt); diff != "" {
			t.Error(diff)
		}
	}
}

// func TestProjectRepository_CreateProject(t *testing.T) {
// }

func TestProjectRepository_UpdateProject(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("project_repository_update_project")
	sqlConf := conf.SQLConf()
	h := testutils.SetupSQLHandler(t, sqlConf)
	repo := irepository.NewProjectRepository(h, mock_external_e2e.NewMockPortalAPI())

	project1 := mustMakeProjectDetail(t, repo, nil)
	mustMakeProjectDetail(t, repo, nil)

	arg1 := urepository.UpdateProjectArgs{
		Name:          random.Optional(random.AlphaNumeric()),
		Description:   random.Optional(random.AlphaNumeric()),
		Link:          random.Optional(random.AlphaNumeric()),
		SinceYear:     random.Optional(int64(2100)), // TODO: intでよさそう
		SinceSemester: random.Optional(int64(2)),
		UntilYear:     random.Optional(int64(2100)),
		UntilSemester: random.Optional(int64(2)),
	}

	if v, ok := arg1.Name.V(); ok {
		project1.Name = v
	}
	if v, ok := arg1.Description.V(); ok {
		project1.Description = v
	}
	if v, ok := arg1.Link.V(); ok {
		project1.Link = v
	}
	if sy, ok := arg1.SinceYear.V(); ok {
		if ss, ok := arg1.SinceSemester.V(); ok {
			project1.Duration.Since.Year = int(sy)
			project1.Duration.Since.Semester = int(ss)
		}
	}
	if uy, ok := arg1.UntilYear.V(); ok {
		if us, ok := arg1.UntilSemester.V(); ok {
			project1.Duration.Until.Year = int(uy)
			project1.Duration.Until.Semester = int(us)
		}
	}

	err := repo.UpdateProject(context.Background(), project1.ID, &arg1)
	assert.NoError(t, err)

	got, err := repo.GetProject(context.Background(), project1.ID)
	assert.NoError(t, err)

	opt := cmpopts.EquateEmpty()
	if diff := cmp.Diff(project1, got, opt); diff != "" {
		t.Error(diff)
	}
}

func TestProjectRepository_GetProjectMembers(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("project_repository_get_project_members")
	sqlConf := conf.SQLConf()
	h := testutils.SetupSQLHandler(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
	repo := irepository.NewProjectRepository(h, mock_external_e2e.NewMockPortalAPI())

	project1 := mustMakeProjectDetail(t, repo, nil)
	project2 := mustMakeProjectDetail(t, repo, nil)
	user1 := domain.NewUser(
		mockdata.MockUsers[0].ID,
		mockdata.MockUsers[0].Name,
		mockdata.MockPortalUsers[0].RealName,
		mockdata.MockUsers[0].Check,
	)
	user2 := domain.NewUser(
		mockdata.MockUsers[1].ID,
		mockdata.MockUsers[1].Name,
		mockdata.MockPortalUsers[1].RealName,
		mockdata.MockUsers[1].Check,
	)

	args1 := mustAddProjectMember(t, repo, project1.ID, user1.ID, nil)
	args2 := mustAddProjectMember(t, repo, project1.ID, user2.ID, nil)
	args3 := mustAddProjectMember(t, repo, project2.ID, user2.ID, nil)

	projectMember1 := &domain.UserWithDuration{
		User: *user1,
		Duration: domain.YearWithSemesterDuration{
			Since: domain.YearWithSemester{
				Year:     args1.SinceYear,
				Semester: args1.SinceSemester,
			},
			Until: domain.YearWithSemester{
				Year:     args1.UntilYear,
				Semester: args1.UntilSemester,
			},
		},
	}
	projectMember2 := &domain.UserWithDuration{
		User: *user2,
		Duration: domain.YearWithSemesterDuration{
			Since: domain.YearWithSemester{
				Year:     args2.SinceYear,
				Semester: args2.SinceSemester,
			},
			Until: domain.YearWithSemester{
				Year:     args2.UntilYear,
				Semester: args2.UntilSemester,
			},
		},
	}
	expected1 := []*domain.UserWithDuration{projectMember1, projectMember2}
	users1, err := repo.GetProjectMembers(context.Background(), project1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected1, users1)

	projectMember3 := &domain.UserWithDuration{
		User: *user2,
		Duration: domain.YearWithSemesterDuration{
			Since: domain.YearWithSemester{
				Year:     args3.SinceYear,
				Semester: args3.SinceSemester,
			},
			Until: domain.YearWithSemester{
				Year:     args3.UntilYear,
				Semester: args3.UntilSemester,
			},
		},
	}

	expected2 := []*domain.UserWithDuration{projectMember3}
	users2, err := repo.GetProjectMembers(context.Background(), project2.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected2, users2)
}

// TestGetProjectMembers と似たような内容になるから省略
// func TestProjectRepository_AddProjectMembers(t *testing.T) {
// }

func TestProjectRepository_DeleteProjectMembers(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("project_repository_delete_project_members")
	sqlConf := conf.SQLConf()
	h := testutils.SetupSQLHandler(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
	repo := irepository.NewProjectRepository(h, mock_external_e2e.NewMockPortalAPI())

	project1 := mustMakeProjectDetail(t, repo, nil)
	project2 := mustMakeProjectDetail(t, repo, nil)
	user1 := domain.NewUser(
		mockdata.MockUsers[1].ID,
		mockdata.MockUsers[1].Name,
		mockdata.MockPortalUsers[1].RealName,
		mockdata.MockUsers[1].Check,
	)
	user2 := domain.NewUser(
		mockdata.MockUsers[2].ID,
		mockdata.MockUsers[2].Name,
		mockdata.MockPortalUsers[2].RealName,
		mockdata.MockUsers[2].Check,
	)

	args1 := mustAddProjectMember(t, repo, project1.ID, user1.ID, nil)
	args2 := mustAddProjectMember(t, repo, project1.ID, user2.ID, nil)
	args3 := mustAddProjectMember(t, repo, project2.ID, user2.ID, nil)

	projectMember1 := &domain.UserWithDuration{
		User: *user1,
		Duration: domain.YearWithSemesterDuration{
			Since: domain.YearWithSemester{
				Year:     args1.SinceYear,
				Semester: args1.SinceSemester,
			},
			Until: domain.YearWithSemester{
				Year:     args1.UntilYear,
				Semester: args1.UntilSemester,
			},
		},
	}
	projectMember2 := &domain.UserWithDuration{
		User: *user2,
		Duration: domain.YearWithSemesterDuration{
			Since: domain.YearWithSemester{
				Year:     args2.SinceYear,
				Semester: args2.SinceSemester,
			},
			Until: domain.YearWithSemester{
				Year:     args2.UntilYear,
				Semester: args2.UntilSemester,
			},
		},
	}
	expected1 := []*domain.UserWithDuration{projectMember1, projectMember2}
	users1, err := repo.GetProjectMembers(context.Background(), project1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected1, users1)

	projectMember3 := &domain.UserWithDuration{
		User: *user2,
		Duration: domain.YearWithSemesterDuration{
			Since: domain.YearWithSemester{
				Year:     args3.SinceYear,
				Semester: args3.SinceSemester,
			},
			Until: domain.YearWithSemester{
				Year:     args3.UntilYear,
				Semester: args3.UntilSemester,
			},
		},
	}
	expected2 := []*domain.UserWithDuration{projectMember3}
	users2, err := repo.GetProjectMembers(context.Background(), project2.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected2, users2)

	err = repo.DeleteProjectMembers(context.Background(), project1.ID, []uuid.UUID{projectMember1.User.ID})
	assert.NoError(t, err)
	expected3 := []*domain.UserWithDuration{projectMember2}
	users3, err := repo.GetProjectMembers(context.Background(), project1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected3, users3)

	err = repo.DeleteProjectMembers(context.Background(), project1.ID, []uuid.UUID{projectMember2.User.ID})
	assert.NoError(t, err)
	expected4 := []*domain.UserWithDuration{}
	users4, err := repo.GetProjectMembers(context.Background(), project1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected4, users4)

	err = repo.DeleteProjectMembers(context.Background(), project2.ID, []uuid.UUID{projectMember3.User.ID})
	assert.NoError(t, err)
	expected5 := []*domain.UserWithDuration{}
	users5, err := repo.GetProjectMembers(context.Background(), project2.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected5, users5)
}
