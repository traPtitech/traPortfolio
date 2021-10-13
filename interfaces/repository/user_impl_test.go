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
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util"
)

const isValidDB = true

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
		want      []*domain.User
		setup     func(f fields, want []*domain.User)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			fields: fields{
				sqlhandler: mock_database.NewMockSQLHandler(isValidDB),
				portal:     mock_external.NewMockPortalAPI(),
				traq:       mock_external.NewMockTraQAPI(),
			},
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
			name: "InvalidDB",
			fields: fields{
				sqlhandler: mock_database.NewMockSQLHandler(!isValidDB),
				portal:     mock_external.NewMockPortalAPI(),
				traq:       mock_external.NewMockTraQAPI(),
			},
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
		want      *domain.UserDetail
		setup     func(f fields, args args, want *domain.UserDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			fields: fields{
				sqlhandler: mock_database.NewMockSQLHandler(isValidDB),
				portal:     mock_external.NewMockPortalAPI(),
				traq:       mock_external.NewMockTraQAPI(),
			},
			args: args{ids[0]},
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
			name: "NotFound",
			fields: fields{
				sqlhandler: mock_database.NewMockSQLHandler(isValidDB),
				portal:     mock_external.NewMockPortalAPI(),
				traq:       mock_external.NewMockTraQAPI(),
			},
			args: args{ids[0]},
			want: nil,
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
			name: "InvalidDB",
			fields: fields{
				sqlhandler: mock_database.NewMockSQLHandler(!isValidDB),
				portal:     mock_external.NewMockPortalAPI(),
				traq:       mock_external.NewMockTraQAPI(),
			},
			args:      args{ids[0]},
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
			tt.setup(tt.fields, tt.args, tt.want)
			repo := NewUserRepository(tt.fields.sqlhandler, tt.fields.portal, tt.fields.traq)
			// Assertion
			got, err := repo.GetUser(tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_GetAccounts(t *testing.T) {
	t.Parallel()
	type fields struct {
		sqlhandler database.SQLHandler
		portal     external.PortalAPI
		traq       external.TraQAPI
	}
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []*domain.Account
		setup     func(f fields, args args, want []*domain.Account)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			fields: fields{
				sqlhandler: mock_database.NewMockSQLHandler(isValidDB),
				portal:     mock_external.NewMockPortalAPI(),
				traq:       mock_external.NewMockTraQAPI(),
			},
			args: args{ids[0]},
			want: []*domain.Account{
				{
					ID:          util.UUID(),
					Type:        domain.HOMEPAGE,
					PrPermitted: true,
				},
			},
			setup: func(f fields, args args, want []*domain.Account) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "type", "check"})
				for _, v := range want {
					rows.AddRow(v.ID, args.userID, v.Type, v.PrPermitted)
				}
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `accounts` WHERE user_id = ?")).
					WithArgs(args.userID).
					WillReturnRows(rows)
			},
			assertion: assert.NoError,
		},
		{
			name: "InvalidDB",
			fields: fields{
				sqlhandler: mock_database.NewMockSQLHandler(!isValidDB),
				portal:     mock_external.NewMockPortalAPI(),
				traq:       mock_external.NewMockTraQAPI(),
			},
			args:      args{ids[0]},
			want:      nil,
			setup:     func(f fields, args args, want []*domain.Account) {},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			tt.setup(tt.fields, tt.args, tt.want)
			repo := NewUserRepository(tt.fields.sqlhandler, tt.fields.portal, tt.fields.traq)
			// Assertion
			got, err := repo.GetAccounts(tt.args.userID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_GetAccount(t *testing.T) {
	t.Parallel()
	type fields struct {
		sqlhandler database.SQLHandler
		portal     external.PortalAPI
		traq       external.TraQAPI
	}
	type args struct {
		userID    uuid.UUID
		accountID uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *domain.Account
		setup     func(f fields, args args, want *domain.Account)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			fields: fields{
				sqlhandler: mock_database.NewMockSQLHandler(isValidDB),
				portal:     mock_external.NewMockPortalAPI(),
				traq:       mock_external.NewMockTraQAPI(),
			},
			args: args{
				userID:    ids[0],
				accountID: util.UUID(),
			},
			want: &domain.Account{
				ID:          uuid.Nil, // setupで変更する TODO: もう少しいい方法を取りたい
				Type:        domain.HOMEPAGE,
				PrPermitted: true,
			},
			setup: func(f fields, args args, want *domain.Account) {
				want.ID = args.userID
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(args.accountID, args.userID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "type", "check"}).
							AddRow(want.ID, args.userID, want.Type, want.PrPermitted),
					)
			},
			assertion: assert.NoError,
		},
		{
			name: "InvalidDB",
			fields: fields{
				sqlhandler: mock_database.NewMockSQLHandler(!isValidDB),
				portal:     mock_external.NewMockPortalAPI(),
				traq:       mock_external.NewMockTraQAPI(),
			},
			args:      args{ids[0], util.UUID()},
			want:      nil,
			setup:     func(f fields, args args, want *domain.Account) {},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			tt.setup(tt.fields, tt.args, tt.want)
			repo := NewUserRepository(tt.fields.sqlhandler, tt.fields.portal, tt.fields.traq)
			// Assertion
			got, err := repo.GetAccount(tt.args.userID, tt.args.accountID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	t.Parallel()
	type fields struct {
		sqlhandler database.SQLHandler
		portal     external.PortalAPI
		traq       external.TraQAPI
	}
	type args struct {
		id      uuid.UUID
		changes map[string]interface{}
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		setup     func(f fields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			fields: fields{
				sqlhandler: mock_database.NewMockSQLHandler(isValidDB),
				portal:     mock_external.NewMockPortalAPI(),
				traq:       mock_external.NewMockTraQAPI(),
			},
			args: args{
				id: ids[0],
				changes: map[string]interface{}{
					"description": util.AlphaNumeric(10),
					"check":       true,
				},
			},
			setup: func(f fields, args args) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(args.id), // TODO: もっとちゃんと返したほうがいいかも
					)
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.
					ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `check`=?,`description`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.changes["check"], args.changes["description"], sqlmock.AnyArg(), args.id). // UpdatedAtはsqlmock.AnyArg()でマッチさせる
					WillReturnResult(sqlmock.NewResult(1, 1))
				sqlhandler.Mock.ExpectCommit()
				sqlhandler.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "NotFound",
			fields: fields{
				sqlhandler: mock_database.NewMockSQLHandler(isValidDB),
				portal:     mock_external.NewMockPortalAPI(),
				traq:       mock_external.NewMockTraQAPI(),
			},
			args: args{
				id: ids[0],
				changes: map[string]interface{}{
					"description": util.AlphaNumeric(10),
					"check":       true,
				},
			},
			setup: func(f fields, args args) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnError(repository.ErrNotFound)
				sqlhandler.Mock.ExpectRollback()
				sqlhandler.Mock.ExpectCommit()
			},
			assertion: assert.Error,
		},
		// TODO: トランザクションエラーのテストを書く
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			tt.setup(tt.fields, tt.args)
			repo := NewUserRepository(tt.fields.sqlhandler, tt.fields.portal, tt.fields.traq)
			// Assertion
			tt.assertion(t, repo.Update(tt.args.id, tt.args.changes))
		})
	}
}

func TestUserRepository_CreateAccount(t *testing.T) {
	t.Parallel()
	type fields struct {
		sqlhandler database.SQLHandler
		portal     external.PortalAPI
		traq       external.TraQAPI
	}
	type args struct {
		id   uuid.UUID
		args *repository.CreateAccountArgs
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *domain.Account
		setup     func(f fields, args args, want *domain.Account)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			fields: fields{
				sqlhandler: mock_database.NewMockSQLHandler(isValidDB),
				portal:     mock_external.NewMockPortalAPI(),
				traq:       mock_external.NewMockTraQAPI(),
			},
			args: args{
				id: ids[0],
				args: &repository.CreateAccountArgs{
					ID:          util.UUID(),
					Name:        util.AlphaNumeric(5),
					Type:        domain.HOMEPAGE,
					URL:         util.AlphaNumeric(5),
					PrPermitted: true,
				},
			},
			want: &domain.Account{
				ID:          uuid.Nil, // setupで変更
				Type:        domain.HOMEPAGE,
				PrPermitted: true,
			},
			setup: func(f fields, args args, want *domain.Account) {
				want.ID = args.args.ID
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.
					ExpectExec(regexp.QuoteMeta("INSERT INTO `accounts` (`id`,`type`,`name`,`url`,`user_id`,`check`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(args.args.ID, args.args.Type, args.args.Name, args.args.URL, args.id, args.args.PrPermitted, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				sqlhandler.Mock.ExpectCommit()
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(args.args.ID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "type", "check"}).
							AddRow(args.args.ID, args.args.Type, args.args.PrPermitted),
					)
			},
			assertion: assert.NoError,
		},
		{
			name: "InvalidDB",
			fields: fields{
				sqlhandler: mock_database.NewMockSQLHandler(!isValidDB),
				portal:     mock_external.NewMockPortalAPI(),
				traq:       mock_external.NewMockTraQAPI(),
			},
			args: args{
				id: ids[0],
				args: &repository.CreateAccountArgs{
					ID:          util.UUID(),
					Name:        util.AlphaNumeric(5),
					Type:        domain.HOMEPAGE,
					URL:         util.AlphaNumeric(5),
					PrPermitted: true,
				},
			},
			want:      nil,
			setup:     func(f fields, args args, want *domain.Account) {},
			assertion: assert.Error,
		},
		{
			name: "CreatedButNotFound",
			fields: fields{
				sqlhandler: mock_database.NewMockSQLHandler(isValidDB),
				portal:     mock_external.NewMockPortalAPI(),
				traq:       mock_external.NewMockTraQAPI(),
			},
			args: args{
				id: ids[0],
				args: &repository.CreateAccountArgs{
					ID:          util.UUID(),
					Name:        util.AlphaNumeric(5),
					Type:        domain.HOMEPAGE,
					URL:         util.AlphaNumeric(5),
					PrPermitted: true,
				},
			},
			want: nil,
			setup: func(f fields, args args, want *domain.Account) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.
					ExpectExec(regexp.QuoteMeta("INSERT INTO `accounts` (`id`,`type`,`name`,`url`,`user_id`,`check`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(args.args.ID, args.args.Type, args.args.Name, args.args.URL, args.id, args.args.PrPermitted, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				sqlhandler.Mock.ExpectCommit()
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(args.args.ID).
					WillReturnError(repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			tt.setup(tt.fields, tt.args, tt.want)
			repo := NewUserRepository(tt.fields.sqlhandler, tt.fields.portal, tt.fields.traq)
			// Assertion
			got, err := repo.CreateAccount(tt.args.id, tt.args.args)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_UpdateAccount(t *testing.T) {
	t.Parallel()
	type fields struct {
		sqlhandler database.SQLHandler
		portal     external.PortalAPI
		traq       external.TraQAPI
	}
	type args struct {
		userID    uuid.UUID
		accountID uuid.UUID
		changes   map[string]interface{}
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		setup     func(f fields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			fields: fields{
				sqlhandler: mock_database.NewMockSQLHandler(isValidDB),
				portal:     mock_external.NewMockPortalAPI(),
				traq:       mock_external.NewMockTraQAPI(),
			},
			args: args{
				userID:    ids[0],
				accountID: util.UUID(),
				changes: map[string]interface{}{
					"name":  util.AlphaNumeric(5),
					"url":   util.AlphaNumeric(5),
					"check": true,
					"type":  domain.HOMEPAGE,
				},
			},
			setup: func(f fields, args args) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(args.accountID, args.userID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.accountID))
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.ExpectExec(regexp.QuoteMeta("UPDATE `accounts` SET `check`=?,`name`=?,`type`=?,`url`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.changes["check"], args.changes["name"], args.changes["type"], args.changes["url"], sqlmock.AnyArg(), args.accountID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				sqlhandler.Mock.ExpectCommit()
				sqlhandler.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "NotFound",
			fields: fields{
				sqlhandler: mock_database.NewMockSQLHandler(isValidDB),
				portal:     mock_external.NewMockPortalAPI(),
				traq:       mock_external.NewMockTraQAPI(),
			},
			args: args{
				userID:    ids[0],
				accountID: util.UUID(),
				changes: map[string]interface{}{
					"name":  util.AlphaNumeric(5),
					"url":   util.AlphaNumeric(5),
					"check": true,
					"type":  domain.HOMEPAGE,
				},
			},
			setup: func(f fields, args args) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(args.accountID, args.userID).
					WillReturnError(repository.ErrNotFound)
				sqlhandler.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		// TODO: トランザクションエラーのテストを書く
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			tt.setup(tt.fields, tt.args)
			repo := NewUserRepository(tt.fields.sqlhandler, tt.fields.portal, tt.fields.traq)
			// Assertion
			tt.assertion(t, repo.UpdateAccount(tt.args.userID, tt.args.accountID, tt.args.changes))
		})
	}
}
