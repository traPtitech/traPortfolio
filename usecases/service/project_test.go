package service

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

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
					Duration:    random.Duration(),
					Description: random.AlphaNumeric(rand.Intn(30) + 1),
					Link:        random.RandURLString(),
					Members: []*domain.ProjectMember{
						{
							UserID:   random.UUID(),
							Name:     random.AlphaNumeric(rand.Intn(30) + 1),
							RealName: random.AlphaNumeric(rand.Intn(30) + 1),
							Duration: random.Duration(),
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
				Duration:    random.Duration(),
				Description: random.AlphaNumeric(rand.Intn(30) + 1),
				Link:        random.RandURLString(),
				Members: []*domain.ProjectMember{
					{
						UserID:   random.UUID(),
						Name:     random.AlphaNumeric(rand.Intn(30) + 1),
						RealName: random.AlphaNumeric(rand.Intn(30) + 1),
						Duration: random.Duration(),
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
	var (
		name        = random.AlphaNumeric(rand.Intn(30) + 1)
		description = random.AlphaNumeric(rand.Intn(30) + 1)
		link        = random.RandURLString()
		duration    = random.Duration()
	)

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
					Name:          name,
					Description:   description,
					Link:          optional.NewString(link, true),
					SinceYear:     duration.Since.Year,
					SinceSemester: duration.Since.Semester,
					UntilYear:     duration.Until.Year,
					UntilSemester: duration.Until.Semester,
				},
			},
			want: &domain.Project{
				ID:          random.UUID(),
				Name:        name,
				Description: description,
				Link:        link,
				Duration:    duration,
			},
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args, want *domain.Project) {
				if args.args.Link.Valid {
					want.Link = args.args.Link.String
				}
				repo.EXPECT().CreateProject(gomock.Any()).Return(want, nil) // TODO: CreateProject内でuuid.NewV4するのでテストができない？
			},
			assertion: assert.NoError,
		},
		{
			name: "InvalidDuration",
			args: args{
				ctx: context.Background(),
				args: &repository.CreateProjectArgs{
					Name:          random.AlphaNumeric(rand.Intn(30) + 1),
					Description:   random.AlphaNumeric(rand.Intn(30) + 1),
					Link:          optional.NewString(random.RandURLString(), true),
					SinceYear:     duration.Until.Year,
					SinceSemester: duration.Until.Semester,
					UntilYear:     duration.Since.Year,
					UntilSemester: duration.Since.Semester,
				},
			},
			want: nil,
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args, want *domain.Project) {
			},
			assertion: assert.Error,
		},
		{
			name: "InvalidDB",
			args: args{
				ctx: context.Background(),
				args: &repository.CreateProjectArgs{
					Name:          random.AlphaNumeric(rand.Intn(30) + 1),
					Description:   random.AlphaNumeric(rand.Intn(30) + 1),
					Link:          optional.NewString(random.RandURLString(), true),
					SinceYear:     duration.Since.Year,
					SinceSemester: duration.Since.Semester,
					UntilYear:     duration.Until.Year,
					UntilSemester: duration.Until.Semester,
				},
			},
			want: nil,
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args, want *domain.Project) {
				repo.EXPECT().CreateProject(args.args).Return(nil, gorm.ErrInvalidDB)
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
	duration := random.Duration()

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
					Name:          optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Description:   optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Link:          optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					SinceYear:     optional.NewInt64(int64(duration.Since.Year), true),
					SinceSemester: optional.NewInt64(int64(duration.Since.Semester), true),
					UntilYear:     optional.NewInt64(int64(duration.Until.Year), true),
					UntilSemester: optional.NewInt64(int64(duration.Until.Semester), true),
				},
			},
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args) {
				repo.EXPECT().GetProject(args.id).Return(&domain.Project{
					ID:       args.id,
					Duration: duration,
				}, nil)
				changes := map[string]interface{}{
					"name":           args.args.Name.String,
					"description":    args.args.Description.String,
					"link":           args.args.Link.String,
					"since_year":     args.args.SinceYear.Int64,
					"since_semester": args.args.SinceSemester.Int64,
					"until_year":     args.args.UntilYear.Int64,
					"until_semester": args.args.UntilSemester.Int64,
				}
				repo.EXPECT().UpdateProject(args.id, changes).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "ErrFind",
			args: args{
				ctx:  context.Background(),
				id:   random.UUID(),
				args: nil,
			},
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args) {
				repo.EXPECT().GetProject(args.id).Return(nil, repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
		{
			name: "InvalidDuration",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
				args: &repository.UpdateProjectArgs{
					Name:          optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Description:   optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Link:          optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					SinceYear:     optional.NewInt64(int64(duration.Until.Year), true),
					SinceSemester: optional.NewInt64(int64(duration.Until.Semester), true),
					UntilYear:     optional.NewInt64(int64(duration.Since.Year), true),
					UntilSemester: optional.NewInt64(int64(duration.Since.Semester), true),
				},
			},
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args) {
				repo.EXPECT().GetProject(args.id).Return(&domain.Project{
					ID:       args.id,
					Duration: duration,
				}, nil)
			},
			assertion: assert.Error,
		},
		{
			name: "ErrUpdate",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
				args: &repository.UpdateProjectArgs{
					Name:          optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Description:   optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Link:          optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					SinceYear:     optional.NewInt64(int64(duration.Since.Year), true),
					SinceSemester: optional.NewInt64(int64(duration.Since.Semester), true),
					UntilYear:     optional.NewInt64(int64(duration.Until.Year), true),
					UntilSemester: optional.NewInt64(int64(duration.Until.Semester), true),
				},
			},
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args) {
				repo.EXPECT().GetProject(args.id).Return(&domain.Project{
					ID:       args.id,
					Duration: duration,
				}, nil)
				changes := map[string]interface{}{
					"name":           args.args.Name.String,
					"description":    args.args.Description.String,
					"link":           args.args.Link.String,
					"since_year":     args.args.SinceYear.Int64,
					"since_semester": args.args.SinceSemester.Int64,
					"until_year":     args.args.UntilYear.Int64,
					"until_semester": args.args.UntilSemester.Int64,
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
	duration := random.Duration()

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
						UserID:        random.UUID(),
						SinceYear:     duration.Since.Year,
						SinceSemester: duration.Since.Semester,
						UntilYear:     duration.Until.Year,
						UntilSemester: duration.Until.Semester,
					},
				},
			},
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args) {
				repo.EXPECT().AddProjectMembers(args.projectID, args.args).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "InvalidDuration",
			args: args{
				ctx:       context.Background(),
				projectID: random.UUID(),
				args: []*repository.CreateProjectMemberArgs{
					{
						UserID:        random.UUID(),
						SinceYear:     duration.Until.Year,
						SinceSemester: duration.Until.Semester,
						UntilYear:     duration.Since.Year,
						UntilSemester: duration.Since.Semester,
					},
				},
			},
			setup: func(repo *mock_repository.MockProjectRepository, portal *mock_repository.MockPortalRepository, args args) {
			},
			assertion: assert.Error,
		},
		{
			name: "InvalidDB",
			args: args{
				ctx:       context.Background(),
				projectID: random.UUID(),
				args: []*repository.CreateProjectMemberArgs{
					{
						UserID:        random.UUID(),
						SinceYear:     duration.Since.Year,
						SinceSemester: duration.Since.Semester,
						UntilYear:     duration.Until.Year,
						UntilSemester: duration.Until.Semester,
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
