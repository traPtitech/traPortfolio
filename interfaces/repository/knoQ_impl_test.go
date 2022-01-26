package repository

import (
	"math/rand"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/traPtitech/traPortfolio/domain"
	"github.com/traPtitech/traPortfolio/interfaces/external"
	"github.com/traPtitech/traPortfolio/interfaces/external/mock_external"
	"github.com/traPtitech/traPortfolio/util/random"
)

func convertToEventResponse(t *testing.T, e *domain.KnoQEvent) *external.EventResponse {
	t.Helper()
	return &external.EventResponse{
		ID:          e.ID,
		Name:        e.Name,
		Description: e.Description,
		GroupID:     e.GroupID,
		RoomID:      e.RoomID,
		TimeStart:   e.TimeStart,
		TimeEnd:     e.TimeEnd,
		SharedRoom:  e.SharedRoom,
	}
}

func TestKnoqRepository_GetAll(t *testing.T) {
	t.Parallel()
	type fields struct {
		api external.KnoqAPI
	}
	tests := []struct {
		name      string
		fields    fields
		want      []*domain.KnoQEvent
		setup     func(t *testing.T, f fields, want []*domain.KnoQEvent)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			want: []*domain.KnoQEvent{
				{
					ID:          random.UUID(),
					Name:        random.AlphaNumeric(rand.Intn(30) + 1),
					Description: random.AlphaNumeric(rand.Intn(30) + 1),
					GroupID:     random.UUID(),
					RoomID:      random.UUID(),
					TimeStart:   random.Time(),
					TimeEnd:     random.Time(),
					SharedRoom:  true,
				},
				{
					ID:          random.UUID(),
					Name:        random.AlphaNumeric(rand.Intn(30) + 1),
					Description: random.AlphaNumeric(rand.Intn(30) + 1),
					GroupID:     random.UUID(),
					RoomID:      random.UUID(),
					TimeStart:   random.Time(),
					TimeEnd:     random.Time(),
					SharedRoom:  false,
				},
			},
			setup: func(t *testing.T, f fields, want []*domain.KnoQEvent) {
				k := f.api.(*mock_external.MockKnoqAPI)
				wantRes := make([]*external.EventResponse, len(want))
				for i, e := range want {
					wantRes[i] = convertToEventResponse(t, e)
				}
				k.EXPECT().GetAll().Return(wantRes, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			want: nil,
			setup: func(t *testing.T, f fields, want []*domain.KnoQEvent) {
				k := f.api.(*mock_external.MockKnoqAPI)
				k.EXPECT().GetAll().Return(nil, errUnexpected)
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
				api: mock_external.NewMockKnoqAPI(ctrl),
			}
			tt.setup(t, tt.fields, tt.want)
			repo := NewKnoqRepository(tt.fields.api)
			// Assertion
			got, err := repo.GetAll()
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestKnoqRepository_GetByID(t *testing.T) {
	t.Parallel()
	type fields struct {
		api external.KnoqAPI
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *domain.KnoQEvent
		setup     func(t *testing.T, f fields, args args, want *domain.KnoQEvent)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			want: &domain.KnoQEvent{
				ID:          random.UUID(),
				Name:        random.AlphaNumeric(rand.Intn(30) + 1),
				Description: random.AlphaNumeric(rand.Intn(30) + 1),
				GroupID:     random.UUID(),
				RoomID:      random.UUID(),
				TimeStart:   random.Time(),
				TimeEnd:     random.Time(),
				SharedRoom:  true,
			},
			setup: func(t *testing.T, f fields, args args, want *domain.KnoQEvent) {
				k := f.api.(*mock_external.MockKnoqAPI)
				k.EXPECT().GetByID(args.id).Return(convertToEventResponse(t, want), nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "UnexpectedError",
			want: nil,
			setup: func(t *testing.T, f fields, args args, want *domain.KnoQEvent) {
				k := f.api.(*mock_external.MockKnoqAPI)
				k.EXPECT().GetByID(args.id).Return(nil, errUnexpected)
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
				api: mock_external.NewMockKnoqAPI(ctrl),
			}
			tt.setup(t, tt.fields, tt.args, tt.want)
			repo := NewKnoqRepository(tt.fields.api)
			// Assertion
			got, err := repo.GetByID(tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
