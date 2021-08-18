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
		setup     func(m *MockRepository, args args, want []*domain.User)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{ctx: context.Background()},
			want: []*domain.User{
				{
					ID:       uuid.Must(uuid.NewV4()),
					Name:     util.AlphaNumeric(5),
					RealName: util.AlphaNumeric(5),
				},
			},
			setup: func(m *MockRepository, args args, want []*domain.User) {
				m.user.EXPECT().GetUsers().Return(want, nil)
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

			repo := newMockRepository(ctrl)
			tt.setup(repo, tt.args, tt.want)

			s := NewUserService(repo.user, repo.event)
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
		setup     func(m *MockRepository, args args, want *domain.UserDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  uuid.Must(uuid.NewV4()),
			},
			want: &domain.UserDetail{
				ID:       uuid.Nil, // setupで変更する
				Name:     util.AlphaNumeric(5),
				RealName: util.AlphaNumeric(5),
				State:    0,
				Bio:      util.AlphaNumeric(10),
				Accounts: []*domain.Account{
					{
						ID:          uuid.Must(uuid.NewV4()),
						Type:        0,
						PrPermitted: true,
					},
				},
			},
			setup: func(m *MockRepository, args args, want *domain.UserDetail) {
				want.ID = args.id
				m.user.EXPECT().GetUser(args.id).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "NilID",
			args: args{
				ctx: context.Background(),
				id:  uuid.Nil,
			},
			want: nil,
			setup: func(m *MockRepository, args args, want *domain.UserDetail) {
				m.user.EXPECT().GetUser(args.id).Return(nil, repository.ErrInvalidID)
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

			repo := newMockRepository(ctrl)
			tt.setup(repo, tt.args, tt.want)

			s := NewUserService(repo.user, repo.event)
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
		setup     func(m *MockRepository, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  uuid.Must(uuid.NewV4()),
				args: &repository.UpdateUserArgs{
					Description: optional.StringFrom(util.AlphaNumeric(10)),
					Check:       optional.BoolFrom(true),
				},
			},
			setup: func(m *MockRepository, args args) {
				changes := map[string]interface{}{
					"description": args.args.Description.String,
					"check":       args.args.Check.Bool,
				}
				m.user.EXPECT().Update(args.id, changes).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Success_OnlyDescription",
			args: args{
				ctx: context.Background(),
				id:  uuid.Must(uuid.NewV4()),
				args: &repository.UpdateUserArgs{
					Description: optional.StringFrom(util.AlphaNumeric(10)),
				},
			},
			setup: func(m *MockRepository, args args) {
				changes := map[string]interface{}{
					"description": args.args.Description.String,
				}
				m.user.EXPECT().Update(args.id, changes).Return(nil)
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

			repo := newMockRepository(ctrl)
			tt.setup(repo, tt.args)

			s := NewUserService(repo.user, repo.event)

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
		setup     func(m *MockRepository, args args, want *domain.Account)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				userID:    uuid.Must(uuid.NewV4()),
				accountID: uuid.Must(uuid.NewV4()),
			},
			want: &domain.Account{
				ID:          uuid.Nil, // setupで変更,
				Type:        domain.HOMEPAGE,
				PrPermitted: true,
			},
			setup: func(m *MockRepository, args args, want *domain.Account) {
				want.ID = args.accountID
				m.user.EXPECT().GetAccount(args.userID, args.accountID).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "NilID",
			args: args{
				userID:    uuid.Must(uuid.NewV4()),
				accountID: uuid.Nil,
			},
			want:      nil,
			setup:     func(m *MockRepository, args args, want *domain.Account) {},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := newMockRepository(ctrl)
			tt.setup(repo, tt.args, tt.want)

			s := NewUserService(repo.user, repo.event)
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
		setup     func(m *MockRepository, args args, want []*domain.Account)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{userID: uuid.Must(uuid.NewV4())},
			want: []*domain.Account{
				{
					ID:          uuid.Must(uuid.NewV4()),
					Type:        domain.HOMEPAGE,
					PrPermitted: true,
				},
			},
			setup: func(m *MockRepository, args args, want []*domain.Account) {
				m.user.EXPECT().GetAccounts(args.userID).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name:      "NilID",
			args:      args{userID: uuid.Nil},
			want:      nil,
			setup:     func(m *MockRepository, args args, want []*domain.Account) {},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := newMockRepository(ctrl)
			tt.setup(repo, tt.args, tt.want)

			s := NewUserService(repo.user, repo.event)
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
		setup     func(m *MockRepository, args args, want *domain.Account)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  uuid.Must(uuid.NewV4()),
				account: &repository.CreateAccountArgs{
					ID:          util.AlphaNumeric(5),
					Type:        domain.HOMEPAGE,
					URL:         "https://" + util.AlphaNumeric(10),
					PrPermitted: true,
				},
			},
			want: &domain.Account{
				ID:          uuid.Must(uuid.NewV4()),
				Type:        domain.HOMEPAGE,
				PrPermitted: true,
			},
			setup: func(m *MockRepository, args args, want *domain.Account) {
				m.user.EXPECT().CreateAccount(args.id, args.account).Return(want, nil)
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

			repo := newMockRepository(ctrl)
			tt.setup(repo, tt.args, tt.want)

			s := NewUserService(repo.user, repo.event)
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
		setup     func(m *MockRepository, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:       context.Background(),
				accountID: uuid.Must(uuid.NewV4()),
				userID:    uuid.Must(uuid.NewV4()),
				args: &repository.UpdateAccountArgs{
					ID:          optional.StringFrom(util.AlphaNumeric(5)),
					Type:        optional.Int64From(int64(domain.HOMEPAGE)),
					URL:         optional.StringFrom("https://" + util.AlphaNumeric(10)),
					PrPermitted: optional.BoolFrom(true),
				},
			},
			setup: func(m *MockRepository, args args) {
				changes := map[string]interface{}{
					"id":    args.args.ID.String,
					"url":   args.args.URL.String,
					"check": args.args.PrPermitted.Bool,
					"type":  args.args.Type.Int64,
				}
				m.user.EXPECT().UpdateAccount(args.accountID, args.userID, changes).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "NoChanges",
			args: args{
				ctx:       context.Background(),
				accountID: uuid.Must(uuid.NewV4()),
				userID:    uuid.Must(uuid.NewV4()),
				args:      &repository.UpdateAccountArgs{},
			},
			setup:     func(m *MockRepository, args args) {},
			assertion: assert.NoError,
		},
		{
			name: "NilID",
			args: args{
				ctx:       context.Background(),
				accountID: uuid.Nil,
				userID:    uuid.Must(uuid.NewV4()),
				args: &repository.UpdateAccountArgs{
					ID:          optional.StringFrom(util.AlphaNumeric(5)),
					Type:        optional.Int64From(int64(domain.HOMEPAGE)),
					URL:         optional.StringFrom("https://" + util.AlphaNumeric(10)),
					PrPermitted: optional.BoolFrom(true),
				},
			},
			setup:     func(m *MockRepository, args args) {},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := newMockRepository(ctrl)
			tt.setup(repo, tt.args)

			s := NewUserService(repo.user, repo.event)

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
		setup     func(m *MockRepository, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:       context.Background(),
				accountid: uuid.Must(uuid.NewV4()),
				userid:    uuid.Must(uuid.NewV4()),
			},
			setup: func(m *MockRepository, args args) {
				m.user.EXPECT().DeleteAccount(args.accountid, args.userid).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "NilID",
			args: args{
				ctx:       context.Background(),
				accountid: uuid.Nil,
				userid:    uuid.Must(uuid.NewV4()),
			},
			setup:     func(m *MockRepository, args args) {},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := newMockRepository(ctrl)
			tt.setup(repo, tt.args)

			s := NewUserService(repo.user, repo.event)

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
		setup     func(m *MockRepository, args args, want []*domain.UserProject)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:    context.Background(),
				userID: uuid.Must(uuid.NewV4()),
			},
			want: []*domain.UserProject{
				{
					ID:        uuid.Must(uuid.NewV4()),
					Name:      util.AlphaNumeric(5),
					Since:     time.Now(),
					Until:     time.Now(),
					UserSince: time.Now(),
					UserUntil: time.Now(),
				},
			},
			setup: func(m *MockRepository, args args, want []*domain.UserProject) {
				m.user.EXPECT().GetProjects(args.userID).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "NilID",
			args: args{
				ctx:    context.Background(),
				userID: uuid.Nil,
			},
			want:      nil,
			setup:     func(m *MockRepository, args args, want []*domain.UserProject) {},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := newMockRepository(ctrl)
			tt.setup(repo, tt.args, tt.want)

			s := NewUserService(repo.user, repo.event)
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
		setup     func(m *MockRepository, args args, want []*domain.UserContest)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:    context.Background(),
				userID: uuid.Must(uuid.NewV4()),
			},
			want: []*domain.UserContest{
				{
					ID:          uuid.Must(uuid.NewV4()),
					Name:        util.AlphaNumeric(5),
					Result:      util.AlphaNumeric(5),
					ContestName: util.AlphaNumeric(5),
				},
			},
			setup: func(m *MockRepository, args args, want []*domain.UserContest) {
				m.user.EXPECT().GetContests(args.userID).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "NilID",
			args: args{
				ctx:    context.Background(),
				userID: uuid.Nil,
			},
			want:      nil,
			setup:     func(m *MockRepository, args args, want []*domain.UserContest) {},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := newMockRepository(ctrl)
			tt.setup(repo, tt.args, tt.want)

			s := NewUserService(repo.user, repo.event)
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
		setup     func(m *MockRepository, args args, want []*domain.Event)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:    context.Background(),
				userID: uuid.Must(uuid.NewV4()),
			},
			want: []*domain.Event{
				{
					ID:        uuid.Must(uuid.NewV4()),
					Name:      util.AlphaNumeric(5),
					TimeStart: time.Now(),
					TimeEnd:   time.Now(),
				},
			},
			setup: func(m *MockRepository, args args, want []*domain.Event) {
				m.event.EXPECT().GetUserEvents(args.userID).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "NilID",
			args: args{
				ctx:    context.Background(),
				userID: uuid.Nil,
			},
			want:      nil,
			setup:     func(m *MockRepository, args args, want []*domain.Event) {},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := newMockRepository(ctrl)
			tt.setup(repo, tt.args, tt.want)

			s := NewUserService(repo.user, repo.event)
			got, err := s.GetUserEvents(tt.args.ctx, tt.args.userID)

			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
