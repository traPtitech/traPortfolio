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
	"github.com/traPtitech/traPortfolio/util/optional"
	"github.com/traPtitech/traPortfolio/util/random"
)

func TestContestService_GetContests(t *testing.T) {
	t.Parallel()
	type fields struct {
		repo repository.ContestRepository
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []*domain.Contest
		setup     func(f fields, args args, want []*domain.Contest)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
			},
			want: []*domain.Contest{
				{
					ID:        random.UUID(),
					Name:      random.AlphaNumeric(),
					TimeStart: time.Now(),
					TimeEnd:   time.Now(),
				},
			},
			setup: func(f fields, args args, want []*domain.Contest) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().GetContests(args.ctx).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Error_FindContests",
			args: args{
				ctx: context.Background(),
			},
			want: nil,
			setup: func(f fields, args args, want []*domain.Contest) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().GetContests(args.ctx).Return(nil, repository.ErrForbidden)
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			tt.fields = fields{
				repo: mock_repository.NewMockContestRepository(ctrl),
			}
			tt.setup(tt.fields, tt.args, tt.want)
			s := NewContestService(tt.fields.repo)
			// Assertion
			got, err := s.GetContests(tt.args.ctx)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContestService_GetContest(t *testing.T) {
	cid := random.UUID() // contestID

	t.Parallel()
	type fields struct {
		repo repository.ContestRepository
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *domain.ContestDetail
		setup     func(f fields, args args, want *domain.ContestDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  cid,
			},
			want: &domain.ContestDetail{
				Contest: domain.Contest{
					ID:        cid,
					Name:      random.AlphaNumeric(),
					TimeStart: time.Now(),
					TimeEnd:   time.Now(),
				},
				Link:        random.RandURLString(),
				Description: random.AlphaNumeric(),
				ContestTeams: []*domain.ContestTeam{
					{
						ContestTeamWithoutMembers: domain.ContestTeamWithoutMembers{
							ID:        random.UUID(),
							ContestID: cid,
							Name:      random.AlphaNumeric(),
							Result:    random.AlphaNumeric(),
						},
					},
				},
			},
			setup: func(f fields, args args, want *domain.ContestDetail) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().GetContest(args.ctx, args.id).Return(want, nil)
				repo.EXPECT().GetContestTeams(args.ctx, args.id).Return(want.ContestTeams, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "Success_TeamNotFound",
			args: args{
				ctx: context.Background(),
				id:  cid,
			},
			want: &domain.ContestDetail{
				Contest: domain.Contest{
					ID:        cid,
					Name:      random.AlphaNumeric(),
					TimeStart: time.Now(),
					TimeEnd:   time.Now(),
				},
				Link:         random.RandURLString(),
				Description:  random.AlphaNumeric(),
				ContestTeams: nil,
			},
			setup: func(f fields, args args, want *domain.ContestDetail) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().GetContest(args.ctx, args.id).Return(want, nil)
				repo.EXPECT().GetContestTeams(args.ctx, args.id).Return(nil, repository.ErrNotFound)
			},
			assertion: assert.NoError,
		},
		{
			name: "Error_FindContest",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
			},
			want: nil,
			setup: func(f fields, args args, want *domain.ContestDetail) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().GetContest(args.ctx, args.id).Return(nil, repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
		{
			name: "Error_FindContestTeams",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
			},
			want: nil,
			setup: func(f fields, args args, want *domain.ContestDetail) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().GetContest(args.ctx, args.id).Return(want, nil)
				repo.EXPECT().GetContestTeams(args.ctx, args.id).Return(nil, repository.ErrInvalidID)
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			tt.fields = fields{
				repo: mock_repository.NewMockContestRepository(ctrl),
			}
			tt.setup(tt.fields, tt.args, tt.want)
			s := NewContestService(tt.fields.repo)
			// Assertion
			got, err := s.GetContest(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContestService_CreateContest(t *testing.T) {
	cname := random.AlphaNumeric() // 作成するコンテストのコンテスト名

	t.Parallel()
	type fields struct {
		repo repository.ContestRepository
	}
	type args struct {
		ctx  context.Context
		args *repository.CreateContestArgs
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *domain.ContestDetail
		setup     func(f fields, args args, want *domain.ContestDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				args: &repository.CreateContestArgs{
					Name:        cname,
					Description: random.AlphaNumeric(),
					Link:        random.Optional(random.RandURLString()),
					Since:       time.Now(),
					Until:       random.Optional(random.Time()),
				},
			},
			want: &domain.ContestDetail{
				Contest: domain.Contest{
					ID:        random.UUID(),
					Name:      cname,
					TimeStart: time.Now(),
					TimeEnd:   time.Now(),
				},
				Link:         random.AlphaNumeric(),
				Description:  random.AlphaNumeric(),
				ContestTeams: []*domain.ContestTeam{},
			},
			setup: func(f fields, args args, want *domain.ContestDetail) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().CreateContest(args.ctx, args.args).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "ErrCreate",
			args: args{
				ctx: context.Background(),
				args: &repository.CreateContestArgs{
					Name:        cname,
					Description: random.AlphaNumeric(),
					Link:        random.Optional(random.AlphaNumeric()),
					Since:       time.Now(),
					Until:       random.Optional(random.Time()),
				},
			},
			want: nil,
			setup: func(f fields, args args, want *domain.ContestDetail) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().CreateContest(args.ctx, args.args).Return(nil, repository.ErrInvalidArg)
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			tt.fields = fields{
				repo: mock_repository.NewMockContestRepository(ctrl),
			}
			tt.setup(tt.fields, tt.args, tt.want)
			s := NewContestService(tt.fields.repo)
			// Assertion
			got, err := s.CreateContest(tt.args.ctx, tt.args.args)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContestService_UpdateContest(t *testing.T) {
	t.Parallel()
	type fields struct {
		repo repository.ContestRepository
	}
	type args struct {
		ctx  context.Context
		id   uuid.UUID
		args *repository.UpdateContestArgs
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		setup     func(f fields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
				args: &repository.UpdateContestArgs{
					Name:        random.Optional(random.AlphaNumeric()),
					Description: random.Optional(random.AlphaNumeric()),
					Link:        random.Optional(random.RandURLString()),
					Since:       random.Optional(random.Time()),
					Until:       random.Optional(random.Time()),
				},
			},
			setup: func(f fields, args args) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().UpdateContest(args.ctx, args.id, args.args).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "ErrUpdate",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
				args: &repository.UpdateContestArgs{
					Name:        random.Optional(random.AlphaNumeric()),
					Description: random.Optional(random.AlphaNumeric()),
					Link:        random.Optional(random.RandURLString()),
					Since:       random.Optional(random.Time()),
					Until:       random.Optional(random.Time()),
				},
			},
			setup: func(f fields, args args) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().UpdateContest(args.ctx, args.id, args.args).Return(repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			tt.fields = fields{
				repo: mock_repository.NewMockContestRepository(ctrl),
			}
			tt.setup(tt.fields, tt.args)
			s := NewContestService(tt.fields.repo)
			// Assertion
			tt.assertion(t, s.UpdateContest(tt.args.ctx, tt.args.id, tt.args.args))
		})
	}
}

func TestContestService_DeleteContest(t *testing.T) {
	t.Parallel()
	type fields struct {
		repo repository.ContestRepository
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		setup     func(f fields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
			},
			setup: func(f fields, args args) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().DeleteContest(args.ctx, args.id).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "ErrDelete",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
			},
			setup: func(f fields, args args) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().DeleteContest(args.ctx, args.id).Return(repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			tt.fields = fields{
				repo: mock_repository.NewMockContestRepository(ctrl),
			}
			tt.setup(tt.fields, tt.args)
			s := NewContestService(tt.fields.repo)
			// Assertion
			tt.assertion(t, s.DeleteContest(tt.args.ctx, tt.args.id))
		})
	}
}

func TestContestService_GetContestTeams(t *testing.T) {
	cid := random.UUID() // contestID

	t.Parallel()
	type fields struct {
		repo repository.ContestRepository
	}
	type args struct {
		ctx       context.Context
		contestID uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []*domain.ContestTeam
		setup     func(f fields, args args, want []*domain.ContestTeam)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:       context.Background(),
				contestID: cid,
			},
			want: []*domain.ContestTeam{
				{
					ContestTeamWithoutMembers: domain.ContestTeamWithoutMembers{
						ID:        random.UUID(),
						ContestID: cid,
						Name:      random.AlphaNumeric(),
						Result:    random.AlphaNumeric(),
					},
				},
			},
			setup: func(f fields, args args, want []*domain.ContestTeam) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().GetContestTeams(args.ctx, args.contestID).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "ErrGetByContestID",
			args: args{
				ctx:       context.Background(),
				contestID: cid,
			},
			want: nil,
			setup: func(f fields, args args, want []*domain.ContestTeam) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().GetContestTeams(args.ctx, args.contestID).Return(nil, repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			tt.fields = fields{
				repo: mock_repository.NewMockContestRepository(ctrl),
			}
			tt.setup(tt.fields, tt.args, tt.want)
			s := NewContestService(tt.fields.repo)
			// Assertion
			got, err := s.GetContestTeams(tt.args.ctx, tt.args.contestID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_contestService_GetContestTeam(t *testing.T) {
	cid := random.UUID()
	tid := random.UUID()

	t.Parallel()
	type fields struct {
		repo repository.ContestRepository
	}
	type args struct {
		ctx       context.Context
		contestID uuid.UUID
		teamID    uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *domain.ContestTeamDetail
		setup     func(f fields, args args, want *domain.ContestTeamDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:       context.Background(),
				contestID: cid,
				teamID:    tid,
			},
			want: &domain.ContestTeamDetail{
				ContestTeam: domain.ContestTeam{
					ContestTeamWithoutMembers: domain.ContestTeamWithoutMembers{
						ID:        tid,
						ContestID: cid,
						Name:      random.AlphaNumeric(),
						Result:    random.AlphaNumeric(),
					},
					Members: []*domain.User{
						domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
					},
				},
				Link:        random.RandURLString(),
				Description: random.AlphaNumeric(),
			},
			setup: func(f fields, args args, want *domain.ContestTeamDetail) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().GetContestTeam(args.ctx, args.contestID, args.teamID).Return(&domain.ContestTeamDetail{
					ContestTeam: domain.ContestTeam{
						ContestTeamWithoutMembers: domain.ContestTeamWithoutMembers{
							ID:        args.teamID,
							ContestID: args.contestID,
							Name:      want.Name,
							Result:    want.Result,
						},
					},
					Link:        want.Link,
					Description: want.Description,
				}, nil)
				repo.EXPECT().GetContestTeamMembers(args.ctx, args.contestID, args.teamID).Return(want.Members, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "ErrGetContestTeam",
			args: args{
				ctx:       context.Background(),
				contestID: random.UUID(),
				teamID:    random.UUID(),
			},
			want: nil,
			setup: func(f fields, args args, want *domain.ContestTeamDetail) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().GetContestTeam(args.ctx, args.contestID, args.teamID).Return(nil, repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
		{
			name: "ErrGetContestTeamMembers",
			args: args{
				ctx:       context.Background(),
				contestID: random.UUID(),
				teamID:    random.UUID(),
			},
			want: nil,
			setup: func(f fields, args args, want *domain.ContestTeamDetail) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().GetContestTeam(args.ctx, args.contestID, args.teamID).Return(&domain.ContestTeamDetail{
					ContestTeam: domain.ContestTeam{
						ContestTeamWithoutMembers: domain.ContestTeamWithoutMembers{
							ID:        args.teamID,
							ContestID: args.contestID,
							Name:      random.AlphaNumeric(),
							Result:    random.AlphaNumeric(),
						},
					},
					Link:        random.RandURLString(),
					Description: random.AlphaNumeric(),
				}, nil)
				repo.EXPECT().GetContestTeamMembers(args.ctx, args.contestID, args.teamID).Return(nil, repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			tt.fields = fields{
				repo: mock_repository.NewMockContestRepository(ctrl),
			}
			tt.setup(tt.fields, tt.args, tt.want)
			s := NewContestService(tt.fields.repo)
			// Assertion
			got, err := s.GetContestTeam(tt.args.ctx, tt.args.contestID, tt.args.teamID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContestService_CreateContestTeam(t *testing.T) {
	var (
		cid         = random.UUID()
		tid         = random.UUID()
		name        = random.AlphaNumeric()
		result      = random.AlphaNumeric()
		link        = random.RandURLString()
		description = random.AlphaNumeric()
	)

	t.Parallel()
	type fields struct {
		repo repository.ContestRepository
	}
	type args struct {
		ctx       context.Context
		contestID uuid.UUID
		args      *repository.CreateContestTeamArgs
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *domain.ContestTeamDetail
		setup     func(f fields, args args, want *domain.ContestTeamDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:       context.Background(),
				contestID: cid,
				args: &repository.CreateContestTeamArgs{
					Name:        name,
					Result:      optional.From(result),
					Link:        optional.From(link),
					Description: description,
				},
			},
			want: &domain.ContestTeamDetail{
				ContestTeam: domain.ContestTeam{
					ContestTeamWithoutMembers: domain.ContestTeamWithoutMembers{
						ID:        tid,
						ContestID: cid,
						Name:      name,
						Result:    result,
					},
					Members: nil,
				},
				Link:        link,
				Description: description,
			},
			setup: func(f fields, args args, want *domain.ContestTeamDetail) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().CreateContestTeam(args.ctx, args.contestID, args.args).Return(&domain.ContestTeamDetail{
					ContestTeam: domain.ContestTeam{
						ContestTeamWithoutMembers: domain.ContestTeamWithoutMembers{
							ID:        tid,
							ContestID: args.contestID,
							Name:      args.args.Name,
							Result:    args.args.Result.ValueOrZero(),
						},
						Members: nil,
					},
					Link:        args.args.Link.ValueOrZero(),
					Description: args.args.Description,
				}, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "ErrCreateContestTeam",
			args: args{
				ctx:       context.Background(),
				contestID: random.UUID(),
				args: &repository.CreateContestTeamArgs{
					Name:        random.AlphaNumeric(),
					Result:      random.Optional(random.AlphaNumeric()),
					Link:        random.Optional(random.RandURLString()),
					Description: random.AlphaNumeric(),
				},
			},
			want: nil,
			setup: func(f fields, args args, want *domain.ContestTeamDetail) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().CreateContestTeam(args.ctx, args.contestID, args.args).Return(nil, repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			tt.fields = fields{
				repo: mock_repository.NewMockContestRepository(ctrl),
			}
			tt.setup(tt.fields, tt.args, tt.want)
			s := NewContestService(tt.fields.repo)
			// Assertion
			got, err := s.CreateContestTeam(tt.args.ctx, tt.args.contestID, tt.args.args)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContestService_UpdateContestTeam(t *testing.T) {
	t.Parallel()
	type fields struct {
		repo repository.ContestRepository
	}
	type args struct {
		ctx    context.Context
		teamID uuid.UUID
		args   *repository.UpdateContestTeamArgs
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		setup     func(f fields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:    context.Background(),
				teamID: random.UUID(),
				args: &repository.UpdateContestTeamArgs{
					Name:        random.Optional(random.AlphaNumeric()),
					Result:      random.Optional(random.AlphaNumeric()),
					Link:        random.Optional(random.RandURLString()),
					Description: random.Optional(random.AlphaNumeric()),
				},
			},
			setup: func(f fields, args args) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().UpdateContestTeam(args.ctx, args.teamID, args.args).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "NoChange",
			args: args{
				ctx:    context.Background(),
				teamID: random.UUID(),
				args:   &repository.UpdateContestTeamArgs{},
			},
			setup: func(f fields, args args) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().UpdateContestTeam(args.ctx, args.teamID, args.args).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "ErrUpdateContestTeam",
			args: args{
				ctx:    context.Background(),
				teamID: random.UUID(),
				args: &repository.UpdateContestTeamArgs{
					Name:        random.Optional(random.AlphaNumeric()),
					Result:      random.Optional(random.AlphaNumeric()),
					Link:        random.Optional(random.RandURLString()),
					Description: random.Optional(random.AlphaNumeric()),
				},
			},
			setup: func(f fields, args args) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().UpdateContestTeam(args.ctx, args.teamID, args.args).Return(repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			tt.fields = fields{
				repo: mock_repository.NewMockContestRepository(ctrl),
			}
			tt.setup(tt.fields, tt.args)
			s := NewContestService(tt.fields.repo)
			// Assertion
			tt.assertion(t, s.UpdateContestTeam(tt.args.ctx, tt.args.teamID, tt.args.args))
		})
	}
}

func Test_contestService_DeleteContestTeam(t *testing.T) {
	t.Parallel()
	type fields struct {
		repo repository.ContestRepository
	}
	type args struct {
		ctx       context.Context
		contestID uuid.UUID
		teamID    uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		setup     func(f fields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:       context.Background(),
				contestID: random.UUID(),
				teamID:    random.UUID(),
			},
			setup: func(f fields, args args) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().DeleteContestTeam(args.ctx, args.contestID, args.teamID).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "ErrDeleteContestTeam",
			args: args{
				ctx:       context.Background(),
				contestID: random.UUID(),
				teamID:    random.UUID(),
			},
			setup: func(f fields, args args) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().DeleteContestTeam(args.ctx, args.contestID, args.teamID).Return(repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			tt.fields = fields{
				repo: mock_repository.NewMockContestRepository(ctrl),
			}
			tt.setup(tt.fields, tt.args)
			s := NewContestService(tt.fields.repo)
			// Assertion
			tt.assertion(t, s.DeleteContestTeam(tt.args.ctx, tt.args.contestID, tt.args.teamID))
		})
	}
}

func TestContestService_GetContestTeamMembers(t *testing.T) {
	t.Parallel()
	type fields struct {
		repo repository.ContestRepository
	}
	type args struct {
		ctx       context.Context
		contestID uuid.UUID
		teamID    uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []*domain.User
		setup     func(f fields, args args, want []*domain.User)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:       context.Background(),
				contestID: random.UUID(),
				teamID:    random.UUID(),
			},
			want: []*domain.User{
				domain.NewUser(random.UUID(), random.AlphaNumeric(), random.AlphaNumeric(), random.Bool()),
			},
			setup: func(f fields, args args, want []*domain.User) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().GetContestTeamMembers(args.ctx, args.contestID, args.teamID).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "ErrGetContestTeamMembers",
			args: args{
				ctx:       context.Background(),
				contestID: random.UUID(),
				teamID:    random.UUID(),
			},
			want: nil,
			setup: func(f fields, args args, want []*domain.User) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().GetContestTeamMembers(args.ctx, args.contestID, args.teamID).Return(nil, repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			tt.fields = fields{
				repo: mock_repository.NewMockContestRepository(ctrl),
			}
			tt.setup(tt.fields, tt.args, tt.want)
			s := NewContestService(tt.fields.repo)
			// Assertion
			got, err := s.GetContestTeamMembers(tt.args.ctx, tt.args.contestID, tt.args.teamID)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestContestService_AddContestTeamMembers(t *testing.T) {
	t.Parallel()
	type fields struct {
		repo repository.ContestRepository
	}
	type args struct {
		ctx       context.Context
		teamID    uuid.UUID
		memberIDs []uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		setup     func(f fields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:       context.Background(),
				teamID:    random.UUID(),
				memberIDs: []uuid.UUID{random.UUID(), random.UUID()},
			},
			setup: func(f fields, args args) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().AddContestTeamMembers(args.ctx, args.teamID, args.memberIDs).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "ErrAddContestTeamMembers",
			args: args{
				ctx:       context.Background(),
				teamID:    random.UUID(),
				memberIDs: []uuid.UUID{random.UUID(), random.UUID()},
			},
			setup: func(f fields, args args) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().AddContestTeamMembers(args.ctx, args.teamID, args.memberIDs).Return(repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			tt.fields = fields{
				repo: mock_repository.NewMockContestRepository(ctrl),
			}
			tt.setup(tt.fields, tt.args)
			s := NewContestService(tt.fields.repo)
			// Assertion
			tt.assertion(t, s.AddContestTeamMembers(tt.args.ctx, tt.args.teamID, tt.args.memberIDs))
		})
	}
}

func TestContestService_EditContestTeamMembers(t *testing.T) {
	t.Parallel()
	type fields struct {
		repo repository.ContestRepository
	}
	type args struct {
		ctx       context.Context
		teamID    uuid.UUID
		memberIDs []uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		setup     func(f fields, args args)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx:       context.Background(),
				teamID:    random.UUID(),
				memberIDs: []uuid.UUID{random.UUID(), random.UUID()},
			},
			setup: func(f fields, args args) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().EditContestTeamMembers(args.ctx, args.teamID, args.memberIDs).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "ErrEditContestTeamMembers",
			args: args{
				ctx:       context.Background(),
				teamID:    random.UUID(),
				memberIDs: []uuid.UUID{random.UUID(), random.UUID()},
			},
			setup: func(f fields, args args) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().EditContestTeamMembers(args.ctx, args.teamID, args.memberIDs).Return(repository.ErrNotFound)
			},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			tt.fields = fields{
				repo: mock_repository.NewMockContestRepository(ctrl),
			}
			tt.setup(tt.fields, tt.args)
			s := NewContestService(tt.fields.repo)
			// Assertion
			tt.assertion(t, s.EditContestTeamMembers(tt.args.ctx, tt.args.teamID, tt.args.memberIDs))
		})
	}
}
