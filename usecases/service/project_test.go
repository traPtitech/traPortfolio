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
		setup     func(repo *mock_repository.MockProjectRepository, args args, want []*domain.Project)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{ctx: context.Background()},
			want: []*domain.Project{
				{
					ID:       random.UUID(),
					Name:     random.AlphaNumeric(),
					Duration: random.Duration(),
				},
			},
			setup: func(repo *mock_repository.MockProjectRepository, args args, want []*domain.Project) {
				repo.EXPECT().GetProjects(args.ctx).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "ErrInvalidDB",
			args: args{ctx: context.Background()},
			want: nil,
			setup: func(repo *mock_repository.MockProjectRepository, args args, want []*domain.Project) {
				repo.EXPECT().GetProjects(args.ctx).Return(nil, gorm.ErrInvalidDB)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			repo := mock_repository.NewMockProjectRepository(ctrl)
			tt.setup(repo, tt.args, tt.want)

			s := NewProjectService(repo)
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
		want      *domain.ProjectDetail
		setup     func(repo *mock_repository.MockProjectRepository, args args, want *domain.ProjectDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
			},
			want: &domain.ProjectDetail{
				Project: domain.Project{
					ID:       random.UUID(),
					Name:     random.AlphaNumeric(),
					Duration: random.Duration(),
				},
				Description: random.AlphaNumeric(),
				Link:        random.RandURLString(),
				Members: []*domain.UserWithDuration{
					{
						User:     *domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
						Duration: random.Duration(),
					},
				},
			},
			setup: func(repo *mock_repository.MockProjectRepository, args args, want *domain.ProjectDetail) {
				repo.EXPECT().GetProject(args.ctx, args.id).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "InvalidDB",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
			},
			want: nil,
			setup: func(repo *mock_repository.MockProjectRepository, args args, want *domain.ProjectDetail) {
				repo.EXPECT().GetProject(args.ctx, args.id).Return(nil, gorm.ErrInvalidDB)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			repo := mock_repository.NewMockProjectRepository(ctrl)
			tt.setup(repo, tt.args, tt.want)

			s := NewProjectService(repo)
			got, err := s.GetProject(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProjectService_CreateProject(t *testing.T) {
	var (
		name        = random.AlphaNumeric()
		description = random.AlphaNumeric()
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
		want      *domain.ProjectDetail
		setup     func(repo *mock_repository.MockProjectRepository, args args, want *domain.ProjectDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				args: &repository.CreateProjectArgs{
					Name:          name,
					Description:   description,
					Link:          optional.From(link),
					SinceYear:     duration.Since.Year,
					SinceSemester: duration.Since.Semester,
					UntilYear:     duration.Until.Year,
					UntilSemester: duration.Until.Semester,
				},
			},
			want: &domain.ProjectDetail{
				Project: domain.Project{
					ID:       random.UUID(),
					Name:     name,
					Duration: duration,
				},
				Description: description,
				Link:        link,
			},
			setup: func(repo *mock_repository.MockProjectRepository, args args, want *domain.ProjectDetail) {
				want.Link = args.args.Link.ValueOr(want.Link)
				repo.EXPECT().CreateProject(args.ctx, gomock.Any()).Return(want, nil) // TODO: CreateProject内でuuid.NewV4するのでテストができない？
			},
			assertion: assert.NoError,
		},
		{
			name: "InvalidDuration",
			args: args{
				ctx: context.Background(),
				args: &repository.CreateProjectArgs{
					Name:          random.AlphaNumeric(),
					Description:   random.AlphaNumeric(),
					Link:          random.Optional(random.RandURLString()),
					SinceYear:     duration.Until.Year,
					SinceSemester: duration.Until.Semester,
					UntilYear:     duration.Since.Year,
					UntilSemester: duration.Since.Semester,
				},
			},
			want: nil,
			setup: func(repo *mock_repository.MockProjectRepository, args args, want *domain.ProjectDetail) {
			},
			assertion: assert.Error,
		},
		{
			name: "InvalidDB",
			args: args{
				ctx: context.Background(),
				args: &repository.CreateProjectArgs{
					Name:          random.AlphaNumeric(),
					Description:   random.AlphaNumeric(),
					Link:          random.Optional(random.RandURLString()),
					SinceYear:     duration.Since.Year,
					SinceSemester: duration.Since.Semester,
					UntilYear:     duration.Until.Year,
					UntilSemester: duration.Until.Semester,
				},
			},
			want: nil,
			setup: func(repo *mock_repository.MockProjectRepository, args args, want *domain.ProjectDetail) {
				repo.EXPECT().CreateProject(args.ctx, args.args).Return(nil, gorm.ErrInvalidDB)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			repo := mock_repository.NewMockProjectRepository(ctrl)
			tt.setup(repo, tt.args, tt.want)

			s := NewProjectService(repo)
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
		setup     func(repo *mock_repository.MockProjectRepository, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
				args: &repository.UpdateProjectArgs{
					Name:          random.Optional(random.AlphaNumeric()),
					Description:   random.Optional(random.AlphaNumeric()),
					Link:          random.Optional(random.AlphaNumeric()),
					SinceYear:     optional.From(int64(duration.Since.Year)),
					SinceSemester: optional.From(int64(duration.Since.Semester)),
					UntilYear:     optional.From(int64(duration.Until.Year)),
					UntilSemester: optional.From(int64(duration.Until.Semester)),
				},
			},
			setup: func(repo *mock_repository.MockProjectRepository, args args) {
				repo.EXPECT().GetProject(args.ctx, args.id).Return(&domain.ProjectDetail{
					Project: domain.Project{
						ID:       args.id,
						Duration: duration,
					},
				}, nil)
				repo.EXPECT().UpdateProject(args.ctx, args.id, args.args).Return(nil)
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
			setup: func(repo *mock_repository.MockProjectRepository, args args) {
				repo.EXPECT().GetProject(args.ctx, args.id).Return(nil, repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
		{
			name: "InvalidDuration",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
				args: &repository.UpdateProjectArgs{
					Name:          random.Optional(random.AlphaNumeric()),
					Description:   random.Optional(random.AlphaNumeric()),
					Link:          random.Optional(random.AlphaNumeric()),
					SinceYear:     optional.From(int64(duration.Until.Year)),
					SinceSemester: optional.From(int64(duration.Until.Semester)),
					UntilYear:     optional.From(int64(duration.Since.Year)),
					UntilSemester: optional.From(int64(duration.Since.Semester)),
				},
			},
			setup: func(repo *mock_repository.MockProjectRepository, args args) {
				repo.EXPECT().GetProject(args.ctx, args.id).Return(&domain.ProjectDetail{
					Project: domain.Project{
						ID:       args.id,
						Duration: duration,
					},
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
					Name:          random.Optional(random.AlphaNumeric()),
					Description:   random.Optional(random.AlphaNumeric()),
					Link:          random.Optional(random.AlphaNumeric()),
					SinceYear:     optional.From(int64(duration.Since.Year)),
					SinceSemester: optional.From(int64(duration.Since.Semester)),
					UntilYear:     optional.From(int64(duration.Until.Year)),
					UntilSemester: optional.From(int64(duration.Until.Semester)),
				},
			},
			setup: func(repo *mock_repository.MockProjectRepository, args args) {
				repo.EXPECT().GetProject(args.ctx, args.id).Return(&domain.ProjectDetail{
					Project: domain.Project{
						ID:       args.id,
						Duration: duration,
					},
				}, nil)
				repo.EXPECT().UpdateProject(args.ctx, args.id, args.args).Return(gorm.ErrInvalidDB)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			repo := mock_repository.NewMockProjectRepository(ctrl)
			tt.setup(repo, tt.args)

			s := NewProjectService(repo)

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
		want      []*domain.UserWithDuration
		setup     func(repo *mock_repository.MockProjectRepository, args args, want []*domain.UserWithDuration)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
			},
			want: []*domain.UserWithDuration{
				{
					User: *domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
				},
			},
			setup: func(repo *mock_repository.MockProjectRepository, args args, want []*domain.UserWithDuration) {
				repo.EXPECT().GetProjectMembers(args.ctx, args.id).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "InvalidDB",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
			},
			want: nil,
			setup: func(repo *mock_repository.MockProjectRepository, args args, want []*domain.UserWithDuration) {
				repo.EXPECT().GetProjectMembers(args.ctx, args.id).Return(nil, gorm.ErrInvalidDB)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			repo := mock_repository.NewMockProjectRepository(ctrl)
			tt.setup(repo, tt.args, tt.want)

			s := NewProjectService(repo)
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
		setup     func(repo *mock_repository.MockProjectRepository, args args)
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
			setup: func(repo *mock_repository.MockProjectRepository, args args) {
				repo.EXPECT().AddProjectMembers(args.ctx, args.projectID, args.args).Return(nil)
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
			setup: func(repo *mock_repository.MockProjectRepository, args args) {
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
			setup: func(repo *mock_repository.MockProjectRepository, args args) {
				repo.EXPECT().AddProjectMembers(args.ctx, args.projectID, args.args).Return(gorm.ErrInvalidDB)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			repo := mock_repository.NewMockProjectRepository(ctrl)
			tt.setup(repo, tt.args)

			s := NewProjectService(repo)

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
		setup     func(repo *mock_repository.MockProjectRepository, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:       context.Background(),
				projectID: random.UUID(),
				memberIDs: []uuid.UUID{random.UUID()},
			},
			setup: func(repo *mock_repository.MockProjectRepository, args args) {
				repo.EXPECT().DeleteProjectMembers(args.ctx, args.projectID, args.memberIDs).Return(nil)
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
			setup: func(repo *mock_repository.MockProjectRepository, args args) {
				repo.EXPECT().DeleteProjectMembers(args.ctx, args.projectID, args.memberIDs).Return(gorm.ErrInvalidDB)
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			repo := mock_repository.NewMockProjectRepository(ctrl)
			tt.setup(repo, tt.args)

			s := NewProjectService(repo)

			tt.assertion(t, s.DeleteProjectMembers(tt.args.ctx, tt.args.projectID, tt.args.memberIDs))
		})
	}
}
