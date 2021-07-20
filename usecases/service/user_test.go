package service

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util"
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
