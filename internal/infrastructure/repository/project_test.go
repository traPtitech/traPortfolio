package repository

import (
	"context"
	"testing"

	urepository "github.com/traPtitech/traPortfolio/internal/usecases/repository"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/internal/domain"
	"github.com/traPtitech/traPortfolio/internal/infrastructure/external/mock_external_e2e"
	"github.com/traPtitech/traPortfolio/internal/pkgs/mockdata"
	"github.com/traPtitech/traPortfolio/internal/pkgs/optional"
	"github.com/traPtitech/traPortfolio/internal/pkgs/random"
)

func TestProjectRepository_GetProjects(t *testing.T) {
	t.Parallel()

	db := SetupTestGormDB(t)
	repo := NewProjectRepository(db, mock_external_e2e.NewMockPortalAPI())

	projectNum := 4
	var projects []*domain.Project
	for range projectNum {
		projects = append(projects, mustMakeProject(t, repo, nil))
	}

	got, err := repo.GetProjects(context.Background())
	assert.NoError(t, err)

	assert.ElementsMatch(t, projects, got)
}

func TestProjectRepository_GetProject(t *testing.T) {
	t.Parallel()

	db := SetupTestGormDB(t)
	repo := NewProjectRepository(db, mock_external_e2e.NewMockPortalAPI())

	projectNum := 4
	var projects []*domain.ProjectDetail
	for range projectNum {
		projects = append(projects, mustMakeProjectDetail(t, repo, nil))
	}

	opts := []cmp.Option{
		cmpopts.EquateEmpty(),
		cmp.AllowUnexported(optional.Of[domain.YearWithSemester]{}),
	}
	for _, p := range projects {
		got, err := repo.GetProject(context.Background(), p.ID)
		assert.NoError(t, err)

		if diff := cmp.Diff(p, got, opts...); diff != "" {
			t.Error(diff)
		}
	}
}

// func TestProjectRepository_CreateProject(t *testing.T) {
// }

func TestProjectRepository_UpdateProject(t *testing.T) {
	t.Parallel()

	db := SetupTestGormDB(t)
	repo := NewProjectRepository(db, mock_external_e2e.NewMockPortalAPI())

	tests := []struct {
		name string
		ctx  context.Context
		args *urepository.UpdateProjectArgs
	}{
		{
			name: "all fields",
			ctx:  context.Background(),
			args: random.UpdateProjectArgs(),
		},
		{
			name: "partial fields",
			ctx:  context.Background(),
			args: random.OptUpdateProjectArgs(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project1 := mustMakeProjectDetail(t, repo, nil)
			mustMakeProjectDetail(t, repo, nil)

			arg1 := tt.args

			project1.Name = arg1.Name.ValueOr(project1.Name)
			project1.Description = arg1.Description.ValueOr(project1.Description)
			project1.Link = arg1.Link.ValueOr(project1.Link)
			if sy, ok := arg1.SinceYear.V(); ok {
				if ss, ok := arg1.SinceSemester.V(); ok {
					project1.Duration.Since.Year = int(sy)
					project1.Duration.Since.Semester = int(ss)
				}
			}
			uy, validYear := arg1.UntilYear.V()
			us, validSemester := arg1.UntilSemester.V()
			if validYear == validSemester {
				originUntil, originValid := project1.Duration.Until.V()
				argUntil := domain.YearWithSemester{Year: int(uy), Semester: int(us)}
				if validYear != originValid || (validYear && argUntil != originUntil) {
					project1.Duration.Until = optional.From(domain.YearWithSemester{
						Year:     int(uy),
						Semester: int(us),
					})
				}
			}

			err := repo.UpdateProject(tt.ctx, project1.ID, arg1)
			assert.NoError(t, err)

			got, err := repo.GetProject(tt.ctx, project1.ID)
			assert.NoError(t, err)

			opts := []cmp.Option{
				cmpopts.EquateEmpty(),
				cmp.AllowUnexported(optional.Of[domain.YearWithSemester]{}),
			}
			if diff := cmp.Diff(project1, got, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestProjectRepository_GetProjectMembers(t *testing.T) {
	t.Parallel()

	db := SetupTestGormDB(t)
	err := mockdata.InsertSampleDataToDB(db)
	assert.NoError(t, err)
	repo := NewProjectRepository(db, mock_external_e2e.NewMockPortalAPI())

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

	dur1 := random.DurationBetween(project1.Duration.Since, project1.Duration.Until.ValueOrZero())
	dur2 := random.DurationBetween(project2.Duration.Since, project2.Duration.Until.ValueOrZero())

	mustExistProjectMember(t, repo, project1.ID, project1.Duration, []*urepository.EditProjectMemberArgs{
		{
			UserID:        user1.ID,
			SinceYear:     dur1.Since.Year,
			SinceSemester: dur1.Since.Semester,
			UntilYear:     dur1.Until.ValueOrZero().Year,
			UntilSemester: dur1.Until.ValueOrZero().Semester,
		},
		{
			UserID:        user2.ID,
			SinceYear:     dur1.Since.Year,
			SinceSemester: dur1.Since.Semester,
			UntilYear:     dur1.Until.ValueOrZero().Year,
			UntilSemester: dur1.Until.ValueOrZero().Semester,
		},
	})
	mustExistProjectMember(t, repo, project2.ID, project2.Duration, []*urepository.EditProjectMemberArgs{
		{
			UserID:        user2.ID,
			SinceYear:     dur2.Since.Year,
			SinceSemester: dur2.Since.Semester,
			UntilYear:     dur2.Until.ValueOrZero().Year,
			UntilSemester: dur2.Until.ValueOrZero().Semester,
		},
	})

	projectMember1 := &domain.UserWithDuration{
		User: *user1,
		Duration: domain.NewYearWithSemesterDuration(
			dur1.Since.Year,
			dur1.Since.Semester,
			dur1.Until.ValueOrZero().Year,
			dur1.Until.ValueOrZero().Semester,
		),
	}
	projectMember2 := &domain.UserWithDuration{
		User: *user2,
		Duration: domain.NewYearWithSemesterDuration(
			dur1.Since.Year,
			dur1.Since.Semester,
			dur1.Until.ValueOrZero().Year,
			dur1.Until.ValueOrZero().Semester,
		),
	}
	expected1 := []*domain.UserWithDuration{projectMember1, projectMember2}
	users1, err := repo.GetProjectMembers(context.Background(), project1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected1, users1)

	projectMember3 := &domain.UserWithDuration{
		User: *user2,
		Duration: domain.NewYearWithSemesterDuration(
			dur2.Since.Year,
			dur2.Since.Semester,
			dur2.Until.ValueOrZero().Year,
			dur2.Until.ValueOrZero().Semester,
		),
	}

	expected2 := []*domain.UserWithDuration{projectMember3}
	users2, err := repo.GetProjectMembers(context.Background(), project2.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected2, users2)
}

// TestGetProjectMembers と似たような内容になるから省略
// func TestProjectRepository_AddProjectMembers(t *testing.T) {
// }

func TestProjectRepository_EditProjectMembers(t *testing.T) {
	t.Parallel()

	db := SetupTestGormDB(t)
	err := mockdata.InsertSampleDataToDB(db)
	assert.NoError(t, err)
	repo := NewProjectRepository(db, mock_external_e2e.NewMockPortalAPI())

	project1 := mustMakeProjectDetail(t, repo, nil)
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
	user3 := domain.NewUser(
		mockdata.MockUsers[2].ID,
		mockdata.MockUsers[2].Name,
		mockdata.MockPortalUsers[2].RealName,
		mockdata.MockUsers[2].Check,
	)

	dur1 := random.DurationBetween(project1.Duration.Since, project1.Duration.Until.ValueOrZero())
	dur2 := random.DurationBetween(project1.Duration.Since, project1.Duration.Until.ValueOrZero())

	mustExistProjectMember(t, repo, project1.ID, project1.Duration, []*urepository.EditProjectMemberArgs{
		{
			UserID:        user1.ID,
			SinceYear:     dur1.Since.Year,
			SinceSemester: dur1.Since.Semester,
			UntilYear:     dur1.Until.ValueOrZero().Year,
			UntilSemester: dur1.Until.ValueOrZero().Semester,
		},
		{
			UserID:        user2.ID,
			SinceYear:     dur1.Since.Year,
			SinceSemester: dur1.Since.Semester,
			UntilYear:     dur1.Until.ValueOrZero().Year,
			UntilSemester: dur1.Until.ValueOrZero().Semester,
		},
	})

	projectMember1 := &domain.UserWithDuration{
		User: *user1,
		Duration: domain.NewYearWithSemesterDuration(
			dur1.Since.Year,
			dur1.Since.Semester,
			dur1.Until.ValueOrZero().Year,
			dur1.Until.ValueOrZero().Semester,
		),
	}
	projectMember2 := &domain.UserWithDuration{
		User: *user2,
		Duration: domain.NewYearWithSemesterDuration(
			dur1.Since.Year,
			dur1.Since.Semester,
			dur1.Until.ValueOrZero().Year,
			dur1.Until.ValueOrZero().Semester,
		),
	}

	expected1 := []*domain.UserWithDuration{projectMember1, projectMember2}
	users1, err := repo.GetProjectMembers(context.Background(), project1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected1, users1)

	projectMember3 := &domain.UserWithDuration{
		User: *user2,
		Duration: domain.NewYearWithSemesterDuration(
			dur2.Since.Year,
			dur2.Since.Semester,
			dur2.Until.ValueOrZero().Year,
			dur2.Until.ValueOrZero().Semester,
		),
	}

	projectMember4 := &domain.UserWithDuration{
		User: *user3,
		Duration: domain.NewYearWithSemesterDuration(
			dur1.Since.Year,
			dur1.Since.Semester,
			dur1.Until.ValueOrZero().Year,
			dur1.Until.ValueOrZero().Semester,
		),
	}

	err = repo.EditProjectMembers(context.Background(), project1.ID, []*urepository.EditProjectMemberArgs{
		{
			UserID:        user2.ID,
			SinceYear:     dur2.Since.Year,
			SinceSemester: dur2.Since.Semester,
			UntilYear:     dur2.Until.ValueOrZero().Year,
			UntilSemester: dur2.Until.ValueOrZero().Semester,
		},
		{
			UserID:        user3.ID,
			SinceYear:     dur1.Since.Year,
			SinceSemester: dur1.Since.Semester,
			UntilYear:     dur1.Until.ValueOrZero().Year,
			UntilSemester: dur1.Until.ValueOrZero().Semester,
		},
	})
	assert.NoError(t, err)

	expected2 := []*domain.UserWithDuration{projectMember3, projectMember4}
	users2, err := repo.GetProjectMembers(context.Background(), project1.ID)
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected2, users2)
}
