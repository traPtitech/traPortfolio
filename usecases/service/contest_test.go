package service

import (
	"context"
	"math/rand"
	"testing"
	"time"

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
