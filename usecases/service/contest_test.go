package service

import (
	"context"
	"math/rand"
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
					Name:      random.AlphaNumeric(rand.Intn(30) + 1),
					TimeStart: time.Now(),
					TimeEnd:   time.Now(),
				},
			},
			setup: func(f fields, args args, want []*domain.Contest) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().GetContests().Return(want, nil)
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
				repo.EXPECT().GetContests().Return(nil, repository.ErrForbidden)
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
					Name:      random.AlphaNumeric(rand.Intn(30) + 1),
					TimeStart: time.Now(),
					TimeEnd:   time.Now(),
				},
				Link:        random.RandURLString(),
				Description: random.AlphaNumeric(rand.Intn(30) + 1),
				Teams: []*domain.ContestTeam{
					{
						ID:        random.UUID(),
						ContestID: cid,
						Name:      random.AlphaNumeric(rand.Intn(30) + 1),
						Result:    random.AlphaNumeric(rand.Intn(30) + 1),
					},
				},
			},
			setup: func(f fields, args args, want *domain.ContestDetail) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().GetContest(args.id).Return(want, nil)
				repo.EXPECT().GetContestTeams(args.id).Return(want.Teams, nil)
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
					Name:      random.AlphaNumeric(rand.Intn(30) + 1),
					TimeStart: time.Now(),
					TimeEnd:   time.Now(),
				},
				Link:        random.RandURLString(),
				Description: random.AlphaNumeric(rand.Intn(30) + 1),
				Teams:       nil,
			},
			setup: func(f fields, args args, want *domain.ContestDetail) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().GetContest(args.id).Return(want, nil)
				repo.EXPECT().GetContestTeams(args.id).Return(nil, repository.ErrNotFound)
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
				repo.EXPECT().GetContest(args.id).Return(nil, repository.ErrNotFound)
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
				repo.EXPECT().GetContest(args.id).Return(want, nil)
				repo.EXPECT().GetContestTeams(args.id).Return(nil, repository.ErrInvalidID)
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
	cname := random.AlphaNumeric(rand.Intn(30) + 1) // 作成するコンテストのコンテスト名

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
		want      *domain.Contest
		setup     func(f fields, args args, want *domain.Contest)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				args: &repository.CreateContestArgs{
					Name:        cname,
					Description: random.AlphaNumeric(rand.Intn(30) + 1),
					Link:        optional.NewString(random.RandURLString(), true),
					Since:       time.Now(),
					Until:       optional.NewTime(time.Now(), true),
				},
			},
			want: &domain.Contest{
				ID:        random.UUID(),
				Name:      cname,
				TimeStart: time.Now(),
				TimeEnd:   time.Now(),
			},
			setup: func(f fields, args args, want *domain.Contest) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().CreateContest(args.args).Return(want, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "ErrCreate",
			args: args{
				ctx: context.Background(),
				args: &repository.CreateContestArgs{
					Name:        cname,
					Description: random.AlphaNumeric(rand.Intn(30) + 1),
					Link:        optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Since:       time.Now(),
					Until:       optional.NewTime(time.Now(), true),
				},
			},
			want: nil,
			setup: func(f fields, args args, want *domain.Contest) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().CreateContest(args.args).Return(nil, repository.ErrInvalidArg)
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
					Name:        optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Description: optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Link:        optional.NewString(random.RandURLString(), true),
					Since:       optional.NewTime(time.Now(), true),
					Until:       optional.NewTime(time.Now(), true),
				},
			},
			setup: func(f fields, args args) {
				changes := map[string]interface{}{
					"name":        args.args.Name.String,
					"description": args.args.Description.String,
					"link":        args.args.Link.String,
					"since":       args.args.Since.Time,
					"until":       args.args.Until.Time,
				}
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().UpdateContest(args.id, changes).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "ErrUpdate",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
				args: &repository.UpdateContestArgs{
					Name:        optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Description: optional.NewString(random.AlphaNumeric(rand.Intn(30)+1), true),
					Link:        optional.NewString(random.RandURLString(), true),
					Since:       optional.NewTime(time.Now(), true),
					Until:       optional.NewTime(time.Now(), true),
				},
			},
			setup: func(f fields, args args) {
				changes := map[string]interface{}{
					"name":        args.args.Name.String,
					"description": args.args.Description.String,
					"link":        args.args.Link.String,
					"since":       args.args.Since.Time,
					"until":       args.args.Until.Time,
				}
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().UpdateContest(args.id, changes).Return(repository.ErrNotFound)
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
				repo.EXPECT().DeleteContest(args.id).Return(nil)
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
				repo.EXPECT().DeleteContest(args.id).Return(repository.ErrNotFound)
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
					ID:        random.UUID(),
					ContestID: cid,
					Name:      random.AlphaNumeric(rand.Intn(30) + 1),
					Result:    random.AlphaNumeric(rand.Intn(30) + 1),
				},
			},
			setup: func(f fields, args args, want []*domain.ContestTeam) {
				repo := f.repo.(*mock_repository.MockContestRepository)
				repo.EXPECT().GetContestTeams(args.contestID).Return(want, nil)
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
				repo.EXPECT().GetContestTeams(args.contestID).Return(nil, repository.ErrNotFound)
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
