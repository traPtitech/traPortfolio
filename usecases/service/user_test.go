package service

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/usecases/repository/mock_repository"
	"github.com/traPtitech/traPortfolio/util"
	"github.com/traPtitech/traPortfolio/util/optional"
)

func TestUserService_GetUsers(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		want      []*domain.User
		setup     func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.User)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{ctx: context.Background()},
			want: []*domain.User{
				{
					ID:       util.UUID(),
					Name:     util.AlphaNumeric(5),
					RealName: util.AlphaNumeric(5),
				},
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.User) {
				repo.EXPECT().GetUsers().Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Forbidden",
			args: args{ctx: context.Background()},
			want: nil,
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.User) {
				repo.EXPECT().GetUsers().Return(want, repository.ErrForbidden)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockUserRepository(ctrl)
			event := mock_repository.NewMockEventRepository(ctrl)
			tt.setup(repo, event, tt.args, tt.want)

			s := NewUserService(repo, event)
			got, err := s.GetUsers(tt.args.ctx)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.UserDetail
		setup     func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want *domain.UserDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  util.UUID(),
			},
			want: &domain.UserDetail{
				User: domain.User{
					ID:       uuid.Nil, // setupで変更する
					Name:     util.AlphaNumeric(5),
					RealName: util.AlphaNumeric(5),
				},
				State: 0,
				Bio:   util.AlphaNumeric(10),
				Accounts: []*domain.Account{
					{
						ID:          util.UUID(),
						Type:        0,
						PrPermitted: true,
					},
				},
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want *domain.UserDetail) {
				want.ID = args.id
				repo.EXPECT().GetUser(args.id).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Forbidden",
			args: args{
				ctx: context.Background(),
				id:  util.UUID(),
			},
			want: nil,
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want *domain.UserDetail) {
				repo.EXPECT().GetUser(args.id).Return(nil, repository.ErrForbidden)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockUserRepository(ctrl)
			event := mock_repository.NewMockEventRepository(ctrl)
			tt.setup(repo, event, tt.args, tt.want)

			s := NewUserService(repo, event)
			got, err := s.GetUser(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserService_Update(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		id   uuid.UUID
		args *repository.UpdateUserArgs
	}
	tests := []struct {
		name      string
		args      args
		setup     func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  util.UUID(),
				args: &repository.UpdateUserArgs{
					Description: optional.StringFrom(util.AlphaNumeric(10)),
					Check:       optional.BoolFrom(true),
				},
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args) {
				changes := map[string]interface{}{
					"description": args.args.Description.String,
					"check":       args.args.Check.Bool,
				}
				repo.EXPECT().Update(args.id, changes).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Notfound",
			args: args{
				ctx: context.Background(),
				id:  util.UUID(),
				args: &repository.UpdateUserArgs{
					Description: optional.StringFrom(util.AlphaNumeric(10)),
					Check:       optional.BoolFrom(true),
				},
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args) {
				changes := map[string]interface{}{
					"description": args.args.Description.String,
					"check":       args.args.Check.Bool,
				}
				repo.EXPECT().Update(args.id, changes).Return(repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockUserRepository(ctrl)
			event := mock_repository.NewMockEventRepository(ctrl)
			tt.setup(repo, event, tt.args)

			s := NewUserService(repo, event)
			tt.assertion(t, s.Update(tt.args.ctx, tt.args.id, tt.args.args))
		})
	}
}

func TestUserService_GetAccount(t *testing.T) {
	t.Parallel()
	type args struct {
		userID    uuid.UUID
		accountID uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.Account
		setup     func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want *domain.Account)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				userID:    util.UUID(),
				accountID: util.UUID(),
			},
			want: &domain.Account{
				ID:          uuid.Nil, // setupで変更,
				Type:        domain.HOMEPAGE,
				PrPermitted: true,
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want *domain.Account) {
				want.ID = args.accountID
				repo.EXPECT().GetAccount(args.userID, args.accountID).Return(want, nil)
			},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockUserRepository(ctrl)
			event := mock_repository.NewMockEventRepository(ctrl)
			tt.setup(repo, event, tt.args, tt.want)

			s := NewUserService(repo, event)
			got, err := s.GetAccount(tt.args.userID, tt.args.accountID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserService_GetAccounts(t *testing.T) {
	t.Parallel()
	type args struct {
		userID uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      []*domain.Account
		setup     func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.Account)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{userID: util.UUID()},
			want: []*domain.Account{
				{
					ID:          util.UUID(),
					Type:        domain.HOMEPAGE,
					PrPermitted: true,
				},
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.Account) {
				repo.EXPECT().GetAccounts(args.userID).Return(want, nil)
			},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockUserRepository(ctrl)
			event := mock_repository.NewMockEventRepository(ctrl)
			tt.setup(repo, event, tt.args, tt.want)

			s := NewUserService(repo, event)
			got, err := s.GetAccounts(tt.args.userID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserService_CreateAccount(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx     context.Context
		id      uuid.UUID
		account *repository.CreateAccountArgs
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.Account
		setup     func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want *domain.Account)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  util.UUID(),
				account: &repository.CreateAccountArgs{
					ID:          util.UUID(),
					Name:        util.AlphaNumeric(5),
					Type:        domain.HOMEPAGE,
					URL:         "https://" + util.AlphaNumeric(10),
					PrPermitted: true,
				},
			},
			want: &domain.Account{
				ID:          util.UUID(),
				Type:        domain.HOMEPAGE,
				PrPermitted: true,
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want *domain.Account) {
				repo.EXPECT().CreateAccount(args.id, args.account).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "EmptyID",
			args: args{
				ctx: context.Background(),
				id:  util.UUID(),
				account: &repository.CreateAccountArgs{
					ID:          uuid.Nil,
					Name:        "",
					Type:        domain.HOMEPAGE,
					URL:         "https://" + util.AlphaNumeric(10),
					PrPermitted: true,
				},
			},
			want: nil,
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want *domain.Account) {
			},
			assertion: assert.Error,
		},
		{
			name: "InvalidAccountType",
			args: args{
				ctx: context.Background(),
				id:  util.UUID(),
				account: &repository.CreateAccountArgs{
					ID:          util.UUID(),
					Name:        util.AlphaNumeric(5),
					Type:        10000,
					URL:         "https://" + util.AlphaNumeric(10),
					PrPermitted: true,
				},
			},
			want: nil,
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want *domain.Account) {
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockUserRepository(ctrl)
			event := mock_repository.NewMockEventRepository(ctrl)
			tt.setup(repo, event, tt.args, tt.want)

			s := NewUserService(repo, event)
			got, err := s.CreateAccount(tt.args.ctx, tt.args.id, tt.args.account)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserService_EditAccount(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx       context.Context
		accountID uuid.UUID
		userID    uuid.UUID
		args      *repository.UpdateAccountArgs
	}
	tests := []struct {
		name      string
		args      args
		setup     func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:       context.Background(),
				accountID: util.UUID(),
				userID:    util.UUID(),
				args: &repository.UpdateAccountArgs{
					ID:          optional.StringFrom(util.AlphaNumeric(5)),
					Type:        optional.Int64From(int64(domain.HOMEPAGE)),
					URL:         optional.StringFrom("https://" + util.AlphaNumeric(10)),
					PrPermitted: optional.BoolFrom(true),
				},
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args) {
				changes := map[string]interface{}{
					"id":    args.args.ID.String,
					"url":   args.args.URL.String,
					"check": args.args.PrPermitted.Bool,
					"type":  args.args.Type.Int64,
				}
				repo.EXPECT().UpdateAccount(args.accountID, args.userID, changes).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Notfound",
			args: args{
				ctx:       context.Background(),
				accountID: util.UUID(),
				userID:    util.UUID(),
				args: &repository.UpdateAccountArgs{
					ID:          optional.StringFrom(util.AlphaNumeric(5)),
					Type:        optional.Int64From(int64(domain.HOMEPAGE)),
					URL:         optional.StringFrom("https://" + util.AlphaNumeric(10)),
					PrPermitted: optional.BoolFrom(true),
				},
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args) {
				changes := map[string]interface{}{
					"id":    args.args.ID.String,
					"url":   args.args.URL.String,
					"check": args.args.PrPermitted.Bool,
					"type":  args.args.Type.Int64,
				}
				repo.EXPECT().UpdateAccount(args.accountID, args.userID, changes).Return(repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockUserRepository(ctrl)
			event := mock_repository.NewMockEventRepository(ctrl)
			tt.setup(repo, event, tt.args)

			s := NewUserService(repo, event)
			tt.assertion(t, s.EditAccount(tt.args.ctx, tt.args.accountID, tt.args.userID, tt.args.args))
		})
	}
}

func TestUserService_DeleteAccount(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx       context.Context
		accountid uuid.UUID
		userid    uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		setup     func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:       context.Background(),
				accountid: util.UUID(),
				userid:    util.UUID(),
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args) {
				repo.EXPECT().DeleteAccount(args.accountid, args.userid).Return(nil)
			},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockUserRepository(ctrl)
			event := mock_repository.NewMockEventRepository(ctrl)
			tt.setup(repo, event, tt.args)

			s := NewUserService(repo, event)
			tt.assertion(t, s.DeleteAccount(tt.args.ctx, tt.args.accountid, tt.args.userid))
		})
	}
}

func TestUserService_GetUserProjects(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx    context.Context
		userID uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      []*domain.UserProject
		setup     func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.UserProject)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:    context.Background(),
				userID: util.UUID(),
			},
			want: []*domain.UserProject{
				{
					ID:        util.UUID(),
					Name:      util.AlphaNumeric(5),
					Since:     time.Now(),
					Until:     time.Now(),
					UserSince: time.Now(),
					UserUntil: time.Now(),
				},
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.UserProject) {
				repo.EXPECT().GetProjects(args.userID).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Notfound",
			args: args{
				ctx:    context.Background(),
				userID: util.UUID(),
			},
			want: nil,
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.UserProject) {
				repo.EXPECT().GetProjects(args.userID).Return(want, repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockUserRepository(ctrl)
			event := mock_repository.NewMockEventRepository(ctrl)
			tt.setup(repo, event, tt.args, tt.want)

			s := NewUserService(repo, event)
			got, err := s.GetUserProjects(tt.args.ctx, tt.args.userID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserService_GetUserContests(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx    context.Context
		userID uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      []*domain.UserContest
		setup     func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.UserContest)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:    context.Background(),
				userID: util.UUID(),
			},
			want: []*domain.UserContest{
				{
					ID:          util.UUID(),
					Name:        util.AlphaNumeric(5),
					Result:      util.AlphaNumeric(5),
					ContestName: util.AlphaNumeric(5),
				},
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.UserContest) {
				repo.EXPECT().GetContests(args.userID).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Notfound",
			args: args{
				ctx:    context.Background(),
				userID: util.UUID(),
			},
			want: nil,
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.UserContest) {
				repo.EXPECT().GetContests(args.userID).Return(want, repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockUserRepository(ctrl)
			event := mock_repository.NewMockEventRepository(ctrl)
			tt.setup(repo, event, tt.args, tt.want)

			s := NewUserService(repo, event)
			got, err := s.GetUserContests(tt.args.ctx, tt.args.userID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserService_GetUserEvents(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx    context.Context
		userID uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      []*domain.Event
		setup     func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.Event)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:    context.Background(),
				userID: util.UUID(),
			},
			want: []*domain.Event{
				{
					ID:        util.UUID(),
					Name:      util.AlphaNumeric(5),
					TimeStart: time.Now(),
					TimeEnd:   time.Now(),
				},
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.Event) {
				event.EXPECT().GetUserEvents(args.userID).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Notfound",
			args: args{
				ctx:    context.Background(),
				userID: util.UUID(),
			},
			want: nil,
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.Event) {
				event.EXPECT().GetUserEvents(args.userID).Return(want, repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_repository.NewMockUserRepository(ctrl)
			event := mock_repository.NewMockEventRepository(ctrl)
			tt.setup(repo, event, tt.args, tt.want)

			s := NewUserService(repo, event)
			got, err := s.GetUserEvents(tt.args.ctx, tt.args.userID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
