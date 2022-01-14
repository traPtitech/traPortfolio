package repository

import (
	"database/sql/driver"
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

// 0 first semester, 1 second semester
func makeYearWithSemester(s int) domain.YearWithSemester {
	return domain.YearWithSemester{
		Year:     random.Time().Year(),
		Semester: s,
	}
}

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
					ID:   random.UUID(),
					Name: random.AlphaNumeric(rand.Intn(30) + 1),
					Duration: domain.YearWithSemesterDuration{
						Since: makeYearWithSemester(rand.Intn(2)),
						Until: makeYearWithSemester(rand.Intn(2)),
					},
					Description: random.AlphaNumeric(rand.Intn(30) + 1),
					Link:        random.RandURLString(),
					Members:     nil,
				},
			},
			setup: func(f mockProjectRepositoryFields, want []*domain.Project) {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "link", "since_year", "since_semester", "until_year", "until_semester"})
				for _, v := range want {
					d := v.Duration
					rows.AddRow(v.ID, v.Name, v.Description, v.Link, d.Since.Year, d.Since.Semester, d.Until.Year, d.Until.Semester)
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
	pid := random.UUID() // Successで使うprojectID

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
		{
			name: "Success_single",
			args: args{
				id: pid,
			},
			want: &domain.Project{
				ID:          pid,
				Name:        random.AlphaNumeric(rand.Intn(30) + 1),
				Duration:    random.Duration(),
				Description: random.AlphaNumeric(rand.Intn(30) + 1),
				Link:        random.RandURLString(),
				Members: []*domain.ProjectMember{
					{
						UserID: random.UUID(),
						Name:   random.AlphaNumeric(rand.Intn(30) + 1),
						// RealName:
						Duration: random.Duration(),
					},
				},
			},
			setup: func(f mockProjectRepositoryFields, args args, want *domain.Project) {
				h := f.h.(*mock_database.MockSQLHandler)
				wd := want.Duration
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since_year", "since_semester", "until_year", "until_semester"}).
							AddRow(want.ID, want.Name, want.Description, want.Link, wd.Since.Year, wd.Since.Semester, wd.Until.Year, wd.Until.Semester),
					)
				wm := want.Members[0]
				wmd := wm.Duration
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `project_members` WHERE `project_members`.`project_id` = ?")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"user_id", "name", "since_year", "since_semester", "until_year", "until_semester"}).
							AddRow(wm.UserID, wm.Name, wmd.Since.Year, wmd.Since.Semester, wmd.Until.Year, wmd.Until.Semester),
					)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ?")).
					WithArgs(want.Members[0].UserID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name"}).
							AddRow(want.Members[0].UserID, want.Members[0].Name),
					)
			},
			assertion: assert.NoError,
		},
		{
			name: "Success_Multiple",
			args: args{
				id: pid,
			},
			want: &domain.Project{
				ID:          pid,
				Name:        random.AlphaNumeric(rand.Intn(30) + 1),
				Duration:    random.Duration(),
				Description: random.AlphaNumeric(rand.Intn(30) + 1),
				Link:        random.RandURLString(),
				Members: []*domain.ProjectMember{
					{
						UserID: random.UUID(),
						Name:   random.AlphaNumeric(rand.Intn(30) + 1),
						// RealName:
						Duration: random.Duration(),
					},
					{
						UserID: random.UUID(),
						Name:   random.AlphaNumeric(rand.Intn(30) + 1),
						// RealName:
						Duration: random.Duration(),
					},
				},
			},
			setup: func(f mockProjectRepositoryFields, args args, want *domain.Project) {
				h := f.h.(*mock_database.MockSQLHandler)
				wd := want.Duration
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since_year", "since_semester", "until_year", "until_semester"}).
							AddRow(want.ID, want.Name, want.Description, want.Link, wd.Since.Year, wd.Since.Semester, wd.Until.Year, wd.Until.Semester),
					)
				memberRows := sqlmock.NewRows([]string{"user_id", "name", "since_year", "since_semester", "until_year", "until_semester"})
				for _, v := range want.Members {
					d := v.Duration
					memberRows.AddRow(v.UserID, v.Name, d.Since.Year, d.Since.Semester, d.Until.Year, d.Until.Semester)
				}
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `project_members` WHERE `project_members`.`project_id` = ?")).
					WithArgs(args.id).
					WillReturnRows(memberRows)
				userIDs := make([]driver.Value, len(want.Members))
				userRows := sqlmock.NewRows([]string{"id", "name"})
				for i, v := range want.Members {
					userIDs[i] = v.UserID
					userRows.AddRow(v.UserID, v.Name)
				}
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` IN (?,?)")).
					WithArgs(userIDs...).
					WillReturnRows(userRows)
			},
			assertion: assert.NoError,
		},
		{
			name: "ProjectNotFound",
			args: args{
				id: random.UUID(),
			},
			want: nil,
			setup: func(f mockProjectRepositoryFields, args args, want *domain.Project) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "link", "since_year", "since_semester", "until_year", "until_semester"}))
			},
			assertion: assert.Error,
		},
		{
			name: "ProjectMemberUnexpectedError",
			args: args{
				id: random.UUID(),
			},
			want: nil,
			setup: func(f mockProjectRepositoryFields, args args, want *domain.Project) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since_year", "since_semester", "until_year", "until_semester"}).
							AddRow(
								args.id,
								random.AlphaNumeric(rand.Intn(30)+1),
								random.AlphaNumeric(rand.Intn(30)+1),
								random.RandURLString(),
								random.Time().Year(),
								rand.Intn(2),
								random.Time().Year(),
								rand.Intn(2),
							),
					)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `project_members` WHERE `project_members`.`project_id` = ?")).
					WithArgs(args.id).
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
	successProject := &model.Project{
		ID:            random.UUID(),
		Name:          random.AlphaNumeric(rand.Intn(30) + 1),
		Description:   random.AlphaNumeric(rand.Intn(30) + 1),
		Link:          random.RandURLString(),
		SinceYear:     random.Time().Year(),
		SinceSemester: rand.Intn(2),
		UntilYear:     random.Time().Year(),
		UntilSemester: rand.Intn(2),
	} // Successで使うProject

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
		{
			name: "Success",
			args: args{
				project: successProject,
			},
			want: &domain.Project{
				ID:   successProject.ID,
				Name: successProject.Name,
				Duration: domain.YearWithSemesterDuration{
					Since: domain.YearWithSemester{
						Year:     successProject.SinceYear,
						Semester: successProject.SinceSemester,
					},
					Until: domain.YearWithSemester{
						Year:     successProject.UntilYear,
						Semester: successProject.UntilSemester,
					},
				},
				Description: successProject.Description,
				Link:        successProject.Link,
			},
			setup: func(f mockProjectRepositoryFields, args args, want *domain.Project) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				p := args.project
				h.Mock.
					ExpectExec(regexp.QuoteMeta("INSERT INTO `projects` (`id`,`name`,`description`,`link`,`since_year`,`since_semester`,`until_year`,`until_semester`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?,?,?)")).
					WithArgs(p.ID, p.Name, p.Description, p.Link, p.SinceYear, p.SinceSemester, p.UntilYear, p.UntilSemester, anyTime{}, anyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			args: args{
				project: &model.Project{
					ID:            random.UUID(),
					Name:          random.AlphaNumeric(rand.Intn(30) + 1),
					Description:   random.AlphaNumeric(rand.Intn(30) + 1),
					Link:          random.RandURLString(),
					SinceYear:     random.Time().Year(),
					SinceSemester: rand.Intn(2),
					UntilYear:     random.Time().Year(),
					UntilSemester: rand.Intn(2),
				},
			},
			want: nil,
			setup: func(f mockProjectRepositoryFields, args args, want *domain.Project) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				p := args.project
				h.Mock.
					ExpectExec(regexp.QuoteMeta("INSERT INTO `projects` (`id`,`name`,`description`,`link`,`since_year`,`since_semester`,`until_year`,`until_semester`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?,?,?)")).
					WithArgs(p.ID, p.Name, p.Description, p.Link, p.SinceYear, p.SinceSemester, p.UntilYear, p.UntilSemester, anyTime{}, anyTime{}).
					WillReturnError(errUnexpected)
				h.Mock.ExpectRollback()
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
		{
			name: "Success_AllParameterWillBeChanged",
			args: args{
				id: random.UUID(),
				changes: map[string]interface{}{
					"name":           random.AlphaNumeric(rand.Intn(30) + 1),
					"description":    random.AlphaNumeric(rand.Intn(30) + 1),
					"link":           random.RandURLString(),
					"since_year":     random.Time().Year(),
					"since_semester": rand.Intn(2),
					"until_year":     random.Time().Year(),
					"until_semester": rand.Intn(2),
				},
			},
			setup: func(f mockProjectRepositoryFields, args args) {
				var (
					name        = random.AlphaNumeric(rand.Intn(30) + 1)
					description = random.AlphaNumeric(rand.Intn(30) + 1)
					link        = random.RandURLString()
					sinceYear   = random.Time().Year()
					sinceSem    = rand.Intn(2)
					untilYear   = random.Time().Year()
					untilSem    = rand.Intn(2)
				)

				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since_year", "since_semester", "until_year", "until_semester"}).
							AddRow(args.id, name, description, link, sinceYear, sinceSem, untilYear, untilSem),
					)
				h.Mock.
					ExpectExec(regexp.QuoteMeta("UPDATE `projects` SET `description`=?,`link`=?,`name`=?,`since_semester`=?,`since_year`=?,`until_semester`=?,`until_year`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.changes["description"], args.changes["link"], args.changes["name"], args.changes["since_semester"], args.changes["since_year"], args.changes["until_semester"], args.changes["until_year"], anyTime{}, args.id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				name = args.changes["name"].(string)
				description = args.changes["description"].(string)
				link = args.changes["link"].(string)
				sinceYear = args.changes["since_year"].(int)
				sinceSem = args.changes["since_semester"].(int)
				untilYear = args.changes["until_year"].(int)
				untilSem = args.changes["until_semester"].(int)

				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since_year", "since_semester", "until_year", "until_semester"}).
							AddRow(args.id, name, description, link, sinceYear, sinceSem, untilYear, untilSem),
					)
				h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "ErrorFirstOldProject",
			args: args{
				id:      random.UUID(),
				changes: map[string]interface{}{},
			},
			setup: func(f mockProjectRepositoryFields, args args) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnError(errUnexpected)
				h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "ErrorUpdate",
			args: args{
				id: random.UUID(),
				changes: map[string]interface{}{
					"name":        random.AlphaNumeric(rand.Intn(30) + 1),
					"description": random.AlphaNumeric(rand.Intn(30) + 1),
					"link":        random.RandURLString(),
					"since":       time.Now(),
					"until":       time.Now(),
				},
			},
			setup: func(f mockProjectRepositoryFields, args args) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since_year", "since_semester", "until_year", "until_semester"}).
							AddRow(
								args.id,
								random.AlphaNumeric(rand.Intn(30)+1),
								random.AlphaNumeric(rand.Intn(30)+1),
								random.RandURLString(),
								random.Time().Year(),
								rand.Intn(2),
								random.Time().Year(),
								rand.Intn(2),
							),
					)
				h.Mock.
					ExpectExec(regexp.QuoteMeta("UPDATE `projects` SET `description`=?,`link`=?,`name`=?,`since`=?,`until`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.changes["description"], args.changes["link"], args.changes["name"], args.changes["since"], args.changes["until"], anyTime{}, args.id).
					WillReturnError(errUnexpected)
				h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "ErrorFirstNewProject",
			args: args{
				id: random.UUID(),
				changes: map[string]interface{}{
					"name":           random.AlphaNumeric(rand.Intn(30) + 1),
					"description":    random.AlphaNumeric(rand.Intn(30) + 1),
					"link":           random.RandURLString(),
					"since_year":     random.Time().Year(),
					"since_semester": rand.Intn(2),
					"until_year":     random.Time().Year(),
					"until_semester": rand.Intn(2),
				},
			},
			setup: func(f mockProjectRepositoryFields, args args) {
				var (
					name        = random.AlphaNumeric(rand.Intn(30) + 1)
					description = random.AlphaNumeric(rand.Intn(30) + 1)
					link        = random.RandURLString()
					sinceYear   = random.Time().Year()
					sinceSem    = rand.Intn(2)
					untilYear   = random.Time().Year()
					untilSem    = rand.Intn(2)
				)

				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since_year", "since_semester", "until_year", "until_semester"}).
							AddRow(args.id, name, description, link, sinceYear, sinceSem, untilYear, untilSem),
					)
				h.Mock.
					ExpectExec(regexp.QuoteMeta("UPDATE `projects` SET `description`=?,`link`=?,`name`=?,`since`=?,`until`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.changes["description"], args.changes["link"], args.changes["name"], args.changes["since"], args.changes["until"], anyTime{}, args.id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnError(errUnexpected)
				h.Mock.ExpectRollback()
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
		{
			name: "Success_Single",
			args: args{
				id: random.UUID(),
			},
			want: []*domain.User{
				{
					ID:   random.UUID(),
					Name: random.AlphaNumeric(rand.Intn(30) + 1),
					// RealName:
				},
			},
			setup: func(f mockProjectRepositoryFields, args args, want []*domain.User) {
				rows := sqlmock.NewRows([]string{"user_id"})
				for _, u := range want {
					rows.AddRow(u.ID)
				}
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `project_members` WHERE `project_members`.`project_id` = ?")).
					WithArgs(args.id).
					WillReturnRows(rows)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ?")).
					WithArgs(want[0].ID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name"}).
							AddRow(want[0].ID, want[0].Name),
					)
			},
			assertion: assert.NoError,
		},
		{
			name: "Success_Multiple",
			args: args{
				id: random.UUID(),
			},
			want: []*domain.User{
				{
					ID:   random.UUID(),
					Name: random.AlphaNumeric(rand.Intn(30) + 1),
					// RealName:
				},
				{
					ID:   random.UUID(),
					Name: random.AlphaNumeric(rand.Intn(30) + 1),
				},
			},
			setup: func(f mockProjectRepositoryFields, args args, want []*domain.User) {
				rows := sqlmock.NewRows([]string{"user_id"})
				for _, u := range want {
					rows.AddRow(u.ID)
				}
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `project_members` WHERE `project_members`.`project_id` = ?")).
					WithArgs(args.id).
					WillReturnRows(rows)
				userIDs := make([]driver.Value, len(want))
				userRows := sqlmock.NewRows([]string{"id", "name"})
				for i, v := range want {
					userIDs[i] = v.ID
					userRows.AddRow(v.ID, v.Name)
				}
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` IN (?,?)")).
					WithArgs(userIDs...).
					WillReturnRows(userRows)
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			args: args{
				id: random.UUID(),
			},
			want: nil,
			setup: func(f mockProjectRepositoryFields, args args, want []*domain.User) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `project_members` WHERE `project_members`.`project_id` = ?")).
					WithArgs(args.id).
					WillReturnError(errUnexpected)
				h.Mock.ExpectRollback()
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
		{
			name: "Success",
			args: args{
				projectID: random.UUID(),
				projectMembers: []*repository.CreateProjectMemberArgs{
					{
						UserID:        random.UUID(),
						SinceYear:     random.Time().Year(),
						SinceSemester: rand.Intn(2),
						UntilYear:     random.Time().Year(),
						UntilSemester: rand.Intn(2),
					},
					{
						UserID:        random.UUID(),
						SinceYear:     random.Time().Year(),
						SinceSemester: rand.Intn(2),
						UntilYear:     random.Time().Year(),
						UntilSemester: rand.Intn(2),
					},
					{
						UserID:        random.UUID(),
						SinceYear:     random.Time().Year(),
						SinceSemester: rand.Intn(2),
						UntilYear:     random.Time().Year(),
						UntilSemester: rand.Intn(2),
					},
					{
						UserID:        random.UUID(),
						SinceYear:     random.Time().Year(),
						SinceSemester: rand.Intn(2),
						UntilYear:     random.Time().Year(),
						UntilSemester: rand.Intn(2),
					},
				},
			},
			setup: func(f mockProjectRepositoryFields, args args) {
				rows := sqlmock.NewRows([]string{"team_id", "user_id"})
				newUsers := make([]*repository.CreateProjectMemberArgs, 0, len(args.projectMembers))
				for i, u := range args.projectMembers {
					if i%2 == 0 {
						rows.AddRow(args.projectID, u.UserID)
					} else {
						newUsers = append(newUsers, u)
					}
				}
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.projectID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(args.projectID),
					)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `project_members` WHERE `project_members`.`project_id` = ?")).
					WithArgs(args.projectID).
					WillReturnRows(rows)
				h.Mock.ExpectBegin()
				for _, u := range newUsers {
					h.Mock.
						ExpectExec(regexp.QuoteMeta("INSERT INTO `project_members` (`id`,`project_id`,`user_id`,`since_year`,`since_semester`,`until_year`,`until_semester`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?,?)")).
						WithArgs(anyUUID{}, args.projectID, u.UserID, u.SinceYear, u.SinceSemester, u.UntilYear, u.UntilSemester, anyTime{}, anyTime{}).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}
				h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "ZeroMembers",
			args: args{
				projectID:      random.UUID(),
				projectMembers: nil,
			},
			setup:     func(f mockProjectRepositoryFields, args args) {},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_FindProject",
			args: args{
				projectID: random.UUID(),
				projectMembers: []*repository.CreateProjectMemberArgs{
					{
						UserID:        random.UUID(),
						SinceYear:     random.Time().Year(),
						SinceSemester: rand.Intn(2),
						UntilYear:     random.Time().Year(),
						UntilSemester: rand.Intn(2),
					},
					{
						UserID:        random.UUID(),
						SinceYear:     random.Time().Year(),
						SinceSemester: rand.Intn(2),
						UntilYear:     random.Time().Year(),
						UntilSemester: rand.Intn(2),
					},
				},
			},
			setup: func(f mockProjectRepositoryFields, args args) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.projectID).
					WillReturnError(errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_FindProjectMembers",
			args: args{
				projectID: random.UUID(),
				projectMembers: []*repository.CreateProjectMemberArgs{
					{
						UserID:        random.UUID(),
						SinceYear:     random.Time().Year(),
						SinceSemester: rand.Intn(2),
						UntilYear:     random.Time().Year(),
						UntilSemester: rand.Intn(2),
					},
					{
						UserID:        random.UUID(),
						SinceYear:     random.Time().Year(),
						SinceSemester: rand.Intn(2),
						UntilYear:     random.Time().Year(),
						UntilSemester: rand.Intn(2),
					},
				},
			},
			setup: func(f mockProjectRepositoryFields, args args) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.projectID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(args.projectID),
					)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `project_members` WHERE `project_members`.`project_id` = ?")).
					WithArgs(args.projectID).
					WillReturnError(errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_AddProjectMembers",
			args: args{
				projectID: random.UUID(),
				projectMembers: []*repository.CreateProjectMemberArgs{
					{
						UserID:        random.UUID(),
						SinceYear:     random.Time().Year(),
						SinceSemester: rand.Intn(2),
						UntilYear:     random.Time().Year(),
						UntilSemester: rand.Intn(2),
					},
					{
						UserID:        random.UUID(),
						SinceYear:     random.Time().Year(),
						SinceSemester: rand.Intn(2),
						UntilYear:     random.Time().Year(),
						UntilSemester: rand.Intn(2),
					},
				},
			},
			setup: func(f mockProjectRepositoryFields, args args) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.projectID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(args.projectID),
					)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `project_members` WHERE `project_members`.`project_id` = ?")).
					WithArgs(args.projectID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(args.projectID),
					)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectExec(regexp.QuoteMeta("INSERT INTO `project_members` (`id`,`project_id`,`user_id`,`since_year`,`since_semester`,`until_year`,`until_semester`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, args.projectID, anyUUID{}, anyTime{}, anyTime{}, anyTime{}, anyTime{}).
					WillReturnError(errUnexpected)
				h.Mock.ExpectRollback()
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
		{
			name: "Success",
			args: args{
				projectID: random.UUID(),
				members: []uuid.UUID{
					random.UUID(), // 元のメンバーであるため削除対象となる
					random.UUID(), // 元のメンバーでないため削除対象とならない
				},
			},
			setup: func(f mockProjectRepositoryFields, args args) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.projectID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(args.projectID),
					)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectExec(regexp.QuoteMeta("DELETE FROM `project_members` WHERE `project_members`.`project_id` = ? AND user_id IN (?,?)")).
					WithArgs(args.projectID, args.members[0], args.members[1]).
					WillReturnResult(sqlmock.NewResult(0, int64(len(args.members)+1)/2))
				h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "ZeroMembers",
			args: args{
				projectID: random.UUID(),
				members:   []uuid.UUID{},
			},
			setup:     func(f mockProjectRepositoryFields, args args) {},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_FindProject",
			args: args{
				projectID: random.UUID(),
				members: []uuid.UUID{
					random.UUID(),
					random.UUID(),
				},
			},
			setup: func(f mockProjectRepositoryFields, args args) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.projectID).
					WillReturnError(errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_DeleteMembers",
			args: args{
				projectID: random.UUID(),
				members: []uuid.UUID{
					random.UUID(),
					random.UUID(),
				},
			},
			setup: func(f mockProjectRepositoryFields, args args) {
				h := f.h.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.projectID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(args.projectID),
					)
				h.Mock.ExpectBegin()
				h.Mock.
					ExpectExec(regexp.QuoteMeta("DELETE FROM `project_members` WHERE `project_members`.`project_id` = ? AND user_id IN (?,?)")).
					WithArgs(args.projectID, args.members[0], args.members[1]).
					WillReturnError(errUnexpected)
				h.Mock.ExpectRollback()
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
			tt.setup(f, tt.args)
			repo := NewProjectRepository(f.h)
			// Assertion
			tt.assertion(t, repo.DeleteProjectMembers(tt.args.projectID, tt.args.members))
		})
	}
}