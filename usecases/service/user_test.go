package service

import (
	"context"
	"testing"
	"time"

	"github.com/traPtitech/traPortfolio/util/random"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/usecases/repository/mock_repository"
	"github.com/traPtitech/traPortfolio/util/optional"
)

func TestUserService_GetUsers(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		args *repository.GetUsersArgs
	}
	tests := []struct {
		name      string
		args      args
		want      []*domain.User
		setup     func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.User)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success_NoOpts",
			args: args{
				ctx:  context.Background(),
				args: &repository.GetUsersArgs{},
			},
			want: []*domain.User{
				{
					ID:       random.UUID(),
					Name:     random.AlphaNumeric(),
					RealName: random.AlphaNumeric(),
				},
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.User) {
				repo.EXPECT().GetUsers(args.args).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Forbidden",
			args: args{
				ctx:  context.Background(),
				args: &repository.GetUsersArgs{},
			},
			want: nil,
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.User) {
				repo.EXPECT().GetUsers(args.args).Return(want, repository.ErrForbidden)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			repo := mock_repository.NewMockUserRepository(ctrl)
			event := mock_repository.NewMockEventRepository(ctrl)
			tt.setup(repo, event, tt.args, tt.want)

			s := NewUserService(repo, event)
			got, err := s.GetUsers(tt.args.ctx, tt.args.args)
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
				id:  random.UUID(),
			},
			want: &domain.UserDetail{
				User: domain.User{
					ID:       uuid.Nil, // setupで変更する
					Name:     random.AlphaNumeric(),
					RealName: random.AlphaNumeric(),
				},
				State: domain.TraqStateActive,
				Bio:   random.AlphaNumeric(),
				Accounts: []*domain.Account{
					{
						ID:          random.UUID(),
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
				id:  random.UUID(),
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
				id:  random.UUID(),
				args: &repository.UpdateUserArgs{
					Description: optional.NewString(random.AlphaNumeric(), true),
					Check:       optional.NewBool(true, true),
				},
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args) {
				repo.EXPECT().UpdateUser(args.id, args.args).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Notfound",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
				args: &repository.UpdateUserArgs{
					Description: optional.NewString(random.AlphaNumeric(), true),
					Check:       optional.NewBool(true, true),
				},
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args) {
				repo.EXPECT().UpdateUser(args.id, args.args).Return(repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

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
				userID:    random.UUID(),
				accountID: random.UUID(),
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
			args: args{userID: random.UUID()},
			want: []*domain.Account{
				{
					ID:          random.UUID(),
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
				id:  random.UUID(),
				account: &repository.CreateAccountArgs{
					DisplayName: random.AlphaNumeric(),
					Type:        domain.HOMEPAGE,
					URL:         "https://" + random.AlphaNumeric(),
					PrPermitted: true,
				},
			},
			want: &domain.Account{
				ID:          random.UUID(),
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
				id:  random.UUID(),
				account: &repository.CreateAccountArgs{
					DisplayName: "",
					Type:        domain.HOMEPAGE,
					URL:         "https://" + random.AlphaNumeric(),
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
				id:  random.UUID(),
				account: &repository.CreateAccountArgs{
					DisplayName: random.AlphaNumeric(),
					Type:        10000,
					URL:         "https://" + random.AlphaNumeric(),
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
		userID    uuid.UUID
		accountID uuid.UUID
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
				userID:    random.UUID(),
				accountID: random.UUID(),
				args: &repository.UpdateAccountArgs{
					DisplayName: optional.NewString(random.AlphaNumeric(), true),
					Type:        optional.NewInt64(int64(domain.HOMEPAGE), true),
					URL:         optional.NewString(random.RandURLString(), true),
					PrPermitted: optional.NewBool(true, true),
				},
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args) {
				repo.EXPECT().UpdateAccount(args.userID, args.accountID, args.args).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Notfound",
			args: args{
				ctx:       context.Background(),
				userID:    random.UUID(),
				accountID: random.UUID(),
				args: &repository.UpdateAccountArgs{
					DisplayName: optional.NewString(random.AlphaNumeric(), true),
					Type:        optional.NewInt64(int64(domain.HOMEPAGE), true),
					URL:         optional.NewString(random.RandURLString(), true),
					PrPermitted: optional.NewBool(true, true),
				},
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args) {
				repo.EXPECT().UpdateAccount(args.userID, args.accountID, args.args).Return(repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			repo := mock_repository.NewMockUserRepository(ctrl)
			event := mock_repository.NewMockEventRepository(ctrl)
			tt.setup(repo, event, tt.args)

			s := NewUserService(repo, event)
			tt.assertion(t, s.EditAccount(tt.args.ctx, tt.args.userID, tt.args.accountID, tt.args.args))
		})
	}
}

func TestUserService_DeleteAccount(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx       context.Context
		userID    uuid.UUID
		accountID uuid.UUID
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
				userID:    random.UUID(),
				accountID: random.UUID(),
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args) {
				repo.EXPECT().DeleteAccount(args.userID, args.accountID).Return(nil)
			},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			repo := mock_repository.NewMockUserRepository(ctrl)
			event := mock_repository.NewMockEventRepository(ctrl)
			tt.setup(repo, event, tt.args)

			s := NewUserService(repo, event)
			tt.assertion(t, s.DeleteAccount(tt.args.ctx, tt.args.userID, tt.args.accountID))
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
				userID: random.UUID(),
			},
			want: []*domain.UserProject{
				{
					ID:           random.UUID(),
					Name:         random.AlphaNumeric(),
					Duration:     random.Duration(),
					UserDuration: random.Duration(),
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
				userID: random.UUID(),
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
				userID: random.UUID(),
			},
			want: []*domain.UserContest{
				{
					ID:          random.UUID(),
					Name:        random.AlphaNumeric(),
					Result:      random.AlphaNumeric(),
					ContestName: random.AlphaNumeric(),
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
				userID: random.UUID(),
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

func TestUserService_GetGroupsByUserID(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx    context.Context
		userID uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      []*domain.GroupUser
		setup     func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.GroupUser)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:    context.Background(),
				userID: random.UUID(),
			},
			want: []*domain.GroupUser{
				{
					ID:       random.UUID(),
					Name:     random.AlphaNumeric(),
					Duration: random.Duration(),
				},
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.GroupUser) {
				repo.EXPECT().GetGroupsByUserID(args.userID).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Notfound",
			args: args{
				ctx:    context.Background(),
				userID: random.UUID(),
			},
			want: nil,
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.GroupUser) {
				repo.EXPECT().GetGroupsByUserID(args.userID).Return(want, repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			repo := mock_repository.NewMockUserRepository(ctrl)
			event := mock_repository.NewMockEventRepository(ctrl)
			tt.setup(repo, event, tt.args, tt.want)

			s := NewUserService(repo, event)
			got, err := s.GetGroupsByUserID(tt.args.ctx, tt.args.userID)
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
				userID: random.UUID(),
			},
			want: []*domain.Event{
				{
					ID:        random.UUID(),
					Name:      random.AlphaNumeric(),
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
				userID: random.UUID(),
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
