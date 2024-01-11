package repository

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/infrastructure/external"
	"github.com/traPtitech/traPortfolio/infrastructure/external/mock_external"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/optional"
	"github.com/traPtitech/traPortfolio/util/random"
)

type mockProjectRepositoryFields struct {
	h      *MockSQLHandler
	portal *mock_external.MockPortalAPI
}

func newMockProjectRepositoryFields(t *testing.T, ctrl *gomock.Controller) mockProjectRepositoryFields {
	t.Helper()
	return mockProjectRepositoryFields{
		h:      NewMockSQLHandler(),
		portal: mock_external.NewMockPortalAPI(ctrl),
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
					ID:       random.UUID(),
					Name:     random.AlphaNumeric(),
					Duration: random.Duration(),
				},
			},
			setup: func(f mockProjectRepositoryFields, want []*domain.Project) {
				rows := sqlmock.NewRows([]string{"id", "name", "since_year", "since_semester", "until_year", "until_semester"})
				for _, v := range want {
					d := v.Duration
					rows.AddRow(v.ID, v.Name, d.Since.Year, d.Since.Semester, d.Until.Year, d.Until.Semester)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `projects`")).
					WillReturnRows(rows)
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			want: nil,
			setup: func(f mockProjectRepositoryFields, want []*domain.Project) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `projects`")).
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
			ctrl := gomock.NewController(t)
			f := newMockProjectRepositoryFields(t, ctrl)
			tt.setup(f, tt.want)
			repo := NewProjectRepository(f.h.Conn, f.portal)
			// Assertion
			got, err := repo.GetProjects(context.Background())
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
		want      *domain.ProjectDetail
		setup     func(f mockProjectRepositoryFields, args args, want *domain.ProjectDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success_single",
			args: args{
				id: pid,
			},
			want: &domain.ProjectDetail{
				Project: domain.Project{
					ID:       pid,
					Name:     random.AlphaNumeric(),
					Duration: random.Duration(),
				},
				Description: random.AlphaNumeric(),
				Link:        random.RandURLString(),
				Members: []*domain.UserWithDuration{
					{
						User:     *domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), true),
						Duration: random.Duration(),
					},
				},
			},
			setup: func(f mockProjectRepositoryFields, args args, want *domain.ProjectDetail) {
				wd := want.Duration
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since_year", "since_semester", "until_year", "until_semester"}).
							AddRow(want.ID, want.Name, want.Description, want.Link, wd.Since.Year, wd.Since.Semester, wd.Until.Year, wd.Until.Semester),
					)
				wm := want.Members[0]
				wmd := wm.Duration
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `project_members` WHERE `project_members`.`project_id` = ?")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"user_id", "name", "since_year", "since_semester", "until_year", "until_semester"}).
							AddRow(wm.User.ID, wm.User.Name, wmd.Since.Year, wmd.Since.Semester, wmd.Until.Year, wmd.Until.Semester),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ?")).
					WithArgs(wm.User.ID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "check"}).
							AddRow(wm.User.ID, wm.User.Name, wm.User.Check),
					)
				f.portal.EXPECT().GetUsers().Return([]*external.PortalUserResponse{
					{
						TraQID:   wm.User.Name,
						RealName: wm.User.RealName(),
					},
				}, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Success_Multiple",
			args: args{
				id: pid,
			},
			want: &domain.ProjectDetail{
				Project: domain.Project{
					ID:       pid,
					Name:     random.AlphaNumeric(),
					Duration: random.Duration(),
				},
				Description: random.AlphaNumeric(),
				Link:        random.RandURLString(),
				Members: []*domain.UserWithDuration{
					{
						User:     *domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), true),
						Duration: random.Duration(),
					},
					{
						User:     *domain.NewUser(random.UUID(), random.AlphaNumeric(), "", false),
						Duration: random.Duration(),
					},
				},
			},
			setup: func(f mockProjectRepositoryFields, args args, want *domain.ProjectDetail) {
				wd := want.Duration
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since_year", "since_semester", "until_year", "until_semester"}).
							AddRow(want.ID, want.Name, want.Description, want.Link, wd.Since.Year, wd.Since.Semester, wd.Until.Year, wd.Until.Semester),
					)
				memberRows := sqlmock.NewRows([]string{"user_id", "name", "since_year", "since_semester", "until_year", "until_semester"})
				for _, v := range want.Members {
					d := v.Duration
					memberRows.AddRow(v.User.ID, v.User.Name, d.Since.Year, d.Since.Semester, d.Until.Year, d.Until.Semester)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `project_members` WHERE `project_members`.`project_id` = ?")).
					WithArgs(args.id).
					WillReturnRows(memberRows)
				userIDs := make([]driver.Value, len(want.Members))
				userRows := sqlmock.NewRows([]string{"id", "name", "check"})
				for i, v := range want.Members {
					userIDs[i] = v.User.ID
					userRows.AddRow(v.User.ID, v.User.Name, v.User.Check)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` IN (?,?)")).
					WithArgs(userIDs...).
					WillReturnRows(userRows)
				wp := make([]*external.PortalUserResponse, len(want.Members))
				for i, v := range want.Members {
					wp[i] = &external.PortalUserResponse{
						TraQID:   v.User.Name,
						RealName: v.User.RealName(),
					}
				}
				f.portal.EXPECT().GetUsers().Return(wp, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "ProjectNotFound",
			args: args{
				id: random.UUID(),
			},
			want: nil,
			setup: func(f mockProjectRepositoryFields, args args, want *domain.ProjectDetail) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
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
			setup: func(f mockProjectRepositoryFields, args args, want *domain.ProjectDetail) {
				d := random.Duration()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since_year", "since_semester", "until_year", "until_semester"}).
							AddRow(
								args.id,
								random.AlphaNumeric(),
								random.AlphaNumeric(),
								random.RandURLString(),
								d.Since.Year,
								d.Since.Semester,
								d.Until.Year,
								d.Until.Semester,
							),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `project_members` WHERE `project_members`.`project_id` = ?")).
					WithArgs(args.id).
					WillReturnError(errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "PortalError",
			args: args{
				id: random.UUID(),
			},
			want: nil,
			setup: func(f mockProjectRepositoryFields, args args, want *domain.ProjectDetail) {
				d := random.Duration()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "description", "link", "since_year", "since_semester", "until_year", "until_semester"}).
							AddRow(
								args.id,
								random.AlphaNumeric(),
								random.AlphaNumeric(),
								random.RandURLString(),
								d.Since.Year,
								d.Since.Semester,
								d.Until.Year,
								d.Until.Semester,
							),
					)
				uid := random.UUID()
				md := random.Duration()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `project_members` WHERE `project_members`.`project_id` = ?")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"user_id", "name", "since_year", "since_semester", "until_year", "until_semester"}).
							AddRow(
								uid,
								random.AlphaNumeric(),
								md.Since.Year,
								md.Since.Semester,
								md.Until.Year,
								md.Until.Semester,
							),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ?")).
					WithArgs(uid).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name"}).
							AddRow(
								uid,
								random.AlphaNumeric(),
							),
					)
				f.portal.EXPECT().GetUsers().Return(nil, errUnexpected)
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			f := newMockProjectRepositoryFields(t, ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := NewProjectRepository(f.h.Conn, f.portal)
			// Assertion
			got, err := repo.GetProject(context.Background(), tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProjectRepository_CreateProject(t *testing.T) {
	duration := random.Duration()
	successProject := &repository.CreateProjectArgs{
		Name:          random.AlphaNumeric(),
		Description:   random.AlphaNumeric(),
		Link:          optional.From(random.RandURLString()),
		SinceYear:     duration.Since.Year,
		SinceSemester: duration.Since.Semester,
		UntilYear:     duration.Until.Year,
		UntilSemester: duration.Until.Semester,
	} // Successで使うProject

	t.Parallel()
	type args struct {
		project *repository.CreateProjectArgs
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.ProjectDetail
		setup     func(f mockProjectRepositoryFields, args args, want *domain.ProjectDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				project: successProject,
			},
			want: &domain.ProjectDetail{
				Project: domain.Project{
					// ID: 比較しない
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
				},
				Description: successProject.Description,
				Link:        successProject.Link.ValueOrZero(),
			},
			setup: func(f mockProjectRepositoryFields, args args, want *domain.ProjectDetail) {
				f.h.Mock.ExpectBegin()
				p := args.project
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("INSERT INTO `projects` (`id`,`name`,`description`,`link`,`since_year`,`since_semester`,`until_year`,`until_semester`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, p.Name, p.Description, p.Link, p.SinceYear, p.SinceSemester, p.UntilYear, p.UntilSemester, anyTime{}, anyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			args: args{
				project: &repository.CreateProjectArgs{
					Name:          random.AlphaNumeric(),
					Description:   random.AlphaNumeric(),
					Link:          random.Optional(random.RandURLString()),
					SinceYear:     duration.Since.Year,
					SinceSemester: duration.Since.Semester,
					UntilYear:     duration.Until.Year,
					UntilSemester: duration.Until.Semester,
				},
			},
			want: nil,
			setup: func(f mockProjectRepositoryFields, args args, want *domain.ProjectDetail) {
				f.h.Mock.ExpectBegin()
				p := args.project
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("INSERT INTO `projects` (`id`,`name`,`description`,`link`,`since_year`,`since_semester`,`until_year`,`until_semester`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, p.Name, p.Description, p.Link, p.SinceYear, p.SinceSemester, p.UntilYear, p.UntilSemester, anyTime{}, anyTime{}).
					WillReturnError(errUnexpected)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			f := newMockProjectRepositoryFields(t, ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := NewProjectRepository(f.h.Conn, f.portal)
			// Assertion
			got, err := repo.CreateProject(context.Background(), tt.args.project)
			if tt.want != nil && got != nil {
				tt.want.ID = got.ID // 関数内でIDを生成するためここで合わせる
			}
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProjectRepository_UpdateProject(t *testing.T) {
	d := random.Duration()

	t.Parallel()
	type args struct {
		id   uuid.UUID
		args *repository.UpdateProjectArgs
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
				args: &repository.UpdateProjectArgs{
					Name:          optional.From(random.AlphaNumeric()),
					Description:   optional.From(random.AlphaNumeric()),
					Link:          optional.From(random.RandURLString()),
					SinceYear:     optional.From(int64(d.Since.Year)),
					SinceSemester: optional.From(int64(d.Since.Semester)),
					UntilYear:     optional.From(int64(d.Until.Year)),
					UntilSemester: optional.From(int64(d.Until.Semester)),
				},
			},
			setup: func(f mockProjectRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("UPDATE `projects` SET `description`=?,`link`=?,`name`=?,`since_semester`=?,`since_year`=?,`until_semester`=?,`until_year`=?,`updated_at`=? WHERE `projects`.`id` = ?")).
					WithArgs(args.args.Description.ValueOrZero(), args.args.Link.ValueOrZero(), args.args.Name.ValueOrZero(), args.args.SinceSemester.ValueOrZero(), args.args.SinceYear.ValueOrZero(), args.args.UntilSemester.ValueOrZero(), args.args.UntilYear.ValueOrZero(), anyTime{}, args.id).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "ErrorUpdate",
			args: args{
				id: random.UUID(),
				args: &repository.UpdateProjectArgs{
					Name:          optional.From(random.AlphaNumeric()),
					Description:   optional.From(random.AlphaNumeric()),
					Link:          optional.From(random.RandURLString()),
					SinceYear:     optional.From(int64(d.Since.Year)),
					SinceSemester: optional.From(int64(d.Since.Semester)),
					UntilYear:     optional.From(int64(d.Until.Year)),
					UntilSemester: optional.From(int64(d.Until.Semester)),
				},
			},
			setup: func(f mockProjectRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("UPDATE `projects` SET `description`=?,`link`=?,`name`=?,`since_semester`=?,`since_year`=?,`until_semester`=?,`until_year`=?,`updated_at`=? WHERE `projects`.`id` = ?")).
					WithArgs(args.args.Description.ValueOrZero(), args.args.Link.ValueOrZero(), args.args.Name.ValueOrZero(), args.args.SinceSemester.ValueOrZero(), args.args.SinceYear.ValueOrZero(), args.args.UntilSemester.ValueOrZero(), args.args.UntilYear.ValueOrZero(), anyTime{}, args.id).
					WillReturnError(errUnexpected)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			f := newMockProjectRepositoryFields(t, ctrl)
			tt.setup(f, tt.args)
			repo := NewProjectRepository(f.h.Conn, f.portal)
			// Assertion
			tt.assertion(t, repo.UpdateProject(context.Background(), tt.args.id, tt.args.args))
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
		want      []*domain.UserWithDuration
		setup     func(f mockProjectRepositoryFields, args args, want []*domain.UserWithDuration)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success_Single",
			args: args{
				id: random.UUID(),
			},
			want: []*domain.UserWithDuration{
				{
					User: *domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), true),
				},
			},
			setup: func(f mockProjectRepositoryFields, args args, want []*domain.UserWithDuration) {
				rows := sqlmock.NewRows([]string{"user_id"})
				for _, pm := range want {
					rows.AddRow(pm.User.ID)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `project_members` WHERE `project_members`.`project_id` = ?")).
					WithArgs(args.id).
					WillReturnRows(rows)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ?")).
					WithArgs(want[0].User.ID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "check"}).
							AddRow(want[0].User.ID, want[0].User.Name, want[0].User.Check),
					)
				f.portal.EXPECT().GetUsers().Return([]*external.PortalUserResponse{
					{
						TraQID:   want[0].User.Name,
						RealName: want[0].User.RealName(),
					},
				}, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Success_Multiple",
			args: args{
				id: random.UUID(),
			},
			want: []*domain.UserWithDuration{
				{
					User: domain.User{
						ID:   random.UUID(),
						Name: random.AlphaNumeric(),
						// RealName:
					},
				},
				{
					User: domain.User{
						ID:   random.UUID(),
						Name: random.AlphaNumeric(),
						// RealName:
					},
				},
			},
			setup: func(f mockProjectRepositoryFields, args args, want []*domain.UserWithDuration) {
				rows := sqlmock.NewRows([]string{"user_id"})
				for _, pm := range want {
					rows.AddRow(pm.User.ID)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `project_members` WHERE `project_members`.`project_id` = ?")).
					WithArgs(args.id).
					WillReturnRows(rows)
				userIDs := make([]driver.Value, len(want))
				userRows := sqlmock.NewRows([]string{"id", "name"})
				for i, v := range want {
					userIDs[i] = v.User.ID
					userRows.AddRow(v.User.ID, v.User.Name)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` IN (?,?)")).
					WithArgs(userIDs...).
					WillReturnRows(userRows)
				wp := make([]*external.PortalUserResponse, len(want))
				for i, v := range want {
					wp[i] = &external.PortalUserResponse{
						TraQID:   v.User.Name,
						RealName: v.User.RealName(),
					}
				}
				f.portal.EXPECT().GetUsers().Return(wp, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			args: args{
				id: random.UUID(),
			},
			want: nil,
			setup: func(f mockProjectRepositoryFields, args args, want []*domain.UserWithDuration) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `project_members` WHERE `project_members`.`project_id` = ?")).
					WithArgs(args.id).
					WillReturnError(errUnexpected)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "PortalError",
			args: args{
				id: random.UUID(),
			},
			want: nil,
			setup: func(f mockProjectRepositoryFields, args args, want []*domain.UserWithDuration) {
				uid := random.UUID()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `project_members` WHERE `project_members`.`project_id` = ?")).
					WithArgs(args.id).
					WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(uid))
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ?")).
					WithArgs(uid).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name"}).
							AddRow(
								uid,
								random.AlphaNumeric(),
							),
					)
				f.portal.EXPECT().GetUsers().Return(nil, errUnexpected)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			f := newMockProjectRepositoryFields(t, ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := NewProjectRepository(f.h.Conn, f.portal)
			// Assertion
			got, err := repo.GetProjectMembers(context.Background(), tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProjectRepository_AddProjectMembers(t *testing.T) {
	duration := random.Duration()
	duplicatedMemberID := random.UUID()

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
						SinceYear:     duration.Since.Year,
						SinceSemester: duration.Since.Semester,
						UntilYear:     duration.Until.Year,
						UntilSemester: duration.Until.Semester,
					},
					{
						UserID:        random.UUID(),
						SinceYear:     duration.Since.Year,
						SinceSemester: duration.Since.Semester,
						UntilYear:     duration.Until.Year,
						UntilSemester: duration.Until.Semester,
					},
					{
						UserID:        random.UUID(),
						SinceYear:     duration.Since.Year,
						SinceSemester: duration.Since.Semester,
						UntilYear:     duration.Until.Year,
						UntilSemester: duration.Until.Semester,
					},
					{
						UserID:        random.UUID(),
						SinceYear:     duration.Since.Year,
						SinceSemester: duration.Since.Semester,
						UntilYear:     duration.Until.Year,
						UntilSemester: duration.Until.Semester,
					},
				},
			},
			setup: func(f mockProjectRepositoryFields, args args) {
				rows := sqlmock.NewRows([]string{"team_id", "user_id"})
				newUsers := make([]*repository.CreateProjectMemberArgs, 0, len(args.projectMembers))
				for i, pm := range args.projectMembers {
					if i%2 == 0 {
						rows.AddRow(args.projectID, pm.UserID)
					} else {
						newUsers = append(newUsers, pm)
					}
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.projectID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "since_year", "since_semester", "until_year", "until_semester"}).
							AddRow(args.projectID, duration.Since.Year, duration.Since.Semester, duration.Until.Year, duration.Until.Semester),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `project_members` WHERE `project_members`.`project_id` = ?")).
					WithArgs(args.projectID).
					WillReturnRows(rows)
				f.h.Mock.ExpectBegin()
				for _, pm := range newUsers {
					f.h.Mock.
						ExpectExec(makeSQLQueryRegexp("INSERT INTO `project_members` (`id`,`project_id`,`user_id`,`since_year`,`since_semester`,`until_year`,`until_semester`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?,?)")).
						WithArgs(anyUUID{}, args.projectID, pm.UserID, pm.SinceYear, pm.SinceSemester, pm.UntilYear, pm.UntilSemester, anyTime{}, anyTime{}).
						WillReturnResult(sqlmock.NewResult(1, 1))
				}
				f.h.Mock.ExpectCommit()
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
			name: "duplicatedMembers",
			args: args{
				projectID: random.UUID(),
				projectMembers: []*repository.CreateProjectMemberArgs{
					{
						UserID:        duplicatedMemberID,
						SinceYear:     duration.Since.Year,
						SinceSemester: duration.Since.Semester,
						UntilYear:     duration.Until.Year,
						UntilSemester: duration.Until.Semester,
					},
					{
						UserID:        duplicatedMemberID,
						SinceYear:     duration.Since.Year,
						SinceSemester: duration.Since.Semester,
						UntilYear:     duration.Until.Year,
						UntilSemester: duration.Until.Semester,
					},
				},
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
						SinceYear:     duration.Since.Year,
						SinceSemester: duration.Since.Semester,
						UntilYear:     duration.Until.Year,
						UntilSemester: duration.Until.Semester,
					},
					{
						UserID:        random.UUID(),
						SinceYear:     duration.Since.Year,
						SinceSemester: duration.Since.Semester,
						UntilYear:     duration.Until.Year,
						UntilSemester: duration.Until.Semester,
					},
				},
			},
			setup: func(f mockProjectRepositoryFields, args args) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.projectID).
					WillReturnError(errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_DurationExceedMember",
			args: args{
				projectID: random.UUID(),
				// project duration is 2020-0 ~ 2021-1
				projectMembers: []*repository.CreateProjectMemberArgs{
					{
						UserID:        random.UUID(), // OK
						SinceYear:     2020,
						SinceSemester: 0,
						UntilYear:     2021,
						UntilSemester: 1,
					},
					{
						UserID:        random.UUID(), // NG
						SinceYear:     2021,
						SinceSemester: 0,
						UntilYear:     2022,
						UntilSemester: 0,
					},
				},
			},
			setup: func(f mockProjectRepositoryFields, args args) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.projectID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "since_year", "since_semester", "until_year", "until_semester"}).
							AddRow(args.projectID, 2020, 0, 2021, 1),
					)
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
						SinceYear:     duration.Since.Year,
						SinceSemester: duration.Since.Semester,
						UntilYear:     duration.Until.Year,
						UntilSemester: duration.Until.Semester,
					},
					{
						UserID:        random.UUID(),
						SinceYear:     duration.Since.Year,
						SinceSemester: duration.Since.Semester,
						UntilYear:     duration.Until.Year,
						UntilSemester: duration.Until.Semester,
					},
				},
			},
			setup: func(f mockProjectRepositoryFields, args args) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.projectID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(args.projectID),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `project_members` WHERE `project_members`.`project_id` = ?")).
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
						SinceYear:     duration.Since.Year,
						SinceSemester: duration.Since.Semester,
						UntilYear:     duration.Until.Year,
						UntilSemester: duration.Until.Semester,
					},
					{
						UserID:        random.UUID(),
						SinceYear:     duration.Since.Year,
						SinceSemester: duration.Since.Semester,
						UntilYear:     duration.Until.Year,
						UntilSemester: duration.Until.Semester,
					},
				},
			},
			setup: func(f mockProjectRepositoryFields, args args) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.projectID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(args.projectID),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `project_members` WHERE `project_members`.`project_id` = ?")).
					WithArgs(args.projectID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(args.projectID),
					)
				f.h.Mock.ExpectBegin()
				pm := args.projectMembers[0]
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("INSERT INTO `project_members` (`id`,`project_id`,`user_id`,`since_year`,`since_semester`,`until_year`,`until_semester`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, args.projectID, pm.UserID, pm.SinceYear, pm.SinceSemester, pm.UntilYear, pm.UntilSemester, anyTime{}, anyTime{}).
					WillReturnError(errUnexpected)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			f := newMockProjectRepositoryFields(t, ctrl)
			tt.setup(f, tt.args)
			repo := NewProjectRepository(f.h.Conn, f.portal)
			// Assertion
			tt.assertion(t, repo.AddProjectMembers(context.Background(), tt.args.projectID, tt.args.projectMembers))
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
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.projectID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(args.projectID),
					)
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("DELETE FROM `project_members` WHERE `project_members`.`project_id` = ? AND `project_members`.`user_id` IN (?,?)")).
					WithArgs(args.projectID, args.members[0], args.members[1]).
					WillReturnResult(sqlmock.NewResult(0, int64(len(args.members)+1)/2))
				f.h.Mock.ExpectCommit()
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
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
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
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `projects` WHERE `projects`.`id` = ? ORDER BY `projects`.`id` LIMIT 1")).
					WithArgs(args.projectID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(args.projectID),
					)
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("DELETE FROM `project_members` WHERE `project_members`.`project_id` = ? AND `project_members`.`user_id` IN (?,?)")).
					WithArgs(args.projectID, args.members[0], args.members[1]).
					WillReturnError(errUnexpected)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			f := newMockProjectRepositoryFields(t, ctrl)
			tt.setup(f, tt.args)
			repo := NewProjectRepository(f.h.Conn, f.portal)
			// Assertion
			tt.assertion(t, repo.DeleteProjectMembers(context.Background(), tt.args.projectID, tt.args.members))
		})
	}
}
