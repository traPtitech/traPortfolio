package service

import (
	"context"
	"fmt"
	"math/rand"
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
	"gorm.io/gorm"
)

func TestProjectService_GetProjects(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		args      args
		want      []*domain.Project
		setup     func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args, want []*domain.Project)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{ctx: context.Background()},
			want: []*domain.Project{
				{
					ID:          random.UUID(),
					Name:        random.AlphaNumeric(rand.Intn(30) + 1),
					Since:       time.Now(),
					Until:       time.Now(),
					Description: random.AlphaNumeric(rand.Intn(30) + 1),
					Link:        random.RandURLString(),
					Members: []*domain.ProjectMember{
						{
							UserID:   random.UUID(),
							Name:     random.AlphaNumeric(rand.Intn(30) + 1),
							RealName: random.AlphaNumeric(rand.Intn(30) + 1),
							Since:    time.Now(),
							Until:    time.Now(),
						},
					},
				},
			},
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args, want []*domain.Project) {
				repo.EXPECT().GetProjects().Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "ErrInvalidDB",
			args: args{ctx: context.Background()},
			want: nil,
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args, want []*domain.Project) {
				repo.EXPECT().GetProjects().Return(nil, gorm.ErrInvalidDB)
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

			repo := mock_repository.NewMockProjectRepository(ctrl)
			portal := mock_repository.NewMockPortalRepository(ctrl)
			tt.setup(repo, portal, tt.args, tt.want)

			s := NewProjectService(repo, portal)
			got, err := s.GetProjects(tt.args.ctx)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProjectService_GetProject(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.Project
		setup     func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args, want *domain.Project)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
			},
			want: &domain.Project{
				ID:          random.UUID(),
				Name:        random.AlphaNumeric(rand.Intn(30) + 1),
				Since:       time.Now(),
				Until:       time.Now(),
				Description: random.AlphaNumeric(rand.Intn(30) + 1),
				Link:        random.RandURLString(),
				Members: []*domain.ProjectMember{
					{
						UserID:   random.UUID(),
						Name:     random.AlphaNumeric(rand.Intn(30) + 1),
						RealName: random.AlphaNumeric(rand.Intn(30) + 1),
						Since:    time.Now(),
						Until:    time.Now(),
					},
				},
			},
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args, want *domain.Project) {
				repo.EXPECT().GetProject(args.id).Return(want, nil)
				portalWant := make([]*domain.PortalUser, 0)
				for _, v := range want.Members {
					portalWant = append(portalWant, &domain.PortalUser{
						ID:   v.Name,
						Name: v.RealName,
					})
				}
				portal.EXPECT().GetUsers(args.ctx).Return(portalWant, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "PortalForbidden",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
			},
			want: nil,
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args, want *domain.Project) {
				repo.EXPECT().GetProject(args.id).Return(want, nil)
				portal.EXPECT().GetUsers(args.ctx).Return(nil, fmt.Errorf("GET /user failed: %v", "forbidden"))
			},
			assertion: assert.Error,
		},
		{
			name: "InvalidDB",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
			},
			want: nil,
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args, want *domain.Project) {
				repo.EXPECT().GetProject(args.id).Return(nil, gorm.ErrInvalidDB)
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

			repo := mock_repository.NewMockProjectRepository(ctrl)
			portal := mock_repository.NewMockPortalRepository(ctrl)
			tt.setup(repo, portal, tt.args, tt.want)

			s := NewProjectService(repo, portal)
			got, err := s.GetProject(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProjectService_CreateProject(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		args *repository.CreateProjectArgs
	}
	tests := []struct {
		name      string
		args      args
		want      *domain.Project
		setup     func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args, want *domain.Project)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				args: &repository.CreateProjectArgs{
					Name:        random.AlphaNumeric(rand.Intn(30) + 1),
					Description: random.AlphaNumeric(rand.Intn(30) + 1),
					Link:        random.RandURLString(),
					Since:       time.Now(),
					Until:       time.Now(),
				},
			},
			want: &domain.Project{
				ID:          random.UUID(),
				Name:        "",
				Description: "",
				Link:        "",
				Since:       time.Time{},
				Until:       time.Time{},
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args, want *domain.Project) {
				want.Name = args.args.Name
				want.Description = args.args.Description
				want.Link = args.args.Link
				want.Since = args.args.Since
				want.Until = args.args.Until
				repo.EXPECT().CreateProject(gomock.Any()).Return(want, nil) // TODO: CreateProject内でuuid.NewV4するのでテストができない？
			},
			assertion: assert.NoError,
		},
		{
			name: "InvalidDB",
			args: args{
				ctx: context.Background(),
				args: &repository.CreateProjectArgs{
					Name:        random.AlphaNumeric(rand.Intn(30) + 1),
					Description: random.AlphaNumeric(rand.Intn(30) + 1),
					Link:        random.RandURLString(),
					Since:       time.Now(),
					Until:       time.Now(),
				},
			},
			want: nil,
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args, want *domain.Project) {
				repo.EXPECT().CreateProject(gomock.Any()).Return(nil, gorm.ErrInvalidDB) // TODO: CreateProject内でuuid.NewV4するのでテストができない？
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

			repo := mock_repository.NewMockProjectRepository(ctrl)
			portal := mock_repository.NewMockPortalRepository(ctrl)
			tt.setup(repo, portal, tt.args, tt.want)

			s := NewProjectService(repo, portal)
			got, err := s.CreateProject(tt.args.ctx, tt.args.args)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProjectService_UpdateProject(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx  context.Context
		id   uuid.UUID
		args *repository.UpdateProjectArgs
	}
	tests := []struct {
		name      string
		args      args
		setup     func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
				args: &repository.UpdateProjectArgs{
					Name:        optional.StringFrom(random.AlphaNumeric(rand.Intn(30) + 1)),
					Description: optional.StringFrom(random.AlphaNumeric(rand.Intn(30) + 1)),
					Link:        optional.StringFrom(random.AlphaNumeric(rand.Intn(30) + 1)),
					Since:       optional.TimeFrom(time.Now()),
					Until:       optional.TimeFrom(time.Now()),
				},
			},
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args) {
				changes := map[string]interface{}{
					"name":        args.args.Name.String,
					"description": args.args.Description.String,
					"link":        args.args.Link.String,
					"since":       args.args.Since.Time,
					"until":       args.args.Until.Time,
				}
				repo.EXPECT().UpdateProject(args.id, changes).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "InvalidDB",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
				args: &repository.UpdateProjectArgs{
					Name:        optional.StringFrom(random.AlphaNumeric(rand.Intn(30) + 1)),
					Description: optional.StringFrom(random.AlphaNumeric(rand.Intn(30) + 1)),
					Link:        optional.StringFrom(random.AlphaNumeric(rand.Intn(30) + 1)),
					Since:       optional.TimeFrom(time.Now()),
					Until:       optional.TimeFrom(time.Now()),
				},
			},
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args) {
				changes := map[string]interface{}{
					"name":        args.args.Name.String,
					"description": args.args.Description.String,
					"link":        args.args.Link.String,
					"since":       args.args.Since.Time,
					"until":       args.args.Until.Time,
				}
				repo.EXPECT().UpdateProject(args.id, changes).Return(gorm.ErrInvalidDB)
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

			repo := mock_repository.NewMockProjectRepository(ctrl)
			portal := mock_repository.NewMockPortalRepository(ctrl)
			tt.setup(repo, portal, tt.args)

			s := NewProjectService(repo, portal)

			tt.assertion(t, s.UpdateProject(tt.args.ctx, tt.args.id, tt.args.args))
		})
	}
}

func TestProjectService_GetProjectMembers(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		want      []*domain.User
		setup     func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args, want []*domain.User)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
			},
			want: []*domain.User{
				{
					ID:       random.UUID(),
					Name:     random.AlphaNumeric(rand.Intn(30) + 1),
					RealName: random.AlphaNumeric(rand.Intn(30) + 1),
				},
			},
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args, want []*domain.User) {
				repo.EXPECT().GetProjectMembers(args.id).Return(want, nil)
				portalWant := make([]*domain.PortalUser, 0)
				for _, v := range want {
					portalWant = append(portalWant, &domain.PortalUser{
						ID:   v.Name,
						Name: v.RealName,
					})
				}
				portal.EXPECT().GetUsers(args.ctx).Return(portalWant, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "PortalForbidden",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
			},
			want: nil,
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args, want []*domain.User) {
				repo.EXPECT().GetProjectMembers(args.id).Return(want, nil)
				portal.EXPECT().GetUsers(args.ctx).Return(nil, fmt.Errorf("GET /user failed: %v", "forbidden"))
			},
			assertion: assert.Error,
		},
		{
			name: "InvalidDB",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
			},
			want: nil,
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args, want []*domain.User) {
				repo.EXPECT().GetProjectMembers(args.id).Return(nil, gorm.ErrInvalidDB)
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

			repo := mock_repository.NewMockProjectRepository(ctrl)
			portal := mock_repository.NewMockPortalRepository(ctrl)
			tt.setup(repo, portal, tt.args, tt.want)

			s := NewProjectService(repo, portal)
			got, err := s.GetProjectMembers(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProjectService_AddProjectMembers(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx       context.Context
		projectID uuid.UUID
		args      []*repository.CreateProjectMemberArgs
	}
	tests := []struct {
		name      string
		args      args
		setup     func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:       context.Background(),
				projectID: random.UUID(),
				args: []*repository.CreateProjectMemberArgs{
					{
						UserID: random.UUID(),
						Since:  time.Now(),
						Until:  time.Now(),
					},
				},
			},
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args) {
				repo.EXPECT().AddProjectMembers(args.projectID, args.args).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "InvalidDB",
			args: args{
				ctx:       context.Background(),
				projectID: random.UUID(),
				args: []*repository.CreateProjectMemberArgs{
					{
						UserID: random.UUID(),
						Since:  time.Now(),
						Until:  time.Now(),
					},
				},
			},
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args) {
				repo.EXPECT().AddProjectMembers(args.projectID, args.args).Return(gorm.ErrInvalidDB)
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

			repo := mock_repository.NewMockProjectRepository(ctrl)
			portal := mock_repository.NewMockPortalRepository(ctrl)
			tt.setup(repo, portal, tt.args)

			s := NewProjectService(repo, portal)

			tt.assertion(t, s.AddProjectMembers(tt.args.ctx, tt.args.projectID, tt.args.args))
		})
	}
}

func TestProjectService_DeleteProjectMembers(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx       context.Context
		projectID uuid.UUID
		memberIDs []uuid.UUID
	}
	tests := []struct {
		name      string
		args      args
		setup     func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:       context.Background(),
				projectID: random.UUID(),
				memberIDs: []uuid.UUID{random.UUID()},
			},
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args) {
				repo.EXPECT().DeleteProjectMembers(args.projectID, args.memberIDs).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "InvalidDB",
			args: args{
				ctx:       context.Background(),
				projectID: random.UUID(),
				memberIDs: []uuid.UUID{random.UUID()},
			},
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args) {
				repo.EXPECT().DeleteProjectMembers(args.projectID, args.memberIDs).Return(gorm.ErrInvalidDB)
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

			repo := mock_repository.NewMockProjectRepository(ctrl)
			portal := mock_repository.NewMockPortalRepository(ctrl)
			tt.setup(repo, portal, tt.args)

			s := NewProjectService(repo, portal)

			tt.assertion(t, s.DeleteProjectMembers(tt.args.ctx, tt.args.projectID, tt.args.memberIDs))
		})
	}
}
