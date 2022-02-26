package repository_test

import (
	"math/rand"
	"regexp"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/traPtitech/traPortfolio/util/optional"
	"github.com/traPtitech/traPortfolio/util/random"
	"gorm.io/gorm"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/database"
	"github.com/traPtitech/traPortfolio/interfaces/database/mock_database"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external"
	impl "github.com/traPtitech/traPortfolio/interfaces/repository"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type mockUserRepositoryFields struct {
	sqlhandler database.SQLHandler
	portal     external.PortalAPI
	traq       external.TraQAPI
}

func newMockUserRepositoryFields(ctrl *gomock.Controller) mockUserRepositoryFields {
	return mockUserRepositoryFields{
		sqlhandler: mock_database.NewMockSQLHandler(),
		portal:     mock_external.NewMockPortalAPI(ctrl),
		traq:       mock_external.NewMockTraQAPI(ctrl),
	}
}

func TestUserRepository_GetUsers(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		want      []*domain.User
		setup     func(f mockUserRepositoryFields, want []*domain.User)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			want: []*domain.User{
				{
					ID:       random.UUID(),
					Name:     random.AlphaNumeric(rand.Intn(30) + 1),
					RealName: random.AlphaNumeric(rand.Intn(30) + 1),
				},
				{
					ID:       random.UUID(),
					Name:     random.AlphaNumeric(rand.Intn(30) + 1),
					RealName: random.AlphaNumeric(rand.Intn(30) + 1),
				},
				{
					ID:       random.UUID(),
					Name:     random.AlphaNumeric(rand.Intn(30) + 1),
					RealName: random.AlphaNumeric(rand.Intn(30) + 1),
				},
			},
			setup: func(f mockUserRepositoryFields, want []*domain.User) {
				rows := sqlmock.NewRows([]string{"id", "name"})
				for _, v := range want {
					rows.AddRow(v.ID, v.Name)
				}
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
					WillReturnRows(rows)
				p := f.portal.(*mock_external.MockPortalAPI)
				p.EXPECT().GetAll().Return(makePortalUsers(want), nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError_Find",
			want: nil,
			setup: func(f mockUserRepositoryFields, want []*domain.User) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
					WillReturnError(errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "PortalError",
			want: nil,
			setup: func(f mockUserRepositoryFields, want []*domain.User) {
				users := []*domain.User{
					{
						ID:       random.UUID(),
						Name:     random.AlphaNumeric(rand.Intn(30) + 1),
						RealName: random.AlphaNumeric(rand.Intn(30) + 1),
					},
				}
				rows := sqlmock.NewRows([]string{"id", "name"})
				for _, v := range users {
					rows.AddRow(v.ID, v.Name)
				}
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).
					WillReturnRows(rows)
				p := f.portal.(*mock_external.MockPortalAPI)
				p.EXPECT().GetAll().Return(nil, errUnexpected)
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
			f := newMockUserRepositoryFields(ctrl)
			tt.setup(f, tt.want)
			repo := impl.NewUserRepository(f.sqlhandler, f.portal, f.traq)
			// Assertion
			got, err := repo.GetUsers()
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_GetUser(t *testing.T) {
	uid := random.UUID()

	t.Parallel()
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.UserDetail
		setup     func(f mockUserRepositoryFields, args args, want *domain.UserDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{uid},
			want: &domain.UserDetail{
				User: domain.User{
					ID:       uid,
					Name:     random.AlphaNumeric(rand.Intn(30) + 1),
					RealName: random.AlphaNumeric(rand.Intn(30) + 1),
				},
				State: domain.TraqStateActive,
				Bio:   random.AlphaNumeric(rand.Intn(30) + 1),
				Accounts: []*domain.Account{
					{
						ID:          random.UUID(),
						Type:        domain.HOMEPAGE,
						PrPermitted: true,
					},
				},
			},
			setup: func(f mockUserRepositoryFields, args args, want *domain.UserDetail) {
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
					WillReturnRows(rows)
				p := f.portal.(*mock_external.MockPortalAPI)
				p.EXPECT().GetByID(want.User.Name).Return(makePortalUser(&want.User), nil)
				t := f.traq.(*mock_external.MockTraQAPI)
				t.EXPECT().GetByID(args.id).Return(makeTraqUser(want), nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "NotFound",
			args: args{random.UUID()},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want *domain.UserDetail) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}))
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError",
			args: args{random.UUID()},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want *domain.UserDetail) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnError(errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "PortalError",
			args: args{random.UUID()},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want *domain.UserDetail) {
				name := random.AlphaNumeric(rand.Intn(30) + 1)
				h := f.sqlhandler.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "name", "description"}).
							AddRow(args.id, name, random.AlphaNumeric(rand.Intn(30)+1)),
					)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `accounts` WHERE `accounts`.`user_id` = ?")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "type", "check"}).
							AddRow(random.UUID(), args.id, 0, 0),
					)
				p := f.portal.(*mock_external.MockPortalAPI)
				p.EXPECT().GetByID(name).Return(nil, errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "TraqError",
			args: args{random.UUID()},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want *domain.UserDetail) {
				name := random.AlphaNumeric(rand.Intn(30) + 1)
				h := f.sqlhandler.(*mock_database.MockSQLHandler)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "name", "description"}).
							AddRow(args.id, name, random.AlphaNumeric(rand.Intn(30)+1)),
					)
				h.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `accounts` WHERE `accounts`.`user_id` = ?")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "type", "check"}).
							AddRow(random.UUID(), args.id, 0, 0),
					)
				p := f.portal.(*mock_external.MockPortalAPI)
				p.EXPECT().GetByID(name).Return(makePortalUser(&domain.User{Name: name}), nil)
				t := f.traq.(*mock_external.MockTraQAPI)
				t.EXPECT().GetByID(args.id).Return(nil, errUnexpected)
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
			f := newMockUserRepositoryFields(ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := impl.NewUserRepository(f.sqlhandler, f.portal, f.traq)
			// Assertion
			got, err := repo.GetUser(tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_GetAccounts(t *testing.T) {
	t.Parallel()
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      []*domain.Account
		setup     func(f mockUserRepositoryFields, args args, want []*domain.Account)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{random.UUID()},
			want: []*domain.Account{
				{
					ID:          random.UUID(),
					Type:        domain.HOMEPAGE,
					PrPermitted: true,
				},
			},
			setup: func(f mockUserRepositoryFields, args args, want []*domain.Account) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "type", "check"})
				for _, v := range want {
					rows.AddRow(v.ID, args.userID, v.Type, v.PrPermitted)
				}
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `accounts` WHERE `accounts`.`user_id` = ?")).
					WithArgs(args.userID).
					WillReturnRows(rows)
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			args: args{random.UUID()},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want []*domain.Account) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `accounts` WHERE user_id = ?")).
					WithArgs(args.userID).
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
			f := newMockUserRepositoryFields(ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := impl.NewUserRepository(f.sqlhandler, f.portal, f.traq)
			// Assertion
			got, err := repo.GetAccounts(tt.args.userID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_GetAccount(t *testing.T) {
	aid := random.UUID() // Successで使うaccountID

	t.Parallel()
	type args struct {
		userID    uuid.UUID
		accountID uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.Account
		setup     func(f mockUserRepositoryFields, args args, want *domain.Account)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				userID:    random.UUID(),
				accountID: aid,
			},
			want: &domain.Account{
				ID:          aid,
				Type:        domain.HOMEPAGE,
				PrPermitted: true,
			},
			setup: func(f mockUserRepositoryFields, args args, want *domain.Account) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(args.accountID, args.userID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "type", "check"}).
							AddRow(args.accountID, args.userID, want.Type, want.PrPermitted),
					)
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			args: args{
				userID:    random.UUID(),
				accountID: random.UUID(),
			},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want *domain.Account) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(args.accountID, args.userID).
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
			f := newMockUserRepositoryFields(ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := impl.NewUserRepository(f.sqlhandler, f.portal, f.traq)
			// Assertion
			got, err := repo.GetAccount(tt.args.userID, tt.args.accountID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_UpdateUser(t *testing.T) {
	t.Parallel()
	type args struct {
		id   uuid.UUID
		args *repository.UpdateUserArgs
	}
	tests := []struct {
		name      string
		args      args
		setup     func(f mockUserRepositoryFields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				id: random.UUID(),
				args: &repository.UpdateUserArgs{
					Description: optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Check:       optional.NewBool(true, true),
				},
			},
			setup: func(f mockUserRepositoryFields, args args) {
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
					WithArgs(args.args.Check.Bool, args.args.Description.String, anyTime{}, args.id).
					WillReturnResult(sqlmock.NewResult(1, 1))
				sqlhandler.Mock.ExpectCommit()
				sqlhandler.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "NotFound",
			args: args{
				id: random.UUID(),
				args: &repository.UpdateUserArgs{
					Description: optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Check:       optional.NewBool(true, true),
				},
			},
			setup: func(f mockUserRepositoryFields, args args) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnError(gorm.ErrRecordNotFound)
				sqlhandler.Mock.ExpectRollback()
				sqlhandler.Mock.ExpectCommit()
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_Update",
			args: args{
				id: random.UUID(),
				args: &repository.UpdateUserArgs{
					Description: optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Check:       optional.NewBool(true, true),
				},
			},
			setup: func(f mockUserRepositoryFields, args args) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(args.id),
					)
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.
					ExpectExec(regexp.QuoteMeta("UPDATE `users` SET `check`=?,`description`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.args.Check.Bool, args.args.Description.String, anyTime{}, args.id).
					WillReturnError(errUnexpected)
				sqlhandler.Mock.ExpectRollback()
				sqlhandler.Mock.ExpectRollback()
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
			f := newMockUserRepositoryFields(ctrl)
			tt.setup(f, tt.args)
			repo := impl.NewUserRepository(f.sqlhandler, f.portal, f.traq)
			// Assertion
			tt.assertion(t, repo.UpdateUser(tt.args.id, tt.args.args))
		})
	}
}

func TestUserRepository_CreateAccount(t *testing.T) {
	t.Parallel()
	type args struct {
		id   uuid.UUID
		args *repository.CreateAccountArgs
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.Account
		setup     func(f mockUserRepositoryFields, args args, want *domain.Account)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				id: random.UUID(),
				args: &repository.CreateAccountArgs{
					ID:          random.AlphaNumeric(rand.Intn(30) + 1),
					Type:        domain.HOMEPAGE,
					URL:         random.AlphaNumeric(rand.Intn(30) + 1),
					PrPermitted: true,
				},
			},
			want: &domain.Account{
				// ID: 関数内で生成するので比較しない
				Type:        domain.HOMEPAGE,
				PrPermitted: true,
			},
			setup: func(f mockUserRepositoryFields, args args, want *domain.Account) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.
					ExpectExec(regexp.QuoteMeta("INSERT INTO `accounts` (`id`,`type`,`name`,`url`,`user_id`,`check`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, args.args.Type, args.args.ID, args.args.URL, args.id, args.args.PrPermitted, anyTime{}, anyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				sqlhandler.Mock.ExpectCommit()
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(anyUUID{}).
					WillReturnRows(
						sqlmock.NewRows([]string{"type", "check"}).
							AddRow(args.args.Type, args.args.PrPermitted),
					)
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			args: args{
				id: random.UUID(),
				args: &repository.CreateAccountArgs{
					ID:          random.AlphaNumeric(rand.Intn(30) + 1),
					Type:        domain.HOMEPAGE,
					URL:         random.AlphaNumeric(rand.Intn(30) + 1),
					PrPermitted: true,
				},
			},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want *domain.Account) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.
					ExpectExec(regexp.QuoteMeta("INSERT INTO `accounts` (`id`,`type`,`name`,`url`,`user_id`,`check`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, args.args.Type, args.args.ID, args.args.URL, args.id, args.args.PrPermitted, anyTime{}, anyTime{}).
					WillReturnError(errUnexpected)
				sqlhandler.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "CreatedButNotFound",
			args: args{
				id: random.UUID(),
				args: &repository.CreateAccountArgs{
					ID:          random.AlphaNumeric(rand.Intn(30) + 1),
					Type:        domain.HOMEPAGE,
					URL:         random.AlphaNumeric(rand.Intn(30) + 1),
					PrPermitted: true,
				},
			},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want *domain.Account) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.
					ExpectExec(regexp.QuoteMeta("INSERT INTO `accounts` (`id`,`type`,`name`,`url`,`user_id`,`check`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, args.args.Type, args.args.ID, args.args.URL, args.id, args.args.PrPermitted, anyTime{}, anyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				sqlhandler.Mock.ExpectCommit()
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(anyUUID{}).
					WillReturnError(gorm.ErrRecordNotFound)
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
			f := newMockUserRepositoryFields(ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := impl.NewUserRepository(f.sqlhandler, f.portal, f.traq)
			// Assertion
			got, err := repo.CreateAccount(tt.args.id, tt.args.args)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_UpdateAccount(t *testing.T) {
	t.Parallel()
	type args struct {
		userID    uuid.UUID
		accountID uuid.UUID
		args      *repository.UpdateAccountArgs
	}
	tests := []struct {
		name      string
		args      args
		setup     func(f mockUserRepositoryFields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				userID:    random.UUID(),
				accountID: random.UUID(),
				args: &repository.UpdateAccountArgs{
					Name:        optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					URL:         optional.NewString(random.RandURLString(), true),
					PrPermitted: optional.NewBool(true, true),
					Type:        optional.NewInt64(int64(domain.HOMEPAGE), true),
				},
			},
			setup: func(f mockUserRepositoryFields, args args) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(anyUUID{}, args.userID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.accountID))
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.ExpectExec(regexp.QuoteMeta("UPDATE `accounts` SET `check`=?,`name`=?,`type`=?,`url`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.args.PrPermitted.Bool, args.args.Name.String, args.args.Type.Int64, args.args.URL.String, anyTime{}, args.accountID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				sqlhandler.Mock.ExpectCommit()
				sqlhandler.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "NotFound",
			args: args{
				userID:    random.UUID(),
				accountID: random.UUID(),
				args: &repository.UpdateAccountArgs{
					Name:        optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					URL:         optional.NewString(random.RandURLString(), true),
					PrPermitted: optional.NewBool(true, true),
					Type:        optional.NewInt64(int64(domain.HOMEPAGE), true),
				},
			},
			setup: func(f mockUserRepositoryFields, args args) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(args.accountID, args.userID).
					WillReturnError(gorm.ErrRecordNotFound)
				sqlhandler.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_Update",
			args: args{
				userID:    random.UUID(),
				accountID: random.UUID(),
				args: &repository.UpdateAccountArgs{
					Name:        optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					URL:         optional.NewString(random.RandURLString(), true),
					PrPermitted: optional.NewBool(true, true),
					Type:        optional.NewInt64(int64(domain.HOMEPAGE), true),
				},
			},
			setup: func(f mockUserRepositoryFields, args args) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(anyUUID{}, args.userID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.accountID))
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.ExpectExec(regexp.QuoteMeta("UPDATE `accounts` SET `check`=?,`name`=?,`type`=?,`url`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.args.PrPermitted.Bool, args.args.Name.String, args.args.Type.Int64, args.args.URL.String, anyTime{}, args.accountID).
					WillReturnError(errUnexpected)
				sqlhandler.Mock.ExpectRollback()
				sqlhandler.Mock.ExpectRollback()
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
			f := newMockUserRepositoryFields(ctrl)
			tt.setup(f, tt.args)
			repo := impl.NewUserRepository(f.sqlhandler, f.portal, f.traq)
			// Assertion
			tt.assertion(t, repo.UpdateAccount(tt.args.userID, tt.args.accountID, tt.args.args))
		})
	}
}

func TestUserRepository_DeleteAccount(t *testing.T) {
	t.Parallel()
	type args struct {
		accountID uuid.UUID
		userID    uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		setup     func(f mockUserRepositoryFields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				accountID: random.UUID(),
				userID:    random.UUID(),
			},
			setup: func(f mockUserRepositoryFields, args args) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ?")).
					WithArgs(args.accountID, args.userID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				sqlhandler.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			args: args{
				accountID: random.UUID(),
				userID:    random.UUID(),
			},
			setup: func(f mockUserRepositoryFields, args args) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.ExpectBegin()
				sqlhandler.Mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ?")).
					WithArgs(args.accountID, args.userID).
					WillReturnError(errUnexpected)
				sqlhandler.Mock.ExpectRollback()
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
			f := newMockUserRepositoryFields(ctrl)
			tt.setup(f, tt.args)
			repo := impl.NewUserRepository(f.sqlhandler, f.portal, f.traq)
			// Assertion
			tt.assertion(t, repo.DeleteAccount(tt.args.accountID, tt.args.userID))
		})
	}
}

func TestUserRepository_GetProjects(t *testing.T) {
	t.Parallel()
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      []*domain.UserProject
		setup     func(f mockUserRepositoryFields, args args, want []*domain.UserProject)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{userID: random.UUID()},
			want: []*domain.UserProject{
				{
					ID:           random.UUID(),
					Name:         random.AlphaNumeric(rand.Intn(30) + 1),
					Duration:     random.Duration(),
					UserDuration: random.Duration(),
				},
			},
			setup: func(f mockUserRepositoryFields, args args, want []*domain.UserProject) {
				rows := sqlmock.NewRows([]string{"id", "project_id", "user_id", "since_year", "since_semester", "until_year", "until_semester"})
				for _, v := range want {
					ud := v.UserDuration
					rows.AddRow(random.UUID(), v.ID, args.userID, ud.Since.Year, ud.Since.Semester, ud.Until.Year, ud.Until.Semester)
				}
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `project_members` WHERE `project_members`.`user_id` = ?")).
					WithArgs(args.userID).
					WillReturnRows(rows)
				for _, v := range want {
					d := v.Duration
					sqlhandler.Mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `projects` WHERE `projects`.`id` = ?")).
						WithArgs(v.ID).
						WillReturnRows(
							sqlmock.NewRows([]string{"id", "name", "description", "link", "since_year", "since_semester", "until_year", "until_semester", "created_at", "updated_at"}).
								AddRow(v.ID, v.Name, random.AlphaNumeric(rand.Intn(30)+1), random.AlphaNumeric(rand.Intn(30)+1), d.Since.Year, d.Since.Semester, d.Until.Year, d.Until.Semester, time.Now(), time.Now()),
						)
				}
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			args: args{userID: random.UUID()},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want []*domain.UserProject) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.
					ExpectQuery(regexp.QuoteMeta("SELECT * FROM `project_members` WHERE `project_members`.`user_id` = ?")).
					WithArgs(args.userID).
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
			f := newMockUserRepositoryFields(ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := impl.NewUserRepository(f.sqlhandler, f.portal, f.traq)
			// Assertion
			got, err := repo.GetProjects(tt.args.userID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_GetContests(t *testing.T) {
	t.Parallel()
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      []*domain.UserContest
		setup     func(f mockUserRepositoryFields, args args, want []*domain.UserContest)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{userID: random.UUID()},
			want: []*domain.UserContest{
				{
					ID:          random.UUID(),
					Name:        random.AlphaNumeric(rand.Intn(30) + 1),
					Result:      random.AlphaNumeric(rand.Intn(30) + 1),
					ContestName: random.AlphaNumeric(rand.Intn(30) + 1),
				},
			},
			setup: func(f mockUserRepositoryFields, args args, want []*domain.UserContest) {
				rows := sqlmock.NewRows([]string{"team_id"})
				for _, v := range want {
					rows.AddRow(v.ID)
				}
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`user_id` = ?")).
					WithArgs(args.userID).
					WillReturnRows(rows)
				cids := make([]uuid.UUID, len(want))
				for i, v := range want {
					cids[i] = random.UUID()
					sqlhandler.Mock.
						ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ?")).
						WithArgs(v.ID).
						WillReturnRows(
							sqlmock.NewRows([]string{"id", "contest_id", "name", "result", "contest_name"}).
								AddRow(v.ID, cids[i], v.Name, v.Result, v.ContestName),
						)
				}
				for i, v := range want {
					sqlhandler.Mock.
						ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contests` WHERE `contests`.`id` = ?")).
						WithArgs(cids[i]).
						WillReturnRows(
							sqlmock.NewRows([]string{"id", "name"}).
								AddRow(cids[i], v.ContestName),
						)
				}
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			args: args{userID: random.UUID()},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want []*domain.UserContest) {
				sqlhandler := f.sqlhandler.(*mock_database.MockSQLHandler)
				sqlhandler.Mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`user_id` = ?")).
					WithArgs(args.userID).
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
			f := newMockUserRepositoryFields(ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := impl.NewUserRepository(f.sqlhandler, f.portal, f.traq)
			// Assertion
			got, err := repo.GetContests(tt.args.userID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
