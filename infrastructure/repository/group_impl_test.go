package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/util/random"
)

type mockGroupRepositoryFields struct {
	h *MockSQLHandler
}

func newMockGroupRepositoryFields(t *testing.T) mockGroupRepositoryFields {
	t.Helper()
	return mockGroupRepositoryFields{
		h: NewMockSQLHandler(),
	}
}

func TestGroupRepository_GetAllGroups(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		want      []*domain.Group
		setup     func(f mockGroupRepositoryFields, want []*domain.Group)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			want: []*domain.Group{
				{
					ID:   random.UUID(),
					Name: random.AlphaNumeric(),
				},
			},
			setup: func(f mockGroupRepositoryFields, want []*domain.Group) {
				rows := sqlmock.NewRows([]string{"group_id", "name"})
				for _, v := range want {
					rows.AddRow(v.ID, v.Name)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `groups`")).
					WillReturnRows(rows)
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			want: nil,
			setup: func(f mockGroupRepositoryFields, want []*domain.Group) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `groups`")).
					WillReturnError(errUnexpected)
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			f := newMockGroupRepositoryFields(t)
			tt.setup(f, tt.want)
			repo := NewGroupRepository(f.h.Conn)
			got, err := repo.GetGroups(context.Background())
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGroupRepository_GetGroup(t *testing.T) {
	gid := random.UUID()

	t.Parallel()
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.GroupDetail
		setup     func(f mockGroupRepositoryFields, args args, want *domain.GroupDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success_Singles",
			args: args{
				id: gid,
			},
			want: &domain.GroupDetail{
				ID:   gid,
				Name: random.AlphaNumeric(),
				Link: random.RandURLString(),
				Admin: []*domain.User{
					{
						ID: random.UUID(),
						// usecasesで後付けしているのでここでは不要
						// Name:     random.AlphaNumeric(),
						// RealName: random.AlphaNumeric(),
					},
				},
				Members: []*domain.UserWithDuration{
					{
						User: domain.User{
							ID: random.UUID(),
							// Name:     random.AlphaNumeric(),
							// RealName: random.AlphaNumeric(),
						},
						Duration: random.Duration(),
					},
				},
				Description: random.AlphaNumeric(),
			},
			setup: func(f mockGroupRepositoryFields, args args, want *domain.GroupDetail) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `groups` WHERE `groups`.`group_id` = ? ORDER BY `groups`.`group_id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"group_id", "name", "link", "description"}).
							AddRow(want.ID, want.Name, want.Link, want.Description),
					)
				wm := want.Members[0]
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `group_user_belongings` WHERE `group_user_belongings`.`group_id` = ?")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"user_id", "group_id", "since_year", "since_semester", "until_year", "until_semester"}).
							AddRow(wm.User.ID, want.ID, wm.Duration.Since.Year, wm.Duration.Since.Semester, wm.Duration.Until.Year, wm.Duration.Until.Semester),
					)
				wad := want.Admin[0]
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `group_user_admins` WHERE `group_user_admins`.`group_id` = ?")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"user_id", "group_id"}).
							AddRow(wad.ID, want.ID),
					)
			},
			assertion: assert.NoError,
		},
		{
			name: "Success_Multiples",
			args: args{
				id: gid,
			},
			want: &domain.GroupDetail{
				ID:   gid,
				Name: random.AlphaNumeric(),
				Link: random.RandURLString(),
				Admin: []*domain.User{
					{
						ID: random.UUID(),
					},
					{
						ID: random.UUID(),
					},
					{
						ID: random.UUID(),
					},
				},
				Members: []*domain.UserWithDuration{
					{
						User: domain.User{
							ID: random.UUID(),
						},
						Duration: random.Duration(),
					},
					{
						User: domain.User{
							ID: random.UUID(),
						},
						Duration: random.Duration(),
					},
					{
						User: domain.User{
							ID: random.UUID(),
						},
						Duration: random.Duration(),
					},
				},
				Description: random.AlphaNumeric(),
			},
			setup: func(f mockGroupRepositoryFields, args args, want *domain.GroupDetail) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `groups` WHERE `groups`.`group_id` = ? ORDER BY `groups`.`group_id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"group_id", "name", "link", "description"}).
							AddRow(want.ID, want.Name, want.Link, want.Description),
					)
				memberRows := sqlmock.NewRows([]string{"user_id", "group_id", "since_year", "since_semester", "until_year", "until_semester"})
				for _, m := range want.Members {
					memberRows.AddRow(m.User.ID, want.ID, m.Duration.Since.Year, m.Duration.Since.Semester, m.Duration.Until.Year, m.Duration.Until.Semester)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `group_user_belongings` WHERE `group_user_belongings`.`group_id` = ?")).
					WithArgs(args.id).
					WillReturnRows(memberRows)
				adminRows := sqlmock.NewRows([]string{"user_id", "group_id"})
				for _, a := range want.Admin {
					adminRows.AddRow(a.ID, want.ID)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `group_user_admins` WHERE `group_user_admins`.`group_id` = ?")).
					WithArgs(args.id).
					WillReturnRows(adminRows)
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			args: args{
				id: gid,
			},
			want: nil,
			setup: func(f mockGroupRepositoryFields, args args, want *domain.GroupDetail) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `groups` WHERE `groups`.`group_id` = ? ORDER BY `groups`.`group_id` LIMIT 1")).
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
			f := newMockGroupRepositoryFields(t)
			tt.setup(f, tt.args, tt.want)
			repo := NewGroupRepository(f.h.Conn)
			got, err := repo.GetGroup(context.Background(), tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
