package service_test

import (
	"context"
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
	"github.com/traPtitech/traPortfolio/usecases/service"
)

func TestEventService_GetEvents(t *testing.T) {
	t.Parallel()
	type fields struct {
		event repository.EventRepository
		user  repository.UserRepository
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      []*domain.Event
		setup     func(f fields, args args, want []*domain.Event)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
			},
			want: []*domain.Event{
				{
					ID:        random.UUID(),
					Name:      random.AlphaNumeric(rand.Intn(30) + 1),
					TimeStart: time.Now(),
					TimeEnd:   time.Now(),
				},
			},
			setup: func(f fields, args args, want []*domain.Event) {
				e := f.event.(*mock_repository.MockEventRepository)
				e.EXPECT().GetEvents().Return(want, nil)
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			tt.fields = fields{
				event: mock_repository.NewMockEventRepository(ctrl),
				user:  mock_repository.NewMockUserRepository(ctrl),
			}
			tt.setup(tt.fields, tt.args, tt.want)
			s := service.NewEventService(tt.fields.event, tt.fields.user)
			// Assertion
			got, err := s.GetEvents(tt.args.ctx)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEventService_GetEventByID(t *testing.T) {
	t.Parallel()
	type fields struct {
		event repository.EventRepository
		user  repository.UserRepository
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      *domain.EventDetail
		setup     func(f fields, args args, want *domain.EventDetail)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
			},
			want: &domain.EventDetail{
				Event: domain.Event{
					// ID:
					Name:      random.AlphaNumeric(rand.Intn(30) + 1),
					TimeStart: time.Now(),
					TimeEnd:   time.Now(),
				},
				Description: random.AlphaNumeric(rand.Intn(30) + 1),
				Place:       random.AlphaNumeric(rand.Intn(30) + 1),
				Level:       domain.EventLevelAnonymous,
				HostName: []*domain.User{
					{
						ID:       random.UUID(),
						Name:     random.AlphaNumeric(rand.Intn(30) + 1),
						RealName: random.AlphaNumeric(rand.Intn(30) + 1),
					},
				},
				GroupID: random.UUID(),
				RoomID:  random.UUID(),
			},
			setup: func(f fields, args args, want *domain.EventDetail) {
				want.ID = args.id
				e := f.event.(*mock_repository.MockEventRepository)
				u := f.user.(*mock_repository.MockUserRepository)
				e.EXPECT().GetEvent(args.id).Return(&domain.EventDetail{
					Event: domain.Event{
						ID:        args.id,
						Name:      want.Name,
						TimeStart: want.TimeStart,
						TimeEnd:   want.TimeEnd,
					},
					Description: want.Description,
					Place:       want.Place,
					Level:       want.Level,
					HostName:    []*domain.User{{ID: want.HostName[0].ID}},
					GroupID:     want.GroupID,
					RoomID:      want.RoomID,
				}, nil)
				u.EXPECT().GetUsers(&repository.GetUsersArgs{}).Return(want.HostName, nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "KnoqForBidden",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
			},
			want: nil,
			setup: func(f fields, args args, want *domain.EventDetail) {
				e := f.event.(*mock_repository.MockEventRepository)
				e.EXPECT().GetEvent(args.id).Return(nil, repository.ErrForbidden)
			},
			assertion: assert.Error,
		},
		{
			name: "PortalForbidden",
			args: args{
				ctx: context.Background(),
				id:  random.UUID(),
			},
			want: nil,
			setup: func(f fields, args args, want *domain.EventDetail) {
				e := f.event.(*mock_repository.MockEventRepository)
				u := f.user.(*mock_repository.MockUserRepository)
				e.EXPECT().GetEvent(args.id).Return(&domain.EventDetail{
					Event: domain.Event{
						ID:        args.id,
						Name:      random.AlphaNumeric(rand.Intn(30) + 1),
						TimeStart: time.Now(),
						TimeEnd:   time.Now(),
					},
					Description: random.AlphaNumeric(rand.Intn(30) + 1),
					Place:       random.AlphaNumeric(rand.Intn(30) + 1),
					Level:       domain.EventLevelAnonymous,
					HostName:    []*domain.User{{ID: random.UUID()}},
					GroupID:     random.UUID(),
					RoomID:      random.UUID(),
				}, nil)
				u.EXPECT().GetUsers(&repository.GetUsersArgs{}).Return(nil, repository.ErrForbidden)
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
				event: mock_repository.NewMockEventRepository(ctrl),
				user:  mock_repository.NewMockUserRepository(ctrl),
			}
			tt.setup(tt.fields, tt.args, tt.want)
			s := service.NewEventService(tt.fields.event, tt.fields.user)
			// Assertion
			got, err := s.GetEventByID(tt.args.ctx, tt.args.id)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestEventService_UpdateEvent(t *testing.T) {
	t.Parallel()
	type fields struct {
		event repository.EventRepository
		user  repository.UserRepository
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
		arg *repository.UpdateEventLevelArg
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
				arg: &repository.UpdateEventLevelArg{
					Level: domain.EventLevelAnonymous,
				},
			},
			setup: func(f fields, args args) {
				e := f.event.(*mock_repository.MockEventRepository)
				e.EXPECT().UpdateEventLevel(args.id, args.arg).Return(nil)
			},
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// Setup mock
			ctrl := gomock.NewController(t)
			tt.fields = fields{
				event: mock_repository.NewMockEventRepository(ctrl),
				user:  mock_repository.NewMockUserRepository(ctrl),
			}
			tt.setup(tt.fields, tt.args)
			s := service.NewEventService(tt.fields.event, tt.fields.user)
			// Assertion
			tt.assertion(t, s.UpdateEventLevel(tt.args.ctx, tt.args.id, tt.args.arg))
		})
	}
}
