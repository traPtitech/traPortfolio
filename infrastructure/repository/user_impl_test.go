package repository

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/traPtitech/traPortfolio/util/optional"
	"github.com/traPtitech/traPortfolio/util/random"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external"
	"github.com/traPtitech/traPortfolio/usecases/repository"
)

type mockUserRepositoryFields struct {
	h      *MockSQLHandler
	portal *mock_external.MockPortalAPI
	traq   *mock_external.MockTraQAPI
}

func newMockUserRepositoryFields(ctrl *gomock.Controller) mockUserRepositoryFields {
	return mockUserRepositoryFields{
		h:      NewMockSQLHandler(),
		portal: mock_external.NewMockPortalAPI(ctrl),
		traq:   mock_external.NewMockTraQAPI(ctrl),
	}
}

func TestUserRepository_GetUsers(t *testing.T) {
	name := random.AlphaNumeric()

	t.Parallel()
	type args struct {
		args *repository.GetUsersArgs
	}
	tests := []struct {
		name      string
		args      args
		want      []*domain.User
		setup     func(t *testing.T, f mockUserRepositoryFields, args args, want []*domain.User)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success_NoOpts",
			args: args{
				&repository.GetUsersArgs{},
			},
			want: []*domain.User{
				domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), true),
				domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), true),
				domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), true),
			},
			setup: func(t *testing.T, f mockUserRepositoryFields, args args, want []*domain.User) {
				f.traq.EXPECT().GetAll(mustMakeTraqGetAllArgs(t, args.args)).Return(makeTraqUsers(t, want), nil)

				rows := sqlmock.NewRows([]string{"id", "name", "check"})
				for _, v := range want {
					rows.AddRow(v.ID, v.Name, v.Check)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` IN (?,?,?)")).
					WithArgs(want[0].ID, want[1].ID, want[2].ID).
					WillReturnRows(rows)

				f.portal.EXPECT().GetAll().Return(makePortalUsers(t, want), nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Success_NoCheck",
			args: args{
				&repository.GetUsersArgs{},
			},
			want: []*domain.User{
				domain.NewUser(random.UUID(), random.AlphaNumeric(), "", false),
				domain.NewUser(random.UUID(), random.AlphaNumeric(), "", false),
				domain.NewUser(random.UUID(), random.AlphaNumeric(), "", false),
			},
			setup: func(t *testing.T, f mockUserRepositoryFields, args args, want []*domain.User) {
				f.traq.EXPECT().GetAll(mustMakeTraqGetAllArgs(t, args.args)).Return(makeTraqUsers(t, want), nil)

				rows := sqlmock.NewRows([]string{"id", "name", "check"})
				for _, v := range want {
					rows.AddRow(v.ID, v.Name, v.Check)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` IN (?,?,?)")).
					WithArgs(want[0].ID, want[1].ID, want[2].ID).
					WillReturnRows(rows)

				f.portal.EXPECT().GetAll().Return(makePortalUsers(t, want), nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Success_WithOpt_IncludeSuspended",
			args: args{
				&repository.GetUsersArgs{
					IncludeSuspended: optional.From(true),
				},
			},
			want: []*domain.User{
				domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), true),
				domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), true),
				domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), true),
			},
			setup: func(t *testing.T, f mockUserRepositoryFields, args args, want []*domain.User) {
				f.traq.EXPECT().GetAll(mustMakeTraqGetAllArgs(t, args.args)).Return(makeTraqUsers(t, want), nil)

				rows := sqlmock.NewRows([]string{"id", "name", "check"})
				for _, v := range want {
					rows.AddRow(v.ID, v.Name, v.Check)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` IN (?,?,?)")).
					WithArgs(want[0].ID, want[1].ID, want[2].ID).
					WillReturnRows(rows)

				f.portal.EXPECT().GetAll().Return(makePortalUsers(t, want), nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Success_WithOpt_Name",
			args: args{
				&repository.GetUsersArgs{
					Name: optional.From(name),
				},
			},
			want: []*domain.User{
				domain.NewUser(random.UUID(), name, random.AlphaNumeric(), true),
			},
			setup: func(t *testing.T, f mockUserRepositoryFields, args args, want []*domain.User) {
				u := want[0]

				f.traq.EXPECT().GetAll(mustMakeTraqGetAllArgs(t, args.args)).Return(makeTraqUsers(t, want), nil)

				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` IN (?)")).
					WithArgs(u.ID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name", "check"}).
							AddRow(u.ID, u.Name, u.Check),
					)

				f.portal.EXPECT().GetByTraqID(u.Name).Return(makePortalUser(t, want[0]), nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Success_zero_result",
			args: args{
				&repository.GetUsersArgs{},
			},
			want: []*domain.User{},
			setup: func(t *testing.T, f mockUserRepositoryFields, args args, want []*domain.User) {
				f.traq.EXPECT().GetAll(mustMakeTraqGetAllArgs(t, args.args)).Return(makeTraqUsers(t, want), nil)

				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` IN (?)")).
					WithArgs().
					WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))
			},
			assertion: assert.NoError,
		},
		{
			name: "Success_With_Limit10",
			args: args{
				&repository.GetUsersArgs{
					Limit: optional.From(10),
				},
			},
			want: []*domain.User{
				domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
				domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
				domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
			},
			setup: func(t *testing.T, f mockUserRepositoryFields, args args, want []*domain.User) {
				f.traq.EXPECT().GetAll(mustMakeTraqGetAllArgs(t, args.args)).Return(makeTraqUsers(t, want), nil)

				rows := sqlmock.NewRows([]string{"id", "name", "check"})
				for _, v := range want {
					rows.AddRow(v.ID, v.Name, v.Check)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` IN (?,?,?) LIMIT 10")).
					WithArgs(want[0].ID, want[1].ID, want[2].ID).
					WillReturnRows(rows)

				f.portal.EXPECT().GetAll().Return(makePortalUsers(t, want), nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Success_With_Random_Limit",
			args: args{
				&repository.GetUsersArgs{
					Name:  optional.From(random.AlphaNumeric()),
					Limit: optional.From(rand.Intn(200)),
				},
			},
			want: []*domain.User{
				domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
				domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
				domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
			},
			setup: func(t *testing.T, f mockUserRepositoryFields, args args, want []*domain.User) {
				f.traq.EXPECT().GetAll(mustMakeTraqGetAllArgs(t, args.args)).Return(makeTraqUsers(t, want), nil)
				rows := sqlmock.NewRows([]string{"id", "name", "check"})
				for _, v := range want {
					rows.AddRow(v.ID, v.Name, v.Check)
				}
				f.h.Mock.ExpectQuery(makeSQLQueryRegexp(fmt.Sprintf("SELECT * FROM `users` WHERE `users`.`id` IN (?,?,?) LIMIT %d", args.args.Limit.ValueOrZero()))).
					WithArgs(want[0].ID, want[1].ID, want[2].ID).
					WillReturnRows(rows)
				f.portal.EXPECT().GetAll().Return(makePortalUsers(t, want), nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Error_WithMultipleOpts",
			args: args{
				&repository.GetUsersArgs{
					IncludeSuspended: optional.From(random.Bool()),
					Name:             optional.From(random.AlphaNumeric()),
				},
			},
			want:      nil,
			setup:     func(t *testing.T, f mockUserRepositoryFields, args args, want []*domain.User) {},
			assertion: assert.Error,
		},
		{
			name: "TraqError",
			args: args{
				&repository.GetUsersArgs{},
			},
			want: nil,
			setup: func(t *testing.T, f mockUserRepositoryFields, args args, want []*domain.User) {
				f.traq.EXPECT().GetAll(mustMakeTraqGetAllArgs(t, args.args)).Return(nil, errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_Find",
			args: args{
				&repository.GetUsersArgs{},
			},
			want: nil,
			setup: func(t *testing.T, f mockUserRepositoryFields, args args, want []*domain.User) {
				f.traq.EXPECT().GetAll(mustMakeTraqGetAllArgs(t, args.args)).Return(makeTraqUsers(t, want), nil)

				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` IN (?,?,?)")).
					WillReturnError(errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "PortalError_Single",
			args: args{
				&repository.GetUsersArgs{
					Name: optional.From(name),
				},
			},
			want: nil,
			setup: func(t *testing.T, f mockUserRepositoryFields, args args, want []*domain.User) {
				id := random.UUID()

				f.traq.EXPECT().GetAll(mustMakeTraqGetAllArgs(t, args.args)).Return([]*external.TraQUserResponse{{ID: id}}, nil)

				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` IN (?)")).
					WithArgs(id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "name"}).AddRow(id, name),
					)

				f.portal.EXPECT().GetByTraqID(name).Return(nil, errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "PortalError_Multiple",
			args: args{
				&repository.GetUsersArgs{},
			},
			want: nil,
			setup: func(t *testing.T, f mockUserRepositoryFields, args args, want []*domain.User) {
				users := []*domain.User{
					domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
					domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
					domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
				}

				f.traq.EXPECT().GetAll(mustMakeTraqGetAllArgs(t, args.args)).Return(makeTraqUsers(t, users), nil)
				rows := sqlmock.NewRows([]string{"id", "name"})
				for _, v := range users {
					rows.AddRow(v.ID, v.Name)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` IN (?,?,?)")).
					WithArgs(users[0].ID, users[1].ID, users[2].ID).
					WillReturnRows(rows)

				f.portal.EXPECT().GetAll().Return(nil, errUnexpected)
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
			tt.setup(t, f, tt.args, tt.want)
			repo := NewUserRepository(f.h.Conn, f.portal, f.traq)
			// Assertion
			got, err := repo.GetUsers(context.Background(), tt.args.args)
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
				User:  *domain.NewUser(uid, random.AlphaNumeric(), random.AlphaNumeric(), true),
				State: domain.TraqStateActive,
				Bio:   random.AlphaNumeric(),
				Accounts: []*domain.Account{
					{
						ID:          random.UUID(),
						Type:        domain.HOMEPAGE,
						PrPermitted: true,
					},
				},
			},
			setup: func(f mockUserRepositoryFields, args args, want *domain.UserDetail) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "name", "check", "description"}).
							AddRow(want.User.ID, want.User.Name, want.Check, want.Bio),
					)
				rows := sqlmock.NewRows([]string{"id", "user_id", "type", "check"})
				for _, v := range want.Accounts {
					rows.AddRow(v.ID, want.User.ID, v.Type, v.PrPermitted)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`user_id` = ?")).
					WithArgs(args.id).
					WillReturnRows(rows)
				f.portal.EXPECT().GetByTraqID(want.User.Name).Return(makePortalUser(t, &want.User), nil)
				f.traq.EXPECT().GetByUserID(args.id).Return(makeTraqUser(want), nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "NotFound",
			args: args{random.UUID()},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want *domain.UserDetail) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
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
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
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
				name := random.AlphaNumeric()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "name", "description"}).
							AddRow(args.id, name, random.AlphaNumeric()),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`user_id` = ?")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "type", "check"}).
							AddRow(random.UUID(), args.id, 0, 0),
					)
				f.portal.EXPECT().GetByTraqID(name).Return(nil, errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "TraqError",
			args: args{random.UUID()},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want *domain.UserDetail) {
				name := random.AlphaNumeric()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows(
							[]string{"id", "name", "description"}).
							AddRow(args.id, name, random.AlphaNumeric()),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`user_id` = ?")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "user_id", "type", "check"}).
							AddRow(random.UUID(), args.id, 0, 0),
					)
				f.portal.EXPECT().GetByTraqID(name).Return(makePortalUser(t, &domain.User{Name: name}), nil)
				f.traq.EXPECT().GetByUserID(args.id).Return(nil, errUnexpected)
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
			repo := NewUserRepository(f.h.Conn, f.portal, f.traq)
			// Assertion
			got, err := repo.GetUser(context.Background(), tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_CreateUser(t *testing.T) {
	t.Parallel()
	name := random.AlphaNumeric()
	realName := random.AlphaNumeric()
	check := random.Bool()
	description := random.AlphaNumeric()

	type args struct {
		args *repository.CreateUserArgs
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.UserDetail
		setup     func(f mockUserRepositoryFields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				args: &repository.CreateUserArgs{
					Description: description,
					Check:       check,
					Name:        name,
				},
			},
			want: &domain.UserDetail{
				User:     *domain.NewUser(random.UUID(), name, realName, check),
				Bio:      description,
				Accounts: []*domain.Account{},
			},
			setup: func(f mockUserRepositoryFields, args args) {
				f.portal.EXPECT().GetByTraqID(args.args.Name).Return(&external.PortalUserResponse{
					TraQID:   args.args.Name,
					RealName: realName,
				}, nil)

				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("INSERT INTO `users` (`id`,`description`,`check`,`name`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, args.args.Description, args.args.Check, args.args.Name, anyTime{}, anyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "PortalError",
			args: args{
				args: &repository.CreateUserArgs{
					Description: description,
					Check:       check,
					Name:        name,
				},
			},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args) {
				f.portal.EXPECT().GetByTraqID(args.args.Name).Return(nil, errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError",
			args: args{
				args: &repository.CreateUserArgs{
					Description: description,
					Check:       check,
					Name:        name,
				},
			},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args) {
				f.portal.EXPECT().GetByTraqID(args.args.Name).Return(&external.PortalUserResponse{
					TraQID:   args.args.Name,
					RealName: realName,
				}, nil)

				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("INSERT INTO `users` (`id`,`description`,`check`,`name`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, args.args.Description, args.args.Check, args.args.Name, anyTime{}, anyTime{}).
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
			f := newMockUserRepositoryFields(ctrl)
			tt.setup(f, tt.args)
			repo := NewUserRepository(f.h.Conn, f.portal, f.traq)
			// Assertion
			got, err := repo.CreateUser(context.Background(), tt.args.args)
			if tt.want != nil && got != nil {
				tt.want.ID = got.ID // 関数内でIDを生成するためここで合わせる
			}
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
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.userID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.userID))
				rows := sqlmock.NewRows([]string{"id", "user_id", "type", "check"})
				for _, v := range want {
					rows.AddRow(v.ID, args.userID, v.Type, v.PrPermitted)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`user_id` = ?")).
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
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.userID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.userID))
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`user_id` = ?")).
					WithArgs(args.userID).
					WillReturnError(errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "User not found",
			args: args{userID: random.UUID()},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want []*domain.Account) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.userID).
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
			ctrl := gomock.NewController(t)
			f := newMockUserRepositoryFields(ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := NewUserRepository(f.h.Conn, f.portal, f.traq)
			// Assertion
			got, err := repo.GetAccounts(context.Background(), tt.args.userID)
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
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
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
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
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
			repo := NewUserRepository(f.h.Conn, f.portal, f.traq)
			// Assertion
			got, err := repo.GetAccount(context.Background(), tt.args.userID, tt.args.accountID)
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
					Description: optional.From(random.AlphaNumeric()),
					Check:       optional.From(true),
				},
			},
			setup: func(f mockUserRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(args.id), // TODO: もっとちゃんと返したほうがいいかも
					)
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("UPDATE `users` SET `check`=?,`description`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.args.Check.ValueOrZero(), args.args.Description.ValueOrZero(), anyTime{}, args.id).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "NotFound",
			args: args{
				id: random.UUID(),
				args: &repository.UpdateUserArgs{
					Description: optional.From(random.AlphaNumeric()),
					Check:       optional.From(true),
				},
			},
			setup: func(f mockUserRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnError(repository.ErrNotFound)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_Update",
			args: args{
				id: random.UUID(),
				args: &repository.UpdateUserArgs{
					Description: optional.From(random.AlphaNumeric()),
					Check:       optional.From(true),
				},
			},
			setup: func(f mockUserRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.id).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(args.id),
					)
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("UPDATE `users` SET `check`=?,`description`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.args.Check.ValueOrZero(), args.args.Description.ValueOrZero(), anyTime{}, args.id).
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
			f := newMockUserRepositoryFields(ctrl)
			tt.setup(f, tt.args)
			repo := NewUserRepository(f.h.Conn, f.portal, f.traq)
			// Assertion
			tt.assertion(t, repo.UpdateUser(context.Background(), tt.args.id, tt.args.args))
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
					DisplayName: random.AlphaNumeric(),
					Type:        domain.HOMEPAGE,
					URL:         random.AccountURLString(domain.HOMEPAGE),
					PrPermitted: true,
				},
			},
			want: &domain.Account{
				// ID: 関数内で生成するので比較しない
				Type:        domain.HOMEPAGE,
				PrPermitted: true,
			},
			setup: func(f mockUserRepositoryFields, args args, want *domain.Account) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`user_id` = ? AND `accounts`.`type` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(args.id, args.args.Type).
					WillReturnError(repository.ErrNotFound)
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("INSERT INTO `accounts` (`id`,`type`,`name`,`url`,`user_id`,`check`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, args.args.Type, args.args.DisplayName, args.args.URL, args.id, args.args.PrPermitted, anyTime{}, anyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.ExpectCommit()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(anyUUID{}).
					WillReturnRows(
						sqlmock.NewRows([]string{"type", "check"}).
							AddRow(args.args.Type, args.args.PrPermitted),
					)
			},
			assertion: assert.NoError,
		},
		{
			name: "AlreadyExists_FindAccount",
			args: args{
				id: random.UUID(),
				args: &repository.CreateAccountArgs{
					DisplayName: random.AlphaNumeric(),
					Type:        domain.HOMEPAGE,
					URL:         random.AlphaNumeric(),
					PrPermitted: true,
				},
			},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want *domain.Account) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`user_id` = ? AND `accounts`.`type` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(anyUUID{}, domain.HOMEPAGE).
					WillReturnRows(
						sqlmock.NewRows([]string{"type", "check"}).
							AddRow(args.args.Type, args.args.PrPermitted),
					)
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_FindAccount",
			args: args{
				id: random.UUID(),
				args: &repository.CreateAccountArgs{
					DisplayName: random.AlphaNumeric(),
					Type:        domain.HOMEPAGE,
					URL:         random.AlphaNumeric(),
					PrPermitted: true,
				},
			},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want *domain.Account) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`user_id` = ? AND `accounts`.`type` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(anyUUID{}, domain.HOMEPAGE).
					WillReturnError(errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_CreateAccount",
			args: args{
				id: random.UUID(),
				args: &repository.CreateAccountArgs{
					DisplayName: random.AlphaNumeric(),
					Type:        domain.HOMEPAGE,
					URL:         random.AlphaNumeric(),
					PrPermitted: true,
				},
			},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want *domain.Account) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`user_id` = ? AND `accounts`.`type` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(anyUUID{}, domain.HOMEPAGE).
					WillReturnError(repository.ErrNotFound)
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("INSERT INTO `accounts` (`id`,`type`,`name`,`url`,`user_id`,`check`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, args.args.Type, args.args.DisplayName, args.args.URL, args.id, args.args.PrPermitted, anyTime{}, anyTime{}).
					WillReturnError(errUnexpected)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "CreatedButNotFound",
			args: args{
				id: random.UUID(),
				args: &repository.CreateAccountArgs{
					DisplayName: random.AlphaNumeric(),
					Type:        domain.HOMEPAGE,
					URL:         random.AlphaNumeric(),
					PrPermitted: true,
				},
			},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want *domain.Account) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`user_id` = ? AND `accounts`.`type` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(anyUUID{}, domain.HOMEPAGE).
					WillReturnError(repository.ErrNotFound)
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("INSERT INTO `accounts` (`id`,`type`,`name`,`url`,`user_id`,`check`,`created_at`,`updated_at`) VALUES (?,?,?,?,?,?,?,?)")).
					WithArgs(anyUUID{}, args.args.Type, args.args.DisplayName, args.args.URL, args.id, args.args.PrPermitted, anyTime{}, anyTime{}).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.ExpectCommit()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(anyUUID{}).
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
			ctrl := gomock.NewController(t)
			f := newMockUserRepositoryFields(ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := NewUserRepository(f.h.Conn, f.portal, f.traq)
			// Assertion
			got, err := repo.CreateAccount(context.Background(), tt.args.id, tt.args.args)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_UpdateAccount(t *testing.T) {
	aType := optional.From(random.Iotan(domain.AccountLimit))

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
					DisplayName: optional.From(random.AlphaNumeric()),
					URL:         optional.From(random.AccountURLString(aType.ValueOrZero())),
					PrPermitted: optional.From(random.Bool()),
					Type:        aType,
				},
			},
			setup: func(f mockUserRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(anyUUID{}, args.userID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id", "type", "name", "url", "check"}).
							AddRow(args.accountID, (args.args.Type.ValueOrZero()+1)%domain.AccountLimit, args.args.DisplayName.ValueOrZero(), args.args.URL.ValueOrZero(), args.args.PrPermitted.ValueOrZero()),
					)
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`user_id` = ? AND `accounts`.`type` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(args.userID, args.args.Type).
					WillReturnError(repository.ErrNotFound)
				f.h.Mock.ExpectExec(makeSQLQueryRegexp("UPDATE `accounts` SET `check`=?,`name`=?,`type`=?,`url`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.args.PrPermitted.ValueOrZero(), args.args.DisplayName.ValueOrZero(), args.args.Type.ValueOrZero(), args.args.URL.ValueOrZero(), anyTime{}, args.accountID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "NotFound",
			args: args{
				userID:    random.UUID(),
				accountID: random.UUID(),
				args: &repository.UpdateAccountArgs{
					DisplayName: optional.From(random.AlphaNumeric()),
					URL:         optional.From(random.RandURLString()),
					PrPermitted: optional.From(random.Bool()),
					Type:        optional.From(random.Iotan(domain.AccountLimit)),
				},
			},
			setup: func(f mockUserRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(args.accountID, args.userID).
					WillReturnError(repository.ErrNotFound)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "AlreadyExists_FindSameTypeAccount",
			args: args{
				userID:    random.UUID(),
				accountID: random.UUID(),
				args: &repository.UpdateAccountArgs{
					DisplayName: optional.From(random.AlphaNumeric()),
					URL:         optional.From(random.RandURLString()),
					PrPermitted: optional.From(random.Bool()),
					Type:        optional.From(random.Iotan(domain.AccountLimit)),
				},
			},
			setup: func(f mockUserRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(anyUUID{}, args.userID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.accountID))
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`user_id` = ? AND `accounts`.`type` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(args.userID, args.args.Type).
					WillReturnRows(sqlmock.NewRows([]string{"type"}).AddRow(args.args.Type))
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_FindSameTypeAccount",
			args: args{
				userID:    random.UUID(),
				accountID: random.UUID(),
				args: &repository.UpdateAccountArgs{
					DisplayName: optional.From(random.AlphaNumeric()),
					URL:         optional.From(random.RandURLString()),
					PrPermitted: optional.From(random.Bool()),
					Type:        optional.From(random.Iotan(domain.AccountLimit)),
				},
			},
			setup: func(f mockUserRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(anyUUID{}, args.userID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.accountID))
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`user_id` = ? AND `accounts`.`type` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(args.userID, args.args.Type).
					WillReturnError(errUnexpected)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError_Update",
			args: args{
				userID:    random.UUID(),
				accountID: random.UUID(),
				args: &repository.UpdateAccountArgs{
					DisplayName: optional.From(random.AlphaNumeric()),
					URL:         optional.From(random.RandURLString()),
					PrPermitted: optional.From(random.Bool()),
					Type:        optional.From(random.Iotan(domain.AccountLimit)),
				},
			},
			setup: func(f mockUserRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(anyUUID{}, args.userID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.accountID))
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`user_id` = ? AND `accounts`.`type` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(args.userID, args.args.Type).
					WillReturnError(repository.ErrNotFound)
				f.h.Mock.ExpectExec(makeSQLQueryRegexp("UPDATE `accounts` SET `check`=?,`name`=?,`type`=?,`url`=?,`updated_at`=? WHERE `id` = ?")).
					WithArgs(args.args.PrPermitted.ValueOrZero(), args.args.DisplayName.ValueOrZero(), args.args.Type.ValueOrZero(), args.args.URL.ValueOrZero(), anyTime{}, args.accountID).
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
			f := newMockUserRepositoryFields(ctrl)
			tt.setup(f, tt.args)
			repo := NewUserRepository(f.h.Conn, f.portal, f.traq)
			// Assertion
			tt.assertion(t, repo.UpdateAccount(context.Background(), tt.args.userID, tt.args.accountID, tt.args.args))
		})
	}
}

func TestUserRepository_DeleteAccount(t *testing.T) {
	t.Parallel()
	type args struct {
		userID    uuid.UUID
		accountID uuid.UUID
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
			},
			setup: func(f mockUserRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(args.accountID, args.userID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.accountID))
				f.h.Mock.
					ExpectExec(makeSQLQueryRegexp("DELETE FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ?")).
					WithArgs(args.accountID, args.userID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				f.h.Mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "NotFound",
			args: args{
				userID:    random.UUID(),
				accountID: random.UUID(),
			},
			setup: func(f mockUserRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(args.accountID, args.userID).
					WillReturnError(repository.ErrNotFound)
				f.h.Mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name: "UnexpectedError",
			args: args{
				accountID: random.UUID(),
				userID:    random.UUID(),
			},
			setup: func(f mockUserRepositoryFields, args args) {
				f.h.Mock.ExpectBegin()
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ? ORDER BY `accounts`.`id` LIMIT 1")).
					WithArgs(args.accountID, args.userID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.accountID))
				f.h.Mock.ExpectExec(makeSQLQueryRegexp("DELETE FROM `accounts` WHERE `accounts`.`id` = ? AND `accounts`.`user_id` = ?")).
					WithArgs(args.userID, args.accountID).
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
			f := newMockUserRepositoryFields(ctrl)
			tt.setup(f, tt.args)
			repo := NewUserRepository(f.h.Conn, f.portal, f.traq)
			// Assertion
			tt.assertion(t, repo.DeleteAccount(context.Background(), tt.args.userID, tt.args.accountID))
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
					Name:         random.AlphaNumeric(),
					Duration:     random.Duration(),
					UserDuration: random.Duration(),
				},
			},
			setup: func(f mockUserRepositoryFields, args args, want []*domain.UserProject) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.userID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.userID))
				rows := sqlmock.NewRows([]string{"id", "project_id", "user_id", "since_year", "since_semester", "until_year", "until_semester"})
				for _, v := range want {
					ud := v.UserDuration
					rows.AddRow(random.UUID(), v.ID, args.userID, ud.Since.Year, ud.Since.Semester, ud.Until.Year, ud.Until.Semester)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `project_members` WHERE `project_members`.`user_id` = ?")).
					WithArgs(args.userID).
					WillReturnRows(rows)
				for _, v := range want {
					d := v.Duration
					f.h.Mock.ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `projects` WHERE `projects`.`id` = ?")).
						WithArgs(v.ID).
						WillReturnRows(
							sqlmock.NewRows([]string{"id", "name", "description", "link", "since_year", "since_semester", "until_year", "until_semester", "created_at", "updated_at"}).
								AddRow(v.ID, v.Name, random.AlphaNumeric(), random.AlphaNumeric(), d.Since.Year, d.Since.Semester, d.Until.Year, d.Until.Semester, time.Now(), time.Now()),
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
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.userID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.userID))
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `project_members` WHERE `project_members`.`user_id` = ?")).
					WithArgs(args.userID).
					WillReturnError(errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "User not found",
			args: args{userID: random.UUID()},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want []*domain.UserProject) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.userID).
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
			ctrl := gomock.NewController(t)
			f := newMockUserRepositoryFields(ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := NewUserRepository(f.h.Conn, f.portal, f.traq)
			// Assertion
			got, err := repo.GetProjects(context.Background(), tt.args.userID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_GetGroupsByUserID(t *testing.T) {
	t.Parallel()
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      []*domain.UserGroup
		setup     func(f mockUserRepositoryFields, args args, want []*domain.UserGroup)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{userID: random.UUID()},
			want: []*domain.UserGroup{
				{
					ID:       random.UUID(),
					Name:     random.AlphaNumeric(),
					Duration: random.Duration(),
				},
			},
			setup: func(f mockUserRepositoryFields, args args, want []*domain.UserGroup) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.userID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.userID))
				rows := sqlmock.NewRows([]string{"id", "group_id", "user_id", "since_year", "since_semester", "until_year", "until_semester"})
				for _, v := range want {
					d := v.Duration
					rows.AddRow(random.UUID(), v.ID, args.userID, d.Since.Year, d.Since.Semester, d.Until.Year, d.Until.Semester)
				}
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `group_user_belongings` WHERE `group_user_belongings`.`user_id` = ?")).
					WithArgs(args.userID).
					WillReturnRows(rows)
				for _, v := range want {
					f.h.Mock.ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `groups` WHERE `groups`.`group_id` = ?")).
						WithArgs(v.ID).
						WillReturnRows(
							sqlmock.NewRows([]string{"group_id", "name"}).
								AddRow(v.ID, v.Name),
						)
				}
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			args: args{userID: random.UUID()},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want []*domain.UserGroup) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.userID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.userID))
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `group_user_belongings` WHERE `group_user_belongings`.`user_id` = ?")).
					WithArgs(args.userID).
					WillReturnError(errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "User not found",
			args: args{userID: random.UUID()},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want []*domain.UserGroup) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.userID).
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
			ctrl := gomock.NewController(t)
			f := newMockUserRepositoryFields(ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := NewUserRepository(f.h.Conn, f.portal, f.traq)
			// Assertion
			got, err := repo.GetGroupsByUserID(context.Background(), tt.args.userID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_GetContests(t *testing.T) {
	cid := random.UUID()

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
					ID:        cid,
					Name:      random.AlphaNumeric(),
					TimeStart: random.Time(),
					TimeEnd:   random.Time(),
					Teams: []*domain.ContestTeam{
						{
							ID:        random.UUID(),
							ContestID: cid,
							Name:      random.AlphaNumeric(),
							Result:    random.AlphaNumeric(),
						},
					},
				},
			},
			setup: func(f mockUserRepositoryFields, args args, want []*domain.UserContest) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.userID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.userID))
				rows := sqlmock.NewRows([]string{"team_id", "user_id"})
				for _, v := range want {
					rows.AddRow(v.Teams[0].ID, args.userID)
				}
				f.h.Mock.ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`user_id` = ?")).
					WithArgs(args.userID).
					WillReturnRows(rows)
				for _, v := range want {
					f.h.Mock.
						ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` = ?")).
						WithArgs(v.Teams[0].ID).
						WillReturnRows(
							sqlmock.NewRows([]string{"id", "contest_id", "name", "result"}).
								AddRow(v.Teams[0].ID, v.ID, v.Teams[0].Name, v.Teams[0].Result),
						)
				}
				for _, v := range want {
					f.h.Mock.
						ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ?")).
						WithArgs(v.ID).
						WillReturnRows(
							sqlmock.NewRows([]string{"id", "name", "since", "until"}).
								AddRow(v.ID, v.Name, v.TimeStart, v.TimeEnd),
						)
				}
			},
			assertion: assert.NoError,
		},
		{
			name: "Success with multiple teams",
			args: args{userID: random.UUID()},
			want: []*domain.UserContest{
				{
					ID:        cid,
					Name:      random.AlphaNumeric(),
					TimeStart: random.Time(),
					TimeEnd:   random.Time(),
					Teams: []*domain.ContestTeam{
						{
							ID:        random.UUID(),
							ContestID: cid,
							Name:      random.AlphaNumeric(),
							Result:    random.AlphaNumeric(),
						},
						{
							ID:        random.UUID(),
							ContestID: cid,
							Name:      random.AlphaNumeric(),
							Result:    random.AlphaNumeric(),
						},
					},
				},
			},
			setup: func(f mockUserRepositoryFields, args args, want []*domain.UserContest) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.userID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.userID))
				rows := sqlmock.NewRows([]string{"team_id", "user_id"})
				for _, v := range want {
					for _, t := range v.Teams {
						rows.AddRow(t.ID, args.userID)
					}
				}
				f.h.Mock.ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`user_id` = ?")).
					WithArgs(args.userID).
					WillReturnRows(rows)
				for _, v := range want {
					rows = sqlmock.NewRows([]string{"id", "contest_id", "name", "result"})
					for _, t := range v.Teams {
						rows.AddRow(t.ID, t.ContestID, t.Name, t.Result)
					}
					f.h.Mock.
						ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_teams` WHERE `contest_teams`.`id` IN (?,?)")).
						WithArgs(v.Teams[0].ID, v.Teams[1].ID).
						WillReturnRows(rows)
				}
				for _, v := range want {
					f.h.Mock.
						ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contests` WHERE `contests`.`id` = ?")).
						WithArgs(v.ID).
						WillReturnRows(
							sqlmock.NewRows([]string{"id", "name", "since", "until"}).
								AddRow(v.ID, v.Name, v.TimeStart, v.TimeEnd),
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
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.userID).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(args.userID))
				f.h.Mock.ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `contest_team_user_belongings` WHERE `contest_team_user_belongings`.`user_id` = ?")).
					WithArgs(args.userID).
					WillReturnError(errUnexpected)
			},
			assertion: assert.Error,
		},
		{
			name: "User not found",
			args: args{userID: random.UUID()},
			want: nil,
			setup: func(f mockUserRepositoryFields, args args, want []*domain.UserContest) {
				f.h.Mock.
					ExpectQuery(makeSQLQueryRegexp("SELECT * FROM `users` WHERE `users`.`id` = ? ORDER BY `users`.`id` LIMIT 1")).
					WithArgs(args.userID).
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
			ctrl := gomock.NewController(t)
			f := newMockUserRepositoryFields(ctrl)
			tt.setup(f, tt.args, tt.want)
			repo := NewUserRepository(f.h.Conn, f.portal, f.traq)
			// Assertion
			got, err := repo.GetContests(context.Background(), tt.args.userID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
