package service

import (
	"context"
	"testing"

	"github.com/traPtitech/traPortfolio/util/random"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/usecases/repository"
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
		setup     func(group *mock_repository.MockGroupRepository, user *mock_repository.MockUserRepository, args args, want []*domain.Group)
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
			setup: func(group *mock_repository.MockGroupRepository, user *mock_repository.MockUserRepository, args args, want []*domain.Group) {
				group.EXPECT().GetAllGroups().Return(want, nil)
			},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			group := mock_repository.NewMockGroupRepository(ctrl)
			user := mock_repository.NewMockUserRepository(ctrl)
			tt.setup(group, user, tt.args, tt.want)

			s := NewGroupService(group, user)
			got, err := s.GetAllGroups(tt.args.ctx)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGroupService_GetGroup(t *testing.T) {
	t.Parallel()
	groupID := random.UUID()
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.GroupDetail
		setup     func(group *mock_repository.MockGroupRepository, user *mock_repository.MockUserRepository, args args, want *domain.GroupDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{ctx: context.Background(), id: groupID},
			want: &domain.GroupDetail{
				ID:   groupID,
				Name: random.AlphaNumeric(),
				Link: random.AlphaNumeric(),
				Admin: []*domain.User{
					domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
				},
				Members: []*domain.UserWithDuration{
					{
						User:     *domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
						Duration: random.Duration(),
					},
				},
				Description: random.AlphaNumeric(),
			},
			setup: func(group *mock_repository.MockGroupRepository, user *mock_repository.MockUserRepository, args args, want *domain.GroupDetail) {
				group.EXPECT().GetGroup(groupID).Return(&domain.GroupDetail{
					ID:   groupID,
					Name: want.Name,
					Link: want.Link,
					Admin: []*domain.User{
						{
							ID: want.Admin[0].ID,
						},
					},
					Members: []*domain.UserWithDuration{
						{
							User: domain.User{
								ID: want.Members[0].User.ID,
							},
							Duration: want.Members[0].Duration,
						},
					},
					Description: want.Description,
				}, nil)
				user.EXPECT().GetUsers(&repository.GetUsersArgs{}).Return([]*domain.User{
					want.Admin[0],
					&want.Members[0].User,
				}, nil)
			},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			group := mock_repository.NewMockGroupRepository(ctrl)
			user := mock_repository.NewMockUserRepository(ctrl)
			tt.setup(group, user, tt.args, tt.want)

			s := NewGroupService(group, user)
			got, err := s.GetGroup(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
