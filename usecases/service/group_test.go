package service

import (
	"context"
	"testing"

	"github.com/traPtitech/traPortfolio/util/random"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository/mock_repository"
)

func TestGroupService_GetAllGroups(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		want      []*domain.Group
		setup     func(repo *mock_repository.MockGroupRepository, user *mock_repository.MockUserRepository, args args, want []*domain.Group)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{ctx: context.Background()},
			want: []*domain.Group{
				{
					ID:   random.UUID(),
					Name: random.AlphaNumeric(),
				},
			},
			setup: func(repo *mock_repository.MockGroupRepository, user *mock_repository.MockUserRepository, args args, want []*domain.Group) {
				repo.EXPECT().GetAllGroups().Return(want, nil)
			},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			repo := mock_repository.NewMockGroupRepository(ctrl)
			user := mock_repository.NewMockUserRepository(ctrl)
			tt.setup(repo, user, tt.args, tt.want)

			s := NewGroupService(repo, user)
			got, err := s.GetAllGroups(tt.args.ctx)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGroupService_GetGroup(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.GroupDetail
		setup     func(repo *mock_repository.MockGroupRepository, user *mock_repository.MockUserRepository, args args, want *domain.GroupDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{ctx: context.Background(), id: random.UUID()},
			want: &domain.GroupDetail{
				ID:          random.UUID(),
				Name:        random.AlphaNumeric(),
				Link:        random.AlphaNumeric(),
				Admin:       []*domain.User{},
				Members:     []*domain.UserGroup{},
				Description: random.AlphaNumeric(),
			},
			setup: func(repo *mock_repository.MockGroupRepository, user *mock_repository.MockUserRepository, args args, want *domain.GroupDetail) {
				repo.EXPECT().GetGroup(args.id).Return(want, nil)
			},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			repo := mock_repository.NewMockGroupRepository(ctrl)
			user := mock_repository.NewMockUserRepository(ctrl)
			tt.setup(repo, user, tt.args, tt.want)

			s := NewGroupService(repo, user)
			got, err := s.GetGroup(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
