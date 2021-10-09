package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/database/mock_database"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external"
	"github.com/traPtitech/traPortfolio/util"
)

var (
	ids = []uuid.UUID{
		uuid.FromStringOrNil("11111111-1111-1111-1111-111111111111"),
		uuid.FromStringOrNil("22222222-2222-2222-2222-222222222222"),
		uuid.FromStringOrNil("33333333-3333-3333-3333-333333333333"),
	}
)

func TestUserRepository_GetUsers(t *testing.T) {
	t.Parallel()
	type fields struct {
		sqlhandler database.SQLHandler
		portal     external.PortalAPI
		traq       external.TraQAPI
	}
	tests := []struct {
		name      string
		fields    fields
		isValidDB bool
		want      []*domain.User
		setup     func(f fields, want []*domain.User)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Success",
			fields:    fields{},
			isValidDB: true,
			want: []*domain.User{
				{
					ID:       ids[0],
					Name:     "user1",
					RealName: "ユーザー1 ユーザー1",
				},
				{
					ID:       ids[1],
					Name:     "user2",
					RealName: "ユーザー2 ユーザー2",
				},
				{
					ID:       ids[2],
					Name:     "lolico",
					RealName: "東 工子",
				},
			},
			setup: func(f fields, want []*domain.User) {
				rows := sqlmock.NewRows([]string{"id", "name"})
				for _, v := range want {
					rows.AddRow(v.ID, v.Name)
				}
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
					WillReturnRows(rows)
			},
			assertion: assert.NoError,
		},
		{
			name:      "InvalidDB",
			fields:    fields{},
			isValidDB: false,
			want:      nil,
			setup:     func(f fields, want []*domain.User) {},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			tt.fields = fields{
				sqlhandler: mock_database.NewMockSQLHandler(tt.isValidDB),
				portal:     mock_external.NewMockPortalAPI(),
				traq:       mock_external.NewMockTraQAPI(),
			}
			tt.setup(tt.fields, tt.want)
			repo := NewUserRepository(tt.fields.sqlhandler, tt.fields.portal, tt.fields.traq)
			// Assertion
			got, err := repo.GetUsers()
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_GetUser(t *testing.T) {
	t.Parallel()
	type fields struct {
		sqlhandler database.SQLHandler
		portal     external.PortalAPI
		traq       external.TraQAPI
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		isValidDB bool
		want      *domain.UserDetail
		setup     func(f fields, args args, want *domain.UserDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Success",
			fields:    fields{},
			args:      args{ids[0]},
			isValidDB: true,
			want: &domain.UserDetail{
				User: domain.User{
					ID:       ids[0],
					Name:     "user1",
					RealName: "ユーザー1 ユーザー1",
				},
				State: 1,
				Bio:   util.AlphaNumeric(5),
				Accounts: []*domain.Account{
					{
						ID:          util.UUID(),
						Type:        domain.HOMEPAGE,
						PrPermitted: true,
					},
				},
			},
			setup: func(f fields, args args, want *domain.UserDetail) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "name", "description"}).
							AddRow(want.User.ID, want.User.Name, want.Bio),
					)
				rows := sqlmock.NewRows([]string{"id", "user_id", "type", "check"})
				for _, v := range want.Accounts {
					rows.AddRow(v.ID, want.User.ID, v.Type, v.PrPermitted)
				}
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `accounts` WHERE `accounts`.`user_id` = ?")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "type", "check"}).
							AddRow(want.Accounts[0].ID, args.id, want.Accounts[0].Type, want.Accounts[0].PrPermitted),
					)
			},
			assertion: assert.NoError,
		},
		{
			name:      "NotFound",
			fields:    fields{},
			args:      args{ids[0]},
			isValidDB: true,
			want:      nil,
			setup: func(f fields, args args, want *domain.UserDetail) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}))
			},
			assertion: assert.Error,
		},
		{
			name:      "InvalidDB",
			fields:    fields{},
			args:      args{ids[0]},
			isValidDB: false,
			want:      nil,
			setup:     func(f fields, args args, want *domain.UserDetail) {},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			tt.fields = fields{
				sqlhandler: mock_database.NewMockSQLHandler(tt.isValidDB),
				portal:     mock_external.NewMockPortalAPI(),
				traq:       mock_external.NewMockTraQAPI(),
			}
			tt.setup(tt.fields, tt.args, tt.want)
			repo := NewUserRepository(tt.fields.sqlhandler, tt.fields.portal, tt.fields.traq)
			// Assertion
			got, err := repo.GetUser(tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
