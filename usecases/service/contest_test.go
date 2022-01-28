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
