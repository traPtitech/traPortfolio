package repository

import (
	"math/rand"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/database/mock_database"
	"github.com/traPtitech/traPortfolio/interfaces/repository/model"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/random"
)

type mockProjectRepositoryFields struct {
	h database.SQLHandler
}

func newMockProjectRepositoryFields() mockProjectRepositoryFields {
	return mockProjectRepositoryFields{
		h: mock_database.NewMockSQLHandler(),
	}
}

func TestProjectRepository_GetProjects(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		want      []*domain.Project
		setup     func(f mockProjectRepositoryFields, want []*domain.Project)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			want: []*domain.Project{
				{
					ID:          random.UUID(),
					Name:        random.AlphaNumeric(rand.Intn(30) + 1),
					Since:       time.Now(),
					Until:       time.Now(),
					Description: random.AlphaNumeric(rand.Intn(30) + 1),
					Link:        random.RandURLString(),
					Members:     nil,
				},
			},
			setup: func(f mockProjectRepositoryFields, want []*domain.Project) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "link", "since", "until"})
				for _, v := range want {
					rows.AddRow(v.ID, v.Name, v.Description, v.Link, v.Since, v.Until)
				}
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects`")).
					WillReturnRows(rows)
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			want: nil,
			setup: func(f mockProjectRepositoryFields, want []*domain.Project) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects`")).
					WillReturnError(errUnexpected)
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			f := newMockProjectRepositoryFields()
			tt.setup(f, tt.want)
			repo := NewProjectRepository(f.h)
			// Assertion
			got, err := repo.GetProjects()
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProjectRepository_GetProject(t *testing.T) {
	t.Parallel()
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.Project
		setup     func(f mockProjectRepositoryFields, args args, want *domain.Project)
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			f := newMockProjectRepositoryFields()
			tt.setup(f, tt.args, tt.want)
			repo := NewProjectRepository(f.h)
			// Assertion
			got, err := repo.GetProject(tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProjectRepository_CreateProject(t *testing.T) {
	t.Parallel()
	type args struct {
		project *model.Project
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.Project
		setup     func(f mockProjectRepositoryFields, args args, want *domain.Project)
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			f := newMockProjectRepositoryFields()
			tt.setup(f, tt.args, tt.want)
			repo := NewProjectRepository(f.h)
			// Assertion
			got, err := repo.CreateProject(tt.args.project)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProjectRepository_UpdateProject(t *testing.T) {
	t.Parallel()
	type args struct {
		id      uuid.UUID
		changes map[string]interface{}
	}
	tests := []struct {
		name      string
		args      args
		setup     func(f mockProjectRepositoryFields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			f := newMockProjectRepositoryFields()
			tt.setup(f, tt.args)
			repo := NewProjectRepository(f.h)
			// Assertion
			tt.assertion(t, repo.UpdateProject(tt.args.id, tt.args.changes))
		})
	}
}

func TestProjectRepository_GetProjectMembers(t *testing.T) {
	t.Parallel()
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      []*domain.User
		setup     func(f mockProjectRepositoryFields, args args, want []*domain.User)
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			f := newMockProjectRepositoryFields()
			tt.setup(f, tt.args, tt.want)
			repo := NewProjectRepository(f.h)
			// Assertion
			got, err := repo.GetProjectMembers(tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProjectRepository_AddProjectMembers(t *testing.T) {
	t.Parallel()
	type args struct {
		projectID      uuid.UUID
		projectMembers []*repository.CreateProjectMemberArgs
	}
	tests := []struct {
		name      string
		args      args
		setup     func(f mockProjectRepositoryFields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			f := newMockProjectRepositoryFields()
			tt.setup(f, tt.args)
			repo := NewProjectRepository(f.h)
			// Assertion
			tt.assertion(t, repo.AddProjectMembers(tt.args.projectID, tt.args.projectMembers))
		})
	}
}

func TestProjectRepository_DeleteProjectMembers(t *testing.T) {
	t.Parallel()
	type args struct {
		projectID uuid.UUID
		members   []uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		setup     func(f mockProjectRepositoryFields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			f := newMockProjectRepositoryFields()
			tt.setup(f, tt.args)
			repo := NewProjectRepository(f.h)
			// Assertion
			tt.assertion(t, repo.DeleteProjectMembers(tt.args.projectID, tt.args.members))
		})
	}
}
