package repository

import (
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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

func TestProjectRepository_GetProjects(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("project_repository_get_projects")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
	repo := irepository.NewProjectRepository(h, mock_external_e2e.NewMockPortalAPI())

	projectNum := 4
	var projects []*domain.Project
	for i := 0; i < projectNum; i++ {
		projects = append(projects, mustMakeProject(t, repo, nil))
	}

	got, err := repo.GetProjects()
	assert.NoError(t, err)

	assert.ElementsMatch(t, projects, got)
}

func TestProjectRepository_GetProject(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("project_repository_get_project")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
	repo := irepository.NewProjectRepository(h, mock_external_e2e.NewMockPortalAPI())

	projectNum := 4
	var projects []*domain.ProjectDetail
	for i := 0; i < projectNum; i++ {
		projects = append(projects, mustMakeProjectDetail(t, repo, nil))
	}

	opt := cmpopts.EquateEmpty()
	for _, p := range projects {
		got, err := repo.GetProject(p.ID)
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
	h := testutils.SetupDB(t, sqlConf)
	repo := irepository.NewProjectRepository(h, mock_external_e2e.NewMockPortalAPI())

	project1 := mustMakeProjectDetail(t, repo, nil)
	mustMakeProjectDetail(t, repo, nil)

	arg1 := urepository.UpdateProjectArgs{
		Name:          random.OptAlphaNumeric(),
		Description:   random.OptAlphaNumeric(),
		Link:          random.OptAlphaNumeric(),
		SinceYear:     optional.NewInt64(int64(random.Duration().Since.Year), random.Bool()),
		SinceSemester: optional.NewInt64(int64(random.Duration().Since.Semester), random.Bool()),
		UntilYear:     optional.NewInt64(int64(random.Duration().Until.Year), random.Bool()),
		UntilSemester: optional.NewInt64(int64(random.Duration().Until.Semester), random.Bool()),
	}

	if arg1.Name.Valid {
		project1.Name = arg1.Name.String
	}
	if arg1.Description.Valid {
		project1.Description = arg1.Description.String
	}
	if arg1.Link.Valid {
		project1.Link = arg1.Link.String
	}
	if arg1.SinceYear.Valid && arg1.SinceSemester.Valid {
		project1.Duration.Since.Year = int(arg1.SinceYear.Int64)
		project1.Duration.Since.Semester = int(arg1.SinceSemester.Int64)
	}
	if arg1.UntilYear.Valid && arg1.UntilSemester.Valid {
		project1.Duration.Until.Year = int(arg1.UntilYear.Int64)
		project1.Duration.Until.Semester = int(arg1.UntilSemester.Int64)
	}

	err := repo.UpdateProject(project1.ID, &arg1)
	assert.NoError(t, err)

	got, err := repo.GetProject(project1.ID)
	assert.NoError(t, err)

	opt := cmpopts.EquateEmpty()
	if diff := cmp.Diff(project1, got, opt); diff != "" {
		t.Error(diff)
	}
}

func TestProjectRepository_UpdateProject_Concurrent(t *testing.T) {
	t.Parallel()

	conf := testutils.GetConfigWithDBName("project_repository_update_project_concurrent")
	sqlConf := conf.SQLConf()
	h := testutils.SetupDB(t, sqlConf)
	repo := irepository.NewProjectRepository(h, mock_external_e2e.NewMockPortalAPI())

	project1 := mustMakeProjectDetail(t, repo, nil)
	mustMakeProjectDetail(t, repo, nil)

	wg := sync.WaitGroup{}
	for i := 0; i < 3; i++ {

		wg.Add(1)
		go func(t *testing.T) {
			defer wg.Done()
			arg1 := urepository.UpdateProjectArgs{
				Name:          random.OptAlphaNumeric(),
				Description:   random.OptAlphaNumeric(),
				Link:          random.OptAlphaNumeric(),
				SinceYear:     optional.NewInt64(int64(random.Duration().Since.Year), random.Bool()),
				SinceSemester: optional.NewInt64(int64(random.Duration().Since.Semester), random.Bool()),
				UntilYear:     optional.NewInt64(int64(random.Duration().Until.Year), random.Bool()),
				UntilSemester: optional.NewInt64(int64(random.Duration().Until.Semester), random.Bool()),
			}

			if arg1.Name.Valid {
				project1.Name = arg1.Name.String
			}
			if arg1.Description.Valid {
				project1.Description = arg1.Description.String
			}
			if arg1.Link.Valid {
				project1.Link = arg1.Link.String
			}
			if arg1.SinceYear.Valid && arg1.SinceSemester.Valid {
				project1.Duration.Since.Year = int(arg1.SinceYear.Int64)
				project1.Duration.Since.Semester = int(arg1.SinceSemester.Int64)
			}
			if arg1.UntilYear.Valid && arg1.UntilSemester.Valid {
				project1.Duration.Until.Year = int(arg1.UntilYear.Int64)
				project1.Duration.Until.Semester = int(arg1.UntilSemester.Int64)
			}

			err := repo.UpdateProject(project1.ID, &arg1)
			assert.NoError(t, err)
		}(t)
	}

	wg.Wait()

	got, err := repo.GetProject(project1.ID)
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
	h := testutils.SetupDB(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
	repo := irepository.NewProjectRepository(h, mock_external_e2e.NewMockPortalAPI())

	project1 := mustMakeProjectDetail(t, repo, nil)
	project2 := mustMakeProjectDetail(t, repo, nil)
	user1 := &domain.User{
		ID:       mockdata.MockUsers[0].ID,
		Name:     mockdata.MockUsers[0].Name,
		RealName: mockdata.MockPortalUsers[0].RealName,
	}
	user2 := &domain.User{
		ID:       mockdata.MockUsers[1].ID,
		Name:     mockdata.MockUsers[1].Name,
		RealName: mockdata.MockPortalUsers[1].RealName,
	}

	args1 := mustAddProjectMember(t, repo, project1.ID, user1.ID, nil)
	args2 := mustAddProjectMember(t, repo, project1.ID, user2.ID, nil)
	args3 := mustAddProjectMember(t, repo, project2.ID, user2.ID, nil)

	projectMember1 := &domain.ProjectMember{
		UserID:   user1.ID,
		Name:     user1.Name,
		RealName: user1.RealName,
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
	projectMember2 := &domain.ProjectMember{
		UserID:   user2.ID,
		Name:     user2.Name,
		RealName: user2.RealName,
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
	expected1 := []*domain.ProjectMember{projectMember1, projectMember2}
	users1, err := repo.GetProjectMembers(project1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected1, users1)

	projectMember3 := &domain.ProjectMember{
		UserID:   user2.ID,
		Name:     user2.Name,
		RealName: user2.RealName,
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

	expected2 := []*domain.ProjectMember{projectMember3}
	users2, err := repo.GetProjectMembers(project2.ID)
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
	h := testutils.SetupDB(t, sqlConf)
	err := mockdata.InsertSampleDataToDB(h)
	assert.NoError(t, err)
	repo := irepository.NewProjectRepository(h, mock_external_e2e.NewMockPortalAPI())

	project1 := mustMakeProjectDetail(t, repo, nil)
	project2 := mustMakeProjectDetail(t, repo, nil)
	user1 := &domain.User{
		ID:       mockdata.MockUsers[1].ID,
		Name:     mockdata.MockUsers[1].Name,
		RealName: mockdata.MockPortalUsers[1].RealName,
	}
	user2 := &domain.User{
		ID:       mockdata.MockUsers[2].ID,
		Name:     mockdata.MockUsers[2].Name,
		RealName: mockdata.MockPortalUsers[2].RealName,
	}

	args1 := mustAddProjectMember(t, repo, project1.ID, user1.ID, nil)
	args2 := mustAddProjectMember(t, repo, project1.ID, user2.ID, nil)
	args3 := mustAddProjectMember(t, repo, project2.ID, user2.ID, nil)

	projectMember1 := &domain.ProjectMember{
		UserID:   user1.ID,
		Name:     user1.Name,
		RealName: user1.RealName,
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
	projectMember2 := &domain.ProjectMember{
		UserID:   user2.ID,
		Name:     user2.Name,
		RealName: user2.RealName,
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
	expected1 := []*domain.ProjectMember{projectMember1, projectMember2}
	users1, err := repo.GetProjectMembers(project1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected1, users1)

	projectMember3 := &domain.ProjectMember{
		UserID:   user2.ID,
		Name:     user2.Name,
		RealName: user2.RealName,
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
	expected2 := []*domain.ProjectMember{projectMember3}
	users2, err := repo.GetProjectMembers(project2.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected2, users2)

	err = repo.DeleteProjectMembers(project1.ID, []uuid.UUID{projectMember1.UserID})
	assert.NoError(t, err)
	expected3 := []*domain.ProjectMember{projectMember2}
	users3, err := repo.GetProjectMembers(project1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected3, users3)

	err = repo.DeleteProjectMembers(project1.ID, []uuid.UUID{projectMember2.UserID})
	assert.NoError(t, err)
	expected4 := []*domain.ProjectMember{}
	users4, err := repo.GetProjectMembers(project1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected4, users4)

	err = repo.DeleteProjectMembers(project2.ID, []uuid.UUID{projectMember3.UserID})
	assert.NoError(t, err)
	expected5 := []*domain.ProjectMember{}
	users5, err := repo.GetProjectMembers(project2.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected5, users5)
}
