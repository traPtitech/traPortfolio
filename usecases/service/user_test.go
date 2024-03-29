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
				domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.User) {
				repo.EXPECT().GetUsers(args.ctx, args.args).Return(want, nil)
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
				repo.EXPECT().GetUsers(args.ctx, args.args).Return(want, repository.ErrForbidden)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
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
				// IDはsetupで変更する
				User:  *domain.NewUser(uuid.Nil, random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
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
				repo.EXPECT().GetUser(args.ctx, args.id).Return(want, nil)
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
				repo.EXPECT().GetUser(args.ctx, args.id).Return(nil, repository.ErrForbidden)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
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
			name: "Success/AllFields",
			args: args{
				ctx:  context.Background(),
				id:   random.UUID(),
				args: random.UpdateUserArgs(),
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args) {
				repo.EXPECT().UpdateUser(args.ctx, args.id, args.args).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Success/PartialFields",
			args: args{
				ctx:  context.Background(),
				id:   random.UUID(),
				args: random.OptUpdateUserArgs(),
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args) {
				repo.EXPECT().UpdateUser(args.ctx, args.id, args.args).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Notfound",
			args: args{
				ctx:  context.Background(),
				id:   random.UUID(),
				args: random.OptUpdateUserArgs(),
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args) {
				repo.EXPECT().UpdateUser(args.ctx, args.id, args.args).Return(repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
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
		ctx       context.Context
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
				ctx:       context.Background(),
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
				repo.EXPECT().GetAccount(args.ctx, args.userID, args.accountID).Return(want, nil)
			},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			repo := mock_repository.NewMockUserRepository(ctrl)
			event := mock_repository.NewMockEventRepository(ctrl)
			tt.setup(repo, event, tt.args, tt.want)

			s := NewUserService(repo, event)
			got, err := s.GetAccount(tt.args.ctx, tt.args.userID, tt.args.accountID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserService_GetAccounts(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx    context.Context
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
			args: args{
				ctx:    context.Background(),
				userID: random.UUID(),
			},
			want: []*domain.Account{
				{
					ID:          random.UUID(),
					Type:        domain.HOMEPAGE,
					PrPermitted: true,
				},
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.Account) {
				repo.EXPECT().GetAccounts(args.ctx, args.userID).Return(want, nil)
			},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			repo := mock_repository.NewMockUserRepository(ctrl)
			event := mock_repository.NewMockEventRepository(ctrl)
			tt.setup(repo, event, tt.args, tt.want)

			s := NewUserService(repo, event)
			got, err := s.GetAccounts(tt.args.ctx, tt.args.userID)
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
				repo.EXPECT().CreateAccount(args.ctx, args.id, args.account).Return(want, nil)
			},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
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
			name: "Success/AllFields",
			args: args{
				ctx:       context.Background(),
				userID:    random.UUID(),
				accountID: random.UUID(),
				args:      random.UpdateAccountArgs(),
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args) {
				repo.EXPECT().UpdateAccount(args.ctx, args.userID, args.accountID, args.args).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Success/PartialFields",
			args: args{
				ctx:       context.Background(),
				userID:    random.UUID(),
				accountID: random.UUID(),
				args:      random.OptUpdateAccountArgs(),
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args) {
				repo.EXPECT().UpdateAccount(args.ctx, args.userID, args.accountID, args.args).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Notfound",
			args: args{
				ctx:       context.Background(),
				userID:    random.UUID(),
				accountID: random.UUID(),
				args:      random.OptUpdateAccountArgs(),
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args) {
				repo.EXPECT().UpdateAccount(args.ctx, args.userID, args.accountID, args.args).Return(repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
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
				repo.EXPECT().DeleteAccount(args.ctx, args.userID, args.accountID).Return(nil)
			},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
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
				repo.EXPECT().GetProjects(args.ctx, args.userID).Return(want, nil)
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
				repo.EXPECT().GetProjects(args.ctx, args.userID).Return(want, repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
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
					ID:        random.UUID(),
					Name:      random.AlphaNumeric(),
					TimeStart: random.Time(),
					TimeEnd:   random.Time(),
					Teams: []*domain.ContestTeamWithoutMembers{
						{
							ID:     random.UUID(),
							Name:   random.AlphaNumeric(),
							Result: random.AlphaNumeric(),
						},
					},
				},
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.UserContest) {
				repo.EXPECT().GetContests(args.ctx, args.userID).Return(want, nil)
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
				repo.EXPECT().GetContests(args.ctx, args.userID).Return(want, repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
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
		want      []*domain.UserGroup
		setup     func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.UserGroup)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:    context.Background(),
				userID: random.UUID(),
			},
			want: []*domain.UserGroup{
				{
					ID:       random.UUID(),
					Name:     random.AlphaNumeric(),
					Duration: random.Duration(),
				},
			},
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.UserGroup) {
				repo.EXPECT().GetGroupsByUserID(args.ctx, args.userID).Return(want, nil)
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
			setup: func(repo *mock_repository.MockUserRepository, event *mock_repository.MockEventRepository, args args, want []*domain.UserGroup) {
				repo.EXPECT().GetGroupsByUserID(args.ctx, args.userID).Return(want, repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
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
				event.EXPECT().GetUserEvents(args.ctx, args.userID).Return(want, nil)
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
				event.EXPECT().GetUserEvents(args.ctx, args.userID).Return(want, repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
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
