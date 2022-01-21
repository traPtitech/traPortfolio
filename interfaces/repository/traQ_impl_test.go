package repository

import (
	"context"
	"math/rand"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/traPtitech/traPortfolio/usecases/repository"
	"github.com/traPtitech/traPortfolio/util/random"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external"
)

func TestTraQRepository_GetUser(t *testing.T) {
	t.Parallel()
	type fields struct {
		api external.TraQAPI
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *domain.TraQUser
		setup     func(f fields, args args, want *domain.TraQUser)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
			},
			want: &domain.TraQUser{
				State:       domain.TraqStateActive,
				Bot:         false,
				DisplayName: random.AlphaNumeric(rand.Intn(30) + 1),
				Name:        random.AlphaNumeric(rand.Intn(30) + 1),
			},
			setup: func(f fields, args args, want *domain.TraQUser) {
				t := f.api.(*mock_external.MockTraQAPI)
				t.EXPECT().GetByID(args.id).Return(&external.TraQUserResponse{
					State:       want.State,
					Bot:         want.Bot,
					DisplayName: want.DisplayName,
					Name:        want.Name,
				}, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
			},
			want: nil,
			setup: func(f fields, args args, want *domain.TraQUser) {
				t := f.api.(*mock_external.MockTraQAPI)
				t.EXPECT().GetByID(args.id).Return(nil, repository.ErrNotFound)
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
				api: mock_external.NewMockTraQAPI(ctrl),
			}
			tt.setup(tt.fields, tt.args, tt.want)
			repo := NewTraQRepository(tt.fields.api)
			// Assertion
			got, err := repo.GetUser(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
