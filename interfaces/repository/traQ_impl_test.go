package repository

import (
	"context"
	"testing"

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
				id:  uuid.FromStringOrNil("11111111-1111-1111-1111-111111111111"),
			},
			want: &domain.TraQUser{
				State:       domain.TraqStateActive,
				Bot:         false,
				DisplayName: "user1",
				Name:        "user1",
			},
			setup:     func(f fields, args args, want *domain.TraQUser) {},
			assertion: assert.NoError,
		},
		{
			name: "NotFound",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
			},
			want:      nil,
			setup:     func(f fields, args args, want *domain.TraQUser) {},
			assertion: assert.Error,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			tt.fields = fields{
				api: mock_external.NewMockTraQAPI(),
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
